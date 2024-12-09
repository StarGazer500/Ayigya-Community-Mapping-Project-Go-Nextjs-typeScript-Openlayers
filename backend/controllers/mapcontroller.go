package controllers

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/StarGazer500/ayigya/inits/db"
	"github.com/StarGazer500/ayigya/models"

	"github.com/gin-gonic/gin"
)

func MapPageDisplay(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodGet {
		// Define the response structure

		// Render the HTML page and inject the data as a JavaScript variable
		ctx.HTML(http.StatusOK, "map.html", gin.H{})
	}
}

func FeatureLayers(ctx *gin.Context) {

	if ctx.Request.Method == http.MethodGet {

		// Define the response structure
		availfeatlarys := struct {
			BuildingTable     string `json:"buildingTable"`
			OtherPolygonTable string `json:"otherPolygonTable"`
		}{
			BuildingTable:     models.BuildingTable.TableName, // Assuming these are valid
			OtherPolygonTable: models.OtherPolygonTable.TableName,
		}
		fmt.Println(availfeatlarys)
		// Send the JSON response
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    availfeatlarys, // Send the data here
		})
	}
}

func FeatreAttributes(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodPost {
		// Extract the raw request body
		rawData, err := ctx.GetRawData()
		if err != nil {
			// If there's an error reading the body, return a bad request response
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to read the request body",
			})
			return
		}

		// Define a map to hold the raw JSON data
		var requestData map[string]interface{}

		// Unmarshal the raw JSON data into the map
		if err := json.Unmarshal(rawData, &requestData); err != nil {
			// If there's an error unmarshalling the data, return an error response
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			return
		}

		fmt.Println("received data", requestData)
		fmt.Println(models.BuildingTable)

		// Extract selectedLayer from requestData
		selectedLayer, exists := requestData["selectedLayer"].(string)
		if !exists {
			// If no selectedLayer is found, return an error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "selectedLayer field is required",
			})
			return
		}

		// Check if the selectedLayer is "building" or "other_polygon_structure"
		if selectedLayer == models.BuildingTable.TableName {
			// Marshal BuildingTable into JSON
			data, err := json.Marshal(models.BuildingTable)
			if err != nil {
				// If there's an error marshaling BuildingTable, return an error
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error serializing BuildingTable",
				})
				return
			}

			// Send the serialized BuildingTable as JSON response
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "received",
				"data":    json.RawMessage(data), // Directly send the marshaled data as a JSON response
			})
			return
		}

		if selectedLayer == models.OtherPolygonTable.TableName {
			// Marshal OtherPolygonTable into JSON
			data, err := json.Marshal(models.OtherPolygonTable)
			if err != nil {
				// If there's an error marshaling OtherPolygonTable, return an error
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error serializing OtherPolygonTable",
				})
				return
			}

			// Send the serialized OtherPolygonTable as JSON response
			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "received",
				"data":    json.RawMessage(data), // Directly send the marshaled data as a JSON response
			})
			return
		}

		// If selectedLayer is neither "building" nor "other_polygon_structure"
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid layer type. Expected 'building' or 'other_polygon_structure'.",
		})
	} else {
		// If the request method is not POST, return an error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid HTTP method, expected POST",
		})
	}
}

func SelectOperator(ctx *gin.Context) {

	if ctx.Request.Method == http.MethodPost {
		// Extract the raw request body
		rawData, err := ctx.GetRawData()
		if err != nil {
			// If there's an error reading the body, return a bad request response
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to read the request body",
			})
			return
		}

		// Define a map to hold the raw JSON data
		var requestData map[string]interface{}

		// Unmarshal the raw JSON data into the map
		if err := json.Unmarshal(rawData, &requestData); err != nil {
			// If there's an error unmarshalling the data, return an error response
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			return
		}

		fmt.Println("received data", requestData)
		fmt.Println(models.BuildingTable)

		// Extract selectedLayer from requestData
		selectedAtribute, exists := requestData["selectedAttribute"].(string)
		selectedLayer, exists1 := requestData["selectedLayer"].(string)
		if !exists {
			// If no selectedLayer is found, return an error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Feature Attributes  field is required",
			})
			return
		}

		if !exists1 {
			// If no selectedLayer is found, return an error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Feature Feature Layer  field is required",
			})
			return
		}

		datatype, tyerr := models.GetColumnDataType(db.PG.Db, selectedLayer, selectedAtribute)
		if tyerr != nil {
			fmt.Println("query eror", tyerr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Feature feature layer or attribute selected invalid",
			})
			return
		}

		operators := GetValidOperationsForDataType(datatype)
		fmt.Println(selectedLayer, selectedAtribute, operators)

		operationsJSON, operr := json.Marshal(operators)
		if operr != nil {

			fmt.Println("operator eror", operr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Operator Error",
			})

		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Attributes found",
			"data":    json.RawMessage(operationsJSON),
		})

	} else {
		// If the request method is not POST, return an error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid HTTP method, expected POST",
		})
	}

}

