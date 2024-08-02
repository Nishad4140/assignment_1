package customframework

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Server struct {
	routes []*route
	server *http.Server
	mu     sync.Mutex
}

type route struct {
	method     string
	pattern    string
	regex      *regexp.Regexp
	paramNames []string
	handler    RouteHandler
}

type RouteHandler = func(*Request, *Response)

func NewServer() *Server {
	return &Server{
		routes: make([]*route, 0),
	}
}

func (s *Server) Listen(port uint16) error {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s,
	}
	return s.server.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	for _, route := range s.routes {
		if route.method != r.Method && route.method != "ANY" {
			continue
		}
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			params := make(map[string]string)
			for i, name := range route.paramNames {
				if i < len(matches)-1 {
					params[name] = matches[i+1]
				}
			}
			req := &Request{r, params}
			res := &Response{ResponseWriter: w}
			route.handler(req, res)
			return
		}
	}
	http.NotFound(w, r)
}

func (s *Server) registerRoute(method, pattern string, handler RouteHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	paramNames := []string{}
	regexPattern := pattern

	paramRegex := regexp.MustCompile(`\{([^/]+)\}`)
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)
	for _, match := range matches {
		paramNames = append(paramNames, match[1])
		regexPattern = strings.Replace(regexPattern, match[0], "([^/]+)", 1)
	}

	regex := regexp.MustCompile("^" + regexPattern + "$")
	s.routes = append(s.routes, &route{method, pattern, regex, paramNames, handler})
}

func (s *Server) Get(route string, handler RouteHandler) {
	s.registerRoute(http.MethodGet, route, handler)
}

func (s *Server) Post(route string, handler RouteHandler) {
	s.registerRoute(http.MethodPost, route, handler)
}

func (s *Server) Put(route string, handler RouteHandler) {
	s.registerRoute(http.MethodPut, route, handler)
}

func (s *Server) Delete(route string, handler RouteHandler) {
	s.registerRoute(http.MethodDelete, route, handler)
}

func (s *Server) Any(route string, handler RouteHandler) {
	s.registerRoute("ANY", route, handler)
}

func (s *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}
