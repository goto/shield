ALTER TABLE public.old_activities RENAME TO activities;

ALTER INDEX IF EXISTS old_activity_timestamp_index RENAME TO activity_timestamp_index;
ALTER INDEX IF EXISTS old_activity_action_index RENAME TO activity_action_index;
ALTER INDEX IF EXISTS old_activity_actor_index RENAME TO activity_actor_index;
ALTER INDEX IF EXISTS old_activity_data_index RENAME TO activity_data_index;
ALTER INDEX IF EXISTS old_activity_metadata_index RENAME TO activity_metadata_index;
