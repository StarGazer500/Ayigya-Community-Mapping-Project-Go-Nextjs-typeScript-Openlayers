package models

import (
	"github.com/StarGazer500/ayigya/inits/db"
	
	// "bytes"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	"github.com/twpayne/go-geom/encoding/geojson"

	"fmt"
	"strings"

	// "github.com/twpayne/go-geom/wkb"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkb"
	"github.com/twpayne/go-geom/encoding/wkb"
)

func CreateTable(querystring string) {

	tab, err := db.PG.Db.Exec(querystring)

	if err != nil {
		fmt.Println("Error creating table: ", err)
		return
	}

	fmt.Println("Table created successfully (if it didn't already exist).", tab)

}

func InsertOne(db *sql.DB, tablename string, columns []string, args ...interface{}) (sql.Result, error) {
	// Check if columns and arguments match
	if len(columns) != len(args) {
		return nil, fmt.Errorf("the number of columns and values must match")
	}

	// Dynamically build the column part for the query
	cols := strings.Join(columns, ", ")

	// Dynamically build the placeholders for the query (e.g. "$1, $2, $3")
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	placeholderStr := strings.Join(placeholders, ", ")

	// Create the query string
	query := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, tablename, cols, placeholderStr)

	// Execute the query with the provided arguments
	data, err := db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Get the number of rows affected
	rowsAffected, err := data.RowsAffected()
	if err != nil {
		return nil, err
	}

	// Print meaningful details
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	return data, nil
}

// InsertMany inserts multiple rows into the specified table.
func InsertMany(db *sql.DB, tablename string, columns []string, rows [][]interface{}) (sql.Result, error) {
	// Check if columns are provided
	if len(columns) == 0 {
		return nil, fmt.Errorf("columns must not be empty")
	}

	// Ensure that each row has the same number of values as there are columns
	for _, row := range rows {
		if len(row) != len(columns) {
			return nil, fmt.Errorf("each row must have the same number of values as the number of columns")
		}
	}

	// Dynamically build the column part for the query
	cols := strings.Join(columns, ", ")

	// Build the placeholders part for multiple rows (e.g. "($1, $2, $3), ($4, $5, $6), ...")
	var placeholders []string
	for i := range rows {
		// Generate the placeholders for each row
		rowPlaceholders := make([]string, len(columns))
		for j := range rowPlaceholders {
			rowPlaceholders[j] = fmt.Sprintf("$%d", i*len(columns)+j+1)
		}
		placeholders = append(placeholders, fmt.Sprintf("(%s)", strings.Join(rowPlaceholders, ", ")))
	}

	// Create the query string
	placeholderStr := strings.Join(placeholders, ", ")
	query := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES %s`, tablename, cols, placeholderStr)

	// Flatten the rows into a single slice of arguments
	var args []interface{}
	for _, row := range rows {
		args = append(args, row...)
	}

	// Execute the query with the provided arguments
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error inserting data into the database: %v", err)
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %v", err)
	}

	// Print meaningful details
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	return result, nil
}

// type ReturnedData struct {
// 	ColumnName string
// 	Value      interface{}
// }

func FindOne(db *sql.DB, tableName string, queryField string, queryValue string) ([]map[string]interface{}, error) {

	// Correctly format the query string with parameterized query to avoid SQL injection
	querystring := fmt.Sprintf("SELECT * FROM \"%s\" WHERE %s = $1", tableName, queryField)

	// Execute the query with the parameterized query
	rows, err := db.Query(querystring, queryValue)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()

	if err != nil {
		fmt.Println("Error getting columns: ", err)
	}

	var result []map[string]interface{}

	// Initialize the slice with pointers to sql.RawBytes for each column
	rawValues := make([]interface{}, len(columns))
	for i := range rawValues {
		rawValues[i] = new(sql.RawBytes)
	}

	// Loop through each row
	for rows.Next() {
		// Scan the row values into rawValues
		err := rows.Scan(rawValues...)
		if err != nil {
			fmt.Println("Error scanning row: ", err)
		}

		// Create a map to hold the column name and value for this row
		rowMap := make(map[string]interface{})

		// Loop over the columns and populate the map
		for i, col := range columns {
			// Type assert to sql.RawBytes
			rawBytes := rawValues[i].(*sql.RawBytes)

			// Convert []byte to a string for text data (or other types as needed)
			// Here, we handle text columns as strings.
			rowMap[col] = string(*rawBytes) // You can enhance this logic based on column types
		}

		// Append the map to the result slice
		result = append(result, rowMap)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		fmt.Println("Error during row iteration: ", err)
	}

	//

	// fmt.Println("yh", rows.Next())

	// If rows are returned, you can scan the first row (for example, assuming the table has 'email' and 'password' columns)

	// Since we're only looking for one row, return rows as the result
	return result, nil
}

// AddColumnIfNotExists checks if the column exists in the table, and if not, adds it.
func AddColumnIfNotExists(db *sql.DB, tableName, columnName, columnType string) (string, error) {
	// Check if the column already exists in the table
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2;
	`

	var existingColumn string
	err := db.QueryRow(query, tableName, columnName).Scan(&existingColumn)
	if err != nil && err != sql.ErrNoRows {
		// If there's an error other than no rows found, log and return the error
		return "", fmt.Errorf("error checking column existence: %w", err)
	}

	if existingColumn == columnName {
		// If the column exists, no need to add it
		fmt.Printf("Column '%s' already exists in table '%s'.\n", columnName, tableName)
	} else {
		// Add the new column if it doesn't exist
		// Here we have to directly interpolate table and column names into the query
		// but make sure the values are sanitized to avoid SQL injection
		// *** Be cautious with interpolating user inputs directly into SQL ***
		// In this case, we are assuming `tableName` and `columnName` are safe.

		// Build the ALTER TABLE query
		addColumnQuery := fmt.Sprintf(`
			ALTER TABLE "%s"
			ADD COLUMN "%s" %s;
		`, tableName, columnName, columnType)

		// Execute the ALTER TABLE query
		_, err := db.Exec(addColumnQuery)
		if err != nil {
			return "", fmt.Errorf("error adding column '%s' to table '%s': %w", columnName, tableName, err)
		}

		fmt.Printf("Column '%s' added successfully to table '%s'.\n", columnName, tableName)
	}

	// Generate the updated CREATE TABLE SQL syntax with all columns
	createTableSQL, err := generateCreateTableSQL(db, tableName)
	if err != nil {
		return "", fmt.Errorf("error generating CREATE TABLE SQL: %w", err)
	}

	// Return the generated CREATE TABLE SQL
	return createTableSQL, nil
}

