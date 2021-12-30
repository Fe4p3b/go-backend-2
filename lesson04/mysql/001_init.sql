CREATE TABLE IF NOT EXISTS entities (
   id INT PRIMARY KEY,
   data VARCHAR(32)
);

CREATE TABLE IF NOT EXISTS tokens (
   id INT PRIMARY KEY,
   token VARCHAR(32)
);

INSERT INTO tokens VALUES (1, "admin_secret_token");