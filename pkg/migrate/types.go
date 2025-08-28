package migrate

// Migration represents a migration
type Migration struct {
	Name     string // Migrations name (ex: 20250828090000_init)
	Checksum string // Hash SHA256 from file
	Content  []byte // Content .sql
}

// AppliedMigration represents a migration already applied in the bank
type AppliedMigration struct {
	ID       int
	Name     string
	Checksum string
	Success  bool
	Error    *string
}
