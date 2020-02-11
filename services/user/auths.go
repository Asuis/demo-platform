package user

import (
	"demo-platform/model/db"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type SignedData struct {
	Ac int64
	N string
	Ava string
}

func SignedToken(user *db.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ac": user.Id,
		"n": user.LoginName,
	})
	tokenString, err := token.SignedString([]byte("demo_platform"))
	if err != nil {
		return tokenString, err
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*SignedData, error)  {

	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	} else {
		return nil, fmt.Errorf("Token string format is error: %s", tokenString)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("demo_platform"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		return &SignedData{
			Ac:  int64(claims["ac"].(float64)),
			N:   fmt.Sprintf("%v", claims["n"]),
			Ava: "",
		},nil
	} else {
		return nil, err
	}
}