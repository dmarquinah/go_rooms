package crypto

import (
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

func GenerateJWT(id int) *string {
	secretKey := getSecretKey()
	token := jwt.New(jwt.SigningMethodHS256)
	ttl := 24 * time.Hour

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().UTC().Add(ttl).Unix()
	claims["id"] = id

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing the token.", err)
		return nil
	}

	return &tokenString
}

func ValidateJWT(tokenString string, key string) bool {
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

	id := GetIdFromJWT(tokenString)

	switch key {
	case "id":
		return int(claims["id"].(float64)) == *id
	default:
		return false
	}
}

func GetIdFromJWT(tokenString string) *int {
	token, err := jwt.Parse(tokenString, tokenValidator)

	if err != nil {
		return nil
	}

	if !token.Valid {
		return nil
	}

	// When token is valid
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil
	}

	id := int(claims["id"].(float64))

	return &id
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
