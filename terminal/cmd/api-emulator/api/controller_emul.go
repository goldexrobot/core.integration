package api

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync/atomic"
	"time"
)

func (c *Controller) ResetAPI() bool {
	return c.rpc.reset(c.apiImpl())
}

func (c *Controller) HealAPI() {
	atomic.StoreInt32(c.flagModuleBroken, 0)
	atomic.StoreInt32(c.flagHardwareFailure, 0)
	atomic.StoreInt32(c.flagNetworkFailure, 0)
	atomic.StoreInt32(c.flagStorageNoRoom, 0)
	atomic.StoreInt32(c.flagStorageAccessForbidden, 0)
	atomic.StoreInt32(c.flagRejectEval, 0)
	atomic.StoreInt32(c.flagUnstableScale, 0)
	c.ResetAPI()
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

	var weight = math.Floor((1.0+(10*rand.Float64()))*1000) / 1000

	c.evalData = randomEvalData{
		Spectrum: map[string]float64{
			alloy: purity,
		},
		Alloy:      alloy,
		Purity:     purity,
		Millesimal: millesimal,
		Carat:      carat,
		Weight:     weight,
		Confidence: confidence,
		Risky:      confidence < 0.88,
	}
}
