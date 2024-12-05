package controllers

import (
	// "Ayigya-Community-WebMap-go-and-go-template-geoserver-leaflet/inits/db"
	"github.com/StarGazer500/ayigya/inits/db"
	"github.com/StarGazer500/ayigya/middlewares"
	"github.com/StarGazer500/ayigya/models"

	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserDetails struct {
	Firstname string `json:"firstname" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	Password1 string `json:"password1" binding:"required"`
	Password2 string `json:"password2" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type ProfileDeTails struct {
	Firstname string `json:"firstname" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	Password1 string `json:"password1" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

func Profile(ctx *gin.Context) {

	if ctx.Request.Method == http.MethodGet {

		// if err := ctx.ShouldBind(&form); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		// }

		ctx.HTML(http.StatusOK, "profile.html", gin.H{"profilepage": "opened"})
	}

}

func SaveUser(db *sql.DB, user UserDetails) (sql.Result, error) {

	tableName := models.UserTable.TableName

	columns := []string{"firstname", "surname", "email", "password1"}

	data, err := models.InsertOne(db, tableName, columns, user.Firstname, user.Surname, user.Email, user.Password1)

	if err != nil {

		fmt.Println("Insertion Error Occured", err)
		return nil, err
	}

	return data, nil
}

// Binding from JSON
type Login struct {
	Email     string `json:"email" binding:"required"`
	Password1 string `json:"password1" binding:"required"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
}

func Register(ctx *gin.Context) {
	// Default empty data for the registration form
	var form UserDetails
	defaultData := UserDetails{
		Firstname: "",
		Surname:   "",
		Password1: "",
		Password2: "",
		Email:     "",
	}

	if ctx.Request.Method == http.MethodPost {

		if err := ctx.ShouldBind(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}

		hashed, err := hashPassword(form.Password1)
		if err != nil {
			fmt.Println("error when hashing")
		}

		form.Password1 = hashed

		res, err := SaveUser(db.PG.Db, form)
		if err != nil {

			if pqErr, ok := err.(*pq.Error); ok {

				fmt.Println("error name is", pqErr.Code.Name())
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "this emaiil is already registered"})
				return

			}
			fmt.Println("db error", err)
			return
		}
		fmt.Println("saved data", res)

		ctx.JSON(http.StatusOK, gin.H{"status": "You have registered Successfully"})

	} else {
		// If the request is not a POST, show the registration form with empty data
		ctx.HTML(http.StatusOK, "auth.html", gin.H{
			"data": defaultData,
		})
	}
}

// Helper function to hash password
func hashPassword(password string) (string, error) {
	// Implement password hashing here
	// Example using bcrypt:
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func comparePassword(requestPassword string, dbPassword string) error {
	// Compare the provided password with the hashed password
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(requestPassword))
	if err != nil {
		// If there's an error (passwords don't match), return false and the error
		return err
	}
	// If no error (passwords match), return true and nil error
	return nil
}

func LoginUser(ctx *gin.Context) {
	// Default empty data for the login form
	var form Login
	defaultData := Login{
		Email:     "",
		Password1: "",
	}
	// body, err := io.ReadAll(ctx.Request.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(string(body))

	if ctx.Request.Method == http.MethodPost {
		// Bind the JSON request body to the form struct
		if err := ctx.ShouldBindJSON(&form); err != nil {
			// If there's a binding error, return a 400 response and stop executionf
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Query the database to find the user with the given email
		fmt.Println(db.PG.Db, models.UserTable.TableName, form.Email)
		result, err := models.FindOne(db.PG.Db, models.UserTable.TableName, "email", form.Email)
		if err != nil {
			// If there's an error querying the database, return a 400 response
			fmt.Println("Querying error occurred", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
			return
		}

		if len(result) == 0 {
			// If no user is found, return a 400 response
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		// If the user is found, set up the token claims
		for _, user := range result {
			middlewares.Claim = middlewares.TokenClaimStruct{
				MyAuthServer:    "AuthServer",
				AuthUserEmail:   user["email"].(string),
				AuthUserSurname: user["surname"].(string),
				AuthUserId:      user["id"].(string),
			}
		}

		// Get stored hashed password from the database
		hashedPassword := result[0]["password1"].(string)

		passerr := comparePassword(form.Password1, hashedPassword)
		if passerr != nil {
			fmt.Println("Password is Invalid")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password is Invalid"})
			return
		}
		// Generate access and refresh tokens
		access_token, aerr := middlewares.GenerateAccessToken(middlewares.Claim)
		refresh_token, rerr := middlewares.GenerateRefreshToken(middlewares.Claim)
		if aerr != nil || rerr != nil {
			fmt.Println("token generation error")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		// Set the cookies for access and refresh tokens
		ctx.SetCookie("access", access_token, 3600, "/", "localhost", false, true)
		ctx.SetCookie("refresh", refresh_token, 3600, "/", "localhost", false, true)

		ctx.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Login successful!",
			"redirect": "/account/profile", // Include the URL to redirect to
		})
		return
	}

	// If the request is not a POST, show the login form with empty data
	ctx.HTML(http.StatusOK, "auth.html", gin.H{
		"data": defaultData,
	})
}
