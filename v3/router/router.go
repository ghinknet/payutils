package router

import (
	"strings"
)

func Route(upstreamName string, interfaceType string) string {
	return strings.Join([]string{"", upstreamName, interfaceType}, "/")
}
