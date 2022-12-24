package utils

import "github.com/vulcand/oxy/forward"

func ProxyServer() (*forward.Forwarder, error) {
	return forward.New()
}
