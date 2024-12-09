package middlewares

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
)

func CheckTokenExpiration(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {

		expirationTime := time.Unix(int64(exp), 0)

		if time.Now().After(expirationTime) {
			return fmt.Errorf("token has expired at %s", expirationTime)
		}

		fmt.Printf("Token is valid. Expiration time: %s\n", expirationTime)
		return nil
	} else {
		return fmt.Errorf("missing expiration claim in token")
	}
}

func VerifyAccessToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing access token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if err := CheckTokenExpiration(claims); err != nil {
			fmt.Println("Access Expiry has happened")
			return nil, err
		}

		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid access token")
	}
}

func VerifyRefreshToken(tokenString string, secretKey []byte) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing refresh token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if err := CheckTokenExpiration(claims); err != nil {
			fmt.Println("Refresh Expiry has happened")
			return nil, err
		}

		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid refresh token")
	}
}

func LoadKeyAndReturnByte() ([]byte, error) {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		return nil, fmt.Errorf("missing JWT secret key in the environment")
	}
	return []byte(key), nil
}

func GenerateExpiryTime(durationInSeconds int64) int64 {
	expirationTime := time.Now().Add(time.Duration(durationInSeconds) * time.Second)
	return expirationTime.Unix()
}

type TokenClaimStruct struct {
	MyAuthServer    string
	AuthUserEmail   string
	AuthUserSurname string
	AuthUserId      string
	AuthExp         int64
}

var Claim TokenClaimStruct

func GenerateAccessToken(claim TokenClaimStruct) (string, error) {
	key, err := LoadKeyAndReturnByte()
	if err != nil {
		return "", fmt.Errorf("error loading the secret key: %w", err)
	}
	accessTokenExpiry := GenerateExpiryTime(120)
	claims := jwt.MapClaims{
		"iss":     claim.MyAuthServer,    // Issuer
		"sub":     claim.AuthUserEmail,   // Subject (user email)
		"surname": claim.AuthUserSurname, // User surname
		"id":      claim.AuthUserId,      // User ID
		"exp":     accessTokenExpiry,     // for 120 seconds
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error signing the access token: %w", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken(claim TokenClaimStruct) (string, error) {

	key, err := LoadKeyAndReturnByte()
	if err != nil {
		return "", fmt.Errorf("error loading the secret key: %w", err)
	}

	refreshTokenExpiry := GenerateExpiryTime(86400 * 30) // 30 days * 24 hours * 60 minutes * 60 seconds

	claims := jwt.MapClaims{
		"iss":     claim.MyAuthServer,    // Issuer
		"sub":     claim.AuthUserEmail,   // Subject (user email)
		"surname": claim.AuthUserSurname, // User surname
		"id":      claim.AuthUserId,      // User ID
		"exp":     refreshTokenExpiry,    // Expiration time (Unix timestamp)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and get the signed string
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("error signing the refresh token: %w", err)
	}

	return tokenString, nil
}

// CORS middleware function definition
func CorsMiddleware() gin.HandlerFunc {
	// Define allowed origins as a comma-separated string
	originsString := "http://localhost:3000"
	var allowedOrigins []string
	if originsString != "" {
		// Split the originsString into individual origins and store them in allowedOrigins slice
		allowedOrigins = strings.Split(originsString, ",")
	}

	// Return the actual middleware handler function
	return func(c *gin.Context) {
		// Function to check if a given origin is allowed
		isOriginAllowed := func(origin string, allowedOrigins []string) bool {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					fmt.Println(origin, allowedOrigin)
					return true
				}
			}
			return false
		}

		// Get the Origin header from the request
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		if isOriginAllowed(origin, allowedOrigins) {
			// If the origin is allowed, set CORS headers in the response
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
			fmt.Println("header",c.Request.Header)
		}

		

		// Handle preflight OPTIONS requests by aborting with status 204
		if c.Request.Method == "OPTIONS" {
			
			
			c.AbortWithStatus(204)
			return
		}

		

		// Call the next handler
		c.Next()
	}
}
