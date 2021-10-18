FROM golang:alpine AS builder
WORKDIR /go/src/github.com/yongqiangyue/go-httpserver/
# RUN go get -d -v github.com/felixge/httpsnoop
#RUN go env -w GOPROXY=https://goproxy.cn \
#         && go get -d -v github.com/felixge/httpsnoop 
RUN go env -w GOPROXY=https://goproxy.cn
COPY go.mod .
COPY go.sum .
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64/go-httpserver .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/yongqiangyue/go-httpserver/bin/amd64/go-httpserver .
ENV MY_SERVICE_PORT=9000
EXPOSE ${MY_SERVICE_PORT} 
ENTRYPOINT ./go-httpserver -port ${MY_SERVICE_PORT} 
