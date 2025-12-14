package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	connectionStr := `host=localhost user=postgres password=lol dbname=postgres port=5432 sslmode=disable`
	db, err := sql.Open("pgx", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("db: Open %w", err)
	}
	fmt.Println("✅ Database Connected")

	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)

	defer goose.SetBaseFS(nil)

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w\n ", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w\n ", err)
	}
	return nil
}
