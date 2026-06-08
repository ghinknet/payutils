package client

import (
	"net/http"

	"go.gh.ink/payutils/v3/errors"
)

// Callback handles an upstream asynchronous notification manually.
//
// This is the counterpart of the auto-registered callback route: when the user
// does not import an http driver (or wants full control over routing), they can
// route the request themselves and forward the standard *http.Request here.
func (c *Client) Callback(upstream string, w http.ResponseWriter, r *http.Request) error {
	pc, ok := c.PayClient[upstream]
	if !ok {
		return errors.ErrUpstreamNotFound.WithUpstreamName(upstream)
	}
	pc.Callback(w, r)
	return nil
}
