CREATE TABLE IF NOT EXISTS `users` (
  user_id VARCHAR(36) NOT NULL,
  username VARCHAR(30) UNIQUE NOT NULL,
  password VARCHAR(60) DEFAULT NULL,
  f_name VARCHAR(10) DEFAULT NULL,
  l_name VARCHAR(32) DEFAULT NULL,
  nickname VARCHAR(50) DEFAULT NULL,
  gender ENUM('MALE', 'FEMALE', 'OTHER') DEFAULT NULL,
  dob DATE DEFAULT NULL,
  avatar VARCHAR(128) DEFAULT NULL,
  role ENUM('ADMIN', 'EDITOR', 'AUTHOR', 'STUDENT', 'SUPERVISOR', 'STAFF', 'TA') DEFAULT NULL,
  email VARCHAR(64) DEFAULT NULL,
  tel VARCHAR(10) DEFAULT NULL,
  added DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_login DATETIME DEFAULT NULL,
  last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_online BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  added_by VARCHAR(40) DEFAULT NULL,
  ci_session INT DEFAULT NULL,
  session_id VARCHAR(50) DEFAULT NULL,
  PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS `students` (
  stu_id VARCHAR(36) NOT NULL,
  kmitl_id VARCHAR(8) NOT NULL,
  group_id INT DEFAULT NULL,
  note VARCHAR(64) DEFAULT NULL,
  dept_id VARCHAR(36) DEFAULT NULL,
  mid_core FLOAT NOT NULL DEFAULT '0',
  can_submit BOOLEAN NOT NULL DEFAULT TRUE,
  PRIMARY KEY (stu_id),
  UNIQUE KEY user_student_pk (stu_id),
  KEY student_group (group_id),
  KEY stu_department (dept_id),
  CONSTRAINT user_student_ibfk_1 FOREIGN KEY (stu_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `supervisors` (
  supervisor_id VARCHAR(36) NOT NULL,
  dept VARCHAR(40) DEFAULT NULL,
  PRIMARY KEY (supervisor_id),
  CONSTRAINT fk_user_supervisor_user1 FOREIGN KEY (supervisor_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `tas` (
  ta_id VARCHAR(36) NOT NULL,
  group_id INT DEFAULT NULL,
  note VARCHAR(64) DEFAULT NULL,
  dept_id VARCHAR(36) DEFAULT NULL,
  PRIMARY KEY (ta_id),
  KEY fk_user_ta_department1_idx (dept_id),
  KEY fk_user_ta_class_schedule1_idx (group_id),
  CONSTRAINT user_ta_ibfk_1 FOREIGN KEY (ta_id) REFERENCES users (user_id) ON DELETE CASCADE ON UPDATE CASCADE
);