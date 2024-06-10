package api

import (
	"context"
	"database/sql"
	"fmt"
	"github/inglo357/fingrind_backend/api/handlers"
	"github/inglo357/fingrind_backend/api/requests"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"github/inglo357/fingrind_backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Account struct{
	server *Server
}

func (a Account) router(server *Server){
	a.server = server

	serverGroup := server.router.Group("/account", AuthorizationMiddleware())
	serverGroup.POST("create", a.createAccount)
	serverGroup.GET("", a.listAccounts)
	serverGroup.POST("transfer", a.transfer)
	serverGroup.POST("add-money", a.addMoney)
}

func (a *Account) createAccount(ctx *gin.Context){
	userId, err := utils.GetActiveUserId(ctx)
	if err != nil {
		return
	}

	var accountRequest requests.CreateAccountRequest

	eViewer := gValid.Validator(requests.CreateAccountRequest{})

	if err := ctx.ShouldBindJSON(&accountRequest); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, ctx, eViewer)})
		return
	}

	currency, err := a.server.queries.GetCurrencyByID(context.Background(), accountRequest.CurrencyID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var newAccount db.Account
	
	err = a.server.queries.ExecTx(ctx, func(q *db.Queries) error {
		arg := db.CreateAccountParams{
			UserID: userId,
			CurrencyID: currency.ID,
		}
	
		newAccount, err = a.server.queries.CreateAccount(context.Background(), arg)
		if err != nil{
			if pgErr, ok := err.(*pq.Error); ok{
				if pgErr.Code == "23505"{
					return fmt.Errorf("you already have an account with this currency")
				}
			}
			return err
		}

		accountNumber, err := handlers.GenerateAccountNumber(newAccount.ID, currency)
		if err != nil {
			return err
		}

		updateArgs := db.UpdateAccountNumberParams{
			ID: 				newAccount.ID,
			AccountNumber:    	accountNumber,
		}

		_, err = q.UpdateAccountNumber(context.Background(), updateArgs)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, &newAccount)
}

func (a *Account) listAccounts(ctx *gin.Context){
	userId, err := utils.GetActiveUserId(ctx)
	if err != nil {
		return
	}

	accounts, err := a.server.queries.GetAccountByUserID(context.Background(), userId)

	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (a *Account) transfer(ctx *gin.Context){
	userId, err := utils.GetActiveUserId(ctx)
	if err != nil {
		return
	}

	var transferRequest requests.TransactionRequest

	eViewer := gValid.Validator(requests.TransactionRequest{})

	if err := ctx.ShouldBindJSON(&transferRequest); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, ctx, eViewer)})
		return
	}

	fromAccount, err := a.server.queries.GetAccountByID(context.Background(), transferRequest.FromAccountID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userId != fromAccount.UserID{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		return
	}

	toAccount, err := a.server.queries.GetAccountByID(context.Background(), transferRequest.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if toAccount.CurrencyID != fromAccount.CurrencyID{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		return
	}

	if fromAccount.Balance < transferRequest.Amount{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "insufficient funds"})
		return
	}
	
	args := db.CreateTransferParams{
		FromAccountID: transferRequest.FromAccountID,
		ToAccountID: transferRequest.ToAccountID,
		Amount: transferRequest.Amount,
	}

	transfer, err := a.server.queries.TransferTx(context.Background(), args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, transfer)
}


func (a *Account) addMoney(ctx *gin.Context){
	userId, err := utils.GetActiveUserId(ctx)
	if err != nil {
		return
	}

	var addMoneyRequest requests.AddMoneyRequest

	eViewer := gValid.Validator(requests.AddMoneyRequest{})

	if err := ctx.ShouldBindJSON(&addMoneyRequest); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, ctx, eViewer)})
		return
	}


	account, err := a.server.queries.GetAccountByID(context.Background(), addMoneyRequest.ToAccountID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userId != account.UserID{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "account not found"})
		return
	}

	args := db.CreateMoneyRecordParams{
		UserID: userId,
		Reference: addMoneyRequest.Reference,
		Amount: addMoneyRequest.Amount,
		Status: "pending",
	}


	_, err = a.server.queries.CreateMoneyRecord(context.Background(), args)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok{
			if pgErr.Code == "23505"{
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "money record with this reference already exists"})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	balanceArgs := db.UpdateAccountBalanceNewParams{
		ID: addMoneyRequest.ToAccountID,
		Amount: addMoneyRequest.Amount,
	}

	_, err = a.server.queries.UpdateAccountBalanceNew(context.Background(), balanceArgs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Money added successfully"})
}