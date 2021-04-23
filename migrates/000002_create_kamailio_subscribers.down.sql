DROP TABLE IF EXISTS subscriber;
DROP INDEX IF EXISTS subscriber_username_idx;
DELETE FROM version WHERE table_name = 'subscriber';