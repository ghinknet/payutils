package errors

type PayutilsError struct {
	message          string
	upstreamCode     string
	upstreamMessage  string
	upstreamResponse any
}

func (e *PayutilsError) Error() string {
	return e.message
}

func (e *PayutilsError) WithUpstreamCode(code string) *PayutilsError {
	e.upstreamCode = code
	return e
}

func (e *PayutilsError) UpstreamCode() string {
	return e.upstreamCode
}

func (e *PayutilsError) WithUpstreamMessage(message string) *PayutilsError {
	e.upstreamMessage = message
	return e
}

func (e *PayutilsError) UpstreamMessage() string {
	return e.upstreamMessage
}

func (e *PayutilsError) WithUpstreamResponse(response any) *PayutilsError {
	e.upstreamResponse = response
	return e
}

func (e *PayutilsError) UpstreamResponse() any {
	return e.upstreamResponse
}

type Option func(*PayutilsError)

func WithUpstreamCode(code string) Option {
	return func(e *PayutilsError) {
		e.upstreamCode = code
	}
}

func WithUpstreamMessage(message string) Option {
	return func(e *PayutilsError) {
		e.upstreamMessage = message
	}
}

func WithUpstreamResponse(response any) Option {
	return func(e *PayutilsError) {
		e.upstreamResponse = response
	}
}

func New(c string, options ...Option) *PayutilsError {
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
