CREATE TABLE IF NOT EXISTS `group_assignment_chapter_items` (
  group_id VARCHAR(36) NOT NULL,
  chapter_id VARCHAR(36) NOT NULL,
  item_id INT NOT NULL,
  exercise_id_list VARCHAR(1024) DEFAULT NULL,
  full_mark INT NOT NULL DEFAULT '2',
  time_start VARCHAR(8) DEFAULT NULL,
  time_end VARCHAR(8) DEFAULT NULL,
  status ENUM('READY','CLOSED','STOP','OPEN') DEFAULT NULL,
  PRIMARY KEY (group_id, chapter_id, item_id),
  KEY chapter_id (chapter_id),
  CONSTRAINT fk_group_assignment_chapter_item FOREIGN KEY (group_id) REFERENCES class_schedules (group_id) ON DELETE CASCADE ON UPDATE CASCADE
);
