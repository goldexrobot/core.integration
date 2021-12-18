package api

// JSONRPC API
type API interface {
	// Opens inlet window.
	OpenInlet() (err error)
	// Closes inlet window.
	CloseInlet() (err error)
	// Closes outlet window.
	CloseOutlet() (err error)

	// Begins spectral evaluation of an item.
	SpectralEvaluation() (res SpectralEvaluationResult, err error)
	// Begins hydrostatic evaluation of an item.
	HydrostaticEvaluation() (res EvaluationResult, err error)

	// Returns the item back to the customer after successful spectral or hydrostatic evaluation. Opens the outlet window automatically (should be closed manually).
	EvaluationCancel() (err error)
	// Moves successfully and fully evaluated item to the storage.
	EvaluationStore() (err error)

	// Requires storage to extract an item from the specified cell. Opens the outlet window automatically (should be closed manually).
	StorageExtraction() (err error)
}

type SpectralEvaluationResult struct {
	Alloy      string
	Purity     float64
	Millesimal int
	Carat      string
}

type EvaluationResult struct {
	Alloy      string
	Purity     float64
	Millesimal int
	Carat      string
	Weight     float64
	Confidence float64
	Risky      bool
}
