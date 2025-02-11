package functional

import (
	"errors"
	"net"
	"time"
)

// waitHTTPServer waits for the HTTP server to be ready
func waitHTTPServer(listenAddress string, sleepTime time.Duration, retries int) error {
	var err error
	var conn net.Conn

	for i := 0; i < retries; i++ {
		conn, err = net.DialTimeout("tcp", listenAddress, sleepTime)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(sleepTime)
	}

	return errors.New("HTTP server is not ready")
}
