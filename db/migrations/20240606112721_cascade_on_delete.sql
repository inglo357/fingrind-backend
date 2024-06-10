-- Modify "accounts" table
ALTER TABLE "accounts" DROP CONSTRAINT "currency_id", DROP CONSTRAINT "user_id", ADD CONSTRAINT "currency_id" FOREIGN KEY ("currency_id") REFERENCES "currencies" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "user_id" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "entries" table
ALTER TABLE "entries" DROP CONSTRAINT "account_id", ADD CONSTRAINT "account_id" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "transfers" table
ALTER TABLE "transfers" DROP CONSTRAINT "from_account_id", DROP CONSTRAINT "to_account_id", ADD CONSTRAINT "from_account_id" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT "to_account_id" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