// DeleteColumnIfExists checks if the column exists in the table, and if so, deletes it.
func DeleteColumnIfExists(db *sql.DB, tableName, columnName string) (string, error) {
	// Check if the column already exists in the table
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2;
	`

	var existingColumn string
	err := db.QueryRow(query, tableName, columnName).Scan(&existingColumn)
	if err != nil && err != sql.ErrNoRows {
		// If there's an error other than no rows found, log and return the error
		return "", fmt.Errorf("error checking column existence: %w", err)
	}

	if existingColumn == columnName {
		// If the column exists, delete it
		deleteColumnQuery := fmt.Sprintf(`
			ALTER TABLE "%s"
			DROP COLUMN "%s";
		`, tableName, columnName)

		// Execute the ALTER TABLE query to drop the column
		_, err := db.Exec(deleteColumnQuery)
		if err != nil {
			return "", fmt.Errorf("error deleting column '%s' from table '%s': %w", columnName, tableName, err)
		}

		fmt.Printf("Column '%s' deleted successfully from table '%s'.\n", columnName, tableName)
	} else {
		// If the column does not exist, no action needed
		fmt.Printf("Column '%s' does not exist in table '%s'. No action taken.\n", columnName, tableName)
	}

	// Generate the updated CREATE TABLE SQL syntax with all remaining columns
	createTableSQL, err := generateCreateTableSQL(db, tableName)
	if err != nil {
		return "", fmt.Errorf("error generating CREATE TABLE SQL: %w", err)
	}

	// Return the generated CREATE TABLE SQL
	return createTableSQL, nil
}

// generateCreateTableSQL generates the SQL query Syntz to create a table
// based on the columns from an existing table in the database.
func generateCreateTableSQL(db *sql.DB, tableName string) (string, error) {
	// Query to get all columns and their types
	query := `
		SELECT column_name, data_type, character_maximum_length
		FROM information_schema.columns
		WHERE table_name = $1
		ORDER BY ordinal_position;
	`

	rows, err := db.Query(query, tableName)
	if err != nil {
		return "", fmt.Errorf("error retrieving columns from information_schema: %w", err)
	}
	defer rows.Close()

	// Start building the CREATE TABLE SQL query
	var createTableSQL = fmt.Sprintf("CREATE TABLE \"%s\" (\n", tableName)

	// Loop through the columns and construct the SQL for each
	firstColumn := true
	for rows.Next() {
		var columnName, dataType string
		var maxLength *int

		err := rows.Scan(&columnName, &dataType, &maxLength)
		if err != nil {
			return "", fmt.Errorf("error scanning row: %w", err)
		}

		// Add column definition to the CREATE TABLE SQL
		if !firstColumn {
			createTableSQL += ",\n"
		}

		// Handle data types (e.g., VARCHAR with length)
		if dataType == "character varying" && maxLength != nil {
			createTableSQL += fmt.Sprintf("\t%s %s(%d)", columnName, dataType, *maxLength)
		} else {
			createTableSQL += fmt.Sprintf("\t%s %s", columnName, dataType)
		}

		firstColumn = false
	}

	// Check for any errors after the loop
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating over rows: %w", err)
	}

	// Close the CREATE TABLE statement
	createTableSQL += "\n);\n"

	// Return the generated SQL
	return createTableSQL, nil
}

func DeleteRowByColumn(db *sql.DB, tableName, columnName, columnValue string) error {
	// Step 1: Check if the row exists by querying the table for the column and value
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM "%s"
		WHERE "%s" = $1;
	`, tableName, columnName)

	var count int
	err := db.QueryRow(query, columnValue).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking row existence: %w", err)
	}

	if count == 0 {
		// If the row does not exist, return an error message
		return fmt.Errorf("row with %s = '%s' not found in table '%s'", columnName, columnValue, tableName)
	}

	// Step 2: If the row exists, delete it using the DELETE statement
	deleteQuery := fmt.Sprintf(`
		DELETE FROM "%s"
		WHERE "%s" = $1;
	`, tableName, columnName)

	// Execute the DELETE query
	_, err = db.Exec(deleteQuery, columnValue)
	if err != nil {
		return fmt.Errorf("error deleting row where %s = '%s' from table '%s': %w", columnName, columnValue, tableName, err)
	}

	fmt.Printf("Row with %s = '%s' deleted successfully from table '%s'.\n", columnName, columnValue, tableName)
	return nil
}

