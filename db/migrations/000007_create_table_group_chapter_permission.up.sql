CREATE TABLE IF NOT EXISTS `group_chapter_permissions` (
  class_id VARCHAR(36) NOT NULL,
  chapter_id VARCHAR(36) NOT NULL,
  allow_access_type ENUM('DENY','ALWAYS','TIMER','TIMER_PAUSED','DATETIME') NOT NULL DEFAULT 'DENY',
  access_time_start DATETIME DEFAULT NULL,
  access_time_end DATETIME DEFAULT NULL,
  allow_submit_type ENUM('DENY','ALWAYS','TIMER','TIMER_PAUSED','DATETIME') NOT NULL DEFAULT 'DENY',
  submit_time_start DATETIME DEFAULT NULL,
  submit_time_end DATETIME DEFAULT NULL,
  allow_submit BOOLEAN NOT NULL DEFAULT TRUE,
  status ENUM('NA','READY','OPEN','CLOSE','STOP') NOT NULL DEFAULT 'NA',
  allow_access BOOLEAN NOT NULL DEFAULT FALSE,
  time_start VARCHAR(8) DEFAULT NULL,
  time_end VARCHAR(8) DEFAULT NULL,
  PRIMARY KEY (class_id,chapter_id),
  KEY class_id (class_id),
  KEY chapter_id (chapter_id),
  CONSTRAINT fk_group_chapter_permissions_class_schedules FOREIGN KEY (class_id) REFERENCES class_schedules (group_id) ON DELETE CASCADE ON UPDATE CASCADE
);
