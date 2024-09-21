CREATE TABLE IF NOT EXISTS `infinite_loops` (
  group_id VARCHAR(36) NOT NULL,
  stu_id VARCHAR(36) NOT NULL,
  chapter_id VARCHAR(36) NOT NULL,
  item_id INT NOT NULL,
  sequence INT NOT NULL,
  start VARCHAR(20) NOT NULL,
  time VARCHAR(20) NOT NULL,
  timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (group_id,stu_id,chapter_id,item_id,sequence)
);

