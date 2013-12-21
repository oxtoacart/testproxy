package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

var (
	empty = []byte("")
)

func respondBadGateway(resp http.ResponseWriter, req *http.Request, msg string) {
	log.Println(msg)
	resp.WriteHeader(502)
	resp.Write([]byte(fmt.Sprintf("Bad Gateway: %s - %s", req.URL, msg)))
}

func pipe(connIn net.Conn, connOut net.Conn) {
	go func() {
		defer connIn.Close()
		io.Copy(connOut, connIn)
		connIn.Write(empty)
		connOut.Close()
	}()
	go func() {
		defer connOut.Close()
		io.Copy(connIn, connOut)
		connOut.Write(empty)
		connIn.Close()
	}()
}
