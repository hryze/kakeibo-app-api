DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;
USE test_db;

CREATE TABLE users
(
  user_id VARCHAR(10) NOT NULL,
  name VARCHAR(50) NOT NULL,
  email VARCHAR(50) NOT NULL,
  password VARCHAR(60) NOT NULL,
  PRIMARY KEY(user_id),
  UNIQUE uq_users(email)
);

CREATE TABLE group_names
(
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  group_name VARCHAR(20) NOT NULL
);

CREATE TABLE group_users
(
  id INT NOT NULL AUTO_INCREMENT,
  group_id INT NOT NULL,
  user_id VARCHAR(10) NOT NULL,
  color_code VARCHAR(7) NOT NULL,
  PRIMARY KEY(id),
  UNIQUE uq_group_users(group_id, user_id),
  FOREIGN KEY fk_group_id(group_id)
    REFERENCES group_names(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY fk_user_id(user_id)
    REFERENCES users(user_id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE group_unapproved_users
(
  id INT NOT NULL AUTO_INCREMENT,
  group_id INT NOT NULL,
  user_id VARCHAR(10) NOT NULL,
  PRIMARY KEY(id),
  UNIQUE uq_group_unapproved_users(group_id, user_id),
  FOREIGN KEY fk_group_id(group_id)
    REFERENCES group_names(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY fk_user_id(user_id)
    REFERENCES users(user_id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);
