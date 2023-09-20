package migration

import (
	"github.com/mi6ule/goboy/internal/infrastructure/database/persistence"
)

func RunMigration(db *persistence.Database) error {
	err := RunSqlMigrations(db)
	return err
}
