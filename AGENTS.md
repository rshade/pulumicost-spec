# Repository Guidelines

## Project Structure & Modules

- Core protocol in `proto/pulumicost/v1/` with generated Go bindings in `sdk/go/proto/`.
- Domain helpers, validation, and test utilities live in `sdk/go/types/` and `sdk/go/testing/`.
- JSON Schemas in `schemas/` with validation helpers in `validate_examples.js`.
- Cross-cloud examples in `examples/specs/` and sample payloads in `examples/requests/`.
- Docs and prompts reside in `docs/` and top-level `*_GUIDE.md` files.

## Build, Test, and Development

- `make generate` installs a local `buf` binary and regenerates Go protobuf code into `sdk/go/proto/`.
- `make test` or `go test ./...` runs the Go suite (sdk, validation helpers).
- `make validate` runs Go tests, GolangCI-Lint, buf lint, markdownlint, yamllint, and npm schema checks.
- Targeted tasks: `npm run validate:schema`, `npm run validate:examples`, `npm run lint:markdown`, `npm run lint:yaml`,
  `buf lint`.
- Cleanups: `make tidy` (Go deps) and `make clean`/`make clean-all` for generated assets and the local buf binary.

## Coding Style & Naming

- Go 1.24.7 (toolchain 1.25.1); format via `goimports`/`golines` (enforced by GolangCI-Lint). Run from repo root to pick
  up `.golangci.yml`.
- Prefer small, composable functions; keep public API names aligned with protobuf/schema terminology (CostSource,
  PricingSpec, etc.).
- Proto files follow snake_case for fields; generated Go uses standard CamelCase. Add `String()` helpers for new enums to
  match existing patterns.
- Markdown follows `.markdownlint.json`; YAML follows `.yamllint`.

## Testing Guidelines

- Add Go tests alongside packages (e.g., `sdk/go/types/..._test.go`). Use table-driven cases and cover error paths.
- For schema changes, add/adjust fixtures in `examples/` and re-run `npm run validate:examples`.
- For protobuf changes, regenerate with `make generate`, then ensure `buf lint` and `go test ./...` pass.
- Capture expected failure messages in tests when tightening validation to avoid regressions.

## Commit & Pull Request Guidelines

- Use conventional-feeling prefixes seen in history (`feat:`, `fix:`, `chore:`) with concise scope descriptions.
- One logical change per commit; include tests or explain why not.
- PRs should describe intent, enumerate functional changes, list validation commands run, and link issues when applicable.
  Include screenshots only if UI/docs rendering changes.
- Expect CI to run `make validate`; ensure local parity before opening a PR.

## Security & Configuration Tips

- Never commit secrets or real cloud identifiers; sanitize example data.
- Generated code and binaries (`sdk/go/proto/`, `bin/buf`) should not be hand-edited. Update sources (`proto/`,
  `schemas/`) and regenerate instead.
- Keep `go.mod`/`go.sum` tidy and avoid replacing upstream modules unless documented.
