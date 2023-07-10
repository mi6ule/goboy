package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/logging"
)

func RunSqlMigrations(db *persistence.Database) error {
	wd, _ := os.Getwd()
	migrationsDir := filepath.Join(wd, "internal", "infrastructure", "database", "migration", "command")

	// Get the list of migration files
	migrationFiles, err := GetSqlMigrationFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Iterate over the migration files
	for _, file := range migrationFiles {
		// Read the migration file
		migrationContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file: %w", err)
		}

		// Check if the migration has already been executed
		migrationName := GetSqlMigrationName(file)
		if IsMigrationExecuted(db, migrationName) {
			logging.Info((logging.LoggerInput{Message: fmt.Sprintf("Skipping migration: %s (already executed)", migrationName)}))
			continue
		}

		// Execute the migration SQL using the database connection
		_, err = db.ExecuteQuery(string(migrationContent))
		if err != nil {
			return fmt.Errorf("failed to execute migration: %w", err)
		}

		// Insert a new row in the migrations table to mark the migration as executed
		err = InsertMigration(db, migrationName)
		if err != nil {
			return fmt.Errorf("failed to insert migration: %w", err)
		}

		// Print the applied migration
		logging.Info((logging.LoggerInput{Message: fmt.Sprintf("Applied migration: %s\n", migrationName)}))
	}

	return nil
}

func GetSqlMigrationFiles(migrationsDir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return nil, err
	}

	return files, nil
}

func GetSqlMigrationName(filePath string) string {
	fileName := filepath.Base(filePath)
	migrationName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	return migrationName
}

func IsMigrationExecuted(db *persistence.Database, migrationName string) bool {
	// Prepare the query
	query := fmt.Sprintf("SELECT id FROM migrations WHERE migration_name = '%s'", migrationName)

	// Execute the query
	rows, err := db.ExecuteQuery(query)
	if err != nil || len(rows) == 0 {
		// Error occurred or no row found, consider migration not executed
		return false
	}

	// Migration executed (row found)
	return true
}

func InsertMigration(db *persistence.Database, migrationName string) error {
	// Prepare the query
	query := fmt.Sprintf("INSERT INTO migrations (migration_name) VALUES ('%s')", migrationName)

	// Execute the query
	_, err := db.ExecuteQuery(query)
	if err != nil {
		return err
	}

	return nil
}
