// Package main provides a standalone executable for running database migrations.
// It is intended to be invoked via 'make migrate' and ensures that the database
// schema is up to date before the main application starts.
package main

import "github.com/HadesHo3820/ebvn-golang-course/internal/infrastructure"

// main initializes the database connection and triggers the execution of all
// pending UP migrations found in the "./migrations" directory.
func main() {
	_ = infrastructure.CreateSQLDBWithMigration()
}
