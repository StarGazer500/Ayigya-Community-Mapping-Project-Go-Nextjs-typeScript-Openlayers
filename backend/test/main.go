package main

import (
	"github.com/StarGazer500/ayigya/inits/db"
	"github.com/joho/godotenv"

	"github.com/StarGazer500/ayigya/models"
	"fmt"

)


func init() {
	db.InitpgDb()
}

func deinit() {

	db.PG.Db.Close()

}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")

	}
	defer deinit()
	// _, err := models.PerformOperation(db.PG.Db, models.BuildingTable.TableName, "shape__len", "Greater than (>)","5")
	// if err != nil {
	// 	fmt.Println("Error occured",err)
	// }


	_,err := models.SearchAllTables(db.PG.Db,"car park")
	if err!=nil{
		fmt.Println("error occured",err)
	}else{
		// fmt.Println(data)
	}

	
}
