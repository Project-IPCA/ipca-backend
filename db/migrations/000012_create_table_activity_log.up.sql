CREATE TABLE IF NOT EXISTS `activity_logs` (
  log_id VARCHAR(26) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  group_id VARCHAR(36) DEFAULT NULL,
  username VARCHAR(30) NOT NULL,
  remote_ip VARCHAR(15) NOT NULL,
  remote_port INT DEFAULT NULL,
  agent VARCHAR(255) DEFAULT NULL,
  page_name VARCHAR(25) NOT NULL,
  action TEXT NOT NULL,
  ci INT UNSIGNED DEFAULT NULL,
  PRIMARY KEY (log_id)
) 

