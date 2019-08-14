package tools

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	HasuraClaims map[string]interface{} `json:"https://hasura.io/jwt/claims"`
	jwt.StandardClaims
}

func GenerateToken(id string) (string, error) {
	nowTime := time.Now()
	// expireTime := nowTime.Add(3 * time.Hour)
	expireTime := nowTime.Add(10000 * time.Hour)
	a := make([]string, 1)
	a[0] = "user"
	claims := &Claims{
		map[string]interface{}{
			"x-hasura-default-role":  "user",
			"x-hasura-allowed-roles": a,
			"x-hasura-user-id":       id,
		},
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "test",
		},
	}
	fmt.Println(claims.HasuraClaims)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	// fmt.Println("[jwt.Valid]", tokenClaims.Valid)
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		} else {
			fmt.Println("[parse error]", err)
		}
	}

	return nil, err
}
