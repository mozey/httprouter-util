# Copied from
# https://gist.github.com/lantins/e83477d8bccab83f078d

# binary name to kill/restart
PROG_DEV = dev.out

dependencies:
	@command -v fswatch --version >/dev/null 2>&1 || \
	{ printf >&2 "Install fswatch, run: brew install fswatch\n"; exit 1; }

# default targets to run when only running `make`
default: dependencies

# dev...........................................................................
# Local server with live reload
clean.dev:
	go clean

build.dev: dependencies clean.dev
	go build -o ./dev.out ./

# attempt to kill running server
kill.dev:
	@echo kill.dev
	-@killall -9 $(PROG_DEV) 2>/dev/null || true

# attempt to build and start server
restart.dev:
	@echo restart.dev
	@make kill.dev
	@make build.dev; (if [ "$$?" -eq 0 ]; then (./${PROG_DEV} &); fi)

# watch .go files for changes then recompile & try to start server
# will also kill server after ctrl+c
# fswatch includes everything unless an exclusion filter says otherwise
# https://stackoverflow.com/a/37237681/639133
dev: dependencies
	@make restart.dev
	@fswatch -or --exclude ".*" --include "\\.go$$" ./ | \
	xargs -n1 -I{} make restart.dev || make kill.dev






