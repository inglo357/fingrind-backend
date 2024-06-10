package api

import (
	"context"
	"github/inglo357/fingrind_backend/api/handlers"
	"github/inglo357/fingrind_backend/api/requests"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Currency struct{
	server *Server
}

func (c Currency) router(server *Server){
	c.server = server

	serverGroup := server.router.Group("/currency", AuthorizationMiddleware())
	serverGroup.GET("", c.listCurrencies)
	serverGroup.POST("add", c.addCurrency)
	serverGroup.POST("addMultiple", c.addMultipleCurrency)
}

func (c *Currency) listCurrencies(ctx *gin.Context){
	arg := db.ListCurrenciesParams{
		Limit: 10,
		Offset: 0,
	}

	currencies, err := c.server.queries.ListCurrencies(context.Background(), arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, currencies)
}

func (c *Currency) addCurrency(ctx *gin.Context){

	var currencyRequest requests.AddCurrencyRequest

	eViewer := gValid.Validator(requests.AddCurrencyRequest{})

	if err := ctx.ShouldBindJSON(&currencyRequest); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, ctx, eViewer)})
		return
	}

	newCurrency, err := c.server.queries.CreateCurrency(context.Background(), currencyRequest.CurrencyString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, newCurrency)
}

func (c *Currency) addMultipleCurrency(ctx *gin.Context){

	var multipleCurrencyRequest []requests.AddCurrencyRequest

	eViewer := gValid.Validator(requests.AddCurrencyRequest{})

	if err := ctx.ShouldBindJSON(&multipleCurrencyRequest); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, ctx, eViewer)})
		return
	}

	currencyResponse := []db.Currency{}
	var wg sync.WaitGroup

	for _, v := range multipleCurrencyRequest{
		wg.Add(1)
		v := v
		go func()  {
			defer wg.Done()
			currency, err := c.server.queries.CreateCurrency(context.Background(), v.CurrencyString)
		
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
	
			currencyResponse = append(currencyResponse, currency)
		}()
	}	
	
	wg.Wait()

	ctx.JSON(http.StatusCreated, currencyResponse)
}