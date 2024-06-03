ALTER TABLE public.activities RENAME to old_activities;

ALTER INDEX IF EXISTS activity_timestamp_index RENAME TO old_activity_timestamp_index
ALTER INDEX IF EXISTS activity_action_index RENAME TO old_activity_action_index
ALTER INDEX IF EXISTS activity_actor_index RENAME TO old_activity_actor_index
ALTER INDEX IF EXISTS activity_data_index RENAME TO old_activity_data_index
ALTER INDEX IF EXISTS activity_metadata_index RENAME TO old_activity_metadata_index