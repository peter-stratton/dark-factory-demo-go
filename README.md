# Dark Factory Demo: Go Bookmarks API

A simple Go REST API for managing bookmarks, designed as a hands-on tutorial for [Dark Factory](https://godarkfactory.com).

## What you'll learn

This demo walks you through real-world Dark Factory workflows using a Go project:

1. **Basic PR workflow** — Dark Factory creates a branch, opens a PR, and gets it merged
2. **Manual vs auto merge** — when to require human approval vs let changes land automatically
3. **Watch mode** — iterating on changes with `--watch` and seeing live feedback
4. **Mechanical verification: tests** — using `go test` as a merge gate
5. **Mechanical verification: build** — ensuring `go build` passes before merge
6. **CI pipeline integration** — GitHub Actions required status checks with Dark Factory

## Getting started

### 1. Create your own repo from this template

Click the **"Use this template"** button above, then **"Create a new repository"** to get your own copy.

### 2. Clone your repo

```bash
git clone https://github.com/YOUR_USERNAME/dark-factory-demo-go.git
cd dark-factory-demo-go
```

### 3. Verify the project works

```bash
make test    # run tests
make lint    # run linter
make build   # compile the binary
make run     # start the server on :8080
```

### 4. Try the API

```bash
# Health check
curl http://localhost:8080/health

# Create a bookmark
curl -X POST http://localhost:8080/bookmarks \
  -H "Content-Type: application/json" \
  -d '{"url":"https://godarkfactory.com","title":"Dark Factory"}'

# List bookmarks
curl http://localhost:8080/bookmarks
```

## Project structure

```
.
├── main.go                        # Entry point
├── internal/
│   ├── model/
│   │   ├── bookmark.go            # Data types and validation
│   │   └── errors.go              # Domain errors
│   ├── server/
│   │   ├── server.go              # HTTP handlers and routing
│   │   └── server_test.go         # Handler tests
│   └── store/
│       ├── store.go               # In-memory data store
│       └── store_test.go          # Store tests
├── .github/workflows/ci.yml      # GitHub Actions CI
├── .golangci.yml                  # Linter configuration
├── Makefile                       # Build targets
└── go.mod                         # Go module definition
```

## Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [golangci-lint](https://golangci-lint.run/welcome/install/)
- [Dark Factory CLI](https://godarkfactory.com)

## Configuring the Claude model

Dark Factory uses whatever model is configured in your Claude Code project settings. To change the model (e.g. to use Sonnet for faster/cheaper runs, or Opus for more complex tasks), update your Claude Code settings — Dark Factory will pick it up automatically.
