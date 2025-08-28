# go-migrate-test

`go-migrate` is a lightweight and straightforward migration tool written in Go.  
It helps you manage and apply SQL migrations in sequence, ensuring consistent database versioning.

---

## Install

To add `go-migrate` to your project:

- `go get github.com/svvictorelias/go-migrate`

## Functions

- `CreateMigration(path string, name string)`  
  Creates an empty migration file with a timestamp in the specified directory.

- `Run(db *sql.DB, path string, noApply bool)`  
  Executes all local migrations that haven't been applied yet, or retries any that previously failed.

  **Parameters:**

  - `db`: active database connection.
  - `path`: path to the directory where migration files are stored.
  - `noApply`: if `true`, skips executing the SQL statements and only logs or marks them as applied (dry-run mode or when the database already has data).

## Features

- 🔒 Safe migrations with **transaction support**
- 📝 Tracks applied migrations with **checksum validation**
- 🧪 **Dry-run mode** (`noApply`) to validate without applying or using preview data
- ⚡ Zero external dependencies, built with Go standard library

---

## 🏗️ Folder Diagram

```
pkg/
└── migrate/               # Core migration logic
    ├── create.go          # Create new migration files
    ├── executor.go        # Apply and reapply migrations
    ├── loader.go          # Load local and applied migrations
    ├── migrate.go         # Public interface to run migrations
    ├── storage.go         # Manage migrations table in the database
    └── types.go           # Struct definitions (Migration, AppliedMigration)
```
