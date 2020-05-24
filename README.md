# httprouter-example

Example of using [httprouter](https://github.com/julienschmidt/httprouter)
with [zerolog](https://github.com/rs/zerolog)
and gorilla middleware [handlers](https://github.com/gorilla/handlers)


## Quick start

Clone the repos (outside your GOPATH since this is a module)

    git clone https://github.com/mozey/httprouter-example.git
    
    cd httprouter-example

Following the 12 factor app recommendation to
[store config in the environment](https://12factor.net/config).
Configuration is done using [environment variables](https://en.wikipedia.org/wiki/Environment_variable)

Copy script to export dev config

    cp sample.dev.sh dev.sh 

Run dev server (no live reload)

    ./dev.sh app_run
    
Run dev server with live reload
    
    ./dev.sh app
    
**Alternatively**,
use [mozey/config](https://github.com/mozey/config)
to manage env with a flat config.json file.
First setup a helper func for [toggling env](https://github.com/mozey/config#toggling-env).
Then run the commands below
    
    APP_DIR=$(pwd) ./scripts/config.sh
    
    conf && ./make.sh app
    
Or start up services 
in [tmux](https://www.hamvocke.com/blog/a-quick-and-easy-guide-to-tmux)

    cp sample.up.sh up.sh
    
    ./up.sh
    
    ./tmux.sh app
    # ctrl+b d
    
Show running processes
   
    ps ax | grep ${APP_NAME}
    
Shut down 

    cp sample.down.sh down.sh
    
    ./down.sh 
    
   
## Examples
  
Token is required by default    
[http://localhost:8118/token/is/required/by/default](http://localhost:8118/token/is/required/by/default)

Static content skips the token check
- [http://localhost:8118/index.html](http://localhost:8118/index.html)
- [http://localhost:8118](http://localhost:8118)
- [http://localhost:8118/www/data/go.txt](http://localhost:8118/www/data/go.txt)
    
[http://localhost:8118/hello/foo?token=123](http://localhost:8118/hello/foo?token=123)
    
[http://localhost:8118/api?token=123](http://localhost:8118/api?token=123)
    
[http://localhost:8118/panic](http://localhost:8118/panic)
    
[http://localhost:8118/does/not/exist?token=123](http://localhost:8118/does/not/exist?token=123)
    
**TODO**
[http://localhost:8118/proxy](http://localhost:8118/proxy)
    
**NOTE** 
Make requests from the cli with [httpie](https://httpie.org/)


## Reset

Removes all user config

    APP_DIR=$(pwd) ./scripts/reset.sh


## Dependencies

This example aims for a good cross platform experience by depending on 
- [Golang](https://golang.org/) 
- [Bash](https://www.gnu.org/software/bash)
- [fswatch](https://github.com/emcrisostomo/fswatch)
- [xargs](https://github.com/emcrisostomo/fswatch)

On macOS, Linux
- [pgrep](https://en.wikipedia.org/wiki/Pgrep) and kill

On Windows
- **TODO** [taskkill](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/taskkill)

[GNU Make](https://stackoverflow.com/questions/3798562/why-use-make-over-a-shell-script) 
is not needed because Golang is fast to build,
and `fswatch` can be used for live reload.
For this example `main.go` is kept in the project root.
Larger projects might have separate bins in the *"/cmd"* dir

Bash on Windows is easy to setup using 
[msys2](https://www.msys2.org/), MinGW, or native shell on Windows 10.
For other UNIX programs see [gow](https://github.com/bmatzelle/gow/wiki)

**TODO** Instructions for installing deps on Windows