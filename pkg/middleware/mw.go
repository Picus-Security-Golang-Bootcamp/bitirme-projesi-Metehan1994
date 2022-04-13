package mw

import (
	"fmt"
	"time"

	jwtHelper "github.com/Metehan1994/final-project/pkg/jwt"
	"github.com/gin-gonic/gin"

	"net/http"
)

func TokenExpControlMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				c.Set("user", decodedClaims)
				fmt.Println(decodedClaims.Exp)
				jwtTime := time.Unix(decodedClaims.Exp, 0)
				timeNow := time.Now()
				fmt.Println(jwtTime)
				if timeNow.Before(jwtTime) {
					c.Next()
					c.Abort()
					return
				}
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "Your token is not available or expired. You need to login."})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to connect system!"})
			c.Abort()
			return
		}
	}
}

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			if decodedClaims != nil {
				for _, role := range decodedClaims.Roles {
					if role == "admin" {
						c.Next()
						c.Abort()
						return
					}
				}
			}

			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to use this endpoint!"})
			c.Abort()
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to make internal changes!"})
			c.Abort()
			return
		}
	}
}
