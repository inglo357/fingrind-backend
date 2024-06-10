-- Modify "accounts" table
ALTER TABLE "accounts" ADD CONSTRAINT "unique_user_currency" UNIQUE ("user_id") INCLUDE ("currency_id");
