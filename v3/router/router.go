package router

import (
	"strings"
)

// Route builds the in-app route path for an upstream interface, e.g.
// Route("alipay", "callback") -> "/alipay/callback".
func Route(upstreamName string, interfaceType string) string {
	return strings.Join([]string{"", upstreamName, interfaceType}, "/")
}

// Notify builds the absolute callback (notify) URL an upstream should call
// back to. It joins the endpoint with the callback route while tolerating a
// trailing slash on the endpoint, e.g.
// Notify("https://gh.ink/", "alipay") -> "https://gh.ink/alipay/callback".
func Notify(endpoint string, upstreamName string) string {
	return strings.Join([]string{
		strings.TrimSuffix(endpoint, "/"), Route(upstreamName, "callback"),
	}, "")
}
