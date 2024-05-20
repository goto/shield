CREATE TABLE IF NOT EXISTS servicedata_keys
(
    id              uuid        PRIMARY KEY     DEFAULT uuid_generate_v4(),
    urn             varchar     UNIQUE,
    project_id      uuid REFERENCES projects(id),
    key             varchar,
    description     varchar,
    resource_id     uuid REFERENCES resources(id),
    created_at      timestamptz NOT NULL        DEFAULT NOW(),
    updated_at      timestamptz NOT NULL        DEFAULT NOW(),
    deleted_at      timestamptz
);
