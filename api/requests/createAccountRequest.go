package requests

type CreateAccountRequest struct{
	CurrencyID int64   `json:"currency_id" binding:"required"`
}

type AddMoneyRequest struct{
	ToAccountID int64 `json:"to_account_id" binding:"required"`
	Reference string `json:"reference" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}