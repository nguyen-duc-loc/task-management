CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "tasks" (
  "id" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "creator" varchar NOT NULL,
  "deadline" timestamptz NOT NULL,
  "completed" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "tasks" ADD FOREIGN KEY ("creator") REFERENCES "users" ("username");
