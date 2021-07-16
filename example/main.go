package main

import (
	"fmt"

	"github.com/buildingwatsize/cipherPayload"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	aesSecretKeyForEncrypt []byte
	aesSecretIVForEncrypt  []byte
	aesSecretKeyForDecrypt []byte
	aesSecretIVForDecrypt  []byte
)

func init() {
	aesSecretKeyForEncrypt = []byte("12345678901234567890123456789012")
	aesSecretIVForEncrypt = []byte("1234567890123456")
	aesSecretKeyForDecrypt = []byte("67890123456789012345678901234567")
	aesSecretIVForDecrypt = []byte("6789012345678901")

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	// NOTE: For more usage please go to README.md

	// Without cipherPayload middleware
	// ================================
	appPlaintext := fiber.New()
	apiPlaintextGroup := appPlaintext.Group("/api", nextHandlerAPI) // /api
	apiPlaintextGroup.Post("/example", HandlerExamplePlaintext)     // /api/example (FullPath = http://localhost:5000/api/example)

	// With cipherPayload middleware
	// =============================
	appCipher := fiber.New()
	appCipher.Use(cipherPayload.New(cipherPayload.Config{
		KeyPairs: cipherPayload.KeyPairs{
			AESKeyForEncrypt: aesSecretKeyForEncrypt,
			AESIVForEncrypt:  aesSecretIVForEncrypt,
			AESKeyForDecrypt: aesSecretKeyForDecrypt,
			AESIVForDecrypt:  aesSecretIVForDecrypt,
		},
		DebugMode: true,
	}))
	apiCipherGroup := appCipher.Group("/api", nextHandlerAPI) // /api
	apiCipherGroup.Post("/example", HandlerExampleCipher)     // /api/example (FullPath = http://localhost:8000/api/example)

	PortPlaintext := "5000"
	go appPlaintext.Listen(":" + PortPlaintext)
	PortCipher := "8000"
	appCipher.Listen(":" + PortCipher)
}

var nextHandlerAPI = func(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMEApplicationJSONCharsetUTF8)

	log.Info().Msg("Called API")
	return c.Next()
}

var HandlerExamplePlaintext = func(c *fiber.Ctx) error {
	log.Info().Msg("Called HandlerExamplePlaintext")

	var reqBody map[string]interface{}
	_ = c.BodyParser(&reqBody)
	for k, v := range reqBody {
		reqBody[k] = fmt.Sprintf("%v [Modified]", v)
	}

	return c.Status(fiber.StatusOK).JSON(reqBody)
}

// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
// === YOU SEE ? HOW IT DIFFERENCE ? ===
// :: PS. the "Plaintext Way" and the "Cipher Way"
//        can be used the same handler ::
// vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv

var HandlerExampleCipher = func(c *fiber.Ctx) error {
	log.Info().Msg("Called HandlerExampleCipher")

	var reqBody map[string]interface{}
	_ = c.BodyParser(&reqBody)
	for k, v := range reqBody {
		reqBody[k] = fmt.Sprintf("%v [Modified]", v)
	}

	return c.Status(fiber.StatusOK).JSON(reqBody)
}
