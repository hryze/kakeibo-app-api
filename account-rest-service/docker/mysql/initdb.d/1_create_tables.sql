DROP DATABASE IF EXISTS test_db;
CREATE DATABASE test_db;
USE test_db;

CREATE TABLE big_categories
(
  id INT NOT NULL AUTO_INCREMENT,
  category_name VARCHAR(10) NOT NULL,
  transaction_type ENUM('expense', 'income') NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE medium_categories
(
  id INT NOT NULL AUTO_INCREMENT,
  category_name VARCHAR(10) NOT NULL,
  big_category_id INT NOT NULL,
  PRIMARY KEY(id),
  UNIQUE uq_medium_category(category_name, big_category_id),
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
  UNIQUE uq_custom_category(category_name, big_category_id, user_id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  INDEX idx_user_id(user_id, id)
);

CREATE TABLE transactions
(
  id INT NOT NULL AUTO_INCREMENT,
  transaction_type ENUM('expense', 'income') NOT NULL,
  posted_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  transaction_date DATE NOT NULL,
  shop VARCHAR(20) DEFAULT NULL,
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
    ON DELETE SET NULL ON UPDATE CASCADE,
  INDEX idx_user_id(user_id)
);

CREATE TABLE standard_budgets
(
  user_id VARCHAR(10) NOT NULL,
  big_category_id INT NOT NULL,
  budget INT NOT NULL DEFAULT 0,
  PRIMARY KEY(user_id, big_category_id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE custom_budgets
(
  user_id VARCHAR(10) NOT NULL,
  years_months DATE NOT NULL,
  big_category_id INT NOT NULL,
  budget INT NOT NULL,
  PRIMARY KEY(user_id, years_months, big_category_id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE group_custom_categories
(
  id INT NOT NULL AUTO_INCREMENT,
  category_name VARCHAR(50) NOT NULL,
  big_category_id INT NOT NULL,
  group_id INT NOT NULL,
  PRIMARY KEY(id),
  UNIQUE uq_group_custom_category(category_name, big_category_id, group_id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE,
  INDEX idx_group_id(group_id, id)
);

CREATE TABLE group_transactions
(
  id INT NOT NULL AUTO_INCREMENT,
  transaction_type ENUM('expense', 'income') NOT NULL,
  posted_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  transaction_date DATE NOT NULL,
  shop VARCHAR(20) DEFAULT NULL,
  memo VARCHAR(50) DEFAULT NULL,
  amount INT NOT NULL,
  group_id INT NOT NULL,
  posted_user_id VARCHAR(10) NOT NULL,
  updated_user_id VARCHAR(10) DEFAULT NULL,
  payment_user_id VARCHAR(10) NOT NULL,
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
    REFERENCES group_custom_categories(id)
    ON DELETE SET NULL ON UPDATE CASCADE,
  INDEX idx_group_id(group_id)
);

CREATE TABLE group_standard_budgets
(
  group_id INT NOT NULL,
  big_category_id INT NOT NULL,
  budget INT NOT NULL DEFAULT 0,
  PRIMARY KEY(group_id, big_category_id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE group_custom_budgets
(
  group_id INT NOT NULL,
  years_months DATE NOT NULL,
  big_category_id INT NOT NULL,
  budget INT NOT NULL,
  PRIMARY KEY(group_id, years_months, big_category_id),
  FOREIGN KEY fk_big_category_id(big_category_id)
    REFERENCES big_categories(id)
    ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE group_accounts
(
  id INT NOT NULL AUTO_INCREMENT,
  years_months DATE NOT NULL,
  payer_user_id VARCHAR(10) DEFAULT NULL,
  recipient_user_id VARCHAR(10) DEFAULT NULL,
  payment_amount INT DEFAULT NULL,
  payment_confirmation bit(1) NOT NULL DEFAULT b'0',
  receipt_confirmation bit(1) NOT NULL DEFAULT b'0',
  group_id INT NOT NULL,
  PRIMARY KEY(id),
  UNIQUE uq_group_accounts(years_months, payer_user_id, recipient_user_id, group_id),
  INDEX idx_group_id(group_id)
);
