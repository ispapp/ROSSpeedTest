package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	speedtest "github.com/kmoz000/RouterOsSpeedTest/v1"
)

func main() {
	log.Default().SetPrefix("[SPEEDTEST]")
	lport := "8000"
	if err := godotenv.Load(); err != nil {
		log.Default().Fatal(err)
	}
	if port := os.Getenv("PORT"); port != "" {
		lport = port
	}
	speedTest := speedtest.SpeedTest{Actives: make(map[string]speedtest.Test, 10000)}
	srv := &http.Server{
		ReadTimeout:  25 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  0 * time.Second,
		Addr:         fmt.Sprintf(":%s", lport),
		// this is the cost of long polling, exceeding the maximum file descriptor count on the system from invalid requests
		//IdleTimeout: 15 * time.Minute,
		MaxHeaderBytes: 1500,
		//ConnState: ConnStateEvent,
	}
	if l, err := net.Listen("tcp", fmt.Sprintf(":%s", lport)); err == nil {
		log.Default().Printf("%shttp://localhost:%s/speedtest", "Runnning on port: ", lport)
		http.HandleFunc("/speedtest", speedTest.Handler)
		srv.Serve(l)
		defer srv.Close()
	}

}
