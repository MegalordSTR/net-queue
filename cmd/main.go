package main

import (
	"flag"
	"fmt"
	"net-queue/internal"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	port := flag.Int("port", 8081, "port")
	flag.Parse()

	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	controller := internal.NewController(internal.NewMessageQueueMap())

	go func() {
		fmt.Printf("Http server stating on %s\n", addr)
		err := http.ListenAndServe(addr, controller)
		if err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Kill, os.Interrupt, os.Interrupt)
	<-stop
}
