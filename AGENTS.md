# Agent Guidelines for qrcatalog

## Build/Lint/Test Commands

### Frontend (React/TypeScript)
- **Build**: `cd frontend && npm run build` (runs vite build && tsc)
- **Test**: `cd frontend && npm run test` (runs vitest run)
- **Format**: `cd frontend && npm run format` (runs prettier --write .)
- **Dev server**: `cd frontend && npm run dev` (runs vite --port 3000)

### Backend (Go)
- **Build**: `go build -o ./tmp/main .`
- **Dev server**: `air` (hot reload with .air.toml config)
- **Test**: `go test ./...`
- **Run single test**: `go test -run TestName ./path/to/package`

## Code Style Guidelines

### Go
- **Naming**: PascalCase for exported identifiers, camelCase for unexported
- **Error handling**: Return errors from functions, use `log.Printf` for logging
- **Package comments**: Use `// Package name provides...` format
- **Imports**: Standard library → third-party → internal packages
- **Database**: Use pgx/v5 with context for all operations
- **HTTP handlers**: Use `http.HandlerFunc` pattern with proper status codes

### TypeScript/React
- **Formatting**: Prettier with semicolons, double quotes, trailing commas, 4-space tabs, 96 char width
- **TypeScript**: Strict mode enabled, no unused locals/parameters
- **Imports**: Use path mapping `@/*` for src directory
- **Components**: Functional components with hooks
- **Routing**: TanStack Router with generated route tree
- **Styling**: Tailwind CSS with Radix UI components
- **Forms**: React Hook Form with Zod validation

### General
- **Security**: Never log or commit secrets/keys
- **Comments**: No unnecessary comments, focus on self-documenting code
- **Dependencies**: Only use libraries already present in the project