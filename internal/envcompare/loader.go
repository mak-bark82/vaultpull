package envcompare

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadFile reads a .env file and returns a map of key=value pairs.
// Lines starting with '#' and blank lines are skipped.
// Malformed lines (no '=') are silently ignored.
func LoadFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envcompare: open %q: %w", path, err)
	}
	defer f.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		// strip surrounding quotes
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		env[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envcompare: scan %q: %w", path, err)
	}
	return env, nil
}

// CompareFiles loads two .env files and returns a Result.
func CompareFiles(leftPath, rightPath string) (Result, error) {
	left, err := LoadFile(leftPath)
	if err != nil {
		return Result{}, err
	}
	right, err := LoadFile(rightPath)
	if err != nil {
		return Result{}, err
	}
	return Compare(left, right), nil
}
