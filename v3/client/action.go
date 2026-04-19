package client

import (
	"github.com/ghinknet/payutils/v3/errors"
	"github.com/ghinknet/payutils/v3/model"
)

func (c *Client) Status(upstreams []string, tradeID string) (map[string]model.ReturnStatus, error) {
	result := make(map[string]model.ReturnStatus)

	for _, upstream := range upstreams {
		if u, ok := c.PayClient[upstream]; !ok {
			return make(map[string]model.ReturnStatus), errors.ErrUpstreamNotFound.WithUpstreamName(upstream)
		} else {
			status, err := u.Status(tradeID)
			if err != nil {
				return make(map[string]model.ReturnStatus), err
			}
			result[upstream] = status
		}
	}

	return result, nil
}

func (c *Client) Close(upstreams []string, tradeID string) error {
	for _, upstream := range upstreams {
		if u, ok := c.PayClient[upstream]; !ok {
			return errors.ErrUpstreamNotFound.WithUpstreamName(upstream)
		} else {
			if err := u.Close(tradeID); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Client) Refund(
	upstreams []string, tradeID string,
	currency string, refundID string,
	totalAmount int64, refundAmount int64,
	reason string,
) error {
	for _, upstream := range upstreams {
		if u, ok := c.PayClient[upstream]; !ok {
			return errors.ErrUpstreamNotFound.WithUpstreamName(upstream)
		} else {
			if err := u.Refund(
				tradeID,
				currency, refundID,
				totalAmount, refundAmount,
				reason,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
