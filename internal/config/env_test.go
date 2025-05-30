package config

import (
	"os"
	"testing"
	"time"
)

func TestGetEnvString(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue string
		expected     string
	}{
		{
			name:         "environment variable set",
			key:          "TEST_STRING",
			envValue:     "test-value",
			defaultValue: "default",
			expected:     "test-value",
		},
		{
			name:         "environment variable with whitespace",
			key:          "TEST_STRING_WHITESPACE",
			envValue:     "  test-value  ",
			defaultValue: "default",
			expected:     "test-value",
		},
		{
			name:         "environment variable empty",
			key:          "TEST_STRING_EMPTY",
			envValue:     "",
			defaultValue: "default",
			expected:     "default",
		},
		{
			name:         "environment variable not set",
			key:          "TEST_STRING_NOT_SET",
			envValue:     "", // Will not be set
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			defer os.Unsetenv(tt.key)

			// Set environment variable if value provided
			if tt.envValue != "" || tt.name == "environment variable empty" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvString(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvString(%q, %q) = %q, want %q",
					tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		// Truthy values
		{
			name:         "true value",
			key:          "TEST_BOOL_TRUE",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "TRUE value (case insensitive)",
			key:          "TEST_BOOL_TRUE_UPPER",
			envValue:     "TRUE",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "True value (mixed case)",
			key:          "TEST_BOOL_TRUE_MIXED",
			envValue:     "True",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "1 value (numeric true)",
			key:          "TEST_BOOL_ONE",
			envValue:     "1",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "t value (short true)",
			key:          "TEST_BOOL_T",
			envValue:     "t",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "T value (short true upper)",
			key:          "TEST_BOOL_T_UPPER",
			envValue:     "T",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "yes value",
			key:          "TEST_BOOL_YES",
			envValue:     "yes",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "YES value (case insensitive)",
			key:          "TEST_BOOL_YES_UPPER",
			envValue:     "YES",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "y value (short yes)",
			key:          "TEST_BOOL_Y",
			envValue:     "y",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "on value",
			key:          "TEST_BOOL_ON",
			envValue:     "on",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "ON value (case insensitive)",
			key:          "TEST_BOOL_ON_UPPER",
			envValue:     "ON",
			defaultValue: false,
			expected:     true,
		},

		// Falsy values
		{
			name:         "false value",
			key:          "TEST_BOOL_FALSE",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "FALSE value (case insensitive)",
			key:          "TEST_BOOL_FALSE_UPPER",
			envValue:     "FALSE",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "0 value (numeric false)",
			key:          "TEST_BOOL_ZERO",
			envValue:     "0",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "f value (short false)",
			key:          "TEST_BOOL_F",
			envValue:     "f",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "F value (short false upper)",
			key:          "TEST_BOOL_F_UPPER",
			envValue:     "F",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "no value",
			key:          "TEST_BOOL_NO",
			envValue:     "no",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "NO value (case insensitive)",
			key:          "TEST_BOOL_NO_UPPER",
			envValue:     "NO",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "n value (short no)",
			key:          "TEST_BOOL_N",
			envValue:     "n",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "off value",
			key:          "TEST_BOOL_OFF",
			envValue:     "off",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "OFF value (case insensitive)",
			key:          "TEST_BOOL_OFF_UPPER",
			envValue:     "OFF",
			defaultValue: true,
			expected:     false,
		},

		// Edge cases
		{
			name:         "invalid value uses default (true)",
			key:          "TEST_BOOL_INVALID_TRUE",
			envValue:     "invalid",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "invalid value uses default (false)",
			key:          "TEST_BOOL_INVALID_FALSE",
			envValue:     "invalid",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "empty value uses default",
			key:          "TEST_BOOL_EMPTY",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "whitespace around true",
			key:          "TEST_BOOL_WHITESPACE",
			envValue:     "  true  ",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "whitespace around 1",
			key:          "TEST_BOOL_WHITESPACE_ONE",
			envValue:     "  1  ",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "typo falls back to default",
			key:          "TEST_BOOL_TYPO",
			envValue:     "tru", // Common typo
			defaultValue: false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvBool(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvBool(%q, %v) with env value %q = %v, want %v",
					tt.key, tt.defaultValue, tt.envValue, result, tt.expected)
			}
		})
	}
}

