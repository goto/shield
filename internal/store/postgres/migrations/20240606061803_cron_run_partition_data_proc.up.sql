SELECT cron.schedule('call-partition_data_proc', '* 21 * * *', 'CALL partman.partition_data_proc(''public.activities'')');