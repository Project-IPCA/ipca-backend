CREATE TABLE IF NOT EXISTS `student_assignment_chapter_items` (
  stu_id VARCHAR(36) NOT NULL,
  chapter_id VARCHAR(36) NOT NULL,
  item_id VARCHAR(36) NOT NULL,
  exercise_id VARCHAR(36) DEFAULT NULL,
  full_mark INT NOT NULL DEFAULT '0',
  marking INT NOT NULL DEFAULT '0',
  added_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  time_start DATETIME DEFAULT NULL,
  time_end DATETIME DEFAULT NULL,
  PRIMARY KEY (stu_id,chapter_id,item_id),
  KEY exercise_id (exercise_id),
  KEY chapter_id (chapter_id),
  CONSTRAINT fk_student_assignment_chapter_items_students FOREIGN KEY (stu_id) REFERENCES students (stu_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_student_assignment_chapter_items_lab_exercises  FOREIGN KEY (exercise_id) REFERENCES lab_exercises (exercise_id) ON DELETE RESTRICT ON UPDATE RESTRICT
);

