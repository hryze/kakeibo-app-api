DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;
USE test_db;

DROP TABLE IF EXISTS users;
CREATE TABLE users
(
  id VARCHAR(10) NOT NULL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  email VARCHAR(50) NOT NULL,
  password TEXT NOT NULL
);

INSERT INTO users
  (id, name, email, password)
VALUES
  ('tati1', '館ひろし', 'tati@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu');

INSERT INTO users
  (id, name, email, password)
VALUES
  ('go2', '郷ひろみ', 'go@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu');

INSERT INTO users
  (id, name, email, password)
VALUES
  ('saigo3', '西郷隆盛', 'saigo@developer.com', '$2a$10$iMB.JdS3kvyX60R31RbCcO.yBfMpDGlwuLL//7bibRHy89vjXgXQu');
