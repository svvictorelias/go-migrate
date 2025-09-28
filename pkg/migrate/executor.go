package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
)

// Execute applies all local migrations in the correct order
func Execute(db *sql.DB, driver string, local []Migration, applied []AppliedMigration, noApply bool) error {
	sort.Slice(local, func(i, j int) bool {
		pi := strings.SplitN(local[i].Name, "_", 2)[0]
		pj := strings.SplitN(local[j].Name, "_", 2)[0]

		// Convert to int
		var ni, nj int64
		fmt.Sscanf(pi, "%d", &ni)
		fmt.Sscanf(pj, "%d", &nj)

		return ni < nj
	})

	// Create quick map of applied by name
	appliedMap := make(map[string]AppliedMigration)
	for _, a := range applied {
		appliedMap[a.Name] = a
	}

	for _, m := range local {
		if a, ok := appliedMap[m.Name]; ok {
			// Migration already applied
			if a.Checksum != m.Checksum {
				return fmt.Errorf("checksums do not match for migration %s", m.Name)
			}
			if !a.Success {
				log.Printf("Migration %s failed, reapplying...", m.Name)
				if err := applyMigration(db, driver, m, true, noApply); err != nil {
					return err
				}
			} else {
				// Sucess, nothing to do
				continue
			}
		} else {
			// Migration never applied -> apply
			log.Printf("Aplying migration: %s", m.Name)
			if err := applyMigration(db, driver, m, false, noApply); err != nil {
				return err
			}
		}
	}
	if noApply {
		log.Println("No SQL was executed.")
	} else {
		log.Println("All migrations applied.")
	}
	return nil
}

// applyMigration
// isReapply show if it's a reapply
func applyMigration(db *sql.DB, driver string, m Migration, isReapply bool, noApply bool) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if noApply {
		// Delete unsuccessful migration
		if isReapply {
			if err := deleteMigration(tx, driver, m.Name); err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := SaveMigration(tx, driver, m.Name, m.Checksum, true, nil); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		return nil
	}

	// Execute SQL
	if _, err := tx.Exec(string(m.Content)); err != nil {
		tx.Rollback()
		message := err.Error()
		SaveMigration(db, driver, m.Name, m.Checksum, false, &message)
		return fmt.Errorf("error to execute migration %s: %w", m.Name, err)
	}

	// Delete unsuccessful migration
	if isReapply {
		if err := deleteMigration(tx, driver, m.Name); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Save migration success=true
	if err := SaveMigration(tx, driver, m.Name, m.Checksum, true, nil); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if isReapply {
		log.Printf("Migration %s reapplied com sucesso.", m.Name)
	} else {
		log.Printf("Migration %s applied com sucesso.", m.Name)
	}

	return nil
}

// deleteMigration delete the migration record by name
func deleteMigration(tx *sql.Tx, driver, name string) error {

	if driver == "postgres" {
		query := `DELETE FROM migrations WHERE name = $1`
		_, err := tx.Exec(query, name)
		return err
	} else if driver == "mysql" {
		query := `DELETE FROM migrations WHERE name = ?`
		_, err := tx.Exec(query, name)
		return err
	}

	return fmt.Errorf("unsupported driver: %s", driver)
}
