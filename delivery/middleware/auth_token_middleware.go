package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/utils/security"
	"strings"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}
type AuthTokenMiddleware interface {
	RequireToken() gin.HandlerFunc
}
type authTokenMiddleware struct {
	tokenService security.AccessToken
}

func NewTokenValidator(tokenService security.AccessToken) AuthTokenMiddleware {
	return &authTokenMiddleware{
		tokenService: tokenService,
	}
}

func (a *authTokenMiddleware) RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
		if tokenString == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		token, err := a.tokenService.VerifyAccessToken(tokenString)
		fmt.Println("err:", err)
		if err != nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		fmt.Println(token)
		if token != nil {
			c.Next()
		} else {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
	}
}
