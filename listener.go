// Package net implements a net.TCPListener with TCP KeepAlive timeouts
//
// Go's default ListenAndServe funcs [1] will wrap the underlying
// net.TCPListener with net.tcpKeepAliveListener with a hardcoded timeout of
// 3 minutes and now way of changing the timeout.
//
// If a Go server runs behind Google Cloud Loadbalancer for example,
// where the TCP session timeout is 10 minutes [2], this leads to the
// following error, logged by the Google Loadbalancer:
// `backend_connection_closed_before_data_sent_to_client`
//
// Please note that Google Cloud Loadbalancer seems to retry failed GET requests,
// but will not retry failed POST requests, which makes sense (idempotency).
//
// This package mimics the offical net.tcpKeepAliveListener but makes the
// timeout configurable.
//
// [1] https://golang.org/pkg/net/http/#Server.ListenAndServe
// [2] https://cloud.google.com/load-balancing/docs/https/#timeouts_and_retries
package net

import (
	"net"
	"time"
)

// TcpKeepAliveListener is more or less copied from:
// https://github.com/golang/go/blob/release-branch.go1.10/src/net/http/server.go#L3211
type TcpKeepAliveListener struct {
	*net.TCPListener
	Timeout time.Duration
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(ln.Timeout)
	return tc, nil
}
