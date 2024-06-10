package db_tests

import (
	"context"
	db "github/inglo357/fingrind_backend/db/sqlc"
	"github/inglo357/fingrind_backend/utils"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func clean_up() {
	err := testQuery.DeleteAllUsers(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func createRandomUser(t *testing.T) db.User{

	hashedPassword, err := utils.GenerateHashPassword(utils.RandomString(10))

	if err != nil{
		log.Fatal("Could not generate hashed password ", err)
	}

	createUserArgs := db.CreateUserParams{
		Name: "asd",
		Email: utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	user, err := testQuery.CreateUser(context.Background(), createUserArgs)

	assert.NoError(t, err)
	assert.NotEmpty(t, user)
	assert.Equal(t, user.Name, createUserArgs.Name)
	assert.Equal(t, user.Email, createUserArgs.Email)
	assert.Equal(t, user.HashedPassword, createUserArgs.HashedPassword)
	assert.WithinDuration(t, user.CreatedAt, time.Now(), 10*time.Second)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 10*time.Second)

	return user
}

func TestCreateUser(t *testing.T){
	defer clean_up()
	user1 := createRandomUser(t)

	user2, err := testQuery.CreateUser(context.Background(), db.CreateUserParams{
		Name: user1.Name,
		Email: user1.Email,
		HashedPassword: user1.HashedPassword,
	})

	assert.Error(t, err)
	assert.Empty(t, user2)

}

func TestUpdateUserPassword(t *testing.T){
	defer clean_up()
	user := createRandomUser(t)
	newHashedPassword, err := utils.GenerateHashPassword(utils.RandomString(10))

	if err != nil{
		log.Fatal("Could not generate hashed password ", err)
	}

	updateUserPasswordArgs := db.UpdateUserPasswordParams{
		HashedPassword: newHashedPassword,
		UpdatedAt: time.Now(),
		ID: user.ID,
	}

	updatedUser, err := testQuery.UpdateUserPassword(context.Background(), updateUserPasswordArgs)

	assert.NoError(t, err)
	assert.NotEmpty(t, updatedUser)
	assert.Equal(t, updatedUser.HashedPassword, updateUserPasswordArgs.HashedPassword)
	assert.Equal(t, updatedUser.Email, user.Email)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 10*time.Second)
	// assert.WithinDuration(t, updatedUser.UpdatedAt, updateUserPasswordArgs.UpdatedAt, 10*time.Second)
}

func TestUpdateUserName(t *testing.T){
	defer clean_up()
	user := createRandomUser(t)

	updateUserNameArgs := db.UpdateUsernameParams{
		Name: "asd1",
		UpdatedAt: time.Now(),
		ID: user.ID,
	}

	updatedUser, err := testQuery.UpdateUsername(context.Background(), updateUserNameArgs)

	assert.NoError(t, err)
	assert.NotEmpty(t, updatedUser)
	assert.Equal(t, updatedUser.Name, updateUserNameArgs.Name)
	assert.Equal(t, updatedUser.Email, user.Email)
	assert.WithinDuration(t, user.UpdatedAt, time.Now(), 10*time.Second)
	// assert.WithinDuration(t, updatedUser.UpdatedAt, updateUserNameArgs.UpdatedAt, 10*time.Second)
}

func TestGetUserById(t *testing.T){
	defer clean_up()
	user := createRandomUser(t)

	gotUser, err := testQuery.GetUserByID(context.Background(), user.ID)

	assert.NoError(t, err)
	assert.NotEmpty(t, gotUser)
	assert.Equal(t, gotUser.Name, user.Name)
	assert.Equal(t, gotUser.Email, user.Email)
	assert.Equal(t, gotUser.HashedPassword, user.HashedPassword)
}

func TestGetUserByEmail(t *testing.T){
	defer clean_up()
	user := createRandomUser(t)

	gotUser, err := testQuery.GetUserByEmail(context.Background(), user.Email)

	assert.NoError(t, err)
	assert.NotEmpty(t, gotUser)
	assert.Equal(t, gotUser.Name, user.Name)
	assert.Equal(t, gotUser.Email, user.Email)
	assert.Equal(t, gotUser.HashedPassword, user.HashedPassword)
}

func TestDeleteUserById(t *testing.T){
	defer clean_up()
	user := createRandomUser(t)

	deleteErr := testQuery.DeleteUser(context.Background(), user.ID)
	gotUser, gotUserErr := testQuery.GetUserByID(context.Background(), user.ID)

	assert.NoError(t, deleteErr)
	assert.Error(t, gotUserErr)
	assert.Empty(t, gotUser)
}

func TestListUsers(t *testing.T){
	defer clean_up()

	var wg sync.WaitGroup
	for i:=0; i<50; i++{
		wg.Add(1)
		go func() {
			defer wg.Done()
			createRandomUser(t)
		}()
	}

	wg.Wait()

	listUserArgs := db.ListUsersParams{
		Limit: 50,
		Offset: 0,
	}

	users, err := testQuery.ListUsers(context.Background(), listUserArgs)

	assert.NoError(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, len(users), 50)
}