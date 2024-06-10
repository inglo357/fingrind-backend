package api

import (
	"context"
	"database/sql"
	"github/inglo357/fingrind_backend/api/handlers"
	"github/inglo357/fingrind_backend/api/requests"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"github/inglo357/fingrind_backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Auth struct{
	server *Server
}

func (a Auth) router(server *Server){
	a.server = server

	serverGroup := server.router.Group("/auth")
	serverGroup.POST("/login", a.login)
	serverGroup.POST("/register", a.register)
}

type UserParams struct{
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (a *Auth) register(c *gin.Context){

	var user UserParams

	eViewer := gValid.Validator(UserParams{})

	if err := c.ShouldBindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, c, eViewer)})
		return
	}

	hashedPassword, err := utils.GenerateHashPassword(user.Password)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	arg := db.CreateUserParams{
		Email: user.Email,
		HashedPassword: hashedPassword,
		Name: user.Name,
	}

	newUser, err := a.server.queries.CreateUser(context.Background(), arg)

	if err != nil{
		if pgErr, ok := err.(*pq.Error); ok{
			if pgErr.Code == "23505"{
				c.JSON(http.StatusBadRequest, "user already exists")
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{}.toUserResponse(&newUser))
}

func (a *Auth) login(c *gin.Context){

	var userParams requests.UserLoginRequest

	eViewer := gValid.Validator(UserParams{})

	if err := c.ShouldBindJSON(&userParams); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, c, eViewer)})
		return
	}

	user, err := a.server.queries.GetUserByEmail(context.Background(), userParams.Email)
	if err == sql.ErrNoRows{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect Email or Password"})
		return
	} else if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.CompareHashAndPassword(user.HashedPassword, userParams.Password); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := TokenHandler.CreateToken(user.ID)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}