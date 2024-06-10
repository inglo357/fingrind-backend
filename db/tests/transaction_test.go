package db_tests

import (
	"context"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createRandomAccount(userId int64, t *testing.T) db.Account{

	args := db.CreateAccountParams{
		UserID: userId,
		CurrencyID: 1,
		Balance: 500,
	}

	account, err := testQuery.CreateAccount(context.Background(), args)

	assert.NoError(t, err)
	assert.NotEmpty(t, account)
	assert.Equal(t, account.UserID, args.UserID)
	assert.Equal(t, account.CurrencyID, args.CurrencyID)
	assert.Equal(t, account.Balance, args.Balance)
	assert.WithinDuration(t, account.CreatedAt, time.Now(), 10*time.Second)
	assert.WithinDuration(t, account.UpdatedAt, time.Now(), 10*time.Second)

	return account
}

func TestTransferTx(t *testing.T) {	

	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	account1 := createRandomAccount(user1.ID, t)
	account2 := createRandomAccount(user2.ID, t)

	args := db.CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: 10,
	}

	const transactionCount = 5
	transferChan := make(chan db.TransactionTxResponse)
	errorChan := make(chan error)

	go func(){
		for i := 0; i < transactionCount; i++ {
			transfer, err := testQuery.TransferTx(context.Background(), args)
			transferChan <- transfer
			errorChan <- err
		}
	}()

	for i := 0; i < transactionCount; i++ {
		transfer := <-transferChan
		err := <-errorChan

		assert.NoError(t, err)
		assert.NotEmpty(t, transfer)
	
		assert.Equal(t, transfer.Transfer.FromAccountID, args.FromAccountID)
		assert.Equal(t, transfer.Transfer.ToAccountID, args.ToAccountID)
		assert.Equal(t, transfer.Transfer.Amount, args.Amount)
	
		assert.Equal(t, transfer.EntryIn.AccountID, args.ToAccountID)
		assert.Equal(t, transfer.EntryIn.Amount, args.Amount)
	
		assert.Equal(t, transfer.EntryOut.AccountID, args.FromAccountID)
		assert.Equal(t, transfer.EntryOut.Amount, -args.Amount)
	
		// assert.Equal(t, transfer.FromAccount.ID, args.FromAccountID)
		// assert.Equal(t, transfer.ToAccount.ID, args.ToAccountID)
		// assert.Equal(t, transfer.FromAccount.Balance, account1.Balance-args.Amount)
		// assert.Equal(t, transfer.ToAccount.Balance, account2.Balance+args.Amount)
	}

	newAccount1, err := testQuery.GetAccountByID(context.Background(), account1.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccount1)

	newAccount2, err := testQuery.GetAccountByID(context.Background(), account2.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, newAccount2)

	assert.Equal(t, newAccount1.Balance, account1.Balance-(transactionCount*args.Amount))
	assert.Equal(t, newAccount2.Balance, account2.Balance+(transactionCount*args.Amount))
}