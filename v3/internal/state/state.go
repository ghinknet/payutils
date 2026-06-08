package state

import "go.gh.ink/payutils/v3/model"

var PayDrivers = make(map[string]model.PayDriver)
var HttpDrivers = make(map[string]model.HttpDriver)
