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
  PRIMARY KEY(id),
  INDEX idx_user_id(user_id)
);

CREATE TABLE regular_shopping_list
(
  id INT NOT NULL AUTO_INCREMENT,
  posted_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  expected_purchase_date DATE NOT NULL,
  cycle_type ENUM('daily', 'weekly', 'monthly', 'custom') NOT NULL,
  cycle INT DEFAULT NULL,
  purchase VARCHAR(50) NOT NULL,
  shop VARCHAR(20) DEFAULT NULL,
  amount INT DEFAULT NULL,
  big_category_id INT NOT NULL,
  medium_category_id INT DEFAULT NULL,
  custom_category_id INT DEFAULT NULL,
  transaction_auto_add bit(1) NOT NULL DEFAULT b'0',
  user_id VARCHAR(10) NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE shopping_list
(
  id INT NOT NULL AUTO_INCREMENT,
  posted_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  expected_purchase_date DATE NOT NULL,
  complete_flag bit(1) NOT NULL DEFAULT b'0',
  purchase VARCHAR(50) NOT NULL,
  shop VARCHAR(20) DEFAULT NULL,
  amount INT DEFAULT NULL,
  big_category_id INT NOT NULL,
  medium_category_id INT DEFAULT NULL,
  custom_category_id INT DEFAULT NULL,
  regular_shopping_list_id INT DEFAULT NULL,
  user_id VARCHAR(10) NOT NULL,
  transaction_auto_add bit(1) NOT NULL DEFAULT b'0',
  transaction_id INT DEFAULT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY fk_shopping_list_id(regular_shopping_list_id)
    REFERENCES regular_shopping_list(id)
    ON DELETE SET NULL ON UPDATE CASCADE
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
  PRIMARY KEY(id),
  INDEX idx_group_id(group_id)
);

CREATE TABLE group_tasks_users
(
  id INT NOT NULL AUTO_INCREMENT,
  user_id VARCHAR(10) NOT NULL,
  group_id INT NOT NULL,
  PRIMARY KEY(id),
  UNIQUE uq_group_tasks_users(user_id, group_id),
  INDEX idx_group_id(group_id)
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
    ON DELETE SET NULL ON UPDATE CASCADE,
  INDEX idx_group_id(group_id)
);
