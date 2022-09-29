package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
)

func AuthMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString := c.Value
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		claims, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		userID := int(claims.(jwt.MapClaims)["user_id"].(float64))
		fmt.Printf("userID: %v of type %T\n", userID, userID)

		w.Header().Set("user_id", strconv.Itoa(userID))
		next.ServeHTTP(w, r)
	})
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("KSFjH$53KSFjH6745u#uEQQjF349%835hFpzA")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
