package main

import (
	"fmt"
	"strings"
	"time"
)

type demoCase struct {
	name   string
	cmd    string
	input  string
	output string
}

var demoCases = []demoCase{
	{
		name: "hosts file",
		cmd:  "cat /etc/hosts | jsone",
		input: `127.0.0.1	localhost
127.0.1.1	devbox
192.168.1.10	fileserver
192.168.1.20	nas.local
10.0.0.1	gateway`,
		output: `[
  {"ip": "127.0.0.1", "hostname": "localhost"},
  {"ip": "127.0.1.1", "hostname": "devbox"},
  {"ip": "192.168.1.10", "hostname": "fileserver"},
  {"ip": "192.168.1.20", "hostname": "nas.local"},
  {"ip": "10.0.0.1", "hostname": "gateway"}
]`,
	},
	{
		name: "docker ps",
		cmd:  "docker ps | jsone",
		input: `CONTAINER ID   IMAGE          COMMAND                  CREATED        STATUS        PORTS                    NAMES
a1b2c3d4e5f6   nginx:latest   "/docker-entrypoint.…"   2 hours ago    Up 2 hours    0.0.0.0:80->80/tcp       web
f6e5d4c3b2a1   postgres:16    "docker-entrypoint.s…"   3 days ago     Up 3 days     0.0.0.0:5432->5432/tcp   db
1a2b3c4d5e6f   redis:7        "docker-entrypoint.s…"   3 days ago     Up 3 days     0.0.0.0:6379->6379/tcp   cache`,
		output: `[
  {
    "container_id": "a1b2c3d4e5f6",
    "image": "nginx:latest",
    "command": "/docker-entrypoint.…",
    "created": "2 hours ago",
    "status": "Up 2 hours",
    "ports": "0.0.0.0:80->80/tcp",
    "names": "web"
  },
  {
    "container_id": "f6e5d4c3b2a1",
    "image": "postgres:16",
    "command": "docker-entrypoint.s…",
    "created": "3 days ago",
    "status": "Up 3 days",
    "ports": "0.0.0.0:5432->5432/tcp",
    "names": "db"
  },
  {
    "container_id": "1a2b3c4d5e6f",
    "image": "redis:7",
    "command": "docker-entrypoint.s…",
    "created": "3 days ago",
    "status": "Up 3 days",
    "ports": "0.0.0.0:6379->6379/tcp",
    "names": "cache"
  }
]`,
	},
	{
		name: "log grouping",
		cmd:  `cat access.log | jsone "group by status code"`,
		input: `192.168.1.1 - - [17/Feb/2026:10:00:01] "GET /api/users HTTP/1.1" 200 1234
192.168.1.2 - - [17/Feb/2026:10:00:02] "POST /api/login HTTP/1.1" 401 89
192.168.1.1 - - [17/Feb/2026:10:00:03] "GET /api/posts HTTP/1.1" 200 5678
10.0.0.5 - - [17/Feb/2026:10:00:04] "GET /api/admin HTTP/1.1" 403 42
192.168.1.3 - - [17/Feb/2026:10:00:05] "GET /missing HTTP/1.1" 404 0
192.168.1.1 - - [17/Feb/2026:10:00:06] "GET /api/health HTTP/1.1" 200 15
10.0.0.5 - - [17/Feb/2026:10:00:07] "POST /api/upload HTTP/1.1" 500 0`,
		output: `{
  "200": 3,
  "401": 1,
  "403": 1,
  "404": 1,
  "500": 1
}`,
	},
	{
		name: "grep extraction",
		cmd:  `grep -rn TODO . | jsone "file, line, text"`,
		input: `./main.go:42:// TODO: add retry logic
./main.go:87:// TODO: handle timeout
./server.go:12:// TODO: add rate limiting
./utils.go:5:// TODO: write tests for this`,
		output: `[
  {"file": "main.go", "line": 42, "text": "add retry logic"},
  {"file": "main.go", "line": 87, "text": "handle timeout"},
  {"file": "server.go", "line": 12, "text": "add rate limiting"},
  {"file": "utils.go", "line": 5, "text": "write tests for this"}
]`,
	},
}

func runDemo() {
	fmt.Println("\033[1mjsone demo\033[0m -- see what jsone does, no API key needed\n")

	for i, dc := range demoCases {
		// Show command
		fmt.Printf("\033[36m$ %s\033[0m\n", dc.cmd)
		fmt.Println()

		// Show input (dimmed)
		fmt.Println("\033[2m# Input:\033[0m")
		for _, line := range strings.Split(dc.input, "\n") {
			fmt.Printf("\033[2m%s\033[0m\n", line)
		}
		fmt.Println()

		// Simulate thinking
		fmt.Print("\033[33m⠋ Processing...\033[0m")
		time.Sleep(400 * time.Millisecond)
		fmt.Print("\r\033[33m⠙ Processing...\033[0m")
		time.Sleep(300 * time.Millisecond)
		fmt.Print("\r\033[33m⠹ Processing...\033[0m")
		time.Sleep(300 * time.Millisecond)
		fmt.Print("\r\033[K") // Clear line

		// Show output (green)
		fmt.Println("\033[32m# Output:\033[0m")
		fmt.Println(dc.output)

		if i < len(demoCases)-1 {
			fmt.Println("\n" + strings.Repeat("─", 50) + "\n")
		}
	}

	fmt.Println()
	fmt.Println("\033[1mReady to try it yourself?\033[0m")
	fmt.Println("Get a free Gemini API key: https://ai.google.dev/aistudio")
	fmt.Println("Then: export GEMINI_API_KEY=your-key")
	fmt.Println()
	fmt.Println("Or install: brew install SignedAdam/tap/jsone")
	fmt.Println("            go install github.com/SignedAdam/jsone@latest")
}
