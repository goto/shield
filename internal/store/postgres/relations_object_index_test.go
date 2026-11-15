package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/goto/salt/log"
	"github.com/stretchr/testify/require"

	"github.com/goto/shield/internal/store/postgres/migrations"
)

func TestRelationsObjectIndexMigration(t *testing.T) {
	client, pool, resource, err := newTestClient(log.NewNoop())
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.Close())
		require.NoError(t, purgeDocker(pool, resource))
	})
	ctx := context.Background()
	_, err = client.ExecContext(ctx, `
		INSERT INTO namespaces (id, name) VALUES
		('index-fixture-a', 'index-fixture-a'), ('index-fixture-b', 'index-fixture-b');
		INSERT INTO relations (subject_id, object_id, object_namespace_id)
		SELECT 'subject-' || n, 'object-' || n, 'index-fixture-a'
		FROM generate_series(1, 30000) AS n;
		INSERT INTO relations (subject_id, object_id, object_namespace_id) VALUES
		('target-a', 'target', 'index-fixture-a'),
		('target-b', 'target', 'index-fixture-a'),
		('other-namespace', 'target', 'index-fixture-b');
		ANALYZE relations;
	`)
	require.NoError(t, err)

	assertLookup := func() {
		t.Helper()
		var subjects []string
		require.NoError(t, client.SelectContext(ctx, &subjects, `
			SELECT subject_id FROM relations WHERE object_id = $1 AND object_namespace_id = $2`,
			"target", "index-fixture-a"))
		require.ElementsMatch(t, []string{"target-a", "target-b"}, subjects)
	}
	assertIndex := func() {
		t.Helper()
		var valid bool
		require.NoError(t, client.GetContext(ctx, &valid, `
			SELECT indisvalid AND indisready FROM pg_index
			WHERE indexrelid = 'public.relations_object_id_namespace_id_idx'::regclass`))
		require.True(t, valid)
		var raw []byte
		require.NoError(t, client.GetContext(ctx, &raw, `EXPLAIN (FORMAT JSON)
			SELECT created_at, id, object_id, object_namespace_id, role_id,
			       subject_id, subject_namespace_id, updated_at
			FROM relations WHERE object_id = $1 AND object_namespace_id = $2`,
			"target", "index-fixture-a"))
		var plans []struct {
			Plan relationLookupPlan `json:"Plan"`
		}
		require.NoError(t, json.Unmarshal(raw, &plans))
		require.Len(t, plans, 1)
		require.True(t, plans[0].Plan.usesIndex("relations_object_id_namespace_id_idx"), string(raw))
	}

	// Exercise the embedded migrator, rollback and reapply on a disposable database.
	assertLookup()
	assertIndex()
	down, err := migrations.MigrationFs.ReadFile("20261116060000_add_relations_object_index.down.sql")
	require.NoError(t, err)
	_, err = client.ExecContext(ctx, string(down))
	require.NoError(t, err)
	var absent bool
	require.NoError(t, client.GetContext(ctx, &absent, `SELECT to_regclass('public.relations_object_id_namespace_id_idx') IS NULL`))
	require.True(t, absent)
	assertLookup()
	up, err := migrations.MigrationFs.ReadFile("20261116060000_add_relations_object_index.up.sql")
	require.NoError(t, err)
	_, err = client.ExecContext(ctx, string(up))
	require.NoError(t, err)
	assertLookup()
	assertIndex()
}

type relationLookupPlan struct {
	IndexName string               `json:"Index Name"`
	Plans     []relationLookupPlan `json:"Plans"`
}

func (p relationLookupPlan) usesIndex(name string) bool {
	if p.IndexName == name {
		return true
	}
	for _, child := range p.Plans {
		if child.usesIndex(name) {
			return true
		}
	}
	return false
}
