-- Modify "accounts" table
ALTER TABLE "accounts" DROP CONSTRAINT "unique_user_currency", ADD CONSTRAINT "unique_user_currency" UNIQUE ("user_id", "currency_id");
