CREATE INDEX activity_timestamp_index ON activities (timestamp);
CREATE INDEX activity_action_index ON activities (action);
CREATE INDEX activity_actor_index ON activities (actor);
CREATE INDEX activity_data_index ON activities USING gin (data jsonb_path_ops);
CREATE INDEX activity_metadata_index ON activities USING gin (metadata jsonb_path_ops);