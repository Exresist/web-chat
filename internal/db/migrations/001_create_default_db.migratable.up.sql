BEGIN;

CREATE TABLE IF NOT EXISTS users (
  id int,
  username varchar(64),
  hashed_password varchar,
  first_name varchar(64),
  last_name varchar(64),
  email varchar(100),
  photo bytea
);
COMMIT;