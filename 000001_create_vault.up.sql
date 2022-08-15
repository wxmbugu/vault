CREATE TABLE "vault" (
  "id" BIGSERIAL PRIMARY KEY,
  "secret" varchar NOT NULL,
  "duration" varchar NOT NULL,
  "uuid" varchar NOT NULL
);
