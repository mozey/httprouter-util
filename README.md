# httprouter-example

Example of using [httprouter](https://github.com/julienschmidt/httprouter)
with 
- [zerolog](https://github.com/rs/zerolog) for logging
- [Middleware](https://github.com/gorilla/handlers): panic handler, request logging, request ID for tracing, token auth, max bytes handler, gzip
- Graceful shutdown on ctrl+c
- [Swagger](https://github.com/swaggo/swag) docs
- [Caddy](https://caddyserver.com/) as a HTTPS endpoint, API gateway, and reverse proxy


## Quick start

Clone the repo (outside your GOPATH since this is a module)
```bash
git clone https://github.com/mozey/httprouter-example.git
cd httprouter-example # This is the APP_DIR
```

Following the 12 factor app recommendation to [store config in the environment](https://12factor.net/config). Configuration is done using [environment variables](https://en.wikipedia.org/wiki/Environment_variable)

Generate script to set dev config
```bash
APP_DIR=$(pwd) ./make.sh env_sh dev
```

Set dev config
```bash
source ./dev.sh
```

Run dev server (no live reload)
```bash
./make.sh app_run
```

Run dev server with live reload
```bash
./make.sh app
```


## Examples

Make requests from the cli with [curlie](https://github.com/rs/curlie)
  
### Authentication

Token is required by default
[http://localhost:8118/token/is/required/by/default](http://localhost:8118/token/is/required/by/default)

Some routes may [skip the token check](https://github.com/mozey/httprouter-example/blob/connect-go/middleware.go#L119)
- [http://localhost:8118](http://localhost:8118)
- [http://localhost:8118/index.html](http://localhost:8118/index.html)
- [http://localhost:8118/www/data/go.txt](http://localhost:8118/www/data/go.txt)

### Using the token

For static files
[http://localhost:8118/hello/foo?token=123](http://localhost:8118/hello/foo?token=123)

And API endpoints
[http://localhost:8118/api?token=123](http://localhost:8118/api?token=123)

### Error handling

[http://localhost:8118/panic](http://localhost:8118/panic)
    
[http://localhost:8118/does/not/exist?token=123](http://localhost:8118/does/not/exist?token=123)

### Configuration

Use http.MaxBytesReader to limit POST body. Make the [request with specified body size](https://serverfault.com/a/283297), Assuming `MaxBytes` is set to 1 KiB the request below will fail
```bash
dd if=/dev/urandom bs=1 count=1025 | curlie --data-binary @- POST "http://localhost:8118/api?token=123"
```

Settings to protect against malicious clients. **NOTE** The response body for errors below is not JSON, it's not possible to override string response hard-coded in Golang SDK
```bash
# ReadTimeout
gotest -v ./... -run TestReadTimeout

# WriteTimeout
gotest -v ./... -run TestWriteTimeout

# MaxHeaderBytes
gotest -v ./... -run TestMaxHeaderBytes
```

### Proxy

**TODO** Proxy request to external service?
[http://localhost:8118/proxy](http://localhost:8118/proxy)

However, a better architecture is to use something like [Caddy](https://github.com/caddyserver/caddy) as a HTTPS endpoint, API gateway, and reverse proxy. The Caddyfile for configuring this is quite simple, see [#6](https://github.com/mozey/httprouter-example/issues/6)

### Services

**TODO** Define services on the handler, e.g. DB connection
[http://localhost:8118/db?sql=select * from color](http://localhost:8118/db?sql=select%20*%20from%20color)


## Client

Example client with self-update feature.

**TODO** Embed client code generated from protobuf schema?

### Build client

Build the client, download it, and print version
```bash
source dev.sh

APP_CLIENT_VERSION=0.1.0 ./scripts/build-client.sh

./dist/client -version

curlie "http://localhost:8118/client/download?token=123" -o client

chmod u+x client

./client -version
```

### Update client

Create a new build
```bash
APP_CLIENT_VERSION=0.2.0 ./scripts/build-client.sh

./dist/client -version

curlie "http://localhost:8118/client/version?token=123"
```

Update from the server and print new version
```bash
./client -update

./client -version
```

Running update again prints *"already on the latest version"*
```
./client -update
```


## Reset

Removes all user config
```bash
APP_DIR=$(pwd) ./scripts/reset.sh
```


## Dependencies

This example aims for a good cross platform experience by depending on 
- [Golang](https://golang.org/) 
- [Bash](https://www.gnu.org/software/bash)
- [fswatch](https://github.com/emcrisostomo/fswatch)
- xargs

On macOS and Linux
- [pgrep](https://en.wikipedia.org/wiki/Pgrep) and kill

**TODO** On Windows
- [taskkill](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/taskkill)

[GNU Make](https://stackoverflow.com/questions/3798562/why-use-make-over-a-shell-script) 
is not needed because Golang is fast to build,
and `fswatch` can be used for live reload.
For this example `main.go` is kept in the project root.
Larger projects might have separate bins in the *"/cmd"* dir

Bash on Windows is easy to setup using 
[msys2](https://www.msys2.org/), MinGW, or native shell on Windows 10.
For other UNIX programs see [gow](https://github.com/bmatzelle/gow/wiki)

**TODO** Instructions for installing deps on Windows

