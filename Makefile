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

clean:
	rm -rf bin progress.db