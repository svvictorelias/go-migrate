package migrate

import (
	"context"
	"database/sql"
)

// queryExecutor is implemented by *sql.DB e *sql.Tx
type queryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// InitStorage create a table to store applied migrations
func InitStorage(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		checksum TEXT NOT NULL,
		success BOOLEAN NOT NULL DEFAULT FALSE,
		error TEXT,
		applied_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`
	_, err := db.Exec(query)
	return err
}

// LoadApplied returns all applied migrations
func LoadApplied(db *sql.DB) ([]AppliedMigration, error) {
	rows, err := db.Query(`SELECT id, name, checksum, success FROM migrations ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applied []AppliedMigration
	for rows.Next() {
		var m AppliedMigration
		if err := rows.Scan(&m.ID, &m.Name, &m.Checksum, &m.Success); err != nil {
			return nil, err
		}
		applied = append(applied, m)
	}
	return applied, nil
}

// SaveMigration insert or update a migration
func SaveMigration(exec queryExecutor, name, checksum string, success bool, errorMessage *string) error {
	query := `
	INSERT INTO migrations (name, checksum, success, error)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (name) DO UPDATE
	SET success = EXCLUDED.success;
	`
	_, err := exec.ExecContext(context.Background(), query, name, checksum, success, errorMessage)
	return err
}
