package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{}

func init() {
	go runRemote()
}

func runRemote() {
	server := &http.Server{
		Addr:         "localhost:8081",
		Handler:      http.HandlerFunc(handleRemoteRequest),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("About to start test proxy at: %s", "localhost:8081")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Unable to start test proxy: %s", err)
	}
}

func handleRemoteRequest(resp http.ResponseWriter, req *http.Request) {
	host := hostIncludingPort(req)
	if connOut, err := net.Dial("tcp", host); err != nil {
		msg := fmt.Sprintf("Unable to open socket to server: %s", err)
		respondBadGateway(resp, req, msg)
	} else {
		if connIn, _, err := resp.(http.Hijacker).Hijack(); err != nil {
			msg := fmt.Sprintf("Unable to access underlying connection from downstream proxy: %s", err)
			respondBadGateway(resp, req, msg)
		} else {
			if req.Method == "CONNECT" {
				connIn.Write([]byte("HTTP/1.0 200 OK\r\n\r\n"))
			} else {
				req.Write(connOut)
			}
			pipe(connIn, connOut)
		}
	}
}

func hostIncludingPort(req *http.Request) (host string) {
	host = req.Host
	if !strings.Contains(host, ":") {
		if req.Method == "CONNECT" {
			host = host + ":443"
		} else {
			host = host + ":80"
		}
	}
	return
}
