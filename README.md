# goproxy
Web Proxy

## Issues
go net/http can't handle the `CONNECT` http method which is why the proxy handles separately in the `handleTunneling` func. This is quite important since most http/1 requests through a proxy
are done using the `CONNECT` http method.

## TODO
- Add logging instead of just printing to console.
- Record total bytes sent and received.
- Add a config file to set port/log details


## Reference links
- https://www.sobyte.net/post/2021-09/https-proxy-in-golang-in-less-than-100-lines-of-code/
- https://github.com/sipt/shuttle/blob/cf12e39f79ddc71155b769108b338b0227d1ee0e/transport_channel.go#L40
- https://github.com/golang/go/issues/17227
- https://www.sohamkamani.com/golang/channels/