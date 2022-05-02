//go:build debug

package process

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal("[debug] pprof listen error:", err)
	}
	go func() {
		log.Println(http.Serve(listener, nil))
	}()
	log.Println("[debug] pprof serve:", "http://"+listener.Addr().String()+"/debug/pprof/")
}
