package errors

type PayutilsError struct {
	message          string
	upstreamName     string
	upstreamCode     string
	upstreamMessage  string
	upstreamResponse any
	frameworkName    string
}

func (e *PayutilsError) Error() string {
	return e.message
}

func (e *PayutilsError) WithUpstreamName(upstreamName string) *PayutilsError {
	e.upstreamName = upstreamName
	return e
}

func (e *PayutilsError) UpstreamName() string {
	return e.upstreamName
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

func (e *PayutilsError) WithFrameworkName(frameworkName string) *PayutilsError {
	e.frameworkName = frameworkName
	return e
}

func (e *PayutilsError) FrameworkName() string {
	return e.frameworkName
}

type Option func(*PayutilsError)

func WithUpstreamName(upstreamName string) Option {
	return func(e *PayutilsError) {
		e.upstreamName = upstreamName
	}
}

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

func WithFrameworkName(frameworkName string) Option {
	return func(e *PayutilsError) {
		e.frameworkName = frameworkName
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

var ErrDriverNotRegistered = New("driver not registered")