func GetValidOperationsForDataType(datatype string) []string {
	var operations []string
	switch datatype {
	case "numeric", "integer", "bigint", "double precision", "float":
		operations = []string{
			"Equality (=)",
			"Less than (<)",
			"Less than or equal to (<=)",
			"Greater than (>)",
			"Greater than or equal to (>=)",
			"Between",
		}
	case "character varying", "text", "char", "varchar":
		operations = []string{
			"Equality (=)",
			"Like (e.g. LIKE '%pattern%')",
			"Not like",
			"Starts with (e.g. LIKE 'prefix%')",
			"Ends with (e.g. LIKE '%suffix')",
			"Contains (e.g. LIKE '%substring%')",
			"Regular expression match (~)",
			"Case-insensitive like (ILIKE '%pattern%')",
		}
	case "date", "timestamp", "time":
		operations = []string{
			"Equality (=)",
			"Less than (<)",
			"Less than or equal to (<=)",
			"Greater than (>)",
			"Greater than or equal to (>=)",
			"Between",
			"Date operations (e.g., DATEADD, DATE_SUB)",
		}
	case "boolean":
		operations = []string{
			"Equality (=)",
			"Not equal to (!= or <>)",
		}
	default:
		operations = []string{"No defined operations"}
	}
	return operations
}

func MakeQuery(ctx *gin.Context) {
	fmt.Println("controller reached")
	if ctx.Request.Method == http.MethodPost {
		fmt.Println("controller reached")
		// Extract the raw request body
		rawData, err := ctx.GetRawData()
		if err != nil {
			// If there's an error reading the body, return a bad request response
			fmt.Println("rawdata", rawData)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to read the request body",
			})
			return
		}

		// Define a map to hold the raw JSON data
		var requestData map[string]interface{}

		// Unmarshal the raw JSON data into the map
		if err := json.Unmarshal(rawData, &requestData); err != nil {
			// If there's an error unmarshalling the data, return an error response
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			return
		}

		fmt.Println("received data", requestData)
		fmt.Println(models.BuildingTable)

		// Extract selectedLayer from requestData
		selectedAtribute, exists := requestData["selectedAttribute"].(string)
		selectedLayer, exists2 := requestData["selectedLayer"].(string)
		selectedOperator, exists3 := requestData["selectedOperator"].(string)
		searchValuer, exists4 := requestData["searchValue"].(string)
		if !exists || !exists2 || !exists3 || !exists4 {
			// If no selectedLayer is found, return an error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You are required to select featurelayer,attribute,operator and enter search value",
			})
			return
		}

		data, searcherr := models.PerformOperation(db.PG.Db, selectedLayer, selectedAtribute, selectedOperator, searchValuer)

		if searcherr != nil {

			fmt.Println("operation error", searcherr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You Operation not successful",
			})
			return

		}

		fmt.Println(data)

		// operationsJSON, operr := json.Marshal(operators)
		// if operr != nil {

		// 	fmt.Println("operator eror", operr)
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"success": false,
		// 		"message": "Operator Error",
		// 	})

		// }

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query successful",
			"data":    data,
		})

	} else {
		// If the request method is not POST, return an error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid HTTP method, expected POST",
		})
	}

}

// func MakeQuery(ctx *gin.Context){

// }

func SimpleSearch(ctx *gin.Context) {
	fmt.Println("controller reached")
	if ctx.Request.Method == http.MethodPost {
		fmt.Println("controller reached")
		// Extract the raw request body
		rawData, err := ctx.GetRawData()
		if err != nil {
			// If there's an error reading the body, return a bad request response
			fmt.Println("rawdata", rawData)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to read the request body",
			})
			return
		}

		// Define a map to hold the raw JSON data
		var requestData map[string]interface{}

		// Unmarshal the raw JSON data into the map
		if err := json.Unmarshal(rawData, &requestData); err != nil {
			// If there's an error unmarshalling the data, return an error response
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			return
		}

		fmt.Println("received data", requestData)
		fmt.Println(models.BuildingTable)

		// Extract selectedLayer from requestData

		searchValuer, exists := requestData["searchValue"].(string)
		if !exists {
			// If no selectedLayer is found, return an error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You are required to select featurelayer,attribute,operator and enter search value",
			})
			return
		}

		data, searcherr := models.SearchAllTables(db.PG.Db, searchValuer)

		if searcherr != nil {

			fmt.Println("operation error", searcherr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You Operation not successful",
			})
			return

		}

		fmt.Println(data)

		// operationsJSON, operr := json.Marshal(operators)
		// if operr != nil {

		// 	fmt.Println("operator eror", operr)
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"success": false,
		// 		"message": "Operator Error",
		// 	})

		// }

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query successful",
			"data":    data,
		})

	} else {
		// If the request method is not POST, return an error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid HTTP method, expected POST",
		})
	}

}

