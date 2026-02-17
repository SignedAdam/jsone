#!/bin/bash
# jsone demo -- run these to see it in action
# Requires: GEMINI_API_KEY or OPENROUTER_API_KEY set

set -e

echo "=== Example 1: Hosts file ==="
echo "127.0.0.1 localhost
192.168.1.1 router
10.0.0.5 database" | jsone

echo ""
echo "=== Example 2: Table data ==="
echo "NAME          STATUS    AGE
nginx         Running   5d
postgres      Running   12d
redis         CrashLoop 1h" | jsone

echo ""
echo "=== Example 3: Guided extraction ==="
echo "ERROR 2024-01-15 auth failed user=admin ip=1.2.3.4
INFO 2024-01-15 login success user=bob ip=5.6.7.8
ERROR 2024-01-15 auth failed user=root ip=9.10.11.12" | jsone "group by level, count each"

echo ""
echo "=== Example 4: Grep-style output ==="
echo "./main.go:15:// TODO: add validation
./llm.go:42:// TODO: retry logic
./format.go:8:// TODO: support yaml" | jsone "file, line, text"

echo ""
echo "=== Example 5: Key-value pairs ==="
echo 'NAME="Ubuntu"
VERSION="22.04.3 LTS"
ID=ubuntu
PRETTY_NAME="Ubuntu 22.04.3 LTS"' | jsone

echo ""
echo "Done. All examples completed successfully."
