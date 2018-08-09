package http

import (
	"net"
	"net/http"
	"testing"
	"time"
)

var randomAddr = "127.0.0.1:0"

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("/sleepForever", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(60 * 24 * time.Minute)
		w.WriteHeader(500)
	})
}

func TestListenAndServeWithShutdown(t *testing.T) {
	shutdown := make(chan struct{}, 0)
	errChan := make(chan error, 0)

	go func() {
		errChan <- ListenAndServeWithShutdown(randomAddr, http.DefaultServeMux, shutdown)
	}()

	shutdown <- struct{}{}

	if err := <-errChan; err != nil {
		t.Fatal(err)
	}
}

func TestServerListenAndServeWithShutdown_UncleanShutdown(t *testing.T) {
	shutdown := make(chan struct{}, 0)
	errChan := make(chan error, 0)

	server := NewServer(randomAddr, http.DefaultServeMux)

	go func() {
		errChan <- server.ListenAndServeWithShutdown(shutdown)
	}()

	// wait for server to start
	time.Sleep(2 * time.Second)

	// call server and let the request sleep forever
	go func() {
		_, err := http.Get("http://" + server.orig.Addr + "/sleepForever")
		if err != nil {
			t.Fatal(err)
		}
	}()

	// wait for GET call to be made
	time.Sleep(2 * time.Second)

	DefaultShutdownTimeout = 2 * time.Second
	shutdown <- struct{}{}

	if err := <-errChan; err != ErrUncleanShutdown {
		t.Fatal(err)
	}
}

func TestListenAndServeWithShutdown_BeforeShutdownErr(t *testing.T) {
	err := ListenAndServeWithShutdown("bogus addr", http.DefaultServeMux, nil)
	if operr, ok := err.(*net.OpError); !ok {
		t.Fatal("expected type net.OpError")
	} else if operr.Op != "listen" {
		t.Fatal("expected net.OpError with op 'listen'")
	}
}
