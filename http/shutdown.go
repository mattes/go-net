package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var (
	DefaultShutdownTimeout = 60 * time.Second

	ErrUncleanShutdown = fmt.Errorf("unclean shutdown: context deadline exceeded")
)

func (s *Server) ListenAndServeWithShutdown(shutdown chan struct{}) error {
	err := make(chan error, 1)
	go func() {
		err <- s.ListenAndServe()
	}()

	return handleShutdownErrs(s, err, shutdown)
}

func ListenAndServeWithShutdown(addr string, handler http.Handler, shutdown chan struct{}) error {
	s := NewServer(addr, handler)

	err := make(chan error, 1)
	go func() {
		err <- s.ListenAndServe()
	}()

	return handleShutdownErrs(s, err, shutdown)
}

func (s *Server) ListenAndServeTLSWithShutdown(certFile, keyFile string, shutdown chan struct{}) error {
	err := make(chan error, 1)
	go func() {
		err <- s.ListenAndServeTLS(certFile, keyFile)
	}()

	return handleShutdownErrs(s, err, shutdown)
}

func ListenAndServeTLSWithShutdown(addr, certFile, keyFile string, handler http.Handler, shutdown chan struct{}) error {
	s := NewServer(addr, handler)

	err := make(chan error, 1)
	go func() {
		err <- s.ListenAndServeTLS(certFile, keyFile)
	}()

	return handleShutdownErrs(s, err, shutdown)
}

func NewShutdownChan() chan struct{} {
	return make(chan struct{}, 1)
}

func handleShutdownErrs(s *Server, listenErr chan error, shutdown chan struct{}) error {
	shutdownErr := make(chan error, 1)

	// Before shutdown request is received:
beforeShutdown:
	for {
		select {
		case err := <-listenErr:
			// We haven't received a shutdown request so far,
			// immediately return any error fromm ListenAndServe
			return err

		case <-shutdown:
			// Shutdown request received
			go func() {
				shutdownErr <- s.ShutdownWithTimeout(
					context.Background(),
					DefaultShutdownTimeout)
			}()
			break beforeShutdown
		}
	}

	// After shutdown request was received:
	// Ignore certain errors from ListenAndServe,
	// and either return error from ListenAndServe
	// or return error or nil from Shutdown
	for {
		select {
		case err := <-listenErr:
			// Ignore ErrServerClosed, because we are actually closing the Server. Duh.
			// All other errors we return immediately though.
			if err != nil && err != http.ErrServerClosed {
				return err
			}

		case err := <-shutdownErr:
			// Non-nil error means that shutdown was not clean.
			// Connections were dropped.
			if err == context.DeadlineExceeded {
				return ErrUncleanShutdown
			}
			return err
		}
	}
}
