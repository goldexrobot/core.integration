package api

import (
	"encoding/json"

	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
)

type RPC struct {
	api *apiv1.Impl
}

func (r *RPC) InletOpen(_ interface{}, _ *interface{}) (err error) {
	return r.api.InletOpen()
}

func (r *RPC) InletClose(_ interface{}, _ *interface{}) (err error) {
	return r.api.InletClose()
}

func (r *RPC) OutletClose(_ interface{}, _ *interface{}) (err error) {
	return r.api.OutletClose()
}

func (r *RPC) EvalNew(_ interface{}, res *apiv1.EvalNewResult) (err error) {
	v, err := r.api.EvalNew()
	*res = v
	return
}

func (r *RPC) EvalSpectrum(_ interface{}, res *apiv1.EvalSpectrumResult) (err error) {
	v, err := r.api.EvalSpectrum()
	*res = v
	return
}

func (r *RPC) EvalHydro(_ interface{}, res *apiv1.EvalHydroResult) (err error) {
	v, err := r.api.EvalHydro()
	*res = v
	return
}

func (r *RPC) EvalReturn(_ interface{}, _ *interface{}) (err error) {
	return r.api.EvalReturn()
}

func (r *RPC) EvalStore(req *apiv1.EvalStoreRequest, res *apiv1.EvalStoreResult) (err error) {
	v, err := r.api.EvalStore(*req)
	*res = v
	return
}

func (r *RPC) StorageExtract(req *apiv1.StorageExtractRequest, res *apiv1.StorageExtractResult) (err error) {
	v, err := r.api.StorageExtract(*req)
	*res = v
	return
}

func (r *RPC) Status(_ json.RawMessage, res *apiv1.StatusResult) (err error) {
	v, err := r.api.Status()
	*res = v
	return
}

func (r *RPC) Call(req *apiv1.CallRequest, res *apiv1.CallResult) (err error) {
	v, err := r.api.Call(*req)
	*res = v
	return
}

func (r *RPC) Hardware(req *apiv1.HardwareRequest, res *apiv1.HardwareResult) (err error) {
	v, err := r.api.Hardware(*req)
	*res = v
	return
}
