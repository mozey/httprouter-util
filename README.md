# httprouter-example

Minimal example of using [httprouter](https://github.com/julienschmidt/httprouter)
with [zerolog](https://github.com/rs/zerolog)
and middleware [handlers](https://github.com/gorilla/handlers)

# Quick start

Run dev server 

    go get github.com/mozey/httprouter-example
    
    cd ${GOPATH}/src/github.com/mozey/httprouter-example
    
    make dev
    
Make requests
    
    http localhost:8080
    
    http localhost:8080/panic
    
    http localhost:8080/hello/foo?token=123
    
    http localhost:8080/does/not/exist?token=123
