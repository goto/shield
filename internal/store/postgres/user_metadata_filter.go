package postgres

import (
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
)

// metadataFilterExpr builds a correlated predicate that filters users by their
// servicedata metadata, mirroring Guardian's appeal_details* generic JSONB filter
// (guardian internal/store/postgres/utils.go applyJSONBPathsLikeAndInFilter).
//
// Shield stores each metadata key as a separate servicedata row, so a path's first
// segment selects the servicedata key (row) and the remaining segments navigate
// into that row's JSON value, e.g. "functional_roles.roles". The predicate
// correlates on "users"."id", so it can be applied to any query whose FROM is the
// unaliased users table (the base list query and the servicedata-join pagination
// subquery).
//
// Positive filters (Metadatas / MetadataStartsWith / MetadataEndsWith /
// MetadataContains) are OR-ed together across paths. Negative filters
// (NotMetadatas / MetadataNot*) are AND-ed together. Equality/IN uses
// array-or-scalar containment, so a JSON array value matches when it contains the
// value; LIKE filters compare the value as text (COALESCE'd to the literal 'null'
// when the path is absent), matching guardian. Returns nil when nothing to filter.
func metadataFilterExpr(projectID string, flt user.Filter) (goqu.Expression, error) {
	if projectID == "" {
		return nil, nil
	}

	type pathParts struct {
		key     string
		subpath string // postgres text[] literal, e.g. "{roles}" or "{}"
	}
	seen := make(map[string]struct{})
	var paths []pathParts
	for _, p := range flt.MetadataPaths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		segments := strings.Split(p, ".")
		key := segments[0]
		if key == "" {
			continue
		}
		paths = append(paths, pathParts{
			key:     key,
			subpath: "{" + strings.Join(segments[1:], ",") + "}",
		})
	}
	if len(paths) == 0 {
		return nil, nil
	}

	hasPositive := len(flt.Metadatas) > 0 || flt.MetadataStartsWith != "" ||
		flt.MetadataEndsWith != "" || flt.MetadataContains != ""
	hasNegative := len(flt.NotMetadatas) > 0 || flt.MetadataNotStartsWith != "" ||
		flt.MetadataNotEndsWith != "" || flt.MetadataNotContains != ""
	if !hasPositive && !hasNegative {
		return nil, nil
	}

	// Validation mirrors guardian: contains cannot combine with starts/ends.
	if (flt.MetadataStartsWith != "" || flt.MetadataEndsWith != "") && flt.MetadataContains != "" {
		return nil, fmt.Errorf("invalid filter: metadata_contains cannot be used together with metadata_starts_with or metadata_ends_with")
	}
	if (flt.MetadataNotStartsWith != "" || flt.MetadataNotEndsWith != "") && flt.MetadataNotContains != "" {
		return nil, fmt.Errorf("invalid filter: metadata_not_contains cannot be used together with metadata_not_starts_with or metadata_not_ends_with")
	}

	var conditions []goqu.Expression

	// ---- POSITIVE (OR across paths) ----
	if hasPositive {
		var positives []goqu.Expression
		for _, p := range paths {
			var matchers []goqu.Expression
			if flt.MetadataStartsWith != "" {
				matchers = append(matchers, jsonbTextLikeExpr(p.subpath, flt.MetadataStartsWith+"%"))
			}
			if flt.MetadataEndsWith != "" {
				matchers = append(matchers, jsonbTextLikeExpr(p.subpath, "%"+flt.MetadataEndsWith))
			}
			if flt.MetadataContains != "" {
				matchers = append(matchers, jsonbTextLikeExpr(p.subpath, "%"+flt.MetadataContains+"%"))
			}
			for _, v := range flt.Metadatas {
				matchers = append(matchers, jsonbContainsExpr(p.subpath, v))
			}
			if len(matchers) > 0 {
				positives = append(positives, metadataKeyExistsExpr(projectID, p.key, goqu.Or(matchers...), false))
			}
		}
		if len(positives) > 0 {
			conditions = append(conditions, goqu.Or(positives...))
		}
	}

	// ---- NEGATIVE (AND across paths/matchers) ----
	if hasNegative {
		for _, p := range paths {
			if flt.MetadataNotStartsWith != "" {
				conditions = append(conditions, metadataKeyExistsExpr(projectID, p.key, jsonbTextLikeExpr(p.subpath, flt.MetadataNotStartsWith+"%"), true))
			}
			if flt.MetadataNotEndsWith != "" {
				conditions = append(conditions, metadataKeyExistsExpr(projectID, p.key, jsonbTextLikeExpr(p.subpath, "%"+flt.MetadataNotEndsWith), true))
			}
			if flt.MetadataNotContains != "" {
				conditions = append(conditions, metadataKeyExistsExpr(projectID, p.key, jsonbTextLikeExpr(p.subpath, "%"+flt.MetadataNotContains+"%"), true))
			}
			for _, v := range flt.NotMetadatas {
				conditions = append(conditions, metadataKeyExistsExpr(projectID, p.key, jsonbContainsExpr(p.subpath, v), true))
			}
		}
	}

	if len(conditions) == 0 {
		return nil, nil
	}
	return goqu.And(conditions...), nil
}

// metadataKeyExistsExpr wraps an inner predicate in a correlated (NOT) EXISTS over
// the servicedata row for a given metadata key, scoped to the project and the user
// principal namespace.
func metadataKeyExistsExpr(projectID, key string, inner goqu.Expression, negate bool) goqu.Expression {
	keyword := "EXISTS"
	if negate {
		keyword = "NOT EXISTS"
	}
	return goqu.L(
		keyword+` (SELECT 1 FROM "`+TABLE_SERVICE_DATA+`" AS "fr_sd" `+
			`JOIN "`+TABLE_SERVICE_DATA_KEYS+`" AS "fr_sk" ON "fr_sk"."id" = "fr_sd"."key_id" `+
			`WHERE "fr_sd"."entity_id" = CAST("users"."id" AS TEXT) `+
			`AND "fr_sd"."namespace_id" = ? AND "fr_sk"."project_id" = ? AND "fr_sk"."name" = ? AND ?)`,
		schema.UserPrincipal, projectID, key, inner,
	)
}

// jsonbTextLikeExpr compares the value at subpath as text (COALESCE'd to 'null'
// when absent) against a LIKE pattern, matching guardian's buildJSONTextExpr.
func jsonbTextLikeExpr(subpath, pattern string) goqu.Expression {
	return goqu.L(`COALESCE(NULLIF("fr_sd"."value" #>> ?::text[], ''), 'null') LIKE ?`, subpath, pattern)
}

// jsonbContainsExpr matches when the value at subpath equals value, or, when it is
// a JSON array, contains value as an element. Scalars are wrapped in a one-element
// array so a single containment check covers both shapes.
func jsonbContainsExpr(subpath, value string) goqu.Expression {
	return goqu.L(
		`(CASE WHEN jsonb_typeof("fr_sd"."value" #> ?::text[]) = 'array' `+
			`THEN "fr_sd"."value" #> ?::text[] `+
			`ELSE jsonb_build_array("fr_sd"."value" #> ?::text[]) END) @> to_jsonb(?::text)`,
		subpath, subpath, subpath, value,
	)
}
