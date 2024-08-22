CREATE TABLE IF NOT EXISTS `lab_class_infos` (
  chapter_id VARCHAR(36) NOT NULL,
  chapter_index INT NOT NULL,
  name VARCHAR(256) NOT NULL,
  fullmark INT NOT NULL,
  no_items INT NOT NULL DEFAULT '5',
  PRIMARY KEY (chapter_id)
);

INSERT INTO `lab_class_infos` (chapter_id, chapter_index, name, fullmark, no_items) VALUES
(UUID(), 1, 'Introduction', 10, 5),
(UUID(), 2, 'Variables Expression Statement', 10, 5),
(UUID(), 3, 'Conditional Execution', 10, 5),
(UUID(), 4, 'Loop while', 10, 5),
(UUID(), 5, 'Loop for', 10, 5),
(UUID(), 6, 'List', 10, 5),
(UUID(), 7, 'String', 10, 5),
(UUID(), 8, 'Function', 10, 5),
(UUID(), 9, 'Dictionary', 10, 5),
(UUID(), 10, 'Files', 10, 5),
(UUID(), 11, 'Best Practice 1', 10, 5),
(UUID(), 12, 'Best Practice 2', 10, 5),
(UUID(), 13, 'Quiz #1 , chapter 01 - 03)', 10, 5),
(UUID(), 14, 'Quiz #2 , chapter 01 - 06)', 10, 5),
(UUID(), 15, 'Quiz #3 , chapter 01 - 09)', 10, 5),
(UUID(), 16, 'Reserve #1', 10, 5),
(UUID(), 17, 'Reserve #2', 10, 5);
