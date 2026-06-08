package client

import (
	"encoding/json"

	"go.gh.ink/payutils/v3/errors"
	"go.gh.ink/payutils/v3/internal/action"
	"go.gh.ink/payutils/v3/internal/state"
	"go.gh.ink/payutils/v3/model"
	"go.gh.ink/payutils/v3/router"
	"go.gh.ink/toolbox/expr"
	"go.gh.ink/toolbox/pointer"
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

	// Check endpoint
	if config.Endpoint == "" {
		return nil, errors.ErrMissEndpoint
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
		payClients[upstreamName], err = upstream.NewClient(
			model.PayDriverClientParam{
				// Debug
				Debug: config.Debug,
				// Pay client credential
				Credential: upstreamCredential,
				// Other customised settings
				Endpoint:      config.Endpoint,
				TradeIDPrefix: config.TradeIDPrefix,
				TradeIDSuffix: config.TradeIDSuffix,
				NoNewPaymentWindows: expr.Ternary(
					config.NoNewPaymentWindows != nil,
					pointer.SafeDeref(config.NoNewPaymentWindows),
					model.NoNewPaymentWindows,
				),
				SafetyMargin: expr.Ternary(
					config.SafetyMargin != nil,
					pointer.SafeDeref(config.SafetyMargin),
					model.SafetyMargin,
				),
				// Contract
				StatusUpdater: action.StatusUpdaterConstructor(config),
				ErrorHandler:  config.ErrorHandler,
				// JSON
				Marshal:   config.Marshal,
				Unmarshal: config.Unmarshal,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// Register http routers.
	//
	// Instances are optional: a user may either import an http driver and let
	// payutils register the callback route automatically, or omit Instances and
	// invoke (*Client).Callback manually with a standard *http.Request.
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

			// Register callback router only. Create is no longer a route; it is
			// exposed as (*Client).Create for the user to call directly.
			unifiedInstance.Post(router.Route(upstreamName, "callback"), payClient.Callback)
		}
	}

	return &Client{
		PayClient: payClients,
		Config:    config,
	}, nil
}
