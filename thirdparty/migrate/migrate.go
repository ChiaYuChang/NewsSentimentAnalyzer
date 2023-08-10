package migrate

import (
	"github.com/golang-migrate/migrate/v4/database"
	pgx5 "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const SRC_URL = "file://relative/path"
const DB_URL = "pgx5://user:password@host:port/dbname?query"

type Instance struct {
	SourceName       string
	SourceInstance   source.Driver
	DatabaseName     string
	DatabaseInstance database.Driver
}

func init() {
	db := pgx5.Postgres{}
	database.Register("pgx5", &db)
}

// func NewMigrate() (*migrate.Migrate, error) {
// 	m, err := migrate.New(SRC_URL, DB_URL)
// 	if err != nil {
// 		return nil, err
// 	}

// }
