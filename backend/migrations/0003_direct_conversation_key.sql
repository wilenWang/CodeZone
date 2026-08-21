SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'conversations' AND COLUMN_NAME = 'direct_pair_key') = 0,
  'ALTER TABLE conversations ADD COLUMN direct_pair_key VARCHAR(50) NULL',
  'DO 0'
);
PREPARE migration_stmt FROM @sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;

UPDATE conversations c
JOIN (
  SELECT
    cm.conversation_id,
    CONCAT(LEAST(MIN(cm.user_id), MAX(cm.user_id)), ':', GREATEST(MIN(cm.user_id), MAX(cm.user_id))) AS direct_pair_key
  FROM conversation_members cm
  JOIN conversations direct_conversation ON direct_conversation.id = cm.conversation_id
  WHERE direct_conversation.type = 'direct'
  GROUP BY cm.conversation_id
  HAVING COUNT(*) = 2
) pairs ON pairs.conversation_id = c.id
SET c.direct_pair_key = pairs.direct_pair_key
WHERE c.type = 'direct';

SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'conversations' AND INDEX_NAME = 'conversations_workspace_direct_pair_uq') = 0,
  'ALTER TABLE conversations ADD UNIQUE KEY conversations_workspace_direct_pair_uq (workspace_id, direct_pair_key)',
  'DO 0'
);
PREPARE migration_stmt FROM @sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;
