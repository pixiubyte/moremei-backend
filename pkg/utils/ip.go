package utils

import (
	"net"
	"net/http"
)

const (
	xRealIP       = "X-Real-IP"
	xForwardedFor = "X-Forwarded-For"
)

func GetRemoteIp(r *http.Request) string {
	var remoteIp string
	if ip := r.Header.Get(xRealIP); len(ip) > 0 {
		remoteIp = ip
	} else if ip = r.Header.Get(xForwardedFor); len(ip) > 0 {
		remoteIp = ip
	} else {
		remoteIp, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	if remoteIp == "" || remoteIp == "::1" {
		remoteIp = "127.0.0.1"
	}

	return remoteIp
}