func DeleteRowByID(db *sql.DB, tableName string, id int) error {
	// Step 1: Check if the row exists by querying the table for the id
	query := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM "%s"
		WHERE id = $1;
	`, tableName)

	var count int
	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking row existence: %w", err)
	}

	if count == 0 {
		// If the row does not exist, return a message
		return fmt.Errorf("row with id %d not found in table '%s'", id, tableName)
	}

	// Step 2: If the row exists, delete it using the DELETE statement
	deleteQuery := fmt.Sprintf(`
		DELETE FROM "%s"
		WHERE id = $1;
	`, tableName)

	// Execute the DELETE query
	_, err = db.Exec(deleteQuery, id)
	if err != nil {
		return fmt.Errorf("error deleting row with id %d from table '%s': %w", id, tableName, err)
	}

	fmt.Printf("Row with id %d deleted successfully from table '%s'.\n", id, tableName)
	return nil
}

// UpdateOne updates a single row in the specified table with the provided column values.
func UpdateOne(db *sql.DB, tablename string, columns []string, args []interface{}, whereColumn string, whereValue interface{}) (sql.Result, error) {
	// Check if columns and arguments match
	if len(columns) != len(args) {
		return nil, fmt.Errorf("the number of columns and values must match")
	}

	// Dynamically build the column=value part for the query (e.g. "col1 = $1, col2 = $2")
	setClauses := make([]string, len(columns))
	for i := range columns {
		setClauses[i] = fmt.Sprintf("\"%s\" = $%d", columns[i], i+1)
	}
	setClauseStr := strings.Join(setClauses, ", ")

	// Create the query string with the WHERE clause for the row to update
	query := fmt.Sprintf(`UPDATE "%s" SET %s WHERE "%s" = $%d`, tablename, setClauseStr, whereColumn, len(columns)+1)

	// Combine the arguments (the column values and the WHERE condition)
	args = append(args, whereValue)

	// Execute the query with the provided arguments
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error updating data in the database: %v", err)
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %v", err)
	}

	// Print meaningful details
	fmt.Printf("Rows affected: %d\n", rowsAffected)
	return result, nil
}

func GetColumnDataType(db *sql.DB, table, column string) (string, error) {
	// Construct query to get column data type from information_schema.columns
	query := `
		SELECT data_type
		FROM information_schema.columns
		WHERE table_name = $1 AND column_name = $2
	`
	var dataType string
	err := db.QueryRow(query, table, column).Scan(&dataType)
	if err != nil {
		return "", fmt.Errorf("failed to determine data type: %v", err)
	}
	return dataType, nil
}

func PerformOperation(db *sql.DB, table, column, operator, value string) ([]map[string]interface{}, error) {

	// Step 4: Build the SQL query based on the operator
	var query string
	var args []interface{}

	// Handle different operators based on the column type
	switch operator {
	case "Equality (=)":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", table, column)
		args = append(args, value)
	case "Less than (<)":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s < $1", table, column)
		args = append(args, value)
	case "Less than or equal to (<=)":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s <= $1", table, column)
		args = append(args, value)
	case "Greater than (>)":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s > $1", table, column)
		args = append(args, value)
	case "Greater than or equal to (>=)":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s >= $1", table, column)
		args = append(args, value)
	case "ILIKE":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s ILIKE $1", table, column)
		args = append(args, "%"+value+"%")
	case "LIKE":
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s LIKE $1", table, column)
		args = append(args, "%"+value+"%")
	case "BETWEEN":
		// For "BETWEEN", assume value is in the form "start_value AND end_value"
		parts := strings.Split(value, " AND ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid value for BETWEEN operation")
		}
		query = fmt.Sprintf("SELECT * FROM %s WHERE %s BETWEEN $1 AND $2", table, column)
		args = append(args, parts[0], parts[1])
	default:
		return nil, fmt.Errorf("operator '%s' is not supported", operator)
	}

	// Step 5: Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Step 6: Fetch the results
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of empty interfaces to hold the row values
		columnPointers := make([]interface{}, len(columns))
		for i := range columnPointers {
			columnPointers[i] = new([]byte)
		}

		// Scan the row into the column pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Convert row to map
		result := make(map[string]interface{})
		for i, col := range columns {

			value := columnPointers[i].(*[]byte)

			if col == "shape__len" || col == "shape__are" {
				// fmt.Println("shape", (string(*value)))
				tofloat, err := strconv.ParseFloat(string(*value), 64)
				if err != nil {
					fmt.Println("conversion failed", err)
				} else {
					result[col] = tofloat
				}

			} else if col == "geom" {
				// fmt.Println("geom", (string(*value)))
				geodata, err := processGeometryData(*value)
				if err != nil {
					fmt.Println("error", err)
				}

				result[col] = json.RawMessage(geodata)

			}else{
				result[col] =string(*value)
			}

			// 
		}

		results = append(results, result)
	}
	fmt.Println(results)
	return results, nil
}

func cleanAndParseWKB(rawData []byte) ([]byte, error) {
	// Try decoding from hex
	decodedHex, err := hex.DecodeString(string(rawData))
	if err == nil {
		// fmt.Println("Successfully decoded from hex")
		return decodedHex, nil
	}

	// Try base64 decoding
	decodedBase64, err := base64.StdEncoding.DecodeString(string(rawData))
	if err == nil {
		// fmt.Println("Successfully decoded from base64")
		return decodedBase64, nil
	}

	// If direct hex string (ASCII representation)
	decodedASCIIHex, err := hex.DecodeString(strings.ReplaceAll(string(rawData), " ", ""))
	if err == nil {
		// fmt.Println("Successfully decoded from ASCII hex")
		return decodedASCIIHex, nil
	}

	return nil, fmt.Errorf("could not decode WKB data")
}

func processGeometryData(rawData []byte) ([]byte, error) {
	// Clean and decode the data
	cleanData, err := cleanAndParseWKB(rawData)
	if err != nil {
		return nil, fmt.Errorf("data decoding failed: %v", err)
	}

	// Now try parsing with various methods
	// Option 1: Try EWKB

	geometry, err := ewkb.Unmarshal(cleanData)
	if err != nil {
		// Option 2: Try standard WKB
		geometry, err = wkb.Unmarshal(cleanData)
		if err != nil {
			return nil, fmt.Errorf("geometry parsing failed: %v", err)
		}
	}

	// fmt.Println(ConvertGeometryToXY(geometry))

	// Convert to GeoJSON if needed
	geoJSONGeom, err := geojson.Marshal(ConvertGeometryToXY(geometry))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to GeoJSON: %v", err)
	}

	//  fmt.Println(geoJSONGeom)

	return geoJSONGeom, nil
}
func ConvertGeometryToXY(geometry geom.T) geom.T {
	// Create a new MultiPolygon with a 2D layout (geom.XY)
	multiPolygon := geom.NewMultiPolygon(geom.XY)

	// Perform type assertion inside the function
	switch g := geometry.(type) {
	case *geom.MultiPolygon:
		// Iterate through each Polygon in the MultiPolygon using Coords()
		for _, polygon := range g.Coords() {
			var coordinates []geom.Coord

			// Iterate through each ring of the Polygon
			for _, ring := range polygon {
				// Extract each point (X, Y) from the ring and append to coordinates
				for i := 0; i < len(ring); i++ {
					// Extract only X and Y, ignoring Z and M
					x, y := ring[i][0], ring[i][1]

					coordinates = append(coordinates, geom.Coord{x, y})
				}
			}

			// Create a new Polygon with the 2D layout (XY)
			newPolygon := geom.NewPolygon(geom.XY)
			// Set the coordinates for the polygon (set only the 2D points)
			newPolygon.SetCoords([][]geom.Coord{coordinates})

			// Add the polygon to the MultiPolygon
			multiPolygon.Push(newPolygon) // Push to MultiPolygon (not append)
		}

	default:
		// Handle unsupported geometry types
		fmt.Println("Unsupported geometry type:", g)
		return nil
	}

	multiPolygon.SetSRID(4326)

	return multiPolygon
}




func SearchAllTables(db *sql.DB, searchTerm string) ([]map[string]interface{}, error) {
    // Query to get all table names and column names from the information schema
    query := `
        SELECT table_name, column_name
        FROM information_schema.columns
        WHERE table_schema = 'public'
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Loop through each row in the result
    var results []map[string]interface{}
    var processedTables = make(map[string]bool)

    for rows.Next() {
        var tableName, columnName string
        if err := rows.Scan(&tableName, &columnName); err != nil {
            return nil, err
        }

        // Skip if table already processed or not in desired tables
        if processedTables[tableName] || (tableName != "building" && tableName != "other_polygon_structure") {
            continue
        }

        // Find all columns in the table to search
        columnsQuery := fmt.Sprintf(`
            SELECT column_name 
            FROM information_schema.columns 
            WHERE table_name = '%s' 
            AND (data_type LIKE '%%character%%' OR data_type LIKE '%%text%%')
        `, tableName)

        columnRows, err := db.Query(columnsQuery)
        if err != nil {
            fmt.Println("Error retrieving searchable columns for table:", tableName, err)
            continue
        }
        defer columnRows.Close()

        var searchColumns []string
        for columnRows.Next() {
            var colName string
            if err := columnRows.Scan(&colName); err != nil {
                fmt.Println("Error scanning column name:", err)
                continue
            }
            searchColumns = append(searchColumns, colName)
        }

        // Construct dynamic search query
        var searchConditions []string
        for _, col := range searchColumns {
            searchConditions = append(searchConditions, fmt.Sprintf(`"%s" ILIKE $1`, col))
        }

        if len(searchConditions) == 0 {
            continue
        }

        querystring := fmt.Sprintf(
            `SELECT * FROM "%s" WHERE %s`, 
            tableName, 
            strings.Join(searchConditions, " OR "),
        )

        // Execute the query with the parameterized search term
        tableRows, err := db.Query(querystring, "%"+searchTerm+"%")
        if err != nil {
            fmt.Println("Error executing query for table:", tableName, err)
            continue
        }
        defer tableRows.Close()

        columns, err := tableRows.Columns()
        if err != nil {
            fmt.Println("Error retrieving columns for table:", tableName, err)
            continue
        }

        // Loop through the rows of the current table
        for tableRows.Next() {
            // Create a slice of empty interfaces to hold the row values
            columnPointers := make([]interface{}, len(columns))
            for i := range columnPointers {
                columnPointers[i] = new([]byte) // Allocate space for the values
            }

            // Scan the row into the column pointers
            if err := tableRows.Scan(columnPointers...); err != nil {
                fmt.Println("Error scanning row for table:", tableName, err)
                continue
            }

            // Convert the row to a map
            result := make(map[string]interface{})
            for i, col := range columns {
                value := columnPointers[i].(*[]byte) // Dereference the value pointer

                // Handle specific columns (e.g., shapelen, shapeare, geom)
                if col == "shapelen" || col == "shapeare" {
                    tofloat, err := strconv.ParseFloat(string(*value), 64)
                    if err != nil {
                        fmt.Println("Error parsing float for", col, err)
                    } else {
                        result[col] = tofloat
                    }
                } else if col == "geom" {
                    geodata, err := processGeometryData(*value)
                    if err != nil {
                        fmt.Println("Error processing geometry data:", err)
                    }
                    result[col] = json.RawMessage(geodata)
                } else {
                    result[col] = string(*value)
                }
            }

            // Append the result to the results slice
            results = append(results, result)
        }

        // Mark table as processed
        processedTables[tableName] = true

        // Check if there was an error during row iteration
        if err := tableRows.Err(); err != nil {
            fmt.Println("Error during row iteration for table:", tableName, err)
        }
    }

    // Check if there was an error during column iteration
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Error during column iteration: %v", err)
    }

    // Return the results after iterating through all tables
    fmt.Println(results)
    return results, nil
}

