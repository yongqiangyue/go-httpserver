FROM ubuntu
ENV MY_SERVICE_PORT=9000

ADD bin/amd64/go-httpserver /go-httpserver
EXPOSE ${MY_SERVICE_PORT} 
ENTRYPOINT /go-httpserver -port ${MY_SERVICE_PORT} 
