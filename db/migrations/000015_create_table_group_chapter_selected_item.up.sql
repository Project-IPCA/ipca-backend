CREATE TABLE IF NOT EXISTS `group_chapter_selected_items` (
  group_id VARCHAR(36) NOT NULL,
  chapter_id VARCHAR(36) NOT NULL,
  item_id INT NOT NULL,
  exercise_id VARCHAR(36) NOT NULL,
  PRIMARY KEY (group_id,chapter_id, exercise_id),
  KEY group_id (group_id),
  KEY chapter_id (chapter_id),
  KEY exercise_id (exercise_id),
  CONSTRAINT `fk_group_chapter_selected_item_class_schedules` FOREIGN KEY (group_id) REFERENCES class_schedules (group_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_group_chapter_selected_item_lab_class_infos` FOREIGN KEY (chapter_id) REFERENCES lab_class_infos (chapter_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_group_chapter_selected_item_lab_exercises` FOREIGN KEY (exercise_id) REFERENCES lab_exercises (exercise_id) ON DELETE CASCADE ON UPDATE CASCADE
);