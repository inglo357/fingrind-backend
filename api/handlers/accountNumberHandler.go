package handlers

import (
	"fmt"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"time"
)

func GenerateAccountNumber(accountId int64, currency db.Currency) (string, error) {

	activeTime := time.Now().Format("20060102150405")
	initialValue := fmt.Sprintf("%v%d", currency.Starter, accountId)

	finalValue := ""
	reminder := 10 - len(initialValue)

	if reminder > 0 {
		finalValue = activeTime[:reminder]
	}

	accountNumber := fmt.Sprintf("%s%s", initialValue, finalValue)

	return accountNumber, nil
}