func TestGetEnvStringSlice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		separator    string
		defaultValue []string
		expected     []string
	}{
		{
			name:         "comma separated values",
			key:          "TEST_SLICE_COMMA",
			envValue:     "app.com,admin.com,api.com",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com", "api.com"},
		},
		{
			name:         "comma separated with whitespace",
			key:          "TEST_SLICE_WHITESPACE",
			envValue:     "app.com, admin.com , api.com ",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com", "api.com"},
		},
		{
			name:         "empty items filtered out",
			key:          "TEST_SLICE_EMPTY_ITEMS",
			envValue:     "app.com,,admin.com,",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com"},
		},
		{
			name:         "single value",
			key:          "TEST_SLICE_SINGLE",
			envValue:     "app.com",
			separator:    ",",
			defaultValue: nil,
			expected:     []string{"app.com"},
		},
		{
			name:         "empty value uses default",
			key:          "TEST_SLICE_DEFAULT",
			envValue:     "",
			separator:    ",",
			defaultValue: []string{"default1", "default2"},
			expected:     []string{"default1", "default2"},
		},
		{
			name:         "semicolon separator",
			key:          "TEST_SLICE_SEMICOLON",
			envValue:     "app.com;admin.com;api.com",
			separator:    ";",
			defaultValue: nil,
			expected:     []string{"app.com", "admin.com", "api.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvStringSlice(tt.key, tt.separator, tt.defaultValue)

			// Compare slices
			if len(result) != len(tt.expected) {
				t.Errorf("GetEnvStringSlice(%q, %q, %v) length = %d, want %d",
					tt.key, tt.separator, tt.defaultValue, len(result), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("GetEnvStringSlice(%q, %q, %v)[%d] = %q, want %q",
						tt.key, tt.separator, tt.defaultValue, i, result[i], expected)
				}
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "valid integer",
			key:          "TEST_INT_VALID",
			envValue:     "42",
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "zero value",
			key:          "TEST_INT_ZERO",
			envValue:     "0",
			defaultValue: 10,
			expected:     0,
		},
		{
			name:         "negative integer",
			key:          "TEST_INT_NEGATIVE",
			envValue:     "-5",
			defaultValue: 10,
			expected:     -5,
		},
		{
			name:         "invalid value uses default",
			key:          "TEST_INT_INVALID",
			envValue:     "not-a-number",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "empty value uses default",
			key:          "TEST_INT_EMPTY",
			envValue:     "",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "whitespace around number",
			key:          "TEST_INT_WHITESPACE",
			envValue:     "  42  ",
			defaultValue: 10,
			expected:     42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvInt(%q, %d) = %d, want %d",
					tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "seconds duration",
			key:          "TEST_DURATION_SECONDS",
			envValue:     "30s",
			defaultValue: 10 * time.Second,
			expected:     30 * time.Second,
		},
		{
			name:         "minutes duration",
			key:          "TEST_DURATION_MINUTES",
			envValue:     "5m",
			defaultValue: 1 * time.Minute,
			expected:     5 * time.Minute,
		},
		{
			name:         "hours duration",
			key:          "TEST_DURATION_HOURS",
			envValue:     "2h",
			defaultValue: 1 * time.Hour,
			expected:     2 * time.Hour,
		},
		{
			name:         "mixed duration",
			key:          "TEST_DURATION_MIXED",
			envValue:     "1h30m",
			defaultValue: 1 * time.Hour,
			expected:     90 * time.Minute,
		},
		{
			name:         "invalid duration uses default",
			key:          "TEST_DURATION_INVALID",
			envValue:     "not-a-duration",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
		{
			name:         "empty value uses default",
			key:          "TEST_DURATION_EMPTY",
			envValue:     "",
			defaultValue: 10 * time.Second,
			expected:     10 * time.Second,
		},
		{
			name:         "whitespace around duration",
			key:          "TEST_DURATION_WHITESPACE",
			envValue:     "  45s  ",
			defaultValue: 10 * time.Second,
			expected:     45 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Unsetenv(tt.key)

			if tt.envValue != "" || tt.name == "empty value uses default" {
				os.Setenv(tt.key, tt.envValue)
			}

			result := GetEnvDuration(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("GetEnvDuration(%q, %v) = %v, want %v",
					tt.key, tt.defaultValue, result, tt.expected)
			}
		})
	}
}
