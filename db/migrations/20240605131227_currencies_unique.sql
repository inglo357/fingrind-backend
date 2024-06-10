-- Modify "currencies" table
ALTER TABLE "currencies" ADD CONSTRAINT "unique_currency" UNIQUE ("currency_string");
