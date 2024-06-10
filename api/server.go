package api

import (
	"database/sql"
	"fmt"
	"github/inglo357/fingrind_backend/api/handlers"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"github/inglo357/fingrind_backend/utils"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golodash/galidator/v2"
	_ "github.com/lib/pq"
)

type Server struct{
	queries *db.Store
	router 	*gin.Engine
	config 	*utils.Config
}

var TokenHandler *handlers.JWTToken
var gValid = galidator.New()

func myCorsHandler() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	return cors.New(config)
}

func NewServer(envPath string) *Server{
	config, err := utils.LoadConfig(envPath)

	if err != nil{
		panic(fmt.Sprintf("Could not load configuration from env %v", err))
	}

	conn, err := sql.Open(config.DB_driver, config.DB_source + config.DB_name + "?sslmode=disable")
	
	if err != nil{
		panic(fmt.Sprintf("Could not connect to db %v", err))
	}

	TokenHandler = handlers.NewJWTToken(config)

	q := db.NewStore(conn)

	g := gin.Default()

	g.Use(myCorsHandler())

	return &Server{
		queries: 		q,
		router: 		g,
		config: 		config,
	}
}

func (s *Server) Start(port int){
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "dzialczy"})
	})

	User{}.router(s)
	Auth{}.router(s)
	Currency{}.router(s)
	Account{}.router(s)

	s.router.Run(fmt.Sprintf(":%v", port))
}