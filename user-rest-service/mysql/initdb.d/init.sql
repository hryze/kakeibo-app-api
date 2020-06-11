DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;
USE test_db;

DROP TABLE IF EXISTS users;
CREATE TABLE users
(
  user_id VARCHAR(10) NOT NULL PRIMARY KEY,
  user_name VARCHAR(50) NOT NULL,
  user_mail VARCHAR(50) NOT NULL
);

INSERT INTO users
  (user_id, user_name, user_mail)
VALUES
  ('tati1', '館 ひろし', 'tati@developer.com');

INSERT INTO users
  (user_id, user_name, user_mail)
VALUES
  ('go2', '郷 ひろみ', 'go@developer.com');

INSERT INTO users
  (user_id, user_name, user_mail)
VALUES
  ('saigo3', '西郷 隆盛', 'saigo@developer.com');