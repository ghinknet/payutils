package model

type Config struct {
	// Payutils
	AllowOrigins []string
	Endpoint     string
	Prefix       string
	Suffix       string
	// Http
	Instances map[string]any
	// Pay
	Credentials C
	// JSON
	Unmarshal func(data []byte, v any) error
	Marshal   func(v any) ([]byte, error)
}
