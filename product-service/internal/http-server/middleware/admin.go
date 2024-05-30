package middleware

import (
	"flower-shop/product-service/internal/config"
	"flower-shop/product-service/internal/lib/jwt"
	"log"
	"net/http"
	"strings"
)

func AdminMiddleware(next http.Handler, cfg config.JwtConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		log.Println("token: ", authHeader)

		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		if !jwt.Verify(tokenString, cfg) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		role := jwt.GetRole(tokenString, cfg)
		if role != "admin" {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
