CREATE TABLE IF NOT EXISTS `exercise_testcases` (
  testcase_id VARCHAR(36) NOT NULL,
  exercise_id VARCHAR(36) NOT NULL,
  is_ready VARCHAR(3) NOT NULL DEFAULT 'yes',
  testcase_content VARCHAR(1024) NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  is_show_student BOOLEAN DEFAULT TRUE,
  testcase_note VARCHAR(1024) DEFAULT NULL,
  testcase_output MEDIUMTEXT,
  testcase_error VARCHAR(4096) DEFAULT NULL,
  PRIMARY KEY (testcase_id),
  KEY exercise_id (exercise_id),
  CONSTRAINT fk_exercise_testcases_lab_exercises FOREIGN KEY (exercise_id) REFERENCES lab_exercises (exercise_id) ON DELETE CASCADE ON UPDATE CASCADE
);  
