package cipherPayload

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

type AES interface {
	Encrypt(plainText string) (string, error)
	Decrypt(encryptedText string) (string, error)
	pkcs7Unpad(b []byte, blocksize int) ([]byte, error)
	pkcs7Pad(b []byte, blocksize int) ([]byte, error)
}

type DefaultAESEncryption struct {
	keys   KeyPairs
	logger logger
}

func NewAESEncryption(keys KeyPairs, debugMode ...bool) AES {
	isShowLog := false
	if len(debugMode) > 0 {
		isShowLog = bool(debugMode[0])
	}
	logger := newLogger(isShowLog)
	return &DefaultAESEncryption{
		keys:   keys,
		logger: logger,
	}
}

func (d *DefaultAESEncryption) Encrypt(plainText string) (string, error) {
	serviceName := "Encryption - AES - Encrypt"
	d.logger.printf("info", serviceName)

	if plainText == "" {
		return "", nil
	}

	byteText := []byte(plainText)
	block, err := aes.NewCipher(d.keys.AESKeyForEncrypt)
	if err != nil {
		d.logger.printf("error", serviceName, err)
		return "", err
	}

	defer recoveryCatch()
	byteText, _ = d.pkcs7Pad(byteText, block.BlockSize())

	ciphertext := make([]byte, len(byteText))

	bm := cipher.NewCBCEncrypter(block, d.keys.AESIVForEncrypt)
	bm.CryptBlocks(ciphertext, byteText)

	cipherText := base64.StdEncoding.EncodeToString(ciphertext)
	d.logger.printf("debug", "Encrypted:", cipherText)
	return cipherText, nil
}

func (d *DefaultAESEncryption) Decrypt(encryptedText string) (string, error) {
	serviceName := "Encryption - AES - Decrypt"
	d.logger.printf("info", serviceName)

	if encryptedText == "" {
		return "", nil
	}
	ciphertext, _ := base64.StdEncoding.DecodeString(encryptedText)
	block, err := aes.NewCipher(d.keys.AESKeyForDecrypt)
	if err != nil {
		d.logger.printf("error", serviceName, err)
		return "", err
	}

	defer recoveryCatch()
	bm := cipher.NewCBCDecrypter(block, d.keys.AESIVForDecrypt)
	bm.CryptBlocks(ciphertext, ciphertext)

	out, err := d.pkcs7Unpad(ciphertext, aes.BlockSize)
	if err != nil {
		d.logger.printf("error", serviceName, err)
		return "", err
	}

	plainText := string(out)
	d.logger.printf("debug", "Decrypted:", plainText)
	return plainText, nil
}

func (d DefaultAESEncryption) pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1] // the last character
	n := int(c)      // convert the last char into a number
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

func (d DefaultAESEncryption) pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}
