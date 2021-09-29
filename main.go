package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
)

type State struct {
	Health string `json:"health"`
	Status uint   `json:"code"`
}

type HTTPReqInfo struct {
	// GET etc.
	method  string
	url     string
	referer string
	ipaddr  string
	// response code, like 200, 404
	code int
	// number of bytes of the response sent
	size int64
	// how long did it take to
	duration  time.Duration
	userAgent string
}

func healthz(w http.ResponseWriter, req *http.Request) {
	status := State{Health: "OK", Status: 200}
	js, _ := json.Marshal(status)
	w.Header().Set("content-type", "text/json")
	w.Header().Set("VERSION", os.Getenv("VERSION"))
	for i, v := range req.Header {
		w.Header().Set(i, v[0])
	}
	_, err := io.WriteString(w, string(js))
	if err != nil {
		log.Panic(err)
	}
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	addr := s[:idx]
	if addr == "[::1]" {
		addr = "127.0.0.1"
	}

	return addr
}

func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		log.Println(r.RemoteAddr)
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For is potentially a list of addresses separated with ","
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: should return first non-local address
		return parts[0]
	}
	log.Println("parts")
	return hdrRealIP
}

func logHTTPReq(req *HTTPReqInfo) {
	log.Printf("Client IP: %s, HTTP result code: %d", req.ipaddr, req.code)
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ri := &HTTPReqInfo{
			method:    r.Method,
			url:       r.URL.String(),
			referer:   r.Header.Get("Referer"),
			userAgent: r.Header.Get("User-Agent"),
		}

		ri.ipaddr = requestGetRemoteAddress(r)

		// this runs handler h and captures information about
		// HTTP request
		m := httpsnoop.CaptureMetrics(h, w, r)

		ri.code = m.Code
		ri.size = m.Written
		ri.duration = m.Duration
		logHTTPReq(ri)
	}
	return http.HandlerFunc(fn)
}

func makeHTTPServer(adrr string) *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/healthz", healthz)
	var handler http.Handler = mux

	handler = logRequestHandler(handler)

	srv := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second, // introduced in Go 1.8
		Handler:      handler,
	}
	srv.Addr = adrr
	return srv
}

func main() {
	log.Println("start httpserver...")
	port := flag.Int("port", 8888, "listen port")
	flag.Parse()
	log.Println("listen port :", *port)
	httpSrv := makeHTTPServer(fmt.Sprintf(":%d", *port))
	err := httpSrv.ListenAndServe()
	// http.HandleFunc("/healthz", healthz)
	// err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
