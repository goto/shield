-- Speed up resolving a service data key by its name (names can repeat across projects).
CREATE INDEX IF NOT EXISTS servicedata_keys_name_index ON servicedata_keys (name);

-- Narrow servicedata rows to a given key within a namespace (e.g. the "shield/user" namespace)
-- before extracting distinct jsonb values.
CREATE INDEX IF NOT EXISTS servicedata_key_id_namespace_id_index ON servicedata (key_id, namespace_id);
