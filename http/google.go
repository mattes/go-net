package http

import (
	"time"
)

// SetDefaultsForGoogleCloudLoadBalancer sets the timeout defaults
// for servers running behind the Google Cloud Load Balancer
// see https://cloud.google.com/load-balancing/docs/https/#timeouts_and_retries
func SetDefaultsForGoogleCloudLoadBalancer() {
	DefaultTcpKeepAlivePeriod = 620 * time.Second
}
