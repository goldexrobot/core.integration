package api

import (
	"sync/atomic"

	apiv1 "github.com/goldexrobot/core.integration/terminal/api/v1"
)

type RPC struct {
	api     *apiv1.Impl
	pending *int32
}

func (r *RPC) add() {
	atomic.AddInt32(r.pending, 1)
}

func (r *RPC) sub() {
	atomic.AddInt32(r.pending, -1)
}

func (r *RPC) reset(a *apiv1.Impl) bool {
	if atomic.LoadInt32(r.pending) == 0 {
		r.api = a
		return true
	}
	return false
}

func (r *RPC) InletOpen(_ interface{}, _ *interface{}) (err error) {
	r.add()
	defer r.sub()
	return r.api.InletOpen()
}

func (r *RPC) InletClose(_ interface{}, _ *interface{}) (err error) {
	r.add()
	defer r.sub()
	return r.api.InletClose()
}

func (r *RPC) OutletClose(_ interface{}, _ *interface{}) (err error) {
	r.add()
	defer r.sub()
	return r.api.OutletClose()
}

func (r *RPC) EvalNew(_ interface{}, res *apiv1.EvalNewResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.EvalNew()
	*res = v
	return
}

func (r *RPC) EvalSpectrum(_ interface{}, res *apiv1.EvalSpectrumResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.EvalSpectrum()
	*res = v
	return
}

func (r *RPC) EvalHydro(_ interface{}, res *apiv1.EvalHydroResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.EvalHydro()
	*res = v
	return
}

func (r *RPC) EvalReturn(_ interface{}, _ *interface{}) (err error) {
	r.add()
	defer r.sub()
	return r.api.EvalReturn()
}

func (r *RPC) EvalStore(req *apiv1.EvalStoreRequest, res *apiv1.EvalStoreResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.EvalStore(*req)
	*res = v
	return
}

func (r *RPC) StorageExtract(req *apiv1.StorageExtractRequest, res *apiv1.StorageExtractResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.StorageExtract(*req)
	*res = v
	return
}

func (r *RPC) Status(_ interface{}, res *apiv1.StatusResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.Status()
	*res = v
	return
}

func (r *RPC) Backend(req *apiv1.BackendRequest, res *apiv1.BackendResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.Backend(*req)
	*res = v
	return
}

func (r *RPC) Hardware(req *apiv1.HardwareRequest, res *apiv1.HardwareResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.Hardware(*req)
	*res = v
	return
}

func (r *RPC) CameraFrontal(_ interface{}, res *apiv1.CameraResult) (err error) {
	r.add()
	defer r.sub()
	v, err := r.api.CameraFrontal()
	*res = v
	return
}
