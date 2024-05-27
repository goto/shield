ALTER TABLE servicedata
ADD CONSTRAINT servicedata_namespace_id_entity_id_key_id_key UNIQUE (namespace_id, entity_id, key_id);