package requests

type TransactionRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required"`
	ToAccountID   int64 `json:"to_account_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
}