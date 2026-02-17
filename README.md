<p align="center">
  <h1 align="center">jsone</h1>
  <p align="center"><strong>Pipe anything, get JSON.</strong></p>
  <p align="center">
    <a href="#install">Install</a> · 
    <a href="#examples">Examples</a> · 
    <a href="#how-it-works">How it works</a> · 
    <a href="#for-ai-agents">For AI agents</a>
  </p>
</p>

---

`jsone` is a CLI tool that reads stdin and outputs structured JSON. It uses a fast LLM (Gemini Flash) to infer structure from any text input -- tables, logs, config files, grep output, whatever.

```bash
cat /etc/hosts | jsone
```
```json
[
  {"ip": "127.0.0.1", "hostname": "localhost"},
  {"ip": "192.168.1.1", "hostname": "router"}
]
```

Zero config to start. Sub-second for small inputs. Built for shell pipelines.

## Install

### Go install (recommended)

```bash
go install github.com/SignedAdam/jsone@latest
```

### From source

```bash
git clone https://github.com/SignedAdam/jsone.git
cd jsone
go build -o jsone .
mv jsone /usr/local/bin/
```

### Set your API key

```bash
# Option 1: Gemini API (fastest, free tier available)
export GEMINI_API_KEY="your-key-here"

# Option 2: Dedicated jsone key (same as Gemini, just namespaced)
export JSONE_API_KEY="your-key-here"

# Option 3: OpenRouter (supports many models, slightly slower)
export OPENROUTER_API_KEY="your-key-here"
```

Get a free Gemini API key in 30 seconds at [ai.google.dev](https://ai.google.dev).

## Examples

### Auto-detect structure (no args)

```bash
# Hosts file
cat /etc/hosts | jsone
```
```json
[{"ip": "127.0.0.1", "hostname": "localhost"}, ...]
```

```bash
# Any table output
docker ps | jsone
```
```json
[{"container_id": "abc123", "image": "nginx", "status": "Up 2 hours", "ports": "80/tcp"}, ...]
```

```bash
# Key-value configs
cat /etc/os-release | jsone
```
```json
{"name": "Ubuntu", "version": "22.04", "id": "ubuntu", ...}
```

### Guided extraction (positional arg)

```bash
# Group and count
cat access.log | jsone "group by status code"
```
```json
{"200": 1523, "404": 47, "500": 3}
```

```bash
# Extract specific fields
grep -r TODO . | jsone "file, line, text"
```
```json
[
  {"file": "./main.go", "line": 42, "text": "TODO: add validation"},
  {"file": "./llm.go", "line": 15, "text": "TODO: retry logic"}
]
```

```bash
# Natural language instructions
ps aux | jsone "top 5 by memory usage, include pid and command"
```
```json
[
  {"pid": 1234, "memory_percent": 12.3, "command": "chrome"},
  {"pid": 5678, "memory_percent": 8.1, "command": "node"}
]
```

### Pipeline composition

```bash
# Chain with jq
docker ps | jsone | jq '.[] | select(.status | contains("Up"))'

# Compact output for piping
cat data.csv | jsone --raw | jq '.[] | select(.active == true)'

# Feed into other tools
kubectl get pods | jsone | jq -r '.[].name' | xargs kubectl describe pod
```

## Usage

```
command | jsone [instruction] [flags]
```

| Flag | Description |
|------|-------------|
| `--model MODEL` | Override LLM model (default: `gemini-2.0-flash`) |
| `--raw` | Compact JSON, no pretty-printing |
| `--version` | Print version |
| `--help` | Show help |

### API key resolution order

1. `JSONE_API_KEY` -- Gemini API (preferred)
2. `GEMINI_API_KEY` -- Gemini API (fallback)
3. `OPENROUTER_API_KEY` -- OpenRouter (any model)

## How it works

1. Reads all of stdin into a buffer (up to 100KB; truncates with a warning)
2. Sends to Gemini Flash with `response_mime_type: application/json` (native JSON mode -- guaranteed valid output)
3. Pretty-prints to stdout

No prompt engineering for JSON validity. Gemini's structured output mode handles that natively.

## For AI agents

`jsone` is designed to be used by AI agents in shell pipelines. Key properties:

- **Deterministic interface.** stdin in, JSON stdout, errors on stderr. Exit code 0 on success, 1 on failure.
- **Structured errors.** All errors go to stderr with a `jsone:` prefix. stdout is always either valid JSON or empty.
- **Composable.** Chain with `jq`, `fx`, or any JSON tool. Use `--raw` for compact output in pipes.
- **No interactive prompts.** Fully non-interactive. Safe for automation.
- **Idempotent.** Same input + same instruction = same structure (content may vary slightly due to LLM, but schema is stable).
- **Fast enough for scripting.** Sub-second with Gemini Flash for inputs under 10KB. Not suitable for tight loops (each call is an API request).

### Agent usage patterns

```bash
# Parse command output into structured data for decision-making
docker ps | jsone --raw  # feed into your context

# Extract actionable items from logs
journalctl --since "1 hour ago" | jsone "errors only, with timestamp and message" --raw

# Convert human-readable output to machine-parseable
git log --oneline -10 | jsone "hash, message" --raw
```

### Cost awareness

Each invocation makes one LLM API call. Gemini Flash pricing is very low (~$0.0001 per typical call), but avoid using jsone in tight loops over thousands of items. Batch your input instead:

```bash
# Good: one call
find . -name "*.go" -exec grep -l TODO {} \; | jsone

# Bad: many calls
find . -name "*.go" | while read f; do grep TODO "$f" | jsone; done
```

## Supported backends

| Backend | Env var | JSON mode | Speed |
|---------|---------|-----------|-------|
| Gemini (direct) | `GEMINI_API_KEY` | Native (`response_mime_type`) | Fastest |
| OpenRouter | `OPENROUTER_API_KEY` | `json_object` format | +~1s latency |

## Roadmap

- [ ] `--schema FILE` -- enforce output against a JSON Schema
- [ ] `~/.config/jsone/config.json` -- persistent config
- [ ] Shell completions (bash, zsh, fish)
- [ ] `-o yaml` / `-o csv` output formats
- [ ] Streaming for large inputs

## License

MIT
