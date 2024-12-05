// data, err := models.AddColumnIfNotExists(db.PG.Db, models.UserModel.TableName, "phone2", "VARCHAR(255)")
// if err != nil {
// 	fmt.Println("adding table error occured", err)
// }
// fmt.Println(data)

// data, err :=models.DeleteColumnIfExists(db.PG.Db, models.UserModel.TableName, "phone1")
// 	if err != nil {
// 		fmt.Println("adding table error occured", err)
// 	}
// 	fmt.Println(data)

// err := models.DeleteRowByID(db.PG.Db, models.UserModel.TableName, 1)
// 	if err != nil {
// 		fmt.Println("adding table error occured", err)
// 	}
// 	fmt.Println(err)

// err := models.DeleteRowByColumn(db.PG.Db, models.UserModel.TableName, "email", "ab@gmail.com")

// 	if err != nil {
// 		fmt.Println("adding table error occured", err)
// 	}
// 	fmt.Println(err)

// tableName := models.UserModel.TableName

// 	columns := []string{"firstname", "surname", "email", "password1"}

// 	data, err := models.InsertOne(db, tableName, columns, user.Firstname, user.Surname, user.Email, user.Password1)

// createTableQuery := `
// 		CREATE TABLE IF NOT EXISTS "UserSQLModel" (
// 			id SERIAL PRIMARY KEY,
// 			firstname VARCHAR(100),
// 			surname VARCHAR(100),
// 			password1 VARCHAR(100),
// 			email VARCHAR(100)
// 		);
// 	`

// 	UserModel = &UserTable{TableName: "UserSQLModel"}
// 	CreateTable(createTableQuery)

// columns := []string{"firstname", "surname", "email", "password1"}
// values := []interface{}{"Mawuli", "Kwadwo", "email@gmail.com", "1111"}

// data, err := models.UpdateOne(db.PG.Db, models.UserModel.TableName, columns, values, "email", "ab@gmail.com")
// if err != nil {
// 	fmt.Println("error occured", err)
// }

// fmt.Println("update result",data)

// columns := []string{"firstname", "surname", "email", "password1"}
// rows := [][]interface{}{
// {"John", "Doe", "john.doe@example.com", "password123"},
// {"Jane", "Doe", "jane.doe@example.com", "password456"},
// {"Alice", "Smith", "alice.smith@example.com", "password789"},
// }

// result, err := models.InsertMany(db.PG.Db, models.UserModel.TableName, columns, rows)
// if err != nil {
// fmt.Println("Error inserting multiple rows: %v", err)
// } else {
// fmt.Println("Rows inserted successfully.",result)
// }

// data,err := models.FindOne(db.PG.Db, models.OtherPolygonModel.TableName, "exact_use", "Car park")
// if err!=nil{
// 	fmt.Println("Error occured",err)
// }

// fmt.Println(data)
