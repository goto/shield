package postgres

import (
	"strings"
	"testing"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/user"
)

// These tests render the SQL for the generic metadata filter without a database,
// so they can run in unit mode (no Docker/Postgres required). The repository uses
// goqu in non-prepared mode, so literals are inlined (and escaped) by goqu.

func renderMetadataFilter(t *testing.T, projectID string, flt user.Filter) (string, error) {
	t.Helper()
	expr, err := metadataFilterExpr(projectID, flt)
	if err != nil {
		return "", err
	}
	if expr == nil {
		return "", nil
	}
	sql, _, sqlErr := dialect.From(goqu.T(TABLE_USERS)).
		Select(goqu.I("id")).
		Where(expr).
		ToSQL()
	if sqlErr != nil {
		t.Fatalf("unexpected ToSQL error: %s", sqlErr)
	}
	return sql, nil
}

func TestMetadataFilterExpr_NilWhenNothingToFilter(t *testing.T) {
	cases := []struct {
		name string
		proj string
		flt  user.Filter
	}{
		{name: "no project", proj: "", flt: user.Filter{MetadataPaths: []string{"functional_roles.roles"}, Metadatas: []string{"X"}}},
		{name: "no paths", proj: "p1", flt: user.Filter{Metadatas: []string{"X"}}},
		{name: "paths but no matchers", proj: "p1", flt: user.Filter{MetadataPaths: []string{"functional_roles.roles"}}},
		{name: "blank path", proj: "p1", flt: user.Filter{MetadataPaths: []string{"   "}, Metadatas: []string{"X"}}},
		// Regression guards: the shapes existing callers use (no metadata fields at
		// all) must produce no predicate, so the List query is unchanged.
		{name: "zero value filter", proj: "", flt: user.Filter{}},
		{name: "typical existing list (project + keyword, no metadata)", proj: "p1", flt: user.Filter{ProjectID: "p1", Keyword: "john", Limit: 10, Page: 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := metadataFilterExpr(tc.proj, tc.flt)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if expr != nil {
				t.Fatalf("expected nil expression, got %#v", expr)
			}
		})
	}
}

