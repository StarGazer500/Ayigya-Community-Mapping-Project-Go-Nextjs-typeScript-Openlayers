package main

import (
	"github.com/StarGazer500/ayigya/inits/db"
	 "context"
	 "os"
	 "time"
	 "log"
	// "Ayigya-Community-WebMap-go-and-go-template-geoserver-leaflet/models"
	"fmt"
	"github.com/pressly/goose/v3"
	"github.com/StarGazer500/ayigya/routers"
	"github.com/StarGazer500/ayigya/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	
)

func init() {
	db.InitpgDb()
}

func deinit() {

	db.PG.Db.Close()

}

func main() {
	// Initialize Gin engine

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")

	}

	// Set a context with a timeout for database migrations
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

	// 	// Ensure the database connection is closed when the main function exits
	defer deinit()
	// //
	// 	// Set the Goose migration dialect for PostgreSQL
		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatalf("Failed to set Goose dialect: %v", err)
		}

		// Validate the command-line arguments
		if len(os.Args) < 2 {
			log.Fatal("No Goose Command ")
		}

		// Parse the migration command and arguments
		command := os.Args[1]
		migrationDir := "migrations"

		// Execute the Goose migration command with the provided context and arguments
		if err := goose.RunContext(ctx, command, db.PG.Db, migrationDir, os.Args[2:]...); err != nil {
			log.Fatalf("Goose command failed Invalid Goose Command: %v", err)
		}

		engine := gin.Default()
		engine.Use(middlewares.CorsMiddleware())

		// Set up account routes
		accountGroup := engine.Group("/account")
		routers.UserRoutes(accountGroup)

		mapGroup := engine.Group("/map")
		routers.MapRoutes(mapGroup)

		// Load HTML templates (make sure the template path is correct)
	

		// Start the server
		engine.Run(":8080") // This starts the server on http://localhost:8080

	// 	data,err := models.FindOne(db.PG.Db, models.BuildingTable.TableName, "name", "Salon")
	// 	if err!=nil{
	// 		fmt.Println("Error occured",err)
	// 	}

	// 	fmt.Println(data)

	// _, err := models.PerformOperation(db.PG.Db, models.BuildingTable.TableName, "shape__len", "Less than (<)","110")
	// if err != nil {
	// 	fmt.Println("Error occured",err)
	// }

// 	data,err := models.FindOne(db.PG.Db,models.OtherPolygonTable.TableName, "exact_use", "Car park")
// if err!=nil{
// 	fmt.Println("Error occured",err)
// }

// fmt.Println(data)


	// fmt.Println("results", data)

}
