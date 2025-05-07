package middleware

import (
	"fmt"
	"smart-scene-app-api/common"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

func authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "", common.TokenNotFound, nil))
			return
		}

		arr := strings.Split(authorization, "Bearer ")
		if len(arr) < 2 {
			c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "", common.TokenNotFound, nil))
			return
		}

		tokenString := arr[1]
		if tokenString == "" {
			c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "0", common.TokenNotFound, nil))
			return
		}

		var claims common.UserJWTProfile
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, err.Error(), common.TokenUnAuthorized, nil))
			return
		}
		if claims, ok := token.Claims.(*common.UserJWTProfile); ok && token.Valid {
			//ex := claims.Exp
			if !claims.AdminAccess {
				c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "Not allow action", common.TokenUnAuthorized, nil))
			}
			if claims.Exp > time.Now().Unix() {
				c.Set(common.USER_JWT_KEY, claims)
				c.Next()
			} else {
				c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "Token expired", common.TokenUnAuthorized, nil))
			}
		} else {
			c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "Internal err", common.TokenUnAuthorized, nil))
		}

	}
}

func AdminAuthenticate() gin.HandlerFunc {
	return authenticate()
}
