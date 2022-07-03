package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Proxy started listening on port 8080...")

	handler := &proxy{}

	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