func SearchAllTables1(db *sql.DB) ([]map[string]interface{}, error) {
    // Query to get all table names and column names from the information schema
    query := `
        SELECT table_name
        FROM information_schema.tables
        WHERE table_schema = 'public' AND table_name IN ('building', 'other_polygon_structure')
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Loop through each row in the result
    var results []map[string]interface{}
    for rows.Next() {
        var tableName string
        if err := rows.Scan(&tableName); err != nil {
            return nil, err // Log the error and continue processing other tables
        }

        // Skip tables that you don't want to process
        if tableName != "building" && tableName != "other_polygon_structure" {
            // fmt.Println("Skipping table:", tableName) // Log which table is being skipped
            continue // Skip this iteration and go to the next table
        }

        // Modify query to select all data from the table
        querystring := fmt.Sprintf(`SELECT * FROM "%s"`, tableName)

        // Execute the query to get all data
        tableRows, err := db.Query(querystring)
        if err != nil {
            fmt.Println("Error executing query for table:", tableName, err)
            continue // Skip to the next table if this query fails
        }
        defer tableRows.Close()

        columns, err := tableRows.Columns()
        if err != nil {
            fmt.Println("Error retrieving columns for table:", tableName, err)
            continue // Skip to the next table if columns retrieval fails
        }

        // Loop through the rows of the current table
        for tableRows.Next() {
            // Create a slice of empty interfaces to hold the row values
            columnPointers := make([]interface{}, len(columns))
            for i := range columnPointers {
                columnPointers[i] = new([]byte) // Allocate space for the values
            }

            // Scan the row into the column pointers
            if err := tableRows.Scan(columnPointers...); err != nil {
                fmt.Println("Error scanning row for table:", tableName, err)
                continue // Skip this row if there's an error scanning it
            }

            // Convert the row to a map
            result := make(map[string]interface{})
            for i, col := range columns {
                value := columnPointers[i].(*[]byte) // Dereference the value pointer

                // Handle specific columns (e.g., shapelen, shapeare, geom)
                if col == "shapelen" || col == "shapeare" {
                    tofloat, err := strconv.ParseFloat(string(*value), 64)
                    if err != nil {
                        fmt.Println("Error parsing float for", col, err)
                    } else {
                        result[col] = tofloat
                    }
                } else if col == "geom" {
                    geodata, err := processGeometryData(*value)
                    if err != nil {
                        fmt.Println("Error processing geometry data:", err)
                    }
                    result[col] = json.RawMessage(geodata)
                } else {
                    result[col] = string(*value)
                }
            }

            // Append the result to the results slice
            results = append(results, result)
        }

        // Check if there was an error during row iteration
        if err := tableRows.Err(); err != nil {
            fmt.Println("Error during row iteration for table:", tableName, err)
        }
    }

    // Check if there was an error during column iteration
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Error during column iteration: %v", err)
    }

    // Return the results after iterating through all tables
    // fmt.Println(results)
    return results, nil
}


func SearchByColumn(db *sql.DB, table, column string) ([]map[string]interface{}, error) {
	// Step 5: Execute the query to get all rows where the column has any non-null value
	// Using the column name directly for exact matching
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s IS NOT NULL", table, column)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Step 6: Fetch the results
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of empty interfaces to hold the row values
		columnPointers := make([]interface{}, len(columns))
		for i := range columnPointers {
			columnPointers[i] = new([]byte)
		}

		// Scan the row into the column pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Convert row to map
		result := make(map[string]interface{})
		for i, col := range columns {
			value := columnPointers[i].(*[]byte)

			if col == "shape__len" || col == "shape__are" {
				// Convert shape__len and shape__are to float
				tofloat, err := strconv.ParseFloat(string(*value), 64)
				if err != nil {
					fmt.Println("conversion failed for column", col, err)
					result[col] = nil // Return nil for failed conversion
				} else {
					result[col] = tofloat
				}
			} else if col == "geom" {
				// Process geometry data
				geodata, err := processGeometryData(*value)
				if err != nil {
					fmt.Println("error processing geometry data:", err)
					result[col] = nil // Return nil for failed geometry processing
				} else {
					result[col] = json.RawMessage(geodata)
				}
			} else {
				// For other columns, just convert the []byte to string
				result[col] = string(*value)
			}
		}

		// Append the result to the results slice
		results = append(results, result)
	}

	// Check if there was an error during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	// Return the results

	// fmt.Println(results)
	return results, nil
}


func SearchByTable(db *sql.DB, table string) ([]map[string]interface{}, error) {
	// Step 5: Execute the query to get all rows where the column has any non-null value
	// Using the column name directly for exact matching
	query := fmt.Sprintf("SELECT * FROM %s", table)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Step 6: Fetch the results
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get column names: %v", err)
	}

	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of empty interfaces to hold the row values
		columnPointers := make([]interface{}, len(columns))
		for i := range columnPointers {
			columnPointers[i] = new([]byte)
		}

		// Scan the row into the column pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		// Convert row to map
		result := make(map[string]interface{})
		for i, col := range columns {
			value := columnPointers[i].(*[]byte)

			if col == "shape__len" || col == "shape__are" {
				// Convert shape__len and shape__are to float
				tofloat, err := strconv.ParseFloat(string(*value), 64)
				if err != nil {
					fmt.Println("conversion failed for column", col, err)
					result[col] = nil // Return nil for failed conversion
				} else {
					result[col] = tofloat
				}
			} else if col == "geom" {
				// Process geometry data
				geodata, err := processGeometryData(*value)
				if err != nil {
					fmt.Println("error processing geometry data:", err)
					result[col] = nil // Return nil for failed geometry processing
				} else {
					result[col] = json.RawMessage(geodata)
				}
			} else {
				// For other columns, just convert the []byte to string
				result[col] = string(*value)
			}
		}

		// Append the result to the results slice
		results = append(results, result)
	}

	// Check if there was an error during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	// Return the results

	// fmt.Println(results)
	return results, nil
}

