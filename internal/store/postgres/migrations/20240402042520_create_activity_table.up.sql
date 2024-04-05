CREATE TABLE IF NOT EXISTS activities (
    timestamp   timestamptz,
    action      varchar,
    actor       uuid,
    data        jsonb,
    metadata    jsonb
);