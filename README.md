# CipherPayload

CipherPayload middleware for Fiber that use AES Algorithm for encrypt and decrypt payload in request and response body.

## Table of Contents

- [CipherPayload](#cipherpayload)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
  - [Signatures](#signatures)
  - [Examples](#examples)
  - [Config](#config)
  - [Default Config](#default-config)
  - [Default Response](#default-response)
  - [Payload Template](#payload-template)
  - [Example Usage](#example-usage)

## Installation

```bash
  go get -u github.com/buildingwatsize/cipherPayload
```

## Signatures

```go
func New(config ...Config) fiber.Handler
```

## Examples

Import the middleware package that is part of the Fiber web framework

```go
import (
  "github.com/gofiber/fiber/v2"
  "github.com/buildingwatsize/cipherPayload"
)
```

After you initiate your Fiber app, you can use the following possibilities:

```go
// Default middleware config
app.Use(cipherPayload.New(cipherPayload.Config{
  AESKey:    []byte("AES_KEY"),
  AESIV:     []byte("AES_IV"),
}))

// Or extend your config for customization
app.Use(cipherPayload.New(cipherPayload.Config{
  AESKey:    []byte("AES_KEY"),
  AESIV:     []byte("AES_IV"),
  AllowMethod: []string{"POST", "OPTIONS"},
  DebugMode: true,
}))
```

## Config

```go
// Config defines the config for middleware.
type Config struct {
  // Next defines a function to skip this middleware when returned true.

  // Optional. Default: nil
  Next func(c *fiber.Ctx) bool

  // Required. Default: nil
  AESKey []byte

  // Required. Default: nil
  AESIV []byte

  // Optional. Default: ["OPTIONS", "POST", "PUT", "DELETE"]
  AllowMethod []string

  // Optional. Default: false
  DebugMode bool

  // Optional. Default: true
  IncludeHealthAPI bool

  // Optional. Default: BadRequestResponse
  FailResponse func(c *fiber.Ctx, msg string) error

  // Optional. Default: InternalServerErrorResponse
  ErrorResponse func(c *fiber.Ctx, msg string) error
}
```

## Default Config

```go
var ConfigDefault = Config{
  Next:   nil,
  AESKey: nil,
  AESIV:  nil,
  AllowMethod: []string{
    fiber.MethodOptions,
    fiber.MethodPost,
    fiber.MethodPut,
    fiber.MethodDelete,
  },
  DebugMode:        false,
  IncludeHealthAPI: true,
  FailResponse:     BadRequestResponse,
  ErrorResponse:    InternalServerErrorResponse,
}
```

## Default Response

```go
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
```

## Payload Template

```json
{
  "payload": "tpkWPEI6F/nfgUjjtwyKSUHhNwBF3fZilleEukZ5GRazPN4rbfuqOasHeNN3OpDG"
}
```

Note: Using

- AES Key: `12345678901234567890123456789012`
- AES IV: `1234567890123456`

Which should be equal to

```json
{
  "firstname": "Chinnawat",
  "lastname": "Chimdee"
}
```

## Example Usage

Please go to [example/README.md](./example/README.md)
