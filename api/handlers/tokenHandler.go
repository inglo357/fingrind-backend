package handlers

import (
	"fmt"
	"github/inglo357/fingrind_backend/utils"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTToken struct{
	config *utils.Config
}

type jwtClaim struct{
	jwt.StandardClaims
	UserID 	int64 `json:"user_id"`
	Exp 	int64 `json:"exp"`
}

func NewJWTToken(config *utils.Config) *JWTToken{
	return &JWTToken{config: config}
}

func (j *JWTToken) CreateToken(userId int64) (string, error){
	claims := jwtClaim{
		UserID: userId,
		Exp: time.Now().Add(time.Minute * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(j.config.Signing_key))

	if err != nil{
		return "", err
	}

	return string(tokenString), nil
}

func (j *JWTToken) VerifyToken(tokenString string) (int64, error){
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid authentication token")
		}
		return []byte(j.config.Signing_key), nil
	})

	if err != nil{
		return 0, fmt.Errorf("invalid authtentication token case2: %v", err)
	}

	claims, ok := token.Claims.(*jwtClaim)

	if !ok{
		return 0, fmt.Errorf("invalid authtentication token case3: %v", j.config.Signing_key)
	}

	if claims.Exp < time.Now().Unix(){
		return 0, fmt.Errorf("authtentication token expired")
	}

	return claims.UserID, nil
}