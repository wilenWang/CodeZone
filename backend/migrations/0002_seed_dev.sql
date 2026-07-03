INSERT INTO workspaces (id, name)
VALUES (1, 'Default')
ON DUPLICATE KEY UPDATE name = VALUES(name);

INSERT INTO users (id, workspace_id, username, display_name, avatar_url, user_type, password_hash)
VALUES
  (1, 1, 'alice', 'Alice Chen', NULL, 'human', NULL),
  (2, 1, 'bob', 'Bob Lin', NULL, 'human', NULL),
  (3, 1, 'carol', 'Carol Wu', NULL, 'human', NULL),
  (10, 1, 'mock-agent', 'Mock Agent', NULL, 'agent', NULL)
ON DUPLICATE KEY UPDATE
  display_name = VALUES(display_name),
  user_type = VALUES(user_type);

INSERT INTO agent_profiles (user_id, kind, config_json, enabled)
VALUES (10, 'mock', JSON_OBJECT('replyPrefix', 'Mock Agent received:'), TRUE)
ON DUPLICATE KEY UPDATE
  config_json = VALUES(config_json),
  enabled = VALUES(enabled);

INSERT INTO conversations (id, workspace_id, type, title, created_by)
VALUES
  (1, 1, 'direct', NULL, 1),
  (2, 1, 'group', 'Project Room', 1),
  (3, 1, 'direct', NULL, 1)
ON DUPLICATE KEY UPDATE title = VALUES(title);

INSERT INTO conversation_members (conversation_id, user_id, role)
VALUES
  (1, 1, 'owner'),
  (1, 2, 'member'),
  (2, 1, 'owner'),
  (2, 2, 'member'),
  (2, 3, 'member'),
  (3, 1, 'owner'),
  (3, 10, 'member')
ON DUPLICATE KEY UPDATE role = VALUES(role);
