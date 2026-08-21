SET @sql := IF(
  (SELECT DATA_TYPE FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'avatar_url') <> 'text',
  'ALTER TABLE users MODIFY COLUMN avatar_url TEXT NULL',
  'DO 0'
);
PREPARE migration_stmt FROM @sql;
EXECUTE migration_stmt;
DEALLOCATE PREPARE migration_stmt;
