package persistence

import (
	"database/sql"
	"fmt"
	"log"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/util"
)

type Database struct {
	db *sql.DB
}

func NewSqlDatabaseConn(driver string, connectionConfig config.DatabaseConfig) (*Database, error) {
	connectionString, err := util.CreateConnectionString(driver, connectionConfig)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(driver, connectionString)
	if err != nil {
		return nil, err
	}
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
func ExampleMySql() {
	// Create a new SQL database connection
	db, _ := NewSqlDatabaseConn("mysql", config.DatabaseConfig{ConnectionString: "mysql://user:password@localhost:3306/database"})
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

// Example usage
func ExamplePostgres() {
	// Create a new SQL database connection
	db, _ := NewSqlDatabaseConn("postgres", config.DatabaseConfig{ConnectionString: "postgres://postgres:123@localhost:5432/golangdb?sslmode=disable"})
	defer db.Close()

	// Execute a query
	query := "SELECT * FROM customer"
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
