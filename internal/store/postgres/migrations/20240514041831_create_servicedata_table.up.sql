CREATE TABLE IF NOT EXISTS servicedata
(
    id              uuid        PRIMARY KEY     DEFAULT uuid_generate_v4(),
    namespace_id    varchar,
    entity_id       varchar,
    key_id          uuid        REFERENCES servicedata_keys(id),
    value           varchar, 
    created_at      timestamptz NOT NULL        DEFAULT NOW(),
    updated_at      timestamptz NOT NULL        DEFAULT NOW(),
    UNIQUE (namespace_id, entity_id)
);
