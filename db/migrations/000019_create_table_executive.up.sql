CREATE TABLE IF NOT EXISTS `executives` (
  executive_id VARCHAR(36) NOT NULL,
  PRIMARY KEY (executive_id),
  CONSTRAINT fk_user_executive_user1 FOREIGN KEY (executive_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
);