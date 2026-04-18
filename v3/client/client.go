package client

import (
	"encoding/json"

	"github.com/ghinknet/payutils/v3/errors"
	"github.com/ghinknet/payutils/v3/internal/state"
	"github.com/ghinknet/payutils/v3/model"
	"github.com/ghinknet/payutils/v3/router"
)

func NewClient(config model.Config) (client *model.Client, err error) {
	// Prepare JSON
	if config.Marshal == nil {
		config.Marshal = json.Marshal
	}
	if config.Unmarshal == nil {
		config.Unmarshal = json.Unmarshal
	}

	// Register pay clients
	payClients := make(map[string]model.PayClient)

	for upstreamName, upstreamCredential := range config.Credentials {
		// Check driver registered
		upstream, ok := state.PayDrivers[upstreamName]
		if !ok {
			return nil, errors.ErrDriverNotRegistered.WithUpstreamName(upstreamName)
		}

		// Create driver client
		payClients[upstreamName], err = upstream.NewClient(model.PayDriverClientParam{
			Credential: upstreamCredential,
			// JSON
			Marshal:   config.Marshal,
			Unmarshal: config.Unmarshal,
		})
		if err != nil {
			return nil, err
		}
	}

	// Register http routers
	for instanceName, instance := range config.Instances {
		// Check driver registered
		framework, ok := state.HttpDrivers[instanceName]
		if !ok {
			return nil, errors.ErrDriverNotRegistered.WithFrameworkName(instanceName)
		}

		for upstreamName, payClient := range payClients {
			// Create unified instance
			unifiedInstance, err := framework.NewInstance(instance)
			if err != nil {
				return nil, err
			}

			// Register router
			unifiedInstance.Post(router.Route(upstreamName, "create"), payClient.Create)
			unifiedInstance.Post(router.Route(upstreamName, "callback"), payClient.Callback)
		}
	}

	return &model.Client{
		PayClient: payClients,
	}, nil
}
