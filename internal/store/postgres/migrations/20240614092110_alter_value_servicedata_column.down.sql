ALTER TABLE servicedata ALTER COLUMN value TYPE varchar USING (value::varchar);

UPDATE servicedata SET value = SUBSTRING(value, 2, LENGTH(value) - 2) WHERE value LIKE '"%"' AND (value NOT LIKE '[%' AND value NOT LIKE '{%');
