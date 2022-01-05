package api

import (
	"sync"

	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
	"github.com/goldexrobot/core.integration/terminal/cmd/api-emulator/backend"
	"github.com/sirupsen/logrus"
)

const hwBusinessMultDivider = 1000

// Implements bot emulator console commands and interfaces required by terminal API
type Controller struct {
	logger  *logrus.Entry
	rpc     *RPC
	backend backend.Backender

	hwBusinessMult *uint64

	flagModuleBroken                 *int32
	flagHardwareFailure              *int32
	flagNetworkFailure               *int32
	flagStorageNoRoom                *int32
	flagStorageAccessForbidden       *int32
	flagRejectEval                   *int32
	flagUnstableScale                *int32
	flagEvaluationAlloySilver        *int32
	flagEvaluationFinenessMillesimal *int32

	evalDataMutex sync.Mutex
	evalData      evalData
}

type evalData struct {
	EvalID    uint64
	Cell      string
	Spectrum  map[string]float64
	DryWeight float64
	WetWeight float64
}

func NewController(b backend.Backender, logger *logrus.Entry) *Controller {
	hwbm := new(uint64)
	*hwbm = 1 * hwBusinessMultDivider
	c := &Controller{
		logger:                           logger,
		backend:                          b,
		hwBusinessMult:                   hwbm,
		flagModuleBroken:                 new(int32),
		flagHardwareFailure:              new(int32),
		flagNetworkFailure:               new(int32),
		flagStorageNoRoom:                new(int32),
		flagStorageAccessForbidden:       new(int32),
		flagRejectEval:                   new(int32),
		flagUnstableScale:                new(int32),
		flagEvaluationAlloySilver:        new(int32),
		flagEvaluationFinenessMillesimal: new(int32),
	}
	c.rpc = &RPC{
		api:     c.apiImpl(),
		pending: new(int32),
	}

	return c
}

// Returns RPC handler
func (c *Controller) RPC() *RPC {
	return c.rpc
}

// Creates a new API implementation instance (called on APi reset)
func (c *Controller) apiImpl() *apiv1.Impl {
	return apiv1.NewImpl(c, c, c.logger.WithField("api", "impl"))
}
