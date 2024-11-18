package crypto

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func getSecretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func tokenValidator(t *jwt.Token) (interface{}, error) {
	secretKey := getSecretKey()

	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("sign method not valid")
	}
	return secretKey, nil
}

func GenerateJWT(id string, targetRole string) *string {
	secretKey := getSecretKey()
	token := jwt.New(jwt.SigningMethodHS256)
	ttl := 24 * time.Hour

	claims := token.Claims.(jwt.MapClaims)
	// Add expiration time to JWT
	claims["exp"] = time.Now().UTC().Add(ttl).Unix()
	// Set ID to identify client
	claims["id"] = id
	claims["role"] = targetRole

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing the token.", err)
		return nil
	}

	return &tokenString
}

func ValidateJWT(tokenString string) bool {
	token, err := jwt.Parse(tokenString, tokenValidator)

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	// When token is valid
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return false
	}
	// For now we validate if those two values exists
	return claims["id"] != nil && claims["role"] != nil
}

func GetFieldFromJWT(tokenString string, field string) (string, error) {
	token, err := jwt.Parse(tokenString, tokenValidator)

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", err
	}

	// When token is valid
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", errors.New("error parsing auth token")
	}

	data := claims[field]

	if data == nil {
		return "", errors.New("invalid auth token")
	}

	return data.(string), nil
}

func GetJWTFromRequest(w http.ResponseWriter, r *http.Request) *string {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		http.Error(w, "'Authorization' Header missing.", http.StatusUnauthorized)
		return nil
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "'Authorization' Header format incorrect.", http.StatusUnauthorized)
		return nil
	}

	return &parts[1]
}
