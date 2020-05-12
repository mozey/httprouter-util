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
dev.clean:
	@echo dev.clean
	go clean

# Build dev server
dev.build: dependencies dev.clean
	@echo dev.build
	go build -o ./dev.out ./

# Attempt to kill running server
dev.kill:
	@echo dev.kill
	-@killall -9 $(PROG_DEV) 2>/dev/null || true

# Restart server
dev.restart:
	@echo dev.restart
	@make dev.kill
	@make dev.build; (if [ "$$?" -eq 0 ]; then (./${PROG_DEV} &); fi)

# Run dev server (no live reload)
run:
	@echo dev.run
	@make dev.build
	./${PROG_DEV}

# Run dev server with live reload
# Watch .go files for changes then recompile & try to start server
# will also kill server after ctrl+c
# fswatch includes everything unless an exclusion filter says otherwise
# https://stackoverflow.com/a/37237681/639133
reload: dependencies
	@make dev.restart
	@fswatch -or --exclude ".*" --include "\\.go$$" ./ | \
	xargs -n1 -I{} make dev.restart || make dev.kill




