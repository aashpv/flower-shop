package jwt

import (
	"flower-shop/user-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

func NewToken(id int, role string, cfg config.JwtConfig) (string, error) {
	claims := jwt.MapClaims{
		"uid":  id,
		"role": role,
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// not used here
func GetData(tokenString string, config config.JwtConfig) string {
	secretKey := []byte(config.SecretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return ""
	}

	if !token.Valid {
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	uid := claims["uid"].(float64)
	return strconv.Itoa(int(uid))
}

func Verify(tokenString string, config config.JwtConfig) bool {
	secretKey := []byte(config.SecretKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}
	return true
}
