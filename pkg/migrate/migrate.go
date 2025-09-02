package migrate

import (
	"database/sql"
	"fmt"
)

// Execute runs all local migrations that haven't been applied yet
// or retries any that previously failed.
//
// Parameters:
//   - db: active database connection.
//   - local: list of migrations found locally in the directory.
//   - applied: list of migrations already applied to the database.
//   - noApply: if true, skips executing the SQL statements and only
//     logs or marks them as applied (dry-run mode or db with data).
func Run(db *sql.DB, driver string, path string, noApply bool) error {
	if err := InitStorage(db, driver); err != nil {
		return fmt.Errorf("error to inicializer storage: %w", err)
	}

	local, err := LoadLocal(path)
	if err != nil {
		return fmt.Errorf("error to load local migrations: %w", err)
	}

	applied, err := LoadApplied(db)
	if err != nil {
		return fmt.Errorf("error to load applied migrations: %w", err)
	}

	return Execute(db, local, applied, noApply)
}
