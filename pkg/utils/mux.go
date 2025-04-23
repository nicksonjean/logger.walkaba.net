package utils

import (
	"net/http"
	"sync"
)

type CountingServeMux struct {
	*http.ServeMux
	handlersCount int
	mu            sync.Mutex
}

func NewCountingServeMux() *CountingServeMux {
	return &CountingServeMux{
		ServeMux:      http.NewServeMux(),
		handlersCount: 0,
	}
}

func (m *CountingServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.mu.Lock()
	m.handlersCount++
	m.mu.Unlock()
	m.ServeMux.HandleFunc(pattern, handler)
}

func (m *CountingServeMux) Handle(pattern string, handler http.Handler) {
	m.mu.Lock()
	m.handlersCount++
	m.mu.Unlock()
	m.ServeMux.Handle(pattern, handler)
}

func (m *CountingServeMux) HandlersCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.handlersCount
}
