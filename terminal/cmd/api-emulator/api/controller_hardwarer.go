package api

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"time"

	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
	"github.com/google/uuid"
)

var storageCellsOrder = []string{
	"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9",
	"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9",
	"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9",
	"D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9",
	"E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8", "E9",
	"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9",
	"G1", "G2", "G3", "G4", "G5", "G6", "G7", "G8", "G9",
	"H1", "H2", "H3", "H4", "H5", "H6", "H7", "H8", "H9",
	"I1", "I2", "I3", "I4", "I5", "I6", "I7", "I8", "I9",
	"J1", "J2", "J3", "J4", "J5", "J6", "J7", "J8", "J9",
}

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// get occupied cells on backend side
	occupiedCells, xerr := c.backend.OccupiedCells(ctx)
	if xerr != nil {
		c.logger.WithError(xerr).Errorf("Failed calling backend to get occupied cells")
		failNet = true
		return
	}

	// next non-occupied cell
	for _, c := range storageCellsOrder {
		occupied := false
		for _, domainCells := range occupiedCells {
			for _, cc := range domainCells {
				if c == cc {
					occupied = true
					break
				}
			}
			if occupied {
				break
			}
		}
		if !occupied {
			cell = c
			break
		}
	}
	if cell == "" {
		failNoRoom = true
		return
	}

	// call backend
	evalID, xerr = c.backend.NewEval(ctx)
	if xerr != nil {
		c.logger.WithError(xerr).Errorf("Failed calling backend to begin a new evaluation")
		failNet = true
		return
	}

	var alloy = "au"
	if atomic.LoadInt32(c.flagEvaluationAlloySilver) != 0 {
		alloy = "ag"
	}
	c.generateEvaluationData(evalID, cell, alloy, int(atomic.LoadInt32(c.flagEvaluationFinenessMillesimal)))
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
	defer c.evalDataMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// call backend
	fin, rejectedEval, xerr := c.backend.EvaluateSpectrum(ctx, c.evalData.EvalID, c.evalData.Spectrum)
	if xerr != nil {
		c.logger.WithError(xerr).Errorf("Failed calling backend to evaluate spectral data")
		netFail = true
		return
	}
	if rejectedEval {
		return
	}

	eval = apiv1.ImplSpectralData{
		Alloy:      fin.Alloy,
		Purity:     fin.Purity,
		Millesimal: fin.Millesimal,
		Carat:      fin.Carat,
		Spectrum:   c.evalData.Spectrum,
	}
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
	defer c.evalDataMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// call backend
	rejectedEval, xerr := c.backend.EvaluateHydro(ctx, c.evalData.EvalID, c.evalData.DryWeight, c.evalData.WetWeight)
	if xerr != nil {
		c.logger.WithError(xerr).Errorf("Failed calling backend to evaluate hydrostatic data")
		netFail = true
		return
	}
	if rejectedEval {
		return
	}

	eval = apiv1.ImplHydroData{
		DryWeight: c.evalData.DryWeight,
		WetWeight: c.evalData.WetWeight,
	}
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
	defer c.evalDataMutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// call backend
	fin, rejectedEval, xerr := c.backend.FinalizeEvaluation(ctx, c.evalData.EvalID)
	if xerr != nil {
		c.logger.WithError(xerr).Errorf("Failed calling backend to finalize evaluation data")
		netFail = true
		return
	}
	if rejectedEval {
		return
	}

	fineness = apiv1.ImplFinenessData{
		Alloy:      fin.Alloy,
		Purity:     fin.Purity,
		Millesimal: fin.Millesimal,
		Carat:      fin.Carat,
		Weight:     c.evalData.DryWeight,
		Confidence: fin.Confidence,
		Risky:      fin.Risky,
	}
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

	c.evalDataMutex.Lock()
	defer c.evalDataMutex.Unlock()
	cell = c.evalData.Cell
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// call backend
	forbidden, reason, xerr := c.backend.OccupyStorageCell(ctx, domain, cell, tx)
	if xerr != nil {
		c.logger.WithError(xerr).Errorf("Failed calling backend to occupy storage cell %v under domain %v", cell, domain)
		netFail = true
		return
	}
	if forbidden {
		c.logger.Errorf("Forbidden to occupy storage cell %v under domain %v: %v", cell, domain, reason)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// call backend
	forbidden, reason, xerr := c.backend.ReleaseStorageCell(ctx, domain, cell, tx, false)
	if xerr != nil {
		c.logger.WithError(xerr).Errorf("Failed calling backend to release storage cell %v under domain %v", cell, domain)
		netFail = true
		return
	}
	if forbidden {
		c.logger.Errorf("Forbidden to release storage cell %v under domain %v: %v", cell, domain, reason)
		return
	}

	return
}

func (c *Controller) IntegrationUIMethod(method string, body map[string]interface{}) (httpStatus int, response map[string]interface{}, err error) {
	if !c.accessNetwork() {
		err = errors.New("network failure")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	// call backend
	response, httpStatus, err = c.backend.IntegrationUIMethod(ctx, method, body)
	if err != nil {
		c.logger.WithError(err).Errorf("Failed calling backend UI method %q", method)
		return
	}
	return
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

func (c *Controller) UploadFrontalCameraPhotoForEval() {}

func (c *Controller) UploadFrontalCameraPhoto() (fileID string, err error) {
	return strings.ReplaceAll(uuid.NewString(), "-", ""), nil
}

func (c *Controller) InternetConnectivity() (ok bool) {
	return c.accessNetwork()
}

func (c *Controller) HasStorage() (ok bool) {
	return true
}

func (c *Controller) HasPositionalStorage() (ok bool) {
	return true
}
