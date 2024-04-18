CREATE TABLE IF NOT EXISTS activities (
    timestamp   timestamptz,
    action      varchar,
    actor       varchar,
    data        jsonb,
    metadata    jsonb
);