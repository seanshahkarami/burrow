package gopher

import (
	"fmt"
	"io"
	"net"
	"net/url"
)

func Open(addr string) (io.ReadCloser, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	return OpenURL(u)
}

func stringOrDefault(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func OpenURL(URL *url.URL) (io.ReadCloser, error) {
	addr := fmt.Sprintf("%s:%s", URL.Hostname(), stringOrDefault(URL.Port(), "70"))

	// connect to server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	// send request to server
	if _, err := fmt.Fprintf(conn, "%s\r\n", URL.Path); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}
