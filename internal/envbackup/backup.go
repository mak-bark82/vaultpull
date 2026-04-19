package envbackup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Backup creates a timestamped backup of the given .env file.
// Returns the path of the backup file, or an error.
func Backup(envPath string) (string, error) {
	src, err := os.Open(envPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // nothing to back up
		}
		return "", fmt.Errorf("open source file: %w", err)
	}
	defer src.Close()

	dir := filepath.Dir(envPath)
	base := filepath.Base(envPath)
	timestamp := time.Now().Format("20060102T150405")
	backupName := fmt.Sprintf("%s.%s.bak", base, timestamp)
	backupPath := filepath.Join(dir, backupName)

	dst, err := os.OpenFile(backupPath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return "", fmt.Errorf("create backup file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("copy to backup: %w", err)
	}

	return backupPath, nil
}
