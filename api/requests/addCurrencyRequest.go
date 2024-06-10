package requests

type AddCurrencyRequest struct{
	CurrencyString string `json:"currency_string" binding:"required"`
}