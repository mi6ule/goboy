package persistence

import (
	"database/sql"
	"fmt"

	"gitlab.avakatan.ir/boilerplates/go-boiler/config"
	constants "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/constant"
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
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Message: "Error closing the database connection", Err: err})
	logging.Info(logging.LoggerInput{Message: "Database closed"})
}

// ExecuteQuery executes the specified SQL query and returns the result
func (db *Database) ExecuteQuery(query string) ([]map[string]any, error) {
	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, 0)

	for rows.Next() {
		values := make([]any, len(columns))
		columnPointers := make([]any, len(columns))
		for i := range columns {
			columnPointers[i] = &values[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
		}

		rowData := make(map[string]any)
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

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *Database) QueryRow(query string, args ...any) *sql.Row {
	logging.Info(logging.LoggerInput{Data: map[string]any{"query": query, "args": args}})
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
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, ErrType: "Fatal", Code: constants.ERROR_CODE_100012})

	// Process the query result
	for _, row := range result {
		// Access row data using column names
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("ID: %v", row["id"])})
		logging.Info(logging.LoggerInput{Message: fmt.Sprintf("Name: %v", row["name"])})
	}
}
