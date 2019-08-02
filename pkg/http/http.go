package http

// TODO: Implement more..

import (
	"net/http"
)

// NewClient configures a basic http Client with Timeout, TLS (Validate, ...) etc.
// so its directly usable and the Configuration of the Client is the same everywhere
// in the Code
func NewClient() http.Client {
	return http.Client{
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       30,
	}
}
