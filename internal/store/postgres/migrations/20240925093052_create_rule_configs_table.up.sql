CREATE TABLE IF NOT EXISTS rule_configs 
(
    id          BIGSERIAL       NOT NULL PRIMARY KEY,
    name        varchar         NOT NULL UNIQUE,
    config      jsonb           NOT NULL,
    created_at  timestamptz     NOT NULL    DEFAULT NOW(),
    updated_at  timestamptz     NOT NULL    DEFAULT NOW()
)