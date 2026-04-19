package errors

import "github.com/ghinknet/toolbox/pointer"

type PayutilsError struct {
	message          string
	upstreamName     string
	upstreamCode     string
	upstreamMessage  string
	upstreamResponse any
	frameworkName    string

	raw error
}

func (e *PayutilsError) Error() string {
	return e.message
}

func (e *PayutilsError) Is(err error) bool {
	return e.raw == err
}

func (e *PayutilsError) Unwrap() error {
	return e.raw
}

func (e *PayutilsError) WithUpstreamName(upstreamName string) *PayutilsError {
	ne := pointer.Copy(e)
	ne.upstreamName = upstreamName
	return ne
}

func (e *PayutilsError) UpstreamName() string {
	return e.upstreamName
}

func (e *PayutilsError) WithUpstreamCode(code string) *PayutilsError {
	ne := pointer.Copy(e)
	ne.upstreamCode = code
	return ne
}

func (e *PayutilsError) UpstreamCode() string {
	return e.upstreamCode
}

func (e *PayutilsError) WithUpstreamMessage(message string) *PayutilsError {
	ne := pointer.Copy(e)
	ne.upstreamMessage = message
	return ne
}

func (e *PayutilsError) UpstreamMessage() string {
	return e.upstreamMessage
}

func (e *PayutilsError) WithUpstreamResponse(response any) *PayutilsError {
	ne := pointer.Copy(e)
	ne.upstreamResponse = response
	return ne
}

func (e *PayutilsError) UpstreamResponse() any {
	return e.upstreamResponse
}

func (e *PayutilsError) WithFrameworkName(frameworkName string) *PayutilsError {
	ne := pointer.Copy(e)
	ne.frameworkName = frameworkName
	return ne
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

	err.raw = err

	return err
}

// Config

var ErrMissAllowedOrigin = New("miss allowed origin")
var ErrMissEndpoint = New("miss endpoint")
var ErrMissInstance = New("miss instance")

// Public payment

var ErrNoEnoughTimeToPay = New("no enough time to pay")
var ErrUnsupportedCurrency = New("unsupported currency")
var ErrUnsupportedMethod = New("unsupported method")

var ErrDriverNotRegistered = New("driver not registered")
var ErrUnsupportedInstance = New("unsupported instance")
var ErrUpstreamNotFound = New("upstream not found")
