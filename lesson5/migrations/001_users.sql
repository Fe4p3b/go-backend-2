-- +goose Up
CREATE TABLE "users" (
   "user_id" INT,
   "name" VARCHAR,
   "age" INT,
   "spouse" INT
);
CREATE UNIQUE INDEX "users_user_id" ON "users" ("user_id");

-- +goose Down
DROP TABLE users;
