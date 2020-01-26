package user

import (
	"demo-plaform/model/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type SignedData struct {
	Ac int64
	N string
	Ava string
}

func SignedToken(user *db.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"ac": user.Id,
		"n": user.LoginName,
	})
	tokenString, err := token.SignedString("beta_"+user.Salt)
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*SignedData, error)  {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("my_secret_key"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])

		return &SignedData{
			Ac:  claims["ac"].(int64),
			N:   fmt.Sprintf("%v", claims["n"]),
			Ava: "",
		},nil
	} else {
		return nil, err
	}
}