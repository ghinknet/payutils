package model

// PayutilsError provides a customised error type includes detail error for payutils
type PayutilsError struct {
	message          string
	upstreamCode     int
	upstreamResponse string
	upstreamMessage  string
}

// Error returns basic error message
func (e *PayutilsError) Error() string {
	return e.message
}

// Is returns the compare result of two error value
func Is(err error, target error) bool {
	return err.Error() == target.Error()
}

// UpstreamDetail returns detail of upstream error in error value
func (e *PayutilsError) UpstreamDetail() (int, string, string) {
	return e.upstreamCode, e.upstreamResponse, e.upstreamMessage
}

// Option provides a basic option type
type Option func(*PayutilsError)

// WithUpstreamCode sets upstream code in an error content
func WithUpstreamCode(upstreamCode int) Option {
	return func(e *PayutilsError) {
		e.upstreamCode = upstreamCode
	}
}

// WithUpstreamResponse sets upstream response in an error content
func WithUpstreamResponse(upstreamResponse string) Option {
	return func(e *PayutilsError) {
		e.upstreamResponse = upstreamResponse
	}
}

// WithUpstreamMessage sets upstream message in an error content
func WithUpstreamMessage(upstreamMessage string) Option {
	return func(e *PayutilsError) {
		e.upstreamMessage = upstreamMessage
	}
}

// New returns a error value
func New(c string, options ...Option) error {
	err := &PayutilsError{message: c}

	for _, option := range options {
		option(err)
	}

	return err
}

// Config

var ErrMissAllowedOrigin = New("miss allowed origin")
var ErrMissEndpoint = New("miss endpoint")
var ErrMissHandler = New("miss handler")

// Public payment

var ErrPaymentMethodDisabled = New("payment method is disabled")
var ErrNoEnoughTimeToPay = New("no enough time to pay")
var ErrUnsupportedCurrency = New("unsupported currency")
var ErrUnsupportedMethod = New("unsupported method")
