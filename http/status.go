package http

import (
	"net/http"
	"sync"
)

var (
	DefaultStatusOkCode    = 200 // OK
	DefaultStatusNotOkCode = 503 // Service Unavailable
)

type Status struct {
	mu sync.RWMutex
	ok bool
}

func (s *Status) Ok() {
	s.mu.Lock()
	s.ok = true
	s.mu.Unlock()
}

func (s *Status) NotOk() {
	s.mu.Lock()
	s.ok = false
	s.mu.Unlock()
}

func NewStatusHandler(pattern string, mux *http.ServeMux, okCode, notOkCode int) *Status {
	s := &Status{}

	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		ok := s.ok
		s.mu.RUnlock()

		if ok {
			w.WriteHeader(okCode)
		} else {
			w.WriteHeader(notOkCode)
		}
	})

	return s
}

func NewHealthyHandler(mux *http.ServeMux) *Status {
	return NewStatusHandler("/health", mux, DefaultStatusOkCode, DefaultStatusNotOkCode)
}

func NewReadyHandler(mux *http.ServeMux) *Status {
	return NewStatusHandler("/ready", mux, DefaultStatusOkCode, DefaultStatusNotOkCode)
}
