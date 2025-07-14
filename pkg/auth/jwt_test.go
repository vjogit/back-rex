package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJwt1(t *testing.T) {

	session := uuid.New()
	version := 1

	token1, err := GenerateJWT(123, time.Duration(3600)*time.Second, &session, &version)
	if err != nil {
		t.Errorf("Erreur lors de la génération du JWT: %s", err)
		return
	}

	fmt.Println("JWT généré:", token1)

	token2, err := jwt.Parse(token1, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Erreur lors de la génération du JWT: %s", err)
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if claims, ok := token2.Claims.(jwt.MapClaims); ok && token2.Valid {
		fmt.Println(claims)
	}

}
