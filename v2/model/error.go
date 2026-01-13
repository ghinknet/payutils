package model

import "errors"

// Config

var ErrMissAllowedOrigin = errors.New("miss allowed origin")
var ErrMissEndpoint = errors.New("miss endpoint")
var ErrMissHandler = errors.New("miss handler")

// Public payment

var ErrNoEnoughTimeToPay = errors.New("no enough time to pay")
var ErrUnsupportedCurrency = errors.New("unsupported currency")
