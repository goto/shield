-- Support object-scoped relation lookups without scanning the subject-first index.
-- Keep this as a single statement: CONCURRENTLY cannot run in a transaction.
CREATE INDEX CONCURRENTLY relations_object_id_namespace_id_idx
    ON public.relations (object_id, object_namespace_id);
