/*
Copyright © 2024 Jan Nawrocki jan.nawrocki06@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package repository

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"terminder/repository/migrations"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

const DBDriver = "sqlite"

type Repository struct {
	*sql.DB
	*Queries
}

func NewRepo(terminderDir string) (*Repository, error) {
	dbPath := filepath.Join(terminderDir, "terminder.db")

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, errors.Join(RepositoryFileErr, err)
		}
		file.Close()
	}

	db, err := sql.Open(DBDriver, dbPath)
	if err != nil {
		return nil, errors.Join(DatabaseConnErr, err)
	}
	err = runMigration(db)
	if err != nil {
		return nil, errors.Join(MigrationsErr, err)
	}

	return &Repository{db, New(db)}, nil
}

func runMigration(db *sql.DB) error {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return err
	}

	// Create an embed source for the migration
	embedSource, err := iofs.New(migrations.Files(), ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs", embedSource,
		DBDriver, driver)
	if err != nil {
		return err
	}

	// Run the migration
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
