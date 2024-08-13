CREATE TABLE IF NOT EXISTS `lab_class_infos` (
  chapter_id VARCHAR(36) NOT NULL,
  chapter_name VARCHAR(256) NOT NULL,
  chapter_fullmark INT NOT NULL,
  no_items INT NOT NULL DEFAULT '5',
  PRIMARY KEY (chapter_id)
);

