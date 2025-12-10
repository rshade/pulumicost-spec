# Agent Guidelines for PulumiCost Spec

## Build/Test Commands

- **All tests**: `make test` or `go test ./...`
- **Single test**: `go test -run TestName ./path/to/package`
- **Validate all**: `make validate` (tests + linting + schemas)
- **Generate protobuf**: `make generate`
- **Lint Go**: `make lint` or `golangci-lint run`
- **Validate schemas**: `npm run validate:schema && npm run validate:examples`

## Code Style

- **Go version**: 1.24.10 (toolchain 1.25.4)
- **Formatting**: `goimports` + `golines` (120 char lines)
- **Linting**: 120+ linters via golangci-lint (see `.golangci.yml`)
- **Imports**: Standard library first, then third-party, then local
- **Naming**: CamelCase (protobuf snake_case → Go CamelCase)
- **Functions**: Small, composable; add `String()` for new enums
- **Error handling**: Check all errors; use table-driven tests
- **Tests**: Separate `_test` packages; cover error paths
- **Markdown/YAML**: Follow `.markdownlint.json` and `.yamllint`

## Key Rules

- Never edit generated code (`sdk/go/proto/`, `bin/buf`)
- Run `make validate` before commits
- Use conventional commits: `feat:`, `fix:`, `chore:`
- Sanitize secrets in examples

## Active Technologies

- Go 1.24.10 (toolchain 1.25.4) + gRPC, protobuf, buf v1.32.1 (001-get-budgets-rpc)
- JSON Schema (Draft 2020-12) for PricingSpec and BudgetSpec validation (001-get-budgets-rpc)

## Recent Changes

- 001-get-budgets-rpc: Added Go 1.24.10 (toolchain 1.25.4) + gRPC, protobuf, buf v1.32.1
