package server

import (
	httppkg "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattes/go-net/http"
)

func ListenAndServe(addr string, mux *httppkg.ServeMux) (err error) {
	// Register SIGINT and SIGTERM interrupts
	osShutdown := make(chan os.Signal, 1)
	signal.Notify(osShutdown, syscall.SIGINT, syscall.SIGTERM)

	// Create Healthy and Ready Handlers
	http.NewHealthyHandler(mux).Ok()
	ready := http.NewReadyHandler(mux)

	// Listen and serve ...
	errChan := make(chan error, 1)
	serverShutdown := http.NewShutdownChan()
	go func() {
		errChan <- http.ListenAndServeWithShutdown(addr, mux, serverShutdown)
	}()

	// Ready to process requests
	ready.Ok()

	// Wait for interrupts ...
	<-osShutdown
	ready.NotOk()
	serverShutdown <- struct{}{}
	return <-errChan
}

func ListenAndServeTLS(addr, certFile, keyFile string, mux *httppkg.ServeMux) error {
	// Register SIGINT and SIGTERM interrupts
	osShutdown := make(chan os.Signal, 1)
	signal.Notify(osShutdown, syscall.SIGINT, syscall.SIGTERM)

	// Create Healthy and Ready Handlers
	http.NewHealthyHandler(mux).Ok()
	ready := http.NewReadyHandler(mux)

	// Listen and serve ...
	errChan := make(chan error, 1)
	serverShutdown := http.NewShutdownChan()
	go func() {
		errChan <- http.ListenAndServeTLSWithShutdown(addr, certFile, keyFile, mux, serverShutdown)
	}()

	// Ready to process requests
	ready.Ok()

	// Wait for interrupts ...
	<-osShutdown
	ready.NotOk()
	serverShutdown <- struct{}{}
	return <-errChan
}

func IsBehindGoogleLoadBalancer() {
	http.SetDefaultsForGoogleCloudLoadBalancer()
}
