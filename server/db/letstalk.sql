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

CREATE TABLE fb_auth_data (
  user_id INT NOT NULL PRIMARY KEY,
  fb_user_id VARCHAR(32) NOT NULL,

  FOREIGN KEY (user_id) REFERENCES user(user_id)
);

CREATE TABLE fb_auth_token (
  user_id INT NOT NULL,
  auth_token VARCHAR(50) NOT NULL,
  expiry DATETIME(6) NOT NULL,

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
);

CREATE TABLE matchings (
  matching_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  mentor INT NOT NULL,
  mentee INT NOT NULL,
  UNIQUE KEY (mentor, mentee),
  KEY(mentor),
  KEY(mentee),
  FOREIGN KEY (mentor)
    REFERENCES user(user_id),
  FOREIGN KEY (mentee)
    REFERENCES user(user_id)
);

CREATE TABLE notification_tokens (
  id VARCHAR(64) NOT NULL PRIMARY KEY,
  user_id INT NOT NULL,
  token varchar(255) NOT NULL,
  -- remove the notification when the session doesn't exist anymore
  FOREIGN KEY (id) REFERENCES sessions(session_id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES user(user_id)
);

CREATE TABLE user_vector (
  user_id INT PRIMARY KEY NOT NULL,
  preference_type INT NOT NULL, -- indicate if this is their preferences for a mentee or for a mentor
  sociable  INT NOT NULL,
  hard_working INT NOT NULL,
  ambitious INT NOT NULL,
  energetic INT NOT NULL,
  carefree iNT NOT NULL,
  confident INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES user(user_id)
);
