package backend

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
)

type MockClient struct {
	evalCounter   *uint64
	lastEval      ResultFineness
	occupiedCells map[string]string
}

func NewMockClient() *MockClient {
	return &MockClient{
		evalCounter:   new(uint64),
		occupiedCells: make(map[string]string),
	}
}

func (c *MockClient) OccupiedCells(ctx context.Context) (domains map[string][]string, err error) {
	domains = make(map[string][]string)
	for cell, domain := range c.occupiedCells {
		domains[domain] = append(domains[domain], cell)
	}
	return
}

func (c *MockClient) NewEval(ctx context.Context) (id uint64, err error) {
	id = atomic.AddUint64(c.evalCounter, 1)
	return
}

func (c *MockClient) EvaluateSpectrum(ctx context.Context, evalID uint64, spectrum map[string]float64) (r ResultFineness, rejected bool, err error) {

	var (
		alloy  string
		purity float64
	)
	for k, v := range spectrum {
		alloy = k
		purity = v
	}

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

	c.lastEval = ResultFineness{
		Alloy:      alloy,
		Purity:     purity,
		Millesimal: millesimal,
		Carat:      carat,
		Confidence: 0.98,
		Risky:      false,
	}

	r = c.lastEval
	return
}

func (c *MockClient) EvaluateHydro(ctx context.Context, evalID uint64, dry, wet float64) (rejected bool, err error) {
	return
}

func (c *MockClient) FinalizeEvaluation(ctx context.Context, evalID uint64) (r ResultFineness, rejected bool, err error) {
	r = c.lastEval
	return
}

func (c *MockClient) OccupyStorageCell(ctx context.Context, domain, cell, tx string) (forbidden bool, reason string, err error) {
	if _, ok := c.occupiedCells[cell]; ok {
		forbidden = true
		reason = "already occupied"
		return
	}
	c.occupiedCells[cell] = domain
	return
}

func (c *MockClient) ReleaseStorageCell(ctx context.Context, domain, cell, tx string, strictDomainCheck bool) (forbidden bool, reason string, err error) {
	d, ok := c.occupiedCells[cell]
	if !ok {
		forbidden = true
		reason = "not occupied"
		return
	}

	if strictDomainCheck && domain != d {
		forbidden = true
		reason = "cell is occupied under another domain"
		return
	}

	delete(c.occupiedCells, cell)
	return
}

func (c *MockClient) IntegrationUIMethod(ctx context.Context, method string, kv map[string]interface{}) (result map[string]interface{}, httpStatus int, err error) {
	result = kv
	httpStatus = 200
	return
}
