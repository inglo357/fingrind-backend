-- name: CreateCurrency :one
INSERT INTO currencies (
    currency_string
) VALUES ($1) RETURNING *;

-- name: GetCurrencyByID :one
SELECT * FROM currencies WHERE id = $1;

-- name: GetCurrencyByCurrencyString :one
SELECT * FROM currencies WHERE currency_string = $1;

-- name: ListCurrencies :many
SELECT * FROM currencies ORDER BY id
LIMIT $1 OFFSET $2;

-- name: DeleteAllCurrencies :exec
DELETE FROM currencies;