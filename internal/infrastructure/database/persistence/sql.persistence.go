package persistence

import (
	"database/sql"
	"fmt"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
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
	logging.Info(logging.LoggerInput{Message: fmt.Sprintf("Connected to %s db", driver)})

	return &Database{
		db: db,
	}, nil
}

func (db *Database) Close() {
	err := db.db.Close()
	if err != nil {
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("Error closing the database connection: %v", err)})
	}
	logging.Info(logging.LoggerInput{Message: "Database closed"})
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

func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	logging.Info(logging.LoggerInput{Message: "", Data: map[string]any{"query": query, "args": args}})
	return db.db.QueryRow(query, args...)
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
		errorhandler.ErrorHandler(err, errorhandler.TErrorData{"errType": "Fatal"})
	}

	// Process the query result
	for _, row := range result {
		// Access row data using column names
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("ID: %v", row["id"])})
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("Name: %v", row["name"])})
	}
}
