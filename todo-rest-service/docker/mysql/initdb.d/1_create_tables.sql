DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;
USE test_db;

CREATE TABLE todo_list
(
  id INT NOT NULL AUTO_INCREMENT,
  posted_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  implementation_date DATE NOT NULL,
  due_date DATE NOT NULL,
  todo_content VARCHAR(100) NOT NULL,
  complete_flag bit(1) NOT NULL DEFAULT b'0',
  user_id VARCHAR(10) NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE group_todo_list
(
  id INT NOT NULL AUTO_INCREMENT,
  posted_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  implementation_date DATE NOT NULL,
  due_date DATE NOT NULL,
  todo_content VARCHAR(100) NOT NULL,
  complete_flag bit(1) NOT NULL DEFAULT b'0',
  user_id VARCHAR(10) NOT NULL,
  group_id INT NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE group_tasks_users
(
  id INT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(10) NOT NULL,
  group_id INT NOT NULL,
  PRIMARY KEY(id),
  UNIQUE uq_group_tasks_users(user_id, group_id)
);

CREATE TABLE group_tasks
(
  id INT NOT NULL AUTO_INCREMENT,
  base_date DATETIME DEFAULT NULL,
  cycle_type ENUM('every', 'consecutive', 'none') DEFAULT NULL,
  cycle INT DEFAULT NULL,
  task_name VARCHAR(20) NOT NULL,
  group_id INT NOT NULL,
  group_tasks_users_id INT DEFAULT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY fk_group_tasks_users_id(group_tasks_users_id)
    REFERENCES group_tasks_users(id)
    ON DELETE SET NULL ON UPDATE CASCADE
);
