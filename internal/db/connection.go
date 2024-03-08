package db

import (
	"database/sql"
	"fmt"

	// Need this for PGSQL
	_ "github.com/lib/pq"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

var dbHandle *sql.DB

// EnsureDatabaseConnection just makes sure we are still connected to the
// motherland (tm)
func EnsureDatabaseConnection() {
	if dbHandle == nil {
		ConnectToDatabase()
	}
}

// ConnectToDatabase Connects to ~ the ~ database
func ConnectToDatabase() {
	credentials := utils.DatabaseCredentials()

	InlineConnectToDatabase(
		credentials["user"],
		credentials["pass"],
		credentials["host"],
		credentials["data"],
	)
}

// InlineConnectToDatabase handles connections using credentials from input,
// instead of credentials from global
func InlineConnectToDatabase(user, pass, host, data string) {
	driver := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		user,
		pass,
		host,
		data,
	)

	db, err := sql.Open("postgres", driver)

	handleError("InlineConnectToDatabase", err)

	dbHandle = db
}

// BeginTransaction starts a DB transaction
func BeginTransaction() *sql.Tx {
	tx, err := GetDatabase().Begin()

	handleError("BeginTransaction", err)

	return tx
}

// CommitOrRollback handles the situation of committing or, on error,
// rolling back
func CommitOrRollback(tx *sql.Tx) {
	err := tx.Commit()

	if err != nil {
		utils.NonHaltingError("CommitOrRollback", err)

		if rollErr := tx.Rollback(); rollErr != nil {
			handleError("Rollback", rollErr)
		}
	}
}

// GetDatabase returns the active database handle
func GetDatabase() *sql.DB {
	return dbHandle
}
