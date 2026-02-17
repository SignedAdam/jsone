# Changelog

## [0.2.0] - 2026-02-17

### Added
- `--demo` flag: interactive demo with 4 examples, no API key needed
- MIT license
- This changelog

## [0.1.0] - 2026-02-17

### Added
- Initial release
- Pipe any text through Gemini Flash LLM, get structured JSON back
- Gemini native JSON mode (`response_mime_type: application/json`)
- OpenRouter fallback support
- Positional instruction argument for guided extraction
- `--model` flag for model override
- `--raw` flag for compact output
- `--version` flag
- 100KB input truncation with warning
- Auto-retry on JSON validation failure
