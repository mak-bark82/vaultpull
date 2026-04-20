// Package envschema provides JSON Schema-based validation for environment variable sets.
// It allows callers to define a schema describing expected keys, types, and constraints,
// then validate a map of env vars against that schema before writing to disk.
package envschema

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// FieldType represents the expected data type for an environment variable value.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInteger FieldType = "integer"
	TypeBoolean FieldType = "boolean"
	TypeURL     FieldType = "url"
)

// FieldSchema describes validation rules for a single environment variable.
type FieldSchema struct {
	Type     FieldType `yaml:"type"`
	Required bool      `yaml:"required"`
	Pattern  string    `yaml:"pattern,omitempty"`
	MinLen   int       `yaml:"min_length,omitempty"`
	MaxLen   int       `yaml:"max_length,omitempty"`
}

// Schema maps environment variable names to their field schemas.
type Schema struct {
	Fields map[string]FieldSchema `yaml:"fields"`
}

// ValidationError holds all errors found during schema validation.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("schema validation failed with %d error(s):\n  %s",
		len(e.Errors), strings.Join(e.Errors, "\n  "))
}

// LoadSchema reads and parses a YAML schema file from the given path.
func LoadSchema(path string) (*Schema, error) {
	if path == "" {
		return nil, fmt.Errorf("schema path must not be empty")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading schema file: %w", err)
	}
	var s Schema
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing schema file: %w", err)
	}
	if s.Fields == nil {
		s.Fields = make(map[string]FieldSchema)
	}
	return &s, nil
}

// Validate checks the provided env map against the schema.
// It returns a *ValidationError if any rules are violated, or nil on success.
func (s *Schema) Validate(env map[string]string) error {
	var errs []string

	for key, field := range s.Fields {
		val, exists := env[key]

		if field.Required && (!exists || strings.TrimSpace(val) == "") {
			errs = append(errs, fmt.Sprintf("%s: required but missing or empty", key))
			continue
		}

		if !exists {
			continue
		}

		if err := validateType(key, val, field.Type); err != nil {
			errs = append(errs, err.Error())
		}

		if field.Pattern != "" {
			re, err := regexp.Compile(field.Pattern)
			if err != nil {
				errs = append(errs, fmt.Sprintf("%s: invalid pattern %q: %v", key, field.Pattern, err))
			} else if !re.MatchString(val) {
				errs = append(errs, fmt.Sprintf("%s: value %q does not match pattern %q", key, val, field.Pattern))
			}
		}

		if field.MinLen > 0 && len(val) < field.MinLen {
			errs = append(errs, fmt.Sprintf("%s: value too short (min %d chars)", key, field.MinLen))
		}
		if field.MaxLen > 0 && len(val) > field.MaxLen {
			errs = append(errs, fmt.Sprintf("%s: value too long (max %d chars)", key, field.MaxLen))
		}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}

// validateType checks that val conforms to the expected FieldType.
func validateType(key, val string, t FieldType) error {
	switch t {
	case TypeInteger:
		if _, err := strconv.Atoi(val); err != nil {
			return fmt.Errorf("%s: expected integer, got %q", key, val)
		}
	case TypeBoolean:
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return fmt.Errorf("%s: expected boolean (true/false/1/0), got %q", key, val)
		}
	case TypeURL:
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			return fmt.Errorf("%s: expected URL starting with http:// or https://, got %q", key, val)
		}
	case TypeString, "":
		// no type-level check needed
	}
	return nil
}
