CREATE TABLE IF NOT EXISTS `user` (
  id VARCHAR(36) NOT NULL,
  username VARCHAR(30) UNIQUE NOT NULL,
  password VARCHAR(60) DEFAULT NULL,
  role ENUM('ADMIN', 'EDITOR', 'AUTHOR', 'STUDENT', 'SUPERVISOR', 'STAFF', 'TA') DEFAULT NULL,
  added DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_login DATETIME DEFAULT NULL,
  last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_online BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  added_by VARCHAR(40) DEFAULT NULL,
  ci_session INT DEFAULT NULL,
  session_id VARCHAR(50) DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS `user_student` (
  stu_id VARCHAR(36) NOT NULL,
  stu_stu_id VARCHAR(10) NOT NULL,
  stu_firstname VARCHAR(40) DEFAULT NULL,
  stu_lastname VARCHAR(32) DEFAULT NULL,
  stu_nickname VARCHAR(20) DEFAULT NULL,
  stu_gender ENUM('MALE', 'FEMALE', 'OTHER') DEFAULT NULL,
  stu_dob DATE DEFAULT NULL,
  stu_avatar VARCHAR(128) DEFAULT NULL,
  stu_email VARCHAR(64) DEFAULT NULL,
  stu_tel VARCHAR(10) DEFAULT NULL,
  stu_group INT DEFAULT NULL,
  note VARCHAR(64) DEFAULT NULL,
  stu_dept_id INT DEFAULT NULL,
  mid_core FLOAT NOT NULL DEFAULT '0',
  can_submit VARCHAR(3) NOT NULL DEFAULT 'YES',
  PRIMARY KEY (stu_id),
  UNIQUE KEY user_student_pk (stu_id),
  KEY student_group (stu_group),
  KEY stu_department (stu_dept_id),
  CONSTRAINT user_student_ibfk_1 FOREIGN KEY (stu_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `user_supervisor` (
  supervisor_id VARCHAR(36) NOT NULL,
  supervisor_firstname VARCHAR(50) DEFAULT NULL,
  supervisor_lastname VARCHAR(50) DEFAULT NULL,
  supervisor_nickname VARCHAR(50) DEFAULT NULL,
  supervisor_gender ENUM('MALE', 'FEMALE', 'OTHER') DEFAULT NULL,
  supervisor_dob DATE DEFAULT NULL,
  supervisor_avatar VARCHAR(64) DEFAULT NULL,
  supervisor_email VARCHAR(64) DEFAULT NULL,
  supervisor_tel VARCHAR(10) DEFAULT NULL,
  supervisor_department VARCHAR(40) DEFAULT NULL,
  PRIMARY KEY (supervisor_id),
  CONSTRAINT fk_user_supervisor_user1 FOREIGN KEY (supervisor_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS `user_ta` (
  ta_id VARCHAR(36) NOT NULL,
  ta_gender ENUM('MALE', 'FEMALE', 'OTHER') DEFAULT NULL,
  ta_firstname VARCHAR(40) DEFAULT NULL,
  ta_lastname VARCHAR(32) DEFAULT NULL,
  ta_nickname VARCHAR(20) DEFAULT NULL,
  ta_dob DATE DEFAULT NULL,
  ta_avatar VARCHAR(128) DEFAULT NULL,
  ta_email VARCHAR(64) DEFAULT NULL,
  ta_tel VARCHAR(10) DEFAULT NULL,
  ta_group INT DEFAULT NULL,
  note VARCHAR(64) DEFAULT NULL,
  ta_dept_id INT DEFAULT NULL,
  PRIMARY KEY (ta_id),
  KEY fk_user_ta_department1_idx (ta_dept_id),
  KEY fk_user_ta_class_schedule1_idx (ta_group),
  CONSTRAINT user_ta_ibfk_1 FOREIGN KEY (ta_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE
);
