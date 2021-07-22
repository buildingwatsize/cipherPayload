package cipherPayload

import (
	"github.com/gofiber/fiber/v2"
)

type Config struct {
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Required. Default: KeyPairs{}
	KeyPairs KeyPairs

	// Optional. Default: OPTIONS, POST, PUT, DELETE
	AllowMethod []string

	// Optional. Default: false
	DebugMode bool

	// Optional. [Default: false]
	StrictMode bool

	// Optional. Default: true
	ExcludeHealthAPI bool

	// Optional. Default: BadRequestResponse
	FailResponse func(c *fiber.Ctx, msg string) error

	// Optional. Default: InternalServerErrorResponse
	ErrorResponse func(c *fiber.Ctx, msg string) error
}

var ConfigDefault = Config{
	Next:     nil,
	KeyPairs: KeyPairs{},
	AllowMethod: []string{
		fiber.MethodOptions,
		fiber.MethodPost,
		fiber.MethodPut,
		fiber.MethodDelete,
	},
	DebugMode:        false,
	StrictMode:       false,
	ExcludeHealthAPI: true,
	FailResponse:     BadRequestResponse,
	ErrorResponse:    InternalServerErrorResponse,
}

func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// set default values
	if cfg.AllowMethod == nil {
		cfg.AllowMethod = ConfigDefault.AllowMethod
	}

	// Note: cfg.DebugMode: it's false by default.
	// Note: cfg.StrictMode: it's false by default.
	// NOTE: cfg.ExcludeHealthAPI: it's true by default

	// set default values
	if cfg.FailResponse == nil {
		cfg.FailResponse = ConfigDefault.FailResponse
	}

	// set default values
	if cfg.ErrorResponse == nil {
		cfg.ErrorResponse = ConfigDefault.ErrorResponse
	}

	return cfg
}

type PayloadBody struct {
	Payload string `json:"payload"`
}

type KeyPairs struct {
	AESKeyForEncrypt []byte
	AESIVForEncrypt  []byte
	AESKeyForDecrypt []byte
	AESIVForDecrypt  []byte
}

func BadRequestResponse(c *fiber.Ctx, msg string) error { // 400
	if msg == "" {
		msg = "Bad Request"
	}
	res := fiber.Map{
		"status":  "bad_request",
		"message": msg,
	}
	return c.Status(fiber.StatusBadRequest).JSON(res)
}

func InternalServerErrorResponse(c *fiber.Ctx, msg string) error { // 500
	if msg == "" {
		msg = "Internal Server Error"
	}
	res := fiber.Map{
		"status":  "internal_server_error",
		"message": msg,
	}
	return c.Status(fiber.StatusInternalServerError).JSON(res)
}
