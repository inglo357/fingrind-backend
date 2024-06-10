-- Create "currencies" table
CREATE TABLE "currencies" ("id" bigserial NOT NULL, "currency_string" character varying(80) NOT NULL, PRIMARY KEY ("id"));
-- Modify "accounts" table
ALTER TABLE "accounts" ADD CONSTRAINT "currency_id" FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
