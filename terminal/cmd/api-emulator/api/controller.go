package api

import (
	"sync"

	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
	"github.com/sirupsen/logrus"
)

const hwBusinessMultDivider = 1000

type Controller struct {
	logger *logrus.Entry
	rpc    *RPC

	evalCounter *uint64

	hwBusinessMult *uint64

	flagModuleBroken           *int32
	flagHardwareFailure        *int32
	flagNetworkFailure         *int32
	flagStorageNoRoom          *int32
	flagStorageAccessForbidden *int32
	flagRejectEval             *int32
	flagUnstableScale          *int32

	evalDataMutex sync.Mutex
	evalData      randomEvalData
}

type randomEvalData struct {
	Spectrum   map[string]float64
	Alloy      string
	Purity     float64
	Millesimal int
	Carat      string
	Weight     float64
	Confidence float64
	Risky      bool
}

func NewController(logger *logrus.Entry) *Controller {
	hwbm := new(uint64)
	*hwbm = 1 * hwBusinessMultDivider
	c := &Controller{
		logger:                     logger,
		hwBusinessMult:             hwbm,
		evalCounter:                new(uint64),
		flagModuleBroken:           new(int32),
		flagHardwareFailure:        new(int32),
		flagNetworkFailure:         new(int32),
		flagStorageNoRoom:          new(int32),
		flagStorageAccessForbidden: new(int32),
		flagRejectEval:             new(int32),
		flagUnstableScale:          new(int32),
	}
	c.rpc = &RPC{
		api:     c.apiImpl(),
		pending: new(int32),
	}

	return c
}

func (c *Controller) RPC() *RPC {
	return c.rpc
}

func (c *Controller) apiImpl() *apiv1.Impl {
	return apiv1.NewImpl(c, c, c.logger.WithField("api", "impl"))
}
