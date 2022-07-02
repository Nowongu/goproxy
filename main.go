package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	_, err := io.Copy(destination, source)
	if err != nil {
		fmt.Println(err)
	}
}

// Hop-by-hop headers
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

type proxy struct {
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	//connect to remote server
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		fmt.Println(err)
		return
	}

	//todo: is this necessary
	w.WriteHeader(http.StatusOK)

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	//take connection away from http library with hijacker
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		fmt.Println(err)
	}

	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Scheme != HTTP && r.URL.Scheme != HTTPS {
		fmt.Printf("Unsupported protocal %v\n", r.URL)
		return
	}
	//http: Request.RequestURI can't be set in client requests.
	//http://golang.org/src/pkg/net/http/client.go
	r.RequestURI = ""

	//remove hop headers
	for _, h := range hopHeaders {
		r.Header.Del(h)
	}

	client := &http.Client{}
	resp, err := client.Do(r)
	//if there's an err resp will be nil
	if err != nil {
		http.Error(w, "InternalServer", http.StatusInternalServerError)
		fmt.Printf("ServeHTTP: %v\n", err)
		return
	}

	//defer evaluates Body.Close() immediately but the Close() func is only executed when handleHTTP returns.
	defer resp.Body.Close()

	//copy response
	for k, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(k, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v: %v\n", r.Method, r.URL)

	if r.Method == http.MethodConnect {
		handleTunneling(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func main() {
	fmt.Println("Proxy started listening on port 8080...")

	handler := &proxy{}

	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Println("ListenAndServe:", err)
	}
}
