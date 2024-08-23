CREATE TABLE IF NOT EXISTS `group_assignment_exercises` (
  group_id VARCHAR(36) NOT NULL,
  exercise_id VARCHAR(36) NOT NULL,
  selected BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY KEY (group_id,exercise_id),
  KEY `group` (exercise_id),
  CONSTRAINT `fk_group_assignment_exercises_class_schedules` FOREIGN KEY (group_id) REFERENCES class_schedules (group_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_group_assignment_exercises_lab_exercises` FOREIGN KEY (exercise_id) REFERENCES lab_exercises (exercise_id) ON DELETE CASCADE ON UPDATE CASCADE
);


