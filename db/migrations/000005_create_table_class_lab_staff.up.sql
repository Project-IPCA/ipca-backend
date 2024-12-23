CREATE TABLE IF NOT EXISTS `class_lab_staffs` (
  class_id VARCHAR(36) NOT NULL,
  staff_id VARCHAR(36) NOT NULL,
  role ENUM('ADMIN', 'SUPERVISOR', 'TA') DEFAULT NULL,
  PRIMARY KEY (class_id,staff_id),
  KEY class_id (class_id),
  KEY staff_id (staff_id),
  CONSTRAINT fk_class_lab_staff_user FOREIGN KEY (staff_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT fk_class_lab_staff_class_schedules FOREIGN KEY (class_id) REFERENCES class_schedules (group_id) ON DELETE CASCADE ON UPDATE CASCADE
); 