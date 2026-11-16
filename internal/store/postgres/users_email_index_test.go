package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/goto/salt/log"
	"github.com/stretchr/testify/require"

	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/internal/store/postgres/migrations"
)

func TestUsersEmailLowerActiveIndexMigration(t *testing.T) {
	client, pool, resource, err := newTestClient(log.NewZap())
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, client.Close())
		require.NoError(t, purgeDocker(pool, resource))
	})
	ctx := context.Background()

	// Use enough synthetic rows for the planner to prefer a selective index.
	_, err = client.ExecContext(ctx, `
		INSERT INTO users (id, name, email)
		SELECT md5(n::text)::uuid, 'lookup-fixture-' || n,
		       'lookup-fixture-' || n || '@example.invalid'
		FROM generate_series(1, 30000) AS n;
		INSERT INTO users (id, name, email, deleted_at) VALUES
		('00000000-0000-0000-0000-000000000001', 'lookup-active', 'Mixed.Case@example.invalid', NULL),
		('00000000-0000-0000-0000-000000000002', 'lookup-deleted', 'MIXED.CASE@example.invalid', NOW());
		ANALYZE users;
	`)
	require.NoError(t, err)

	assertLookup := func() {
		t.Helper()
		found, err := postgres.NewUserRepository(client).GetByEmail(ctx, "mIxEd.CaSe@example.invalid")
		require.NoError(t, err)
		require.Equal(t, "00000000-0000-0000-0000-000000000001", found.ID)
	}
	assertValidIndex := func() {
		t.Helper()
		var valid bool
		err := client.GetContext(ctx, &valid, `
			SELECT indisvalid AND indisready FROM pg_index
			WHERE indexrelid = 'public.users_email_lower_active_idx'::regclass
		`)
		require.NoError(t, err)
		require.True(t, valid)
	}
	assertIndexPlan := func() {
		t.Helper()
		var raw []byte
		err := client.GetContext(ctx, &raw, `EXPLAIN (FORMAT JSON)
			SELECT * FROM users
			WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL`, "mIxEd.CaSe@example.invalid")
		require.NoError(t, err)
		var plans []struct {
			Plan emailLookupPlan `json:"Plan"`
		}
		require.NoError(t, json.Unmarshal(raw, &plans))
		require.Len(t, plans, 1)
		require.True(t, plans[0].Plan.usesIndex("users_email_lower_active_idx"), string(raw))
	}

	// newTestClient applies the embedded migrations through the real migrator.
	assertValidIndex()
	assertLookup()
	assertIndexPlan()

	down, err := migrations.MigrationFs.ReadFile("20261113110000_add_users_email_lower_active_index.down.sql")
	require.NoError(t, err)
	_, err = client.ExecContext(ctx, string(down))
	require.NoError(t, err)
	var absent bool
	require.NoError(t, client.GetContext(ctx, &absent, `SELECT to_regclass('public.users_email_lower_active_idx') IS NULL`))
	require.True(t, absent)
	assertLookup()

	up, err := migrations.MigrationFs.ReadFile("20261113110000_add_users_email_lower_active_index.up.sql")
	require.NoError(t, err)
	_, err = client.ExecContext(ctx, string(up))
	require.NoError(t, err)
	assertValidIndex()
	assertLookup()
	assertIndexPlan()
}

type emailLookupPlan struct {
	IndexName string            `json:"Index Name"`
	Plans     []emailLookupPlan `json:"Plans"`
}

func (p emailLookupPlan) usesIndex(name string) bool {
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
