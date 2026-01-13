package model

// Config is the public part of config
type Config struct {
	Debug         bool
	AllowedOrigin []string
	Endpoint      string
}
