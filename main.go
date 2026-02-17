package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	maxInputBytes = 100 * 1024 // 100KB
	version       = "0.1.0"
)

func main() {
	// Parse args
	var instruction string
	var model string
	var raw bool
	var showHelp bool
	var showVersion bool

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			showHelp = true
		case "--version", "-v":
			showVersion = true
		case "--raw":
			raw = true
		case "--model":
			if i+1 < len(args) {
				i++
				model = args[i]
			} else {
				fatal("--model requires a value")
			}
		default:
			if strings.HasPrefix(args[i], "-") {
				fatal("unknown flag: " + args[i])
			}
			if instruction == "" {
				instruction = args[i]
			} else {
				instruction = instruction + " " + args[i]
			}
		}
	}

	if showVersion {
		fmt.Println("jsone " + version)
		os.Exit(0)
	}

	if showHelp {
		printUsage()
		os.Exit(0)
	}

	// Check for piped input
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "jsone: no input. Pipe something in: command | jsone")
		os.Exit(1)
	}

	// Read stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fatal("reading stdin: " + err.Error())
	}

	if len(input) == 0 {
		fatal("empty input")
	}

	truncated := false
	if len(input) > maxInputBytes {
		input = input[:maxInputBytes]
		truncated = true
		fmt.Fprintf(os.Stderr, "jsone: input truncated to %dKB\n", maxInputBytes/1024)
	}

	// Get API key (Gemini direct or OpenRouter)
	apiKey := getAPIKey()
	orKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" && orKey == "" {
		fatal("no API key. Set JSONE_API_KEY, GEMINI_API_KEY, or OPENROUTER_API_KEY")
	}

	if model == "" {
		model = "gemini-2.0-flash"
	}

	// Build prompt and call LLM
	result, err := callLLM(apiKey, model, string(input), instruction, truncated)
	if err != nil {
		fatal(err.Error())
	}

	// Output
	if raw {
		fmt.Print(result)
	} else {
		formatted, err := prettyJSON(result)
		if err != nil {
			// If pretty-print fails, output raw
			fmt.Print(result)
		} else {
			fmt.Print(formatted)
		}
	}
	fmt.Println()
}

func getAPIKey() string {
	if k := os.Getenv("JSONE_API_KEY"); k != "" {
		return k
	}
	if k := os.Getenv("GEMINI_API_KEY"); k != "" {
		return k
	}
	return ""
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "jsone: "+msg)
	os.Exit(1)
}

func printUsage() {
	fmt.Println(`jsone - pipe anything, get JSON

Usage: command | jsone [instruction] [flags]

Arguments:
  instruction    Natural language instruction for guided extraction
                 (positional, no flag needed)

Flags:
  --model MODEL  Override LLM model (default: gemini-2.0-flash)
  --raw          Compact JSON output (no pretty-printing)
  --version      Show version
  --help         Show this help

Examples:
  cat /etc/hosts | jsone
  docker ps | jsone
  cat access.log | jsone "group by status code"
  grep -r TODO . | jsone "file, line, text"

Environment:
  JSONE_API_KEY   API key for Gemini (preferred)
  GEMINI_API_KEY  Fallback API key`)
}
