# systems-forge developer/learner commands.
# All verbs forward to ./bin/forge; the binary is (re)built on demand.
ROOT  := $(abspath .)
GO    ?= go
GOFLAGS ?=

.PHONY: setup selfcheck init list show check score selftest tui clean

setup:
	cd forge && $(GO) build $(GOFLAGS) -o ../bin/forge ./cmd/forge

selfcheck: setup
	./tools/selfcheck.sh

# Learner verbs: `make check M11-ex03`, `make show M0-ex01`, …
init list show check score selftest tui: setup
	./bin/forge $@ $(filter-out setup $@, $(MAKECMDGOALS))

# Desktop GUI (Wails v2 + Svelte + CodeMirror). The GUI is additive: it
# reuses the same engine/store as the CLI/TUI and writes to progress.db.
# Requires the wails CLI and a WebKitGTK 4.1 host (see docs/gui-spec.md).
WAILS ?= $(shell command -v wails 2>/dev/null || echo "$(HOME)/go/bin/wails")

.PHONY: gui gui-dev
gui:
	mkdir -p bin
	cd forge/gui && $(WAILS) build && cp build/bin/forge-gui ../../bin/forge-gui

gui-dev:
	cd forge/gui && $(WAILS) dev

clean:
	rm -rf bin progress.db