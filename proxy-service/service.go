package proxyservice

import (
	"net/http"
	"net/url"
)

type ProxyService interface {
	AddRoute(Route) error
	DeleteRoute(Route) error
	GetRoutes() map[string]string
	ProxyRequest(*url.URL) http.HandlerFunc
}

type Route struct {
	Path   string `json:"path"`
	Target string `json:"target"`
}
