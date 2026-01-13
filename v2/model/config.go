package model

// Config is the public part of config
type Config struct {
	Debug        bool
	AllowOrigins []string
	Endpoint     string
	Prefix       string
	Suffix       string
}
