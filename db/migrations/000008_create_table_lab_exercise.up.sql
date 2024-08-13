CREATE TABLE IF NOT EXISTS `lab_exercises` (
  exercise_id VARCHAR(36) NOT NULL,
  lab_chapter INT DEFAULT NULL,
  lab_level ENUM('0','1','2','3','4','5','6') DEFAULT NULL,
  lab_name VARCHAR(1024) DEFAULT NULL,
  lab_content MEDIUMTEXT,
  testcase ENUM('NO_INPUT','YES','UNDEFINED') NOT NULL DEFAULT 'NO_INPUT',
  sourcecode VARCHAR(50) DEFAULT NULL,
  full_mark INT NOT NULL DEFAULT '10',
  added_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_update DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
  user_defined_constraints JSON DEFAULT NULL,
  suggested_constraints JSON DEFAULT NULL,
  added_by VARCHAR(40) DEFAULT NULL,
  created_by INT DEFAULT NULL,
  PRIMARY KEY (exercise_id),
  KEY created_by (created_by),
  KEY lab_chapter (lab_chapter)
) 

