ALTER TABLE conversations
ADD COLUMN created_by UUID NOT NULL;

ALTER TABLE conversations
ADD CONSTRAINT fk_conversations_created_by
FOREIGN KEY (created_by) REFERENCES users(id)
ON DELETE RESTRICT;

CREATE INDEX idx_conversations_created_by ON conversations(created_by);
