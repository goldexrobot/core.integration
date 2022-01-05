package api

import (
	"errors"
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

func (c *Controller) ToggleEvaluationAlloySilver() bool {
	i := c.flagEvaluationAlloySilver
	v := 1 - atomic.LoadInt32(i)
	atomic.StoreInt32(i, v)
	return v == 1
}

func (c *Controller) SetEvaluationFineness(millesimal int) {
	i := c.flagEvaluationFinenessMillesimal
	atomic.StoreInt32(i, int32(millesimal))
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

func (c *Controller) generateEvaluationData(evalID uint64, cell, alloy string, millesimal int) {
	c.evalDataMutex.Lock()
	defer c.evalDataMutex.Unlock()

	var (
		spectrum  float64
		dryWeight float64
		wetWeight float64
	)

	switch alloy {
	case "ag":
		switch millesimal {
		case 900:
			spectrum = 94.72
			dryWeight = 6.640
			wetWeight = 5.927
		case 960:
			spectrum = 95.53
			dryWeight = 20.011
			wetWeight = 18.085
		default:
			millesimal = 999
			spectrum = 98.64
			dryWeight = 31.282
			wetWeight = 28.270
		}
	default:
		switch millesimal {
		case 375:
			spectrum = 36.23
			dryWeight = 2.230
			wetWeight = 2.019
		case 500:
			spectrum = 53.27
			dryWeight = 3.974
			wetWeight = 3.640
		case 585:
			spectrum = 57.77
			dryWeight = 2.590
			wetWeight = 2.400
		case 750:
			spectrum = 82.43
			dryWeight = 4.130
			wetWeight = 3.820
		case 999:
			spectrum = 98.34
			dryWeight = 7.803
			wetWeight = 7.411
		default:
			millesimal = 9999
			spectrum = 99.97
			dryWeight = 7.803
			wetWeight = 7.411
		}
	}

	c.evalData = evalData{
		EvalID: evalID,
		Cell:   cell,
		Spectrum: map[string]float64{
			alloy: spectrum,
		},
		DryWeight: dryWeight,
		WetWeight: wetWeight,
	}
}
