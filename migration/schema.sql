CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(50),
  lastname VARCHAR(50),
  email VARCHAR(255),
  password VARCHAR(255),
  created_at DATETIME,
  updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS boards (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(100),
  sort_order INT UNSIGNED NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS statuses (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  board_id INT UNSIGNED NOT NULL,
  title VARCHAR(50),
  sort_order INT UNSIGNED NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (board_id) REFERENCES boards(id)
);

CREATE TABLE IF NOT EXISTS tickets (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  status_id INT UNSIGNED NOT NULL,
  title VARCHAR(100),
  description TEXT,
  contact VARCHAR(100),
  sort_order INT UNSIGNED NOT NULL,
  created_at DATETIME,
  updated_at DATETIME,
  FOREIGN KEY (status_id) REFERENCES statuses(id)
);

CREATE
OR REPLACE VIEW status_view AS
SELECT
  s.*,
  JSON_ARRAYAGG(
    JSON_OBJECT(
      'id',
      t.id,
      'status_id',
      t.status_id,
      'title',
      t.title,
      'description',
      t.description,
      'contact',
      t.contact,
      'sort_order',
      t.sort_order,
      'created_at',
      t.created_at,
      'updated_at',
      t.updated_at
    )
  ) AS tickets
FROM
  statuses AS s
  LEFT JOIN tickets AS t ON s.id = t.status_id
GROUP BY
  s.id;

CREATE
OR REPLACE VIEW board_view AS
SELECT
  b.*,
  JSON_ARRAYAGG(
    JSON_OBJECT(
      'id',
      sv.id,
      'board_id',
      sv.board_id,
      'title',
      sv.title,
      'sort_order',
      sv.sort_order,
      'created_at',
      sv.created_at,
      'updated_at',
      sv.updated_at,
      'tickets',
      sv.tickets
    )
  ) AS statuses
FROM
  boards AS b
  LEFT JOIN status_view AS sv ON b.id = sv.board_id
GROUP BY
  b.id