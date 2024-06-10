package main

import (
	"github/inglo357/fingrind_backend/api"
	db "github/inglo357/fingrind_backend/db/sqlc"

	"github.com/gin-gonic/gin"
)

type Server struct{
	queries *db.Queries
	router *gin.Engine
}

func main(){
	server := api.NewServer(".")
	server.Start(8000)
}