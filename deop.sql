--DROP TABLE IF EXISTS schema_migrations;
--UPDATE schema_migrations SET dirty = false;
--SELECT * FROM schema_migrations;

update users set role = 'admin' where phone='+998111231234';