func SearchAllFeaturesData(ctx *gin.Context) {
	fmt.Println("controller reached")
	if ctx.Request.Method == http.MethodPost {
		fmt.Println("controller reached")
		// Extract the raw request body
		

		// Define a map to hold the raw JSON data
		

		

		// Extract selectedLayer from requestData


		data, searcherr := models.SearchAllTables1(db.PG.Db)

		if searcherr != nil {

			fmt.Println("operation error", searcherr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You Operation not successful",
			})
			return

		}

		fmt.Println(data)

		// operationsJSON, operr := json.Marshal(operators)
		// if operr != nil {

		// 	fmt.Println("operator eror", operr)
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"success": false,
		// 		"message": "Operator Error",
		// 	})

		// }

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query successful",
			"data":    data,
		})

	} else {
		// If the request method is not POST, return an error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid HTTP method, expected POST",
		})
	}

}


func SearchByColumn(ctx *gin.Context) {
	fmt.Println("controller reached")
	if ctx.Request.Method == http.MethodPost {
		fmt.Println("controller reached")
		// Extract the raw request body
		rawData, err := ctx.GetRawData()
		if err != nil {
			// If there's an error reading the body, return a bad request response
			fmt.Println("rawdata", rawData)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to read the request body",
			})
			return
		}

		// Define a map to hold the raw JSON data
		var requestData map[string]interface{}

		// Unmarshal the raw JSON data into the map
		if err := json.Unmarshal(rawData, &requestData); err != nil {
			// If there's an error unmarshalling the data, return an error response
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			return
		}

		fmt.Println("received data", requestData)
		fmt.Println(models.BuildingTable)

		// Extract selectedLayer from requestData
		selectedAtribute, exists := requestData["selectedAttribute"].(string)
		selectedLayer, exists2 := requestData["selectedLayer"].(string)
		
		if !exists || !exists2{
			// If no selectedLayer is found, return an error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You are required to select featurelayer,attribute,operator and enter search value",
			})
			return
		}

		data, searcherr := models.SearchByColumn(db.PG.Db, selectedLayer, selectedAtribute)

		if searcherr != nil {

			fmt.Println("operation error", searcherr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You Operation not successful",
			})
			return

		}

		fmt.Println(data)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query successful",
			"data":    data,
		})

	} else {
		// If the request method is not POST, return an error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid HTTP method, expected POST",
		})
	}

}



func SearchByFeatureLayer(ctx *gin.Context) {
	fmt.Println("controller reached")
	if ctx.Request.Method == http.MethodPost {
		fmt.Println("controller reached")
		// Extract the raw request body
		rawData, err := ctx.GetRawData()
		if err != nil {
			// If there's an error reading the body, return a bad request response
			fmt.Println("rawdata", rawData)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Failed to read the request body",
			})
			return
		}

		// Define a map to hold the raw JSON data
		var requestData map[string]interface{}

		// Unmarshal the raw JSON data into the map
		if err := json.Unmarshal(rawData, &requestData); err != nil {
			// If there's an error unmarshalling the data, return an error response
			fmt.Println(err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid JSON format",
			})
			return
		}

		fmt.Println("received data", requestData)
		fmt.Println(models.BuildingTable)

		// Extract selectedLayer from requestData
	
		selectedLayer, exists2 := requestData["selectedLayer"].(string)
		
		if !exists2{
			// If no selectedLayer is found, return an error
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You are required to select featurelayer,attribute,operator and enter search value",
			})
			return
		}

		data, searcherr := models.SearchByTable(db.PG.Db, selectedLayer)

		if searcherr != nil {

			fmt.Println("operation error", searcherr)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "You Operation not successful",
			})
			return

		}

		fmt.Println(data)

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "query successful",
			"data":    data,
		})

	} else {
		// If the request method is not POST, return an error
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid HTTP method, expected POST",
		})
	}

}



