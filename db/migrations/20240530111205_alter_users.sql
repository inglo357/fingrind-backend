-- Create "users" table
CREATE TABLE "users" ("id" bigserial NOT NULL, "name" character varying(256) NOT NULL, "email" character varying(256) NOT NULL, "hashed_password" character varying(256) NOT NULL, "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"));
-- Create index "email" to table: "users"
CREATE UNIQUE INDEX "email" ON "users" ("email");
