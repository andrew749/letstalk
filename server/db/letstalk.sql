DROP DATABASE letstalk;
CREATE DATABASE letstalk
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;
USE letstalk;

-- NOTE: max length of indexed VARCHAR is 191 (767 bytes / 4 bytes per utf8 char)

CREATE TABLE id_gen (
  id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  num_id INT NOT NULL
);

INSERT INTO id_gen (num_id) VALUE (1);

CREATE TABLE user (
  user_id INT NOT NULL PRIMARY KEY,
  first_name VARCHAR(128) NOT NULL,
  last_name VARCHAR(128) NOT NULL,
  email VARCHAR(128) NOT NULL,
  gender INT NOT NULL,
  birthdate DATETIME(6) NOT NULL,

  UNIQUE KEY (email)
);

CREATE TABLE authentication_data (
  user_id INT NOT NULL PRIMARY KEY,
  password_hash VARCHAR(128) NOT NULL,

  FOREIGN KEY (user_id) REFERENCES user(user_id)
);

CREATE TABLE program (
  program_id VARCHAR(64) NOT NULL PRIMARY KEY
);

INSERT INTO program (program_id) VALUES
  ('SOFTWARE_ENGINEERING'),
  ('COMPUTER_ENGINEERING'),
  ('UNKNOWN')
;

CREATE TABLE cohort (
  cohort_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  program_id VARCHAR(64) NOT NULL,
  grad_year SMALLINT NOT NULL,
  sequence VARCHAR(32) NOT NULL,

  UNIQUE KEY (program_id, grad_year, sequence),

  FOREIGN KEY (program_id)
    REFERENCES program(program_id)
);

INSERT INTO cohort (program_id, grad_year, sequence) VALUES
  ('SOFTWARE_ENGINEERING', 2019, '8STREAM'),
  ('COMPUTER_ENGINEERING', 2019, '8STREAM'),
  ('COMPUTER_ENGINEERING', 2019, '4STREAM');

CREATE TABLE user_cohort (
  user_id INT NOT NULL PRIMARY KEY,
  cohort_id INT NOT NULL,

  FOREIGN KEY (user_id)
    REFERENCES user(user_id),

  FOREIGN KEY (cohort_id)
    REFERENCES cohort(cohort_id)
);

CREATE TABLE sessions (
  session_id VARCHAR(64) NOT NULL,
  user_id INT NOT NULL,
  expiry_date DATETIME(6) NOT NULL,

  UNIQUE KEY (session_id),

  FOREIGN KEY (user_id)
    REFERENCES user(user_id)
)
