package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

type proxy struct {
}

const (
	HTTP  = "http"
	HTTPS = "https"
)

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v: %v\n", r.Method, r.URL)
	if r.Method == http.MethodConnect {
		handleTunneling(w, r)
	} else if r.URL.Path == "/" {
		handleInternalRouting(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	//connect to remote server
	proxytoDest, err := net.DialTimeout("tcp", r.Host, time.Second*10)
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
	clientToProxy, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		fmt.Println(err)
	}
	go transfer(proxytoDest, clientToProxy)
	go transfer(clientToProxy, proxytoDest)
}

func transfer(writer io.WriteCloser, source io.ReadCloser) {
	defer writer.Close()
	defer source.Close()
	_, err := io.Copy(writer, source)
	if err != nil {
		//example connection closed
		//readfrom tcp 192.168.45.187:56740->111.119.27.33:443: read tcp [::1]:8080->[::1]:56739: use of closed network connection
		//fmt.Println(err)
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

func handleInternalRouting(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	body := "hello"
	data := []byte(body)
	w.Write(data)
}
