# jsone

Pipe anything, get JSON.

```bash
cat /etc/hosts | jsone
docker ps | jsone
cat access.log | jsone "group by status code"
grep -r TODO . | jsone "file, line, text"
```

## What it does

`jsone` reads stdin, sends it to a fast LLM (Gemini Flash by default), and outputs valid JSON to stdout. That's it.

- No args = auto-detect the most obvious structure
- Positional arg = guided extraction
- Output is always valid JSON, nothing else

## Install

```bash
go install github.com/SignedAdam/jsone@latest
```

Or download a binary from releases.

## Usage

```
command | jsone [instruction] [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `--model MODEL` | Override LLM model (default: gemini-2.0-flash) |
| `--raw` | Compact JSON output (no pretty-printing) |
| `--version` | Show version |
| `--help` | Show help |

### API Key

Set one of these environment variables:

| Variable | Backend |
|----------|---------|
| `JSONE_API_KEY` | Gemini API (preferred, fastest) |
| `GEMINI_API_KEY` | Gemini API (fallback) |
| `OPENROUTER_API_KEY` | OpenRouter (supports many models) |

Get a free Gemini API key at [ai.google.dev](https://ai.google.dev).

## Examples

```bash
# Hosts file to JSON
cat /etc/hosts | jsone
# [{"ip": "127.0.0.1", "hostname": "localhost"}, ...]

# Table output to array of objects
docker ps | jsone
# [{"container_id": "abc123", "image": "nginx", ...}, ...]

# Log analysis
cat access.log | jsone "group by status code"
# {"200": 1523, "404": 47, "500": 3}

# Grep results
grep -r TODO . | jsone "file, line, text"
# [{"file": "./main.go", "line": 42, "text": "TODO: fix this"}, ...]

# Compact output for piping
cat data.csv | jsone --raw | jq '.[] | select(.status == "active")'
```

## How it works

1. Reads all of stdin (up to 100KB, truncates with warning)
2. Sends to LLM with JSON-mode enabled (guaranteed valid JSON output)
3. Pretty-prints to stdout

Sub-second for small inputs. Built for shell pipelines.

## License

MIT
