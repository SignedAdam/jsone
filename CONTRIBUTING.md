# Contributing to jsone

## Development

```bash
git clone https://github.com/SignedAdam/jsone.git
cd jsone
go build -o jsone .
```

## Architecture

jsone is intentionally simple. Four files:

| File | Purpose |
|------|---------|
| `main.go` | Entry point, arg parsing, stdin reading |
| `llm.go` | Gemini and OpenRouter API clients |
| `prompt.go` | System prompt and user prompt construction |
| `format.go` | JSON pretty-printing |

## Adding a new backend

1. Add a new `callX()` function in `llm.go` following the pattern of `callGemini()` or `callOpenRouter()`
2. Add the env var check in `callLLM()`
3. Document the env var in README.md

## Design principles

- **One file per concern.** Don't merge unrelated logic.
- **No dependencies.** Standard library only. No third-party packages.
- **Errors to stderr.** stdout is sacred -- only valid JSON or nothing.
- **Sub-second is the target.** Don't add anything that blocks the critical path.