func TestMetadataFilterExpr_EqualsArrayContainment(t *testing.T) {
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths: []string{"functional_roles.roles"},
		Metadatas:     []string{"BIG_BOSS_SAMPLE1"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	for _, want := range []string{
		`EXISTS (SELECT 1 FROM "servicedata" AS "fr_sd"`,
		`JOIN "servicedata_keys" AS "fr_sk" ON "fr_sk"."id" = "fr_sd"."key_id"`,
		`"fr_sd"."entity_id" = CAST("users"."id" AS TEXT)`,
		`"fr_sd"."namespace_id" = 'shield/user'`,
		`"fr_sk"."project_id" = 'p1'`,
		`"fr_sk"."name" = 'functional_roles'`,
		// first path segment is the key, the rest is the jsonb text[] path
		`#> '{roles}'::text[]`,
		`jsonb_build_array`,
		`@> to_jsonb('BIG_BOSS_SAMPLE1'::text)`,
	} {
		if !strings.Contains(sql, want) {
			t.Fatalf("sql missing %q\nfull sql: %s", want, sql)
		}
	}
	if strings.Contains(sql, "NOT EXISTS") {
		t.Fatalf("positive-only filter should not produce NOT EXISTS: %s", sql)
	}
}

func TestMetadataFilterExpr_ScalarEqualsMatchesTextForm(t *testing.T) {
	// employee_details.terminated is a JSON boolean; equality must also compare as
	// text so `terminated=true` (stored boolean true) matches.
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths: []string{"employee_details.terminated"},
		Metadatas:     []string{"true"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	// array/string-scalar containment branch
	if !strings.Contains(sql, `@> to_jsonb('true'::text)`) {
		t.Fatalf("missing containment branch, sql: %s", sql)
	}
	// non-string scalar (bool/number) text-equality branch
	if !strings.Contains(sql, `COALESCE(NULLIF("fr_sd"."value" #>> '{terminated}'::text[], ''), 'null') = 'true'`) {
		t.Fatalf("missing text-equality branch, sql: %s", sql)
	}
}

func TestMetadataFilterExpr_KeyOnlyPathUsesEmptyJSONPath(t *testing.T) {
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths: []string{"employee_details"},
		Metadatas:     []string{"foo"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !strings.Contains(sql, `#> '{}'::text[]`) {
		t.Fatalf("key-only path should navigate the whole value with '{}', sql: %s", sql)
	}
	if !strings.Contains(sql, `"fr_sk"."name" = 'employee_details'`) {
		t.Fatalf("key not resolved from path, sql: %s", sql)
	}
}

func TestMetadataFilterExpr_DeepDynamicPath(t *testing.T) {
	// metadatakey.p1.p2.xxx -> key = metadatakey, json path = {p1,p2,xxx}.
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths: []string{"employee_details.p1.p2.xxx"},
		Metadatas:     []string{"VAL"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !strings.Contains(sql, `"fr_sk"."name" = 'employee_details'`) {
		t.Fatalf("first segment should be the servicedata key, sql: %s", sql)
	}
	if !strings.Contains(sql, `#> '{p1,p2,xxx}'::text[]`) {
		t.Fatalf("remaining segments should form an arbitrary-depth json path, sql: %s", sql)
	}
}

func TestMetadataFilterExpr_MultiPathOrsPositives(t *testing.T) {
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths: []string{"functional_roles.roles", "employee_details.title"},
		Metadatas:     []string{"X"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	// Two per-path EXISTS clauses OR-ed together.
	if strings.Count(sql, "EXISTS (SELECT 1") != 2 {
		t.Fatalf("expected 2 EXISTS clauses, sql: %s", sql)
	}
	if !strings.Contains(sql, ") OR (") && !strings.Contains(sql, " OR ") {
		t.Fatalf("expected OR between path clauses, sql: %s", sql)
	}
}

func TestMetadataFilterExpr_LikeVariants(t *testing.T) {
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths:      []string{"employee_details.title"},
		MetadataStartsWith: "Head",
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !strings.Contains(sql, `COALESCE(NULLIF("fr_sd"."value" #>> '{title}'::text[], ''), 'null') LIKE 'Head%'`) {
		t.Fatalf("starts_with LIKE not rendered, sql: %s", sql)
	}
}

func TestMetadataFilterExpr_NegativeUsesNotExistsAnded(t *testing.T) {
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths: []string{"functional_roles.roles"},
		NotMetadatas:  []string{"SECRET_ROLE"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if !strings.Contains(sql, "NOT EXISTS (SELECT 1") {
		t.Fatalf("negative filter should use NOT EXISTS, sql: %s", sql)
	}
	if !strings.Contains(sql, `@> to_jsonb('SECRET_ROLE'::text)`) {
		t.Fatalf("excluded value not rendered, sql: %s", sql)
	}
}

func TestMetadataFilterExpr_ContainsConflictErrors(t *testing.T) {
	_, err := metadataFilterExpr("p1", user.Filter{
		MetadataPaths:      []string{"employee_details.title"},
		MetadataStartsWith: "Head",
		MetadataContains:   "Engineer",
	})
	if err == nil {
		t.Fatal("expected validation error for contains + starts_with")
	}
	if !strings.Contains(err.Error(), "metadata_contains cannot be used together") {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestMetadataFilterExpr_EscapesLiterals(t *testing.T) {
	sql, err := renderMetadataFilter(t, "p1", user.Filter{
		MetadataPaths: []string{"functional_roles.roles"},
		Metadatas:     []string{"O'Brien"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	// goqu doubles the single quote inside the string literal (no SQL injection).
	if !strings.Contains(sql, `to_jsonb('O''Brien'::text)`) {
		t.Fatalf("value literal not escaped, sql: %s", sql)
	}
}
