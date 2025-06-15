package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtExpiry = 10 * time.Second
)

func validateJWT(tokenString string, jwtSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims["iat"] == nil || time.Unix(int64(claims["iat"].(float64)), 0).Add(jwtExpiry).Before(time.Now()) {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}

func authMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		if token == "" {
			context.JSON(http.StatusUnauthorized, apiResponse{Success: false, Error: "Unauthorized"})
			context.Abort()
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := validateJWT(token, jwtSecret)
		if err != nil {
			context.JSON(http.StatusUnauthorized, apiResponse{Success: false, Error: err.Error()})
			context.Abort()
			return
		}
		context.Set("claims", claims)
		context.Next()
	}
}
