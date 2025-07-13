package middleware

import (
	"fmt"
	"net/http"
	"os"
	"smart-scene-app-api/common"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv(common.ENVJWTSecretKey))

func getRawToken(c *gin.Context) string {
	rawToken := c.GetHeader("Authorization")
	return strings.TrimSpace(strings.Replace(rawToken, "Bearer", "", -1))
}

func validateAndSetTokenToCtx(c *gin.Context, tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &common.JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return false, common.ErrCodeNotAuthorized
	}
	if _, ok := token.Claims.(*common.JWTCustomClaims); ok && token.Valid {
		c.Set("token", token)
		return true, nil
	}
	return false, common.ErrCodeNotAuthorized
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken := getRawToken(c)
		if rawToken == "" {
			if c.FullPath() == "" {
				c.Next()
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.AllErrors.New(common.ErrCodeNotAuthorized, "vi"))
			return
		}
		isValid, err := validateAndSetTokenToCtx(c, rawToken)
		if err != nil || !isValid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.AllErrors.New(common.ErrCodeNotAuthorized, "vi"))
			return
		}
		// Token is valid, continue with the request
		c.Next()
	}
}
