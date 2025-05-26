--DROP TABLE IF EXISTS schema_migrations;


UPDATE schema_migrations SET dirty = false;
SELECT * FROM schema_migrations;