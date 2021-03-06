# goproxy
Web Proxy

## Issues
go net/http [can't handle](https://github.com/golang/go/issues/17227) the `CONNECT` http method which is why the proxy uses a separate `handleTunneling` func. This is important since most http/1.1 requests use the `CONNECT` method for tsl connections to hide data from the proxy. http/2 doesn't suffer from the same problem since it doesn't use the CONNECT verb, however go doesn't have good http2 support at the moment, but development is in progress.


## Reference links
- https://www.sobyte.net/post/2021-09/https-proxy-in-golang-in-less-than-100-lines-of-code/
- https://github.com/sipt/shuttle
- https://www.sohamkamani.com/golang/channels/

## Create certificate for HTTP/2
```
openssl genrsa -out https-server.key 2048
openssl ecparam -genkey -name secp384r1 -out https-server.key
openssl req -new -x509 -sha256 -key https-server.key -out https-server.crt -days 3650
```
