-- Modify "accounts" table
ALTER TABLE "accounts" ALTER COLUMN "balance" TYPE numeric USING ("balance"::numeric), ALTER COLUMN "balance" SET DEFAULT 0.00;
-- Modify "accounts" table
ALTER TABLE "accounts" ALTER COLUMN "balance" TYPE double precision USING ("balance"::double precision), ALTER COLUMN "balance" SET DEFAULT 0.00;
-- Modify "entries" table
ALTER TABLE "entries" ALTER COLUMN "amount" TYPE numeric USING ("amount"::numeric);
-- Modify "entries" table
ALTER TABLE "entries" ALTER COLUMN "amount" TYPE double precision USING ("amount"::double precision);
-- Modify "transfers" table
ALTER TABLE "transfers" ALTER COLUMN "amount" TYPE numeric USING ("amount"::numeric);
-- Modify "transfers" table
ALTER TABLE "transfers" ALTER COLUMN "amount" TYPE double precision USING ("amount"::double precision);
