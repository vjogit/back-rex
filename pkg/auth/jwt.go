package auth

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey []byte

const length = 64

func init() {

	// Ouvrir /dev/urandom
	f, err := os.Open("/dev/urandom")
	if err != nil {
		fmt.Printf("ne peux generer la clef secrete: %s\n", err)
		os.Exit(2)
	}
	defer f.Close()

	// Lire les octets aléatoires
	secretKey = make([]byte, length)
	_, err = io.ReadFull(f, secretKey)
	if err != nil {
		fmt.Printf("ne peux generer la clef secrete: %s\n", err)
		os.Exit(2)
	}

}

func createSubject(userID int) string {
	return "_" + strconv.Itoa(userID)
}

func getSubject(subjet string) (int, error) {
	id, err := strconv.ParseInt(subjet[1:], 10, 64)
	if err != nil {
		return -1, err
	}

	return int(id), nil
}

func GenerateJWT(userID int, duree time.Duration, session *uuid.UUID, version *int) (string, error) {

	var claims jwt.MapClaims = jwt.MapClaims{}

	claims["sub"] = createSubject(userID)
	claims["exp"] = time.Now().Add(duree).Unix() // Expiration.
	claims["iat"] = time.Now().Unix()            // date de creation.

	if session != nil {
		claims["session"] = session.String()
		claims["version"] = version
	}

	// Créer le token avec les revendications et la méthode de signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// Signer le token avec une clé secrète
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type Claims struct {
	UserID  int
	session uuid.UUID
	version int
}

func VerifyJWT(tokenString string) (*Claims, error) {
	var err error

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// S'assurer que la méthode de signature est correcte
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature invalide")
		}
		return secretKey, nil
	})

	// expiration est traite comme une erreur.
	if err != nil {
		return nil, err
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("claims invalide")
	}

	claims := Claims{}

	sub, err := jwtClaims.GetSubject()
	if err != nil {
		return nil, err
	}

	userId, err := getSubject(sub)
	if err != nil {
		return nil, err
	}
	claims.UserID = int(userId)

	session := jwtClaims["session"]
	switch have := session.(type) {
	case string:
		claims.session, err = uuid.Parse(have)
		if err != nil {
			return nil, err
		}
	default:
		claims.session = uuid.UUID{}
	}

	version := jwtClaims["version"]
	switch have := version.(type) {
	case float64:
		claims.version = int(have)
	default:
		claims.version = -1
	}

	return &claims, nil
}

func GenerateJWTCookies(userID int, session *uuid.UUID, version *int) (*http.Cookie, *http.Cookie, error) {
	accessCokie, err := GenerateJWTCookie(userID, "access_token", time.Duration(30)*time.Second, "/", nil, nil)
	if err != nil {
		return nil, nil, err
	}

	refreshCookie, err := GenerateJWTCookie(userID, "refresh_cookie", time.Duration(3600)*time.Second, "/api/v0/auth/refresh", session, version)
	if err != nil {
		return nil, nil, err
	}

	return accessCokie, refreshCookie, nil
}

func GenerateJWTCookie(userID int, name string, duree time.Duration, path string, session *uuid.UUID, version *int) (*http.Cookie, error) {
	token, err := GenerateJWT(userID, duree, session, version)
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{
		Name:     name,
		Path:     path,
		HttpOnly: true,
		//Secure: true,
		Expires:  time.Now().Add(duree),
		Value:    token,
		SameSite: http.SameSiteStrictMode,
	}

	return &cookie, nil
}
