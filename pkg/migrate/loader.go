package migrate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// LoadLocal reads all migrations in the directory, sorts them and calculates checksums
func LoadLocal(path string) ([]Migration, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory not found: %s", path)
	}

	files, err := filepath.Glob(filepath.Join(path, "*.sql"))
	if err != nil {
		return nil, fmt.Errorf("error to list migrations: %w", err)
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no migrations found in %s", path)
	}

	sort.Slice(files, func(i, j int) bool {
		fi := filepath.Base(files[i])
		fj := filepath.Base(files[j])

		pi := strings.SplitN(fi, "_", 2)[0]
		pj := strings.SplitN(fj, "_", 2)[0]

		ni, _ := strconv.ParseInt(pi, 10, 64)
		nj, _ := strconv.ParseInt(pj, 10, 64)

		return ni < nj
	})

	var migrations []Migration
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("error to read file %s: %w", f, err)
		}

		name := strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))
		hash := sha256.Sum256(data)

		migrations = append(migrations, Migration{
			Name:     name,
			Checksum: hex.EncodeToString(hash[:]),
			Content:  data,
		})
	}

	return migrations, nil
}
