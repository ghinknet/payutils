package client

import (
	"context"

	"go.gh.ink/payutils/v3/errors"
)

// Create dispatches the creation param to the registered pay drivers.
//
// payutils does not own any order-creation logic: the caller prepares the
// trade detail and wraps it into a driver-specific param struct (e.g.
// payalipay.CreateParam or paywechat.CreateParam). Each driver asserts whether
// it claims the given param; the first one that claims it handles the request.
//
// result is the driver-specific payload (ready to be marshalled). When no
// driver claims the param, ErrNoDriverClaimed is returned.
func (c *Client) Create(ctx context.Context, param any) (any, error) {
	for _, pc := range c.PayClient {
		claimed, result, err := pc.Create(ctx, param)
		if !claimed {
			continue
		}
		return result, err
	}
	return nil, errors.ErrNoDriverClaimed
}
