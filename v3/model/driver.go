package model

type PayDriver interface {
	NewClient(params PayDriverClientParam) (PayClient, error)
}

type PayDriverClientParam struct {
	// Pay client credential
	Credential map[string]string
	// JSON
	Unmarshal func(data []byte, v any) error
	Marshal   func(v any) ([]byte, error)
}

type HttpDriver interface {
	NewInstance(instance any) (HttpInstance, error)
}
