package utils

import (

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error){
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil{
		return "", err
	}

	return string(hash), nil
}

func CompareHashAndPassword(hashedPassword string, password string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func VerifyPassword(password string) bool{
	if len(password) < 8{
		return false
	}

	return true
}