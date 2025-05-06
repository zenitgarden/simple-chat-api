DROP INDEX idx_conversations_created_by;

ALTER TABLE conversations
DROP CONSTRAINT fk_conversations_created_by;

ALTER TABLE conversations
DROP COLUMN created_by;