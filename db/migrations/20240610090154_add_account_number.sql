-- Modify "accounts" table
ALTER TABLE "accounts" ADD COLUMN "account_number" character varying(20) NOT NULL, ADD CONSTRAINT "unique_account_number" UNIQUE ("account_number");
