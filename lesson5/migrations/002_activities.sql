-- +goose Up
CREATE TABLE "activities" (
   "user_id" INT,
   "date" TIMESTAMP,
   "name" VARCHAR
);

CREATE INDEX "activities_user_id_date" ON "activities" ("user_id", "date");
-- +goose Down
DROP TABLE activities;