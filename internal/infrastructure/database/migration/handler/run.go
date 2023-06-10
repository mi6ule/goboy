package migration

import (
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
)

func RunMigration(db *persistence.Database) error {
	err := RunSqlMigrations(db)
	return err
}
