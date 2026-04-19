package client

import (
	"encoding/json"

	"github.com/ghinknet/payutils/v3/errors"
	"github.com/ghinknet/payutils/v3/internal/action"
	"github.com/ghinknet/payutils/v3/internal/state"
	"github.com/ghinknet/payutils/v3/model"
	"github.com/ghinknet/payutils/v3/router"
)

type Client struct {
	PayClient map[string]model.PayClient
	Config    model.Config
}

func NewClient(config model.Config) (client *Client, err error) {
	// Prepare JSON
	if config.Marshal == nil {
		config.Marshal = json.Marshal
	}
	if config.Unmarshal == nil {
		config.Unmarshal = json.Unmarshal
	}

	// Check AllowOrigins
	if config.AllowOrigins == nil {
		return nil, errors.ErrMissAllowedOrigin
	}

	// Check endpoint
	if config.Endpoint == "" {
		return nil, errors.ErrMissEndpoint
	}

	// Check handler
	if len(config.Instances) == 0 {
		return nil, errors.ErrMissInstance
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
			// Contract
			StatusUpdater: action.StatusUpdaterConstructor(config),
			ErrorHandler:  config.ErrorHandler,
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

	return &Client{
		PayClient: payClients,
		Config:    config,
	}, nil
}
