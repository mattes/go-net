package http

import (
	"context"
	"net"
	"net/http"
	"time"

	netpkg "github.com/mattes/go-net"
)

var (
	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	DefaultReadTimeout = 30 * time.Second

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.
	DefaultReadHeaderTimeout = 30 * time.Second

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	DefaultWriteTimeout = 30 * time.Second

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	DefaultIdleTimeout = 60 * time.Second

	// TcpKeepAlivePeriod is the amount of time a TCP connection
	// is kept open. On Linux, this translates to the following
	// two Socket options being set:
	// TCP_KEEPINTVL The time between individual keepalive probes.
	// TCP_KEEPIDLE  The time the connection needs to remain idle
	//               before TCP starts sending keepalive probes.
	DefaultTcpKeepAlivePeriod = 120 * time.Second

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	DefaultMaxHeaderBytes = http.DefaultMaxHeaderBytes
)

func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

func ListenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServeTLS(certFile, keyFile)
}

type Server struct {
	addr string
	orig *http.Server
}

func NewServer(addr string, handler http.Handler) *Server {
	s := &Server{
		addr: addr,
	}

	orig := &http.Server{
		Handler: handler,

		ReadTimeout:       DefaultReadTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MaxHeaderBytes:    DefaultMaxHeaderBytes,
	}

	s.orig = orig

	return s
}

func (s *Server) Close() error {
	return s.orig.Close()
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	return s.Serve(
		&netpkg.TcpKeepAliveListener{
			listener.(*net.TCPListener),
			DefaultTcpKeepAlivePeriod})
}

func (s *Server) ListenAndServeTLS(certFile, keyFile string) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.orig.Addr = listener.Addr().String()

	return s.ServeTLS(
		&netpkg.TcpKeepAliveListener{
			listener.(*net.TCPListener),
			DefaultTcpKeepAlivePeriod},
		certFile, keyFile)
}

func (s *Server) RegisterOnShutdown(f func()) {
	s.orig.RegisterOnShutdown(f)
}

func (s *Server) Serve(l net.Listener) error {
	s.orig.Addr = l.Addr().String()
	return s.orig.Serve(l)
}

func (s *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	s.orig.Addr = l.Addr().String()
	return s.orig.ServeTLS(l, certFile, keyFile)
}

func (s *Server) SetKeepAlivesEnabled(v bool) {
	s.orig.SetKeepAlivesEnabled(v)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.orig.Shutdown(ctx)
}

func (s *Server) ShutdownWithTimeout(ctx context.Context, timeout time.Duration) error {
	c, _ := context.WithTimeout(ctx, timeout)
	return s.orig.Shutdown(c)
}
