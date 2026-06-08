package driver

import (
	"go.gh.ink/payutils/v3/internal/state"
	"go.gh.ink/payutils/v3/model"
)

func RegisterPay(name string, driver model.PayDriver) {
	state.PayDrivers[name] = driver
}

func RegisterHttp(name string, driver model.HttpDriver) {
	state.HttpDrivers[name] = driver
}
