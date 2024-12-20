package repository

import "errors"

var (
	RepositoryFileErr = errors.New("failed to fetch file")
	DatabaseConnErr   = errors.New("failed to estabilish a database connection")
	MigrationsErr     = errors.New("failed to run migrations")
)
