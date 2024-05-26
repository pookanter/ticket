CREATE TABLE users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50),
  lastname VARCHAR(50),
  email VARCHAR(255),
  password VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE boards (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED,
  title VARCHAR(100),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE statuses (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  board_id INT UNSIGNED,
  title VARCHAR(50),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (board_id) REFERENCES boards(id)
);

CREATE TABLE tickets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  status_id INT UNSIGNED,
  title VARCHAR(100),
  description TEXT,
  contact VARCHAR(100),
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (status_id) REFERENCES statuses(id)
);