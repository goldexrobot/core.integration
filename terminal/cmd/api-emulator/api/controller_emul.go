package api

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

const hwBusinessMultDivider = 1000

type Controller struct {
	logger *logrus.Entry

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
	return &Controller{
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
}

func (c *Controller) HealAPI() {
	atomic.StoreInt32(c.flagModuleBroken, 0)
	atomic.StoreInt32(c.flagHardwareFailure, 0)
	atomic.StoreInt32(c.flagNetworkFailure, 0)
	atomic.StoreInt32(c.flagStorageNoRoom, 0)
	atomic.StoreInt32(c.flagStorageAccessForbidden, 0)
	atomic.StoreInt32(c.flagRejectEval, 0)
	atomic.StoreInt32(c.flagUnstableScale, 0)
}

func (c *Controller) SetHardwareBusinessMult(m float64) {
	if m < 0 {
		m = 0
	}
	if m > 10 {
		m = 10
	}
	atomic.StoreUint64(c.hwBusinessMult, uint64(hwBusinessMultDivider*m))
}

func (c *Controller) ToggleFailOnHardwareAccess() bool {
	i := c.flagHardwareFailure
	v := 1 - atomic.LoadInt32(i)
	atomic.StoreInt32(i, v)
	return v == 1
}

func (c *Controller) ToggleFailOnNetworkAccess() bool {
	i := c.flagNetworkFailure
	v := 1 - atomic.LoadInt32(i)
	atomic.StoreInt32(i, v)
	return v == 1
}

func (c *Controller) ToggleFailOnStorageRoomCheck() bool {
	i := c.flagStorageNoRoom
	v := 1 - atomic.LoadInt32(i)
	atomic.StoreInt32(i, v)
	return v == 1
}

func (c *Controller) ToggleFailOnStorageAccess() bool {
	i := c.flagStorageAccessForbidden
	v := 1 - atomic.LoadInt32(i)
	atomic.StoreInt32(i, v)
	return v == 1
}

func (c *Controller) ToggleEvalRejection() bool {
	i := c.flagRejectEval
	v := 1 - atomic.LoadInt32(i)
	atomic.StoreInt32(i, v)
	return v == 1
}

func (c *Controller) ToggleWeighingScaleUnstable() bool {
	i := c.flagUnstableScale
	v := 1 - atomic.LoadInt32(i)
	atomic.StoreInt32(i, v)
	return v == 1
}

// ---

func (c *Controller) accessHardware(timeUnits uint) error {
	if timeUnits > 0 {
		m := atomic.LoadUint64(c.hwBusinessMult)
		if m > 0 {
			d := time.Duration(timeUnits) * time.Second * time.Duration(m) / time.Duration(hwBusinessMultDivider)
			<-time.After(d)
		}
	}
	if atomic.LoadInt32(c.flagHardwareFailure) != 0 {
		return errors.New("emulated hardware failure")
	}
	return nil
}

func (c *Controller) accessNetwork() bool {
	return atomic.LoadInt32(c.flagNetworkFailure) == 0
}

func (c *Controller) generateEvaluationData() {
	c.evalDataMutex.Lock()
	defer c.evalDataMutex.Unlock()

	rand.Seed(time.Now().UnixMilli())

	var alloy = "au"
	if rand.Int()%2 == 0 {
		alloy = "ag"
	}
	var purity = math.Floor(math.Min(0.9999, 0.375+rand.Float64()*0.625)*10000) / 100

	var millesimal int
	switch {
	case purity >= 99.90:
		millesimal = 9999
	case purity >= 97.00:
		millesimal = 999
	case purity >= 90.00:
		millesimal = 925
	case purity >= 70.00:
		millesimal = 750
	case purity >= 50.00:
		millesimal = 585
	case purity >= 37.00:
		millesimal = 375
	}

	var carat = fmt.Sprintf("%dK", int(math.Ceil(purity*24/100)))

	var confidence = math.Floor((0.5+rand.Float64()*0.5)*1000) / 1000

	c.evalData = randomEvalData{
		Spectrum: map[string]float64{
			alloy: purity,
		},
		Alloy:      alloy,
		Purity:     purity,
		Millesimal: millesimal,
		Carat:      carat,
		Weight:     rand.Float64(),
		Confidence: confidence,
		Risky:      confidence < 0.88,
	}
}
