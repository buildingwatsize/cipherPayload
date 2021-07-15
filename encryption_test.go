package cipherPayload

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	aesKey []byte
	aesIV  []byte
)

func init() {
	setAValidConfig()
}

func setAValidConfig() {
	aesKey = []byte("12345678901234567890123456789012")
	aesIV = []byte("1234567890123456")
}

type testCaseEncryption struct {
	name     string
	input    string
	expected string
}

func TestEncryptionAES(t *testing.T) {
	setAValidConfig()

	tests := []testCaseEncryption{
		{
			name:     "Encryption Success",
			input:    "1234567890121",
			expected: "",
		},
		{
			name:     "Encryption Success - Empty Case",
			input:    "",
			expected: "SOMETHING ELSE",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.input
			expected := test.expected

			d := NewAESEncryption(aesKey, aesIV, false)
			actual, _ := d.Encrypt(input)

			require.NotEqual(t, expected, actual)
		})
	}

	setAValidConfig()
	// Fail Case
	aesSecretKeyBackup := aesKey
	t.Run("Fail Case", func(t *testing.T) {
		aesKey = []byte("")
		input := "smtg"

		d := NewAESEncryption(aesKey, aesIV)
		_, err := d.Encrypt(input)

		require.Error(t, err)
	})
	aesKey = aesSecretKeyBackup
}

func TestDecryptionAES(t *testing.T) {
	setAValidConfig()

	tests := []testCaseEncryption{
		{
			name:     "Decryption Success",
			input:    "y/hNE1N17iIxSqvM0IWidQ==",
			expected: "1234567890121",
		},
		{
			name:     "Decryption Success - Empty Case",
			input:    "",
			expected: "",
		},
		{
			name:     "Decryption Fail - Invalid",
			input:    "y/hNE1N17iI",
			expected: "",
		},
		{
			name:     "Decryption Failed",
			input:    "a",
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := test.input
			expected := test.expected

			d := NewAESEncryption(aesKey, aesIV)
			actual, _ := d.Decrypt(input)

			require.Equal(t, expected, actual)
		})
	}

	setAValidConfig()
	// Fail Case
	aesSecretKeyBackup := aesKey
	t.Run("Fail Case", func(t *testing.T) {
		aesKey = []byte("")
		input := "smtg"

		d := NewAESEncryption(aesKey, aesIV)
		_, err := d.Decrypt(input)

		require.Error(t, err)
	})
	aesKey = aesSecretKeyBackup
}

func TestPkcs7Unpad(t *testing.T) {
	tests := []struct {
		name           string
		inputByte      []byte
		inputBlocksize int
		expected       []byte
		expectedErr    error
	}{
		{
			name:           "Invalid Blocksize - Zero",
			inputByte:      []byte(""),
			inputBlocksize: 0,
			expected:       []byte(nil),
			expectedErr:    ErrInvalidBlockSize,
		},
		{
			name:           "Invalid InputByte - Nil",
			inputByte:      []byte(nil),
			inputBlocksize: 2,
			expected:       []byte(nil),
			expectedErr:    ErrInvalidPKCS7Data,
		},
		{
			name:           "Invalid Padding - Modulo != 0",
			inputByte:      []byte("abc"),
			inputBlocksize: 2,
			expected:       []byte(nil),
			expectedErr:    ErrInvalidPKCS7Padding,
		},
		{
			name:           "Invalid Padding - Trailing with space",
			inputByte:      []byte("abc "),
			inputBlocksize: 4,
			expected:       []byte(nil),
			expectedErr:    ErrInvalidPKCS7Padding,
		},
		{
			name: "Invalid Padding - Trailing with tab space",
			inputByte: []byte("y/hNE1N17iIxSqvM0IWidQ=	"),
			inputBlocksize: 24,
			expected:       []byte(nil),
			expectedErr:    ErrInvalidPKCS7Padding,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputByte := test.inputByte
			inputBlocksize := test.inputBlocksize
			expected := test.expected
			expectedErr := test.expectedErr

			d := NewAESEncryption(aesKey, aesIV)
			actual, actualErr := d.pkcs7Unpad(inputByte, inputBlocksize)

			require.Equal(t, expected, actual)
			require.Equal(t, expectedErr, actualErr)
		})
	}
}

func TestPkcs7Pad(t *testing.T) {
	tests := []struct {
		name           string
		inputByte      []byte
		inputBlocksize int
		expected       []byte
		expectedErr    error
	}{
		{
			name:           "Invalid Blocksize - Lower than Zero",
			inputByte:      []byte(""),
			inputBlocksize: -1,
			expected:       []byte(nil),
			expectedErr:    ErrInvalidBlockSize,
		},
		{
			name:           "Invalid Data - Empty",
			inputByte:      []byte(nil),
			inputBlocksize: 2,
			expected:       []byte(nil),
			expectedErr:    ErrInvalidPKCS7Data,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputByte := test.inputByte
			inputBlocksize := test.inputBlocksize
			expected := test.expected
			expectedErr := test.expectedErr

			d := NewAESEncryption(aesKey, aesIV)
			actual, actualErr := d.pkcs7Pad(inputByte, inputBlocksize)

			require.Equal(t, expected, actual)
			require.Equal(t, expectedErr, actualErr)
		})
	}
}
