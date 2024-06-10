package db

import (
	"context"
)

type TransactionTxResponse struct {
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	EntryIn     Entry   `json:"entry_in"`
	EntryOut    Entry   `json:"entry_out"`
	Transfer    Transfer `json:"transfer"`
}

func (s	*Store) TransferTx(ctx context.Context, params CreateTransferParams) (TransactionTxResponse, error) {

	var transactionResponse TransactionTxResponse
	err := s.ExecTx(ctx, func(q *Queries) error {
		var err error
		transactionResponse.Transfer, err = q.CreateTransfer(context.Background(), CreateTransferParams{
			FromAccountID: params.FromAccountID,
			ToAccountID:   params.ToAccountID,
			Amount:        params.Amount,
		})
		if err != nil {
			return err
		}

		transactionResponse.EntryIn, err = q.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: params.ToAccountID,
			Amount:    params.Amount,
		})
		if err != nil {
			return err
		}

		transactionResponse.EntryOut, err = q.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: params.FromAccountID,
			Amount:    -params.Amount,
		})
		if err != nil {
			return err
		}

		transactionResponse.FromAccount, err = q.UpdateAccountBalanceNew(context.Background(), UpdateAccountBalanceNewParams{
			ID:     params.FromAccountID,
			Amount: -params.Amount,
		})
		if err != nil {
			return err
		}

		transactionResponse.ToAccount, err = q.UpdateAccountBalanceNew(context.Background(), UpdateAccountBalanceNewParams{
			ID:     params.ToAccountID,
			Amount: params.Amount,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return transactionResponse, err
}