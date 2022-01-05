package backend

import "context"

// Provides backend connectivity
type Backender interface {
	OccupiedCells(ctx context.Context) (domains map[string][]string, err error)
	NewEval(ctx context.Context) (id uint64, err error)
	EvaluateSpectrum(ctx context.Context, evalID uint64, spectrum map[string]float64) (r ResultFineness, rejected bool, err error)
	EvaluateHydro(ctx context.Context, evalID uint64, dry, wet float64) (rejected bool, err error)
	FinalizeEvaluation(ctx context.Context, evalID uint64) (r ResultFineness, rejected bool, err error)
	OccupyStorageCell(ctx context.Context, domain, cell, tx string) (forbidden bool, reason string, err error)
	ReleaseStorageCell(ctx context.Context, domain, cell, tx string, strictDomainCheck bool) (forbidden bool, reason string, err error)
	IntegrationUIMethod(ctx context.Context, method string, kv map[string]interface{}) (result map[string]interface{}, httpStatus int, err error)
}
