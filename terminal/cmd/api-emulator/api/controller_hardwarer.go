package api

import (
	"errors"
	"sync/atomic"

	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
)

////// MODULE IMPL //////

func (c *Controller) Operational() (ok bool) {
	return atomic.LoadInt32(c.flagModuleBroken) == 0
}

func (c *Controller) Broken(err error) {
	atomic.StoreInt32(c.flagModuleBroken, 1)
	c.logger.WithError(err).Errorf("Module now is broken")
}

////// HARDWARE IMPL //////

func (c *Controller) OpenInlet() (err error) {
	return c.accessHardware(2)
}

func (c *Controller) CloseInlet() (err error) {
	return c.accessHardware(2)
}

func (c *Controller) CloseOutlet() (err error) {
	return c.accessHardware(4)
}

func (c *Controller) NewEval() (evalID uint64, cell string, failNet, failNoRoom, failHW bool, err error) {
	if c.accessHardware(1) != nil {
		failHW = true
		return
	}
	if !c.accessNetwork() {
		failNet = true
		return
	}
	if atomic.LoadInt32(c.flagStorageNoRoom) != 0 {
		failNoRoom = true
		return
	}
	cell = "A1"
	evalID = atomic.AddUint64(c.evalCounter, 1)
	c.generateEvaluationData()
	return
}

func (c *Controller) SpectralEval() (eval apiv1.ImplSpectralData, netFail, rejectedEval bool, err error) {
	if err = c.accessHardware(10); err != nil {
		return
	}
	if !c.accessNetwork() {
		netFail = true
		return
	}
	if atomic.LoadInt32(c.flagRejectEval) != 0 {
		rejectedEval = true
		return
	}
	c.evalDataMutex.Lock()
	eval = apiv1.ImplSpectralData{
		Alloy:      c.evalData.Alloy,
		Purity:     c.evalData.Purity,
		Millesimal: c.evalData.Millesimal,
		Carat:      c.evalData.Carat,
		Spectrum:   c.evalData.Spectrum,
	}
	c.evalDataMutex.Unlock()
	return
}

func (c *Controller) HydroEval() (eval apiv1.ImplHydroData, netFail, rejectedEval, unstableScale bool, err error) {
	if err = c.accessHardware(10); err != nil {
		return
	}
	if atomic.LoadInt32(c.flagUnstableScale) != 0 {
		unstableScale = true
		return
	}
	if !c.accessNetwork() {
		netFail = true
		return
	}
	if atomic.LoadInt32(c.flagRejectEval) != 0 {
		rejectedEval = true
		return
	}
	c.evalDataMutex.Lock()
	eval = apiv1.ImplHydroData{
		DryWeight: c.evalData.Weight,
		WetWeight: c.evalData.Weight,
	}
	c.evalDataMutex.Unlock()
	return
}

func (c *Controller) FinalizeEval() (fineness apiv1.ImplFinenessData, netFail, rejectedEval bool, err error) {
	if !c.accessNetwork() {
		netFail = true
		return
	}
	if atomic.LoadInt32(c.flagRejectEval) != 0 {
		rejectedEval = true
		return
	}
	c.evalDataMutex.Lock()
	fineness = apiv1.ImplFinenessData{
		Alloy:      c.evalData.Alloy,
		Purity:     c.evalData.Purity,
		Millesimal: c.evalData.Millesimal,
		Carat:      c.evalData.Carat,
		Weight:     c.evalData.Weight,
		Confidence: c.evalData.Confidence,
		Risky:      c.evalData.Risky,
	}
	c.evalDataMutex.Unlock()
	return
}

func (c *Controller) ReturnAfterSpectrumEval(customerChoice bool) (err error) {
	if err = c.accessHardware(6); err != nil {
		return
	}
	return
}

func (c *Controller) ReturnAfterHydroEval(customerChoice bool) (err error) {
	if err = c.accessHardware(6); err != nil {
		return
	}
	return
}

func (c *Controller) StoreAfterHydroEval() (cell string, err error) {
	if err = c.accessHardware(6); err != nil {
		return
	}
	cell = "A1"
	return
}

func (c *Controller) ExtractCellFromStorage(cell string) (err error) {
	if err = c.accessHardware(6); err != nil {
		return
	}
	return
}

func (c *Controller) StorageOccupyCell(cell, domain, tx string) (netFail, forbidden bool, err error) {
	if !c.accessNetwork() {
		netFail = true
		return
	}
	if atomic.LoadInt32(c.flagStorageAccessForbidden) != 0 {
		forbidden = true
		return
	}
	return
}

func (c *Controller) StorageReleaseCell(cell, domain, tx string) (netFail, forbidden bool, err error) {
	if !c.accessNetwork() {
		netFail = true
		return
	}
	if atomic.LoadInt32(c.flagStorageAccessForbidden) != 0 {
		forbidden = true
		return
	}
	return
}

func (c *Controller) IntegrationUIMethod(method string, body map[string]interface{}) (httpStatus int, response map[string]interface{}, err error) {
	if !c.accessNetwork() {
		err = errors.New("network failure")
		return
	}
	httpStatus = 200
	response = body
	return
}

func (c *Controller) UploadEvalCustomerPhoto() {
}

func (c *Controller) InternetConnectivity() (ok bool) {
	return c.accessNetwork()
}

func (c *Controller) OptionalHardwareHealthcheck() (health map[string]bool, err error) {
	health = map[string]bool{
		"my-pos-terminal": true,
		"my-printer":      true,
	}
	return
}

func (c *Controller) OptionalHardwareRPC(module, method string, request map[string]interface{}) (result map[string]interface{}, subError string, err error) {
	if xerr := c.accessHardware(3); xerr != nil {
		subError = xerr.Error()
		return
	}
	result = request
	return
}
