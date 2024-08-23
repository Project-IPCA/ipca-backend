CREATE TABLE IF NOT EXISTS `lab_exercises` (
  exercise_id VARCHAR(36) NOT NULL,
  chapter_id VARCHAR(36) DEFAULT NULL,
  level ENUM('0','1','2','3','4','5','6') DEFAULT NULL,
  name VARCHAR(1024) DEFAULT NULL,
  content MEDIUMTEXT,
  testcase ENUM('NO_INPUT','YES','UNDEFINED') NOT NULL DEFAULT 'NO_INPUT',
  sourcecode VARCHAR(50) DEFAULT NULL,
  full_mark INT NOT NULL DEFAULT '10',
  added_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_update DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  user_defined_constraints JSON DEFAULT NULL,
  suggested_constraints JSON DEFAULT NULL,
  added_by VARCHAR(40) DEFAULT NULL,
  created_by VARCHAR(36) DEFAULT NULL,
  PRIMARY KEY (exercise_id),
  KEY created_by (created_by),
  KEY chapter_id (chapter_id),
  CONSTRAINT fk_lab_class_infos_lab_exercises FOREIGN KEY (chapter_id) REFERENCES lab_class_infos (chapter_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_lab_class_infos_supervisors FOREIGN KEY (created_by) REFERENCES supervisors (supervisor_id) ON DELETE CASCADE ON UPDATE CASCADE
) 

