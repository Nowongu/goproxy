package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	proxyConfig.Load()

	host := proxyConfig.Hostname + ":" + strconv.Itoa(proxyConfig.Port)
	fmt.Printf("Proxy started listening on port %v...", host)

	handler := &proxy{}

	if err := http.ListenAndServe(host, handler); err != nil {
		fmt.Println("Terminated:", err)
	}
}
