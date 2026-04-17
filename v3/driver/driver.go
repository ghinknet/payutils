package driver

import (
	"github.com/ghinknet/payutils/v3/internal/state"
	"github.com/ghinknet/payutils/v3/model"
)

func RegisterPay(name string, driver model.PayDriver) {
	state.PayDrivers[name] = driver
}

func RegisterHttp(name string, driver model.HttpDriver) {
	state.HttpDrivers[name] = driver
}
