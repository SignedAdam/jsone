<p align="center">
  <h1 align="center">jsone</h1>
  <p align="center"><strong>Pipe anything in, get structured JSON out.</strong></p>
  <p align="center">A CLI that uses an LLM to turn any text into structured JSON. Like jq, but for data that doesn't have a schema yet.</p>
  <p align="center">
    <a href="#try-it-now">Try it</a> · 
    <a href="#install">Install</a> · 
    <a href="#examples">Examples</a> · 
    <a href="#how-it-works">How it works</a> · 
    <a href="#for-ai-agents">For AI agents</a>
  </p>
</p>

---

```bash
cat /etc/hosts | jsone
```
```json
[
  {"ip": "127.0.0.1", "hostname": "localhost"},
  {"ip": "192.168.1.1", "hostname": "router"}
]
```

No regex. No parsing code. No schema definition. Just pipe and go.

## Try it now

No API key needed:

```bash
jsone --demo
```

This runs 4 interactive examples showing real transformations -- hosts files, docker output, log grouping, grep extraction -- with zero setup.

Ready to use it for real? Get a free Gemini API key in 30 seconds at [ai.google.dev/aistudio](https://ai.google.dev/aistudio), then:

```bash
export GEMINI_API_KEY="your-key"
```

## Install

### Homebrew (macOS/Linux)

```bash
brew install SignedAdam/tap/jsone
```

### Go install

```bash
go install github.com/SignedAdam/jsone@latest
```

### Binary download

Grab a pre-built binary from [Releases](https://github.com/SignedAdam/jsone/releases) for your platform (macOS, Linux, Windows -- amd64/arm64).

### From source

```bash
git clone https://github.com/SignedAdam/jsone.git
cd jsone
go build -o jsone .
```

## Examples

### Auto-detect structure (no args)

jsone looks at your input and figures out the obvious structure:

```bash
cat /etc/hosts | jsone
```
```json
[{"ip": "127.0.0.1", "hostname": "localhost"}, ...]
```

```bash
docker ps | jsone
```
```json
[{"container_id": "abc123", "image": "nginx", "status": "Up 2 hours", "ports": "80/tcp"}, ...]
```

```bash
cat /etc/os-release | jsone
```
```json
{"name": "Ubuntu", "version": "22.04", "id": "ubuntu", ...}
```

### Guided extraction (just tell it what you want)

```bash
cat access.log | jsone "group by status code"
```
```json
{"200": 1523, "404": 47, "500": 3}
```

```bash
grep -rn TODO . | jsone "file, line, text"
```
```json
[
  {"file": "main.go", "line": 42, "text": "add validation"},
  {"file": "llm.go", "line": 15, "text": "retry logic"}
]
```

```bash
ps aux | jsone "top 5 by memory usage, include pid and command"
```
```json
[
  {"pid": 1234, "memory_percent": 12.3, "command": "chrome"},
  {"pid": 5678, "memory_percent": 8.1, "command": "node"}
]
```

### Works great with jq

jsone turns unstructured data into JSON. jq processes structured JSON. They're complementary:

```bash
docker ps | jsone | jq '.[] | select(.status | contains("Up"))'
kubectl get pods | jsone | jq -r '.[].name' | xargs kubectl describe pod
cat data.csv | jsone --raw | jq '.[] | select(.active == true)'
```

## How is this different from jq?

| | jq | jsone |
|---|---|---|
| **Input** | Must be valid JSON | Anything (logs, tables, configs, prose) |
| **Schema** | You define the structure | LLM infers the structure |
| **Offline** | Yes | No (needs API) |
| **Speed** | Instant | Sub-second (API call) |
| **Use case** | Transform JSON | Create JSON from non-JSON |

They solve different problems. Use jsone to get your data into JSON, then jq to work with it.

## Usage

```
command | jsone [instruction] [flags]
```

| Flag | Description |
|------|-------------|
| `--demo` | Run interactive demo, no API key needed |
| `--model MODEL` | Override LLM model (default: `gemini-2.0-flash`) |
| `--raw` | Compact JSON, no pretty-printing |
| `--version` | Print version |
| `--help` | Show help |

### API key resolution

1. `JSONE_API_KEY` -- Gemini API (preferred)
2. `GEMINI_API_KEY` -- Gemini API (fallback)
3. `OPENROUTER_API_KEY` -- OpenRouter (any model, slightly slower)

## How it works

1. Reads stdin (up to 100KB, truncates with warning)
2. Sends to Gemini Flash with `response_mime_type: application/json` (native JSON mode -- guaranteed valid output)
3. Pretty-prints to stdout

No prompt engineering for JSON validity. Gemini's structured output mode handles that natively. One API call per invocation.

## For AI agents

jsone is designed to work in automated pipelines. Key properties:

- **Deterministic interface.** stdin in, JSON stdout, errors on stderr. Exit 0 on success, 1 on failure.
- **Structured errors.** All errors go to stderr with a `jsone:` prefix. stdout is always valid JSON or empty.
- **Composable.** Chain with jq, fx, or any JSON tool. Use `--raw` for compact output.
- **Non-interactive.** No prompts, no confirmations. Safe for automation.
- **Stable schema.** Same input + same instruction = same JSON structure across calls.

### Cost

Each call is one Gemini Flash API request (~$0.0001). Don't use in tight loops -- batch your input instead:

```bash
# Good: one call
find . -name "*.go" -exec grep -l TODO {} \; | jsone

# Bad: N calls
find . -name "*.go" | while read f; do grep TODO "$f" | jsone; done
```

## Roadmap

- [ ] `--schema FILE` -- enforce output against a JSON Schema
- [ ] `~/.config/jsone/config.json` -- persistent API key config
- [ ] Shell completions (bash, zsh, fish)
- [ ] Streaming for large inputs

## License

[MIT](LICENSE)
