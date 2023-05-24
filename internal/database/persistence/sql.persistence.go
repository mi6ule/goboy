package persistence

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
)

type Database struct {
	db *sql.DB
}

func NewSqlDatabaseConn(driver string, connectionConfig config.DatabaseConfig) (*Database, error) {
	var connectionString string
	if len(connectionConfig.ConnectionString) > 0 {
		connectionString = connectionConfig.ConnectionString
	} else if len(connectionConfig.Host) > 0 {
		connectionString = fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s", driver, connectionConfig.User, connectionConfig.Pwd, connectionConfig.Host, connectionConfig.Port, connectionConfig.Name, connectionConfig.Options)
	} else {
		return nil, nil
	}
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to %s db", driver)

	return &Database{
		db: db,
	}, nil
}

func (db *Database) Close() {
	err := db.db.Close()
	if err != nil {
		log.Printf("Error closing the database connection: %v", err)
	}
	log.Println("Database closed")
}

// ExecuteQuery executes the specified SQL query and returns the result
func (db *Database) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columns {
			columnPointers[i] = &values[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})
		for i, colName := range columns {
			rowData[colName] = values[i]
		}

		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Example usage
func ExampleUsage() {
	// Create a new SQL database connection
	db, err := NewSqlDatabaseConn("mysql", config.DatabaseConfig{ConnectionString: "user:password@tcp(localhost:3306)/database"})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Execute a query
	query := "SELECT * FROM users"
	result, err := db.ExecuteQuery(query)
	if err != nil {
		log.Fatal(err)
	}

	// Process the query result
	for _, row := range result {
		// Access row data using column names
		fmt.Println("ID:", row["id"])
		fmt.Println("Name:", row["name"])
		// ...
	}
}
