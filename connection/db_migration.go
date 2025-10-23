package connection

import (
	"chat/globals"
	"database/sql"
	"strings"
)

func validSqlError(err error) bool {
	if err == nil {
		return false
	}

	content := err.Error()

	// Error 1060: Duplicate column name
	// Error 1050: Table already exists

	return !(strings.Contains(content, "Error 1060") || strings.Contains(content, "Error 1050"))
}

func checkSqlError(_ sql.Result, err error) error {
	if validSqlError(err) {
		return err
	}

	return nil
}

func execSql(db *sql.DB, sql string, args ...interface{}) error {
	return checkSqlError(globals.ExecDb(db, sql, args...))
}

func doMigration(db *sql.DB) error {
	if globals.SqliteEngine {
		return doSqliteMigration(db)
	} else if globals.PostgresEngine {
		// no migration for postgres yet
		return nil
	}

	// v3.10 migration

	// update `quota`, `used` field in `quota` table
	// migrate `DECIMAL(16, 4)` to `DECIMAL(24, 6)`

	if err := execSql(db, `
		ALTER TABLE quota
		MODIFY COLUMN quota DECIMAL(24, 6),
		MODIFY COLUMN used DECIMAL(24, 6);
	`); err != nil {
		return err
	}

	// add new field `is_banned` in `auth` table
	if err := execSql(db, `
		ALTER TABLE auth
		ADD COLUMN is_banned BOOLEAN DEFAULT FALSE;
	`); err != nil {
		return err
	}

	return nil
}

func doSqliteMigration(db *sql.DB) error {
	// v3.10 added sqlite support, no migration needed before this version

	return nil
}
