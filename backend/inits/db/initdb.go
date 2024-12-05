package db

import (
	"fmt"

	"os"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func InitpgDb() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Get the database connection parameters from the environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbSslMode := os.Getenv("DB_SSLMODE")

	error1:=ConnectTODb(dbUser,dbPassword,dbName,dbHost,dbPort,dbSslMode)
	if error1 != nil {
		// If there was an error, log it and exit
		fmt.Println("Error connecting to the database:", err)
		return
	}

	// fmt.Println("This is the PG assessing", PG.Db)

	// If no error, use the db connection (example)

	// fmt.Println("Connected to the database successfully!",value.Db)

	// Don't forget to close the database when don

}
