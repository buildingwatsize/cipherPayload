package cipherPayload

import (
	"bytes"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	defer func() {
		log.Logger = log.Output(os.Stderr)
	}()
	f()
	return buf.String()
}

type testCaseIsExist struct {
	name     string
	inputArr []string
	inputStr string
	expected bool
}

func TestIsExist(t *testing.T) {
	tests := []testCaseIsExist{
		{
			name:     "Empty Case",
			inputArr: []string{},
			inputStr: "",
			expected: false,
		},
		{
			name:     "Valid Case",
			inputArr: []string{"1", "2", "3"},
			inputStr: "1",
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputArr := test.inputArr
			inputStr := test.inputStr
			expected := test.expected

			actual := isExist(inputArr, inputStr)

			require.Equal(t, expected, actual)
		})
	}
}

type testCaseLogger struct {
	name            string
	inputLogType    string
	inputLogMessage string
	expected        string
}

func TestLogger(t *testing.T) {
	tests := []testCaseLogger{
		{
			name:            "Default Case",
			inputLogType:    "",
			inputLogMessage: "",
			expected:        "",
		},
		{
			name:            "Log Error Case",
			inputLogType:    "error",
			inputLogMessage: "test error",
			expected:        "test error",
		},
		{
			name:            "Log Info Case",
			inputLogType:    "info",
			inputLogMessage: "test info",
			expected:        "test info",
		},
		{
			name:            "Log Debug Case",
			inputLogType:    "debug",
			inputLogMessage: "test debug",
			expected:        "test debug",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputLogType := test.inputLogType
			inputLogMessage := test.inputLogMessage
			expected := test.expected

			actual := captureOutput(func() {
				logger := newLogger(true)
				logger.printf(inputLogType, inputLogMessage)
			})

			require.Contains(t, actual, expected)
		})
	}
}
