-- Match GetByEmail's case-insensitive lookup while excluding soft-deleted users.
-- Keep this as a single statement: CONCURRENTLY cannot run in a transaction.
CREATE INDEX CONCURRENTLY users_email_lower_active_idx
    ON public.users (LOWER(email))
    WHERE deleted_at IS NULL;
