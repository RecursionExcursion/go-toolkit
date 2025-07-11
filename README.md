# 🧰 go-toolkit

A personal utility library for Go, offering reusable tools and helpers across a wide range of applications. Designed to simplify everyday tasks such as HTTP request batching, environment loading, JWT handling, MongoDB operations, and more — with zero dependencies unless explicitly used.

---

## 📦 Package Structure

- `core/` — General-purpose, dependency-free utilities
- `jwt/` — JWT token creation and validation (wraps `golang-jwt/jwt/v5`)
- `mongo/` — MongoDB client wrapper (wraps `go.mongodb.org/mongo-driver`)

---

## ✨ Core Utilities

### `fetch_and_map.go`
- `FetchAndMap[T any](urls []string, fn func([]byte) (T, error)) []T`  
  Concurrently fetches a list of URLs and maps their responses into a type.

### `request_batcher.go`
- `RequestBatcher` — Queue and execute HTTP requests in configurable batches.

### `env_loader.go`
- `LoadEnvVars(required []string)`  
  Load required environment variables or panic with a detailed message.

### `compression.go`
- Helpers for `gzip` compression and decompression.

### `logger.go`
- Minimal log wrapper with leveled printing.

### `prettyLog.go`
- `PrettyLog(data any)` — Prints structs/maps in color-coded, readable JSON.

### `converters.go`
- Basic type-safe casting and string conversion helpers.

### `testing_flags.go`
- CLI test-mode toggle for integration tests or conditional logic.

### `timer.go`
- Simple stopwatch-style `Timer` struct with duration reporting.

---

## 🔐 JWT Package

### `jwt/jwt.go`
- `GenerateToken(claims jwt.MapClaims, secret string) (string, error)`
- `ValidateToken(tokenStr, secret string) (*jwt.Token, error)`

> Requires: `github.com/golang-jwt/jwt/v5`

---

## 🍃 Mongo Package

### `mongo/mongo.go`
- `ConnectMongo(uri string) (*mongo.Client, error)`
- `FindOne[T any](collection *mongo.Collection, filter interface{}) (T, error)`
- `InsertOne(collection *mongo.Collection, doc interface{}) (primitive.ObjectID, error)`

> Requires: `go.mongodb.org/mongo-driver`

---

## 📦 Installation

```bash
go get github.com/RecursionExcursion/go-toolkit

Import only what you need:

import "github.com/RecursionExcursion/go-toolkit/core"
import "github.com/RecursionExcursion/go-toolkit/mongo"
import "github.com/RecursionExcursion/go-toolkit/jwt"
```

🧼 Dependency Philosophy

    All core utilities are 100% dependency-free.

    jwt and mongo use third-party packages, but are isolated to avoid polluting the main dependency tree.

    Only importing the subpackages will bring in their respective dependencies.

📄 License

MIT © RecursionExcursion
