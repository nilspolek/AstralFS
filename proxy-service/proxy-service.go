package proxyservice

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

var (
	RouteNotFoundErr = errors.New("route not found")
)

type svc struct {
	routes   map[string]string
	routesMu sync.Mutex
	router   *mux.Router
}

func New(r *mux.Router) ProxyService {
	return &svc{
		routes:   make(map[string]string),
		routesMu: sync.Mutex{},
		router:   r,
	}
}

func (s *svc) AddRoute(r Route) error {
	s.routesMu.Lock()
	s.routes[r.Path] = r.Target
	s.routesMu.Unlock()
	targetURL, err := url.Parse(r.Target)
	if err != nil {
		return err
	}
	s.router.PathPrefix(r.Path).HandlerFunc(s.ProxyRequest(targetURL))
	return nil
}

func (s *svc) DeleteRoute(r Route) error {
	s.routesMu.Lock()
	if _, exists := s.routes[r.Path]; !exists {
		return RouteNotFoundErr
	}
	delete(s.routes, r.Path)
	s.routesMu.Unlock()

	s.router.PathPrefix(r.Path).Handler(nil)
	return nil
}

func (s *svc) GetRoutes() map[string]string {
	s.routesMu.Lock()
	defer s.routesMu.Unlock()
	return s.routes
}

func (s *svc) ProxyRequest(url *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		if realIP := r.Header.Get("X-Forwarded-For"); realIP != "" {
			clientIP = realIP + ", " + clientIP
		}
		r.Header.Set("X-Forwarded-For", clientIP)

		r.Header.Set("X-Forwarded-Proto", r.URL.Scheme)

		r.Header.Set("X-Forwarded-Host", r.Host)

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)
	}
}
