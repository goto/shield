BEGIN TRANSACTION;

  CALL partman.partition_data_proc('public.activities', p_interval := '1 week', p_batch := 100, p_source_table := 'public.old_activities');

COMMIT;