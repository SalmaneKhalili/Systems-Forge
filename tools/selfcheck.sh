#!/usr/bin/env bash
# selfcheck — prove the platform builds clean and every grader fixture scores
# as expected. Fails (exit 1) on any slip so it can gate CI / release.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

# Locate a Go toolchain: PATH first, then the user-local install we provisioned.
if ! command -v go >/dev/null 2>&1; then
	if [ -x "$HOME/.local/go/bin/go" ]; then
		export PATH="$HOME/.local/go/bin:$PATH"
	else
		echo "selfcheck: no go toolchain found. Install go 1.27+ or add it to PATH." >&2
		exit 1
	fi
fi

echo "==> forge version"
go version

echo "==> fmt"
(cd forge && gofmt -l . | tee /dev/stderr | grep -q .) && {
	echo "selfcheck: gofmt found unformatted files" >&2
	exit 1
} || true

echo "==> vet"
(cd forge && go vet ./...)

echo "==> build bin/forge"
rm -rf bin
(cd forge && go build -o ../bin/forge ./cmd/forge)

echo "==> selftest (grader fixtures)"
./bin/forge selftest

echo "==> readings lint (Reading-Ladder standard)"
./bin/forge lint readings

echo "==> cli smoke (list/score)"
./bin/forge list >/dev/null
./bin/forge score >/dev/null

echo
echo "selfcheck: OK"