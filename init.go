package cipherPayload

import (
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var (
	serviceName = "[Middleware] cipherPayload"
)

func New(config ...Config) fiber.Handler {
	// set default config
	cfg := configDefault(config...)

	panicResponseHeader := serviceName + ": Some configuration is missing: "

	// config is required
	if cfg.AESKey == nil || len(cfg.AESKey) == 0 {
		panic(panicResponseHeader + "`AESKey` is required.")
	}

	// config is required
	if cfg.AESIV == nil || len(cfg.AESIV) == 0 {
		panic(panicResponseHeader + "`AESIV` is required.")
	}

	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Don't execute when url is contains "health"
		if strings.Contains(c.OriginalURL(), "health") {
			return c.Next()
		}

		// Don't execute when the method is not matched a list
		if !isExist(cfg.AllowMethod, c.Method()) {
			return c.Next()
		}

		logger := newLogger(cfg.DebugMode)

		var reqBody PayloadBody
		err := c.BodyParser(&reqBody)
		if err != nil || reqBody.Payload == "" {
			errMsg := "Invalid Payload"
			logger.printf("error", errMsg, string(c.Request().Body()))
			return cfg.FailResponse(c, errMsg)
		}
		logger.printf("debug", "Request:", reqBody.Payload)

		// Payload Decrypting
		encrypterDecrypter := NewAESEncryption(cfg.AESKey, cfg.AESIV)
		decryptedPayload, err := encrypterDecrypter.Decrypt(reqBody.Payload)
		logger.printf("debug", "Decrypted:", decryptedPayload)

		if err != nil || decryptedPayload == "" {
			logger.printf("error", serviceName, err)
			errMsg := "Cannot decrypt payload or invalid payload"
			return cfg.FailResponse(c, errMsg)
		}

		jsonRaw := json.RawMessage(decryptedPayload)
		jsonBytes, _ := json.Marshal(jsonRaw)

		// Set plaintext back into request body
		logger.printf("debug", "jsonBytes:", string(jsonBytes))
		c.Request().SetBodyRaw(jsonBytes)

		// Let request to continue execute
		c.Next()

		// Intercept the response body
		interceptBody := string(c.Response().Body())
		logger.printf("debug", "Intercept Body:", interceptBody)

		// Payload Encrypting
		encryptedPayload, err := encrypterDecrypter.Encrypt(interceptBody)
		logger.printf("debug", "Encrypted:", encryptedPayload)

		if err != nil || encryptedPayload == "" {
			logger.printf("error", serviceName, err)
			errMsg := "InternalServerError: Cannot encrypt payload or invalid payload"
			return cfg.ErrorResponse(c, errMsg)
		}

		// Set ciphertext back into response body
		var resBody PayloadBody
		resBody.Payload = encryptedPayload
		return c.JSON(resBody)
	}
}
