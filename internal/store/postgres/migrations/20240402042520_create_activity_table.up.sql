CREATE TABLE IF NOT EXISTS activities (
    timestamp   timestamptz NOT NULL,
    action      varchar     NOT NULL,
    actor       varchar,
    data        jsonb       NOT NULL,
    metadata    jsonb       NOT NULL
);

CREATE INDEX activity_timestamp_index ON activities (timestamp);
CREATE INDEX activity_action_index ON activities (action);
CREATE INDEX activity_actor_index ON activities (actor);
CREATE INDEX activity_data_index ON activities (data);
CREATE INDEX activity_metadata_index ON activities (metadata);