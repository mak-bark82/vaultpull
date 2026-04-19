// Package envreader parses an existing .env file into a key/value map.
package envreader

import (
	"bufio"
	"os"
	"strings"
)

// Read parses a .env file at path and returns a map of key=value pairs.
// If the file does not exist, an empty map is returned without error.
func Read(path string) (map[string]string, error) {
	result := make(map[string]string)

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		result[key] = val
	}

	return result, scanner.Err()
}
