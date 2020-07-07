DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;
USE test_db;

CREATE TABLE big_categories
(
  id INT NOT NULL AUTO_INCREMENT,
  category_name VARCHAR(10) NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE medium_categories
(
  id INT NOT NULL AUTO_INCREMENT,
  category_name VARCHAR(10) NOT NULL,
  big_category_id INT NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE custom_categories
(
  id INT NOT NULL AUTO_INCREMENT,
  category_name VARCHAR(50) NOT NULL,
  big_category_id INT NOT NULL,
  user_id VARCHAR(10) NOT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  INDEX idx_user_id(user_id)
);

CREATE TABLE expenses
(
  id INT NOT NULL AUTO_INCREMENT,
  posted_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  payment_date DATE NOT NULL,
  payment_shop VARCHAR(20) DEFAULT NULL,
  memo VARCHAR(50) DEFAULT NULL,
  amount INT NOT NULL,
  user_id VARCHAR(10) NOT NULL,
  big_category_id INT NOT NULL,
  medium_category_id INT DEFAULT NULL,
  custom_category_id INT DEFAULT NULL,
  PRIMARY KEY(id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY fk_medium_category_id(medium_category_id)
    REFERENCES medium_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY fk_custom_category_id(custom_category_id)
    REFERENCES custom_categories(id)
    ON DELETE SET NULL ON UPDATE CASCADE
);
