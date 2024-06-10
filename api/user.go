package api

import (
	"context"
	"database/sql"
	"github/inglo357/fingrind_backend/api/handlers"
	"github/inglo357/fingrind_backend/api/requests"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"github/inglo357/fingrind_backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct{
	server *Server
}

func (u User) router(server *Server){
	u.server = server

	serverGroup := server.router.Group("/users", AuthorizationMiddleware())
	serverGroup.GET("", u.listUsers)
	serverGroup.GET("me", u.getLoggedInUser)
	serverGroup.PATCH("username", u.updateUsername)
	serverGroup.PATCH("password", u.updatePassword)
}

func (u *User) listUsers(c *gin.Context){
	arg := db.ListUsersParams{
		Limit: 10,
		Offset: 0,
	}

	users, err := u.server.queries.ListUsers(context.Background(), arg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	usersResponse := []UserResponse{}

	for _, v := range users{
		n := UserResponse{}.toUserResponse(&v)
		usersResponse = append(usersResponse, *n)
	}

	c.JSON(http.StatusOK, usersResponse)
}

func (u *User) getLoggedInUser(c *gin.Context){

	userId, exists := c.Get("user_id")

	if !exists{
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	convertedUserId, ok := userId.(int64)

	if !ok{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Encountered an issue"})
		return
	}

	user, err := u.server.queries.GetUserByID(context.Background(), convertedUserId)
	if err == sql.ErrNoRows{
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	} else if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))
}

func (u *User) updateUsername(c *gin.Context){
	userId, err := utils.GetActiveUserId(c)
	if err != nil {
		return
	}

	var usernameRequest requests.UpdateUsernameRequest

	eViewer := gValid.Validator(requests.UpdateUsernameRequest{})
	
	if err := c.ShouldBindJSON(&usernameRequest); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, c, eViewer)})
		return
	}

	user, err := u.server.queries.UpdateUsername(context.Background(), db.UpdateUsernameParams{
		ID:     	userId,
		Name:    	usernameRequest.NewUsername,
		UpdatedAt: 	time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))
}

func (u *User) updatePassword(c *gin.Context){
	userId, err := utils.GetActiveUserId(c)
	if err != nil {
		return
	}

	var passwordRequest requests.UpdatePasswordRequest

	eViewer := gValid.Validator(requests.UpdatePasswordRequest{})

	if err := c.ShouldBindJSON(&passwordRequest); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": handlers.HandleError(err, c, eViewer)})
		return
	}	

	hashedPassword, err := utils.GenerateHashPassword(passwordRequest.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := u.server.queries.UpdateUserPassword(context.Background(), db.UpdateUserPasswordParams{
		ID:     		userId,
		HashedPassword: hashedPassword,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{}.toUserResponse(&user))
}
type UserResponse struct{
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u UserResponse) toUserResponse(user *db.User) *UserResponse{
	return &UserResponse{
		ID:          	user.ID,
		Name:           user.Name,
		Email:          user.Email,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}