package middlewares

import (
	"fmt"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example: Check if user is authenticated using a cookie
		access_token, err := c.Cookie("access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no access token or token expired"})
			c.Abort() // Stop further processing of the request
			return
		}

		verifyaccess, err := VerifyAccessToken(access_token, []byte(os.Getenv("JWT_SECRET_KEY")))
		if err != nil {

			fmt.Println("err", verifyaccess)
			// If user cookie is not set, return a 401 Unauthorized error
			refresh_token, err := c.Cookie("refresh")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token or token expired"})
				c.Abort() // Stop further processing of the request
				return
			}

			verifyrefresh, err := VerifyRefreshToken(refresh_token, []byte(os.Getenv("JWT_SECRET_KEY")))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthor"})
				c.Abort() // Stop further processing of the request
				return

			}
			newaccess, err := GenerateAccessToken(Claim)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "no claim payload"})
				c.Abort() // Stop further processing of the request
			}
			c.SetCookie("access", newaccess, 3600, "/", "localhost", false, true)
			fmt.Println("newaccess", newaccess)
			fmt.Println("newaccess", verifyrefresh)
			c.Next()

		}

		c.Next()

		fmt.Println(verifyaccess)
	}

	// If the user is authenticated, proceed to the next handler

}
