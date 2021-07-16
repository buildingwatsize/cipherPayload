# CipherPayload

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org) [![Go Reference](https://pkg.go.dev/badge/github.com/buildingwatsize/cipherPayload@v0.1.0.svg)](https://pkg.go.dev/github.com/buildingwatsize/cipherPayload@v0.1.0) [![GitHub issues](https://img.shields.io/github/issues/buildingwatsize/cipherPayload)](https://github.com/buildingwatsize/cipherPayload/issues) [![GitHub forks](https://img.shields.io/github/forks/buildingwatsize/cipherPayload)](https://github.com/buildingwatsize/cipherPayload/network) [![GitHub stars](https://img.shields.io/github/stars/buildingwatsize/cipherPayload)](https://github.com/buildingwatsize/cipherPayload/stargazers)

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
  - [KeyPairs Property](#keypairs-property)
  - [Payload Template](#payload-template)
    - [Request](#request)
    - [Response](#response)
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
  KeyPairs: cipherPayload.KeyPairs{
    AESKeyForEncrypt: []byte("AES_KEY_FOR_ENCRYPT"),
    AESIVForEncrypt:  []byte("AES_IV_FOR_ENCRYPT"),
    AESKeyForDecrypt: []byte("AES_KEY_FOR_DECRYPT"),
    AESIVForDecrypt:  []byte("AES_IV_FOR_DECRYPT"),
  },
}))

// Or extend your config for customization
app.Use(cipherPayload.New(cipherPayload.Config{
  KeyPairs: cipherPayload.KeyPairs{
    AESKeyForEncrypt: []byte("AES_KEY_FOR_ENCRYPT"),
    AESIVForEncrypt:  []byte("AES_IV_FOR_ENCRYPT"),
    AESKeyForDecrypt: []byte("AES_KEY_FOR_DECRYPT"),
    AESIVForDecrypt:  []byte("AES_IV_FOR_DECRYPT"),
  },
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

  // Required. Default: KeyPairs{}
  KeyPairs KeyPairs

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
  KeyPairs: KeyPairs{},
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

## KeyPairs Property

```go
type KeyPairs struct {
  AESKeyForEncrypt []byte
  AESIVForEncrypt  []byte
  AESKeyForDecrypt []byte
  AESIVForDecrypt  []byte
}
```

## Payload Template

An example of payload template (see more how to work in [Example](./example))

### Request

```json
{
  "payload": "FDp1Dl31zGx5nRXFNKihB+k3ly/L7HI9tlHycbKVRwhaf3RRdyFGviuntEZqst0/"
}
```

which can be decrypt to:

```json
{
  "firstname": "Chinnawat",
  "lastname": "Chimdee"
}
```

### Response

```json
{
  "payload": "tpkWPEI6F/nfgUjjtwyKSUf1erxPL6rQt8jG3RitQ1KpvRALfR5YAgQ0CXYkrwLfTid6VdK3SNlffuu/kvI7Hj7br0ur01TUFUWxQ9cl+8U="
}
```

which encrypted from:

```json
{
  "firstname": "Chinnawat [Modified]",
  "lastname": "Chimdee [Modified]"
}
```

Note: This payload using

- AESKeyForEncrypt (used in encrypting response body): `67890123456789012345678901234567`
- AESIVForEncrypt (used in encrypting response body): `6789012345678901`
- AESKeyForDecrypt (used in decrypting request body): `12345678901234567890123456789012`
- AESIVForDecrypt (used in decrypting request body): `1234567890123456`

## Example Usage

Please go to [example/README.md](./example/README.md)
