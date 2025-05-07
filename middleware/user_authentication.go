package middleware

import (
	"fmt"
	"smart-scene-app-api/common"
	"smart-scene-app-api/config"
	"smart-scene-app-api/server"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

var secretKey = []byte(config.Config.JwtSecret)

func UserAuthentication() gin.HandlerFunc {
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
			if claims.Exp > time.Now().Unix() {
				c.Set(common.USER_JWT_KEY, claims)
				c.Set(common.UserId, claims.Id)
				c.Next()
			} else {
				c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "Token expired", common.TokenUnAuthorized, nil))
			}
		} else {
			c.AbortWithStatusJSON(common.UNAUTHORIZED_STATUS, common.BaseResponse(common.UNAUTHORIZED_STATUS, "2", common.TokenUnAuthorized, nil))
		}

	}
}

func OptionalUserAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization != "" {
			arr := strings.Split(authorization, "Bearer ")
			if len(arr) >= 2 {
				tokenString := arr[1]
				if tokenString != "" {
					var claims common.UserJWTProfile
					token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
						if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
							return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
						}
						return secretKey, nil
					})

					if err == nil {
						if claims, ok := token.Claims.(*common.UserJWTProfile); ok && token.Valid {
							if claims.Exp > time.Now().Unix() {
								c.Set(common.USER_JWT_KEY, claims)
							}
						}
					}

				}

			}

		}
		c.Next()

	}
}

type Authenticator struct {
	authConfig *server.AuthorizationConfig
}

func NewAuthenticator(authConfig *server.AuthorizationConfig) Authenticator {
	return Authenticator{authConfig: authConfig}
}

func (a Authenticator) ACLAuthentication(actionId string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			common.AbortWithError(c, common.ErrTokenNotFound)
			return
		}

		arr := strings.Split(authorization, "Bearer ")
		if len(arr) < 2 {
			common.AbortWithError(c, common.ErrTokenNotFound)
			return
		}
		tokenString := arr[1]
		if tokenString == "" {
			common.AbortWithError(c, common.ErrTokenNotFound)
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
			common.AbortWithError(c, common.ErrNotAuthorized)
			return
		}
		if claims, ok := token.Claims.(*common.UserJWTProfile); ok && token.Valid {
			if claims.Exp > time.Now().Unix() {
				err := a.authConfig.CheckValidValidRole(claims.Role, claims.ID, actionId)
				if err != nil {
					common.AbortWithError(c, common.ErrActionNotAllowed)
					return
				}
				c.Set(common.USER_JWT_KEY, claims)
				c.Next()
			} else {
				common.AbortWithError(c, common.ErrNotAuthorized)
				return
			}
		} else {
			common.AbortWithError(c, common.ErrNotAuthorized)
			return
		}

	}
}
