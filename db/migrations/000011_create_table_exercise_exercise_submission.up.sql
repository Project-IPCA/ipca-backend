CREATE TABLE IF NOT EXISTS `exercise_submissions` (
  submission_id VARCHAR(36) NOT NULL,
  stu_id VARCHAR(36) NOT NULL,
  exercise_id VARCHAR(36) NOT NULL,
  status ENUM('ACCEPTED','WRONG_ANSWER','PENDING','REJECTED','ERROR') NOT NULL DEFAULT 'PENDING',
  sourcecode_filename VARCHAR(40) NOT NULL,
  marking INT NOT NULL DEFAULT '-1',
  time_submit DATETIME DEFAULT CURRENT_TIMESTAMP,
  is_inf_loop BOOLEAN DEFAULT NULL,
  output TEXT DEFAULT NULL,
  result JSON DEFAULT NULL,
  error_message MEDIUMTEXT,
  PRIMARY KEY (submission_id),
  KEY stu_id (stu_id),
  KEY exercise_id (exercise_id),
  CONSTRAINT fk_exercise_submissions_students FOREIGN KEY (stu_id) REFERENCES students (stu_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_exercise_submissions_lab_exercises FOREIGN KEY (exercise_id) REFERENCES lab_exercises (exercise_id) ON DELETE RESTRICT ON UPDATE RESTRICT
); 
