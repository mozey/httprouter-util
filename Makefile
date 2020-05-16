# Copied from
# https://gist.github.com/lantins/e83477d8bccab83f078d

# binary name to kill/restart
EXE = ${APP_EXE}

fswatch:
	@command -v fswatch --version >/dev/null 2>&1 || \
	{ printf >&2 "Install fswatch, run: brew install fswatch\n"; exit 1; }

# app...........................................................................
#app.clean:
#	@echo app.clean
#	go clean

# Build dev server
app.build.dev:
	@echo app.build.dev
	/usr/bin/env bash -c "scripts/build.dev.sh"

#app.build: app.clean
#	@echo app.build
#	/usr/bin/env bash -c "scripts/build.dev.sh"

# Attempt to kill running server
app.kill:
	@echo app.kill
	-@killall -9 $(EXE) 2>/dev/null || true

# Just run the server, no live reload
app.run:
	@echo app.run
	@make app.build.dev
	/usr/bin/env bash -c ./${EXE}

# Restart server, for use with fswatch
# TODO Run `make app.restart` then ctrl+x and do `ps ax | grep .out`
# Why is app.out still running?
app.restart:
	@echo app.restart
	@make app.kill
	@make app.build.dev; (if [ "$$?" -eq 0 ]; then (./${EXE} &); fi)

# Run app server with live reload
# Watch .go files for changes then recompile & try to start server
# will also kill server after ctrl+c
# fswatch includes everything unless an exclusion filter says otherwise
# https://stackoverflow.com/a/37237681/639133
app: fswatch
	@make app.restart
	@fswatch -or --exclude ".*" --include "\\.go$$" ./ | \
	xargs -n1 -I{} make app.restart || make app.kill




