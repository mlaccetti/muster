#!/usr/bin/env bash
set -euo pipefail
if rg -n 'time\.Sleep|Sleep\(' --glob '*.go' --glob '!*_test.go' .; then
  echo "check-no-sleep: forbidden sleep in non-test Go code" >&2
  exit 1
fi
