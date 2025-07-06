package logger_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/derpartizanen/gophkeeper/internal/logger"
)

func TestNewLogger(t *testing.T) {
	tt := []struct {
		name     string
		level    string
		expected zerolog.Level
	}{
		{
			name:     "Set info level",
			level:    "info",
			expected: zerolog.InfoLevel,
		},
		{
			name:     "Set warning level",
			level:    "warn",
			expected: zerolog.WarnLevel,
		},
		{
			name:     "Set error level",
			level:    "error",
			expected: zerolog.ErrorLevel,
		},
		{
			name:     "Set debug level",
			level:    "debug",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "Set uppercase level",
			level:    "Debug",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "Set capitalized level",
			level:    "DEBUG",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "Empty string",
			level:    "",
			expected: zerolog.InfoLevel,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := logger.New(tc.level)

			require.NoError(t, err)
			require.Equal(t, tc.expected, zerolog.GlobalLevel())
		})
	}
}

func TestNewLoggerFailsOnBadLevel(t *testing.T) {
	_, err := logger.New("xxx")

	require.Error(t, err)
}
