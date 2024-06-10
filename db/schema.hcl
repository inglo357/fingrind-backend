schema "fingrind" {
}

table "users" {
  schema = schema.fingrind
  column "id" {
    null = false
    type = bigserial
  }
  column "name" {
    null = false
    type = character_varying(256)
  }
  column "email" {
    null = false
    type = character_varying(256)
  }
  column "hashed_password" {
    null = false
    type = character_varying(256)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  index "email" {
    unique  = true
    columns = [column.email]
  }
}

table "accounts" {
  schema = schema.fingrind
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id"{
    null = false
    type = bigserial
  }
  column "balance"{
    null = false
    type = float(30)
    default = 0.00
  }
  column "currency_id"{
    null = false
    type = bigserial
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "updated_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  column "account_number"{
    null = false
    type = character_varying(20)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "user_id"{
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  foreign_key "currency_id"{
    columns = [column.currency_id]
    ref_columns = [table.currencies.column.id]
    on_delete = CASCADE
  }
  unique "unique_user_currency"{
    columns = [column.user_id, column.currency_id]
  }
  unique "unique_account_number"{
    columns = [column.account_number]
  }
}

table "entries"{
  schema = schema.fingrind
  column "id" {
    null = false
    type = bigserial
  }
  column "account_id"{
    null = false
    type = bigserial
  }
  column "amount"{
    null = false
    type = float(30)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "account_id"{
    columns = [column.account_id]
    ref_columns = [table.accounts.column.id]
    on_delete = CASCADE
  }
}

table "transfers"{
  schema = schema.fingrind
  column "id" {
    null = false
    type = bigserial
  }
  column "from_account_id"{
    null = false
    type = bigserial
  }
  column "to_account_id"{
    null = false
    type = bigserial
  }
  column "amount"{
    null = false
    type = float(30)
  }
  column "created_at" {
    null    = false
    type    = timestamp
    default = sql("CURRENT_TIMESTAMP")
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "from_account_id"{
    columns = [column.from_account_id]
    ref_columns = [table.accounts.column.id]
    on_delete = CASCADE
  }
  foreign_key "to_account_id"{
    columns = [column.to_account_id]
    ref_columns = [table.accounts.column.id]
    on_delete = CASCADE
  }
}

table "currencies"{
  schema = schema.fingrind
  column "id" {
    null = false
    type = bigserial
  }
  column "currency_string"{
    null = false
    type = character_varying(80)
  }
  column "currency_name"{
    null = true
    type = character_varying(256)
  }
  column "starter"{
    null = true
    type = integer
  }
  primary_key {
    columns = [column.id]
  }
  unique "unique_currency"{
    columns = [column.currency_string]
  }
}

table "money_records"{
  schema = schema.fingrind
  column "id" {
    null = false
    type = bigserial
  }
  column "user_id" {
    null = false
    type = bigserial
  }
  column "reference"{
    null = false
    type = character_varying(80)
  }
  column "status"{
    null = false
    type = character_varying(80)
  }
  column "amount"{
    null = false
    type = float(30)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "user_id"{
    columns = [column.user_id]
    ref_columns = [table.users.column.id]
    on_delete = CASCADE
  }
  unique "unique_reference"{
    columns = [column.reference]
  }
}
