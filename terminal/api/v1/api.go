// Terminal JSONRPC API.
//
// Goldex Robot terminal serves JSONRPC 2 API and accepts Websocket connections on localhost\:80\/ws. Websocket messages are textual, not binary.
//
// How to read the document.
// This doc is generated according to OpenAPI 2, so JSONRPC methods are defined as POST methods. Actual method name is defined after slash symbol.
// For example\:
//
// `"POST /inlet.open"`
//
// means JSONRPC request
//
// `{"version":"2.0","id":1,"method":"inlet.open","params":{...}}`
//
//     Schemes: ws
//     Host: localhost:80
//     BasePath: /ws
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package api

// JSONRPC API
//
type API interface {

	////// INLET/OUTLET WINDOW /////

	// swagger:operation POST /inlet.open InletOpen
	//
	// Open inlet window.
	//
	// Requires hardware to open inlet window. Should be called to receive a customer item before evaluation.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Inlet/outlet window
	// responses:
	//   x-jsonrpc-success:
	//     description: No payload
	//   default:
	//     description: JSONRPC error
	InletOpen() (err error)

	// swagger:operation POST /inlet.close InletClose
	//
	// Close inlet window.
	//
	// Requires hardware to close inlet window. Should be called right before evaluation launch.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Inlet/outlet window
	// responses:
	//   x-jsonrpc-success:
	//     description: No payload
	//   default:
	//     description: JSONRPC error
	InletClose() (err error)

	// swagger:operation POST /outlet.close OutletClose
	//
	// Close outlet window.
	//
	// Requires hardware to close outlet window. Should be called manually after customer item return or storage item extraction.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Inlet/outlet window
	// responses:
	//   x-jsonrpc-success:
	//     description: No payload
	//   default:
	//     description: JSONRPC error
	OutletClose() (err error)

	///// EVALUATION //////

	// swagger:operation POST /eval.new EvalNew
	//
	// New evaluation [I].
	//
	// Prepares a new evaluation operation: check hardware, notify backend server, etc.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Evaluation
	// responses:
	//   x-jsonrpc-success:
	//     description: New evaluation ID or a failure
	//     scheme:
	//       "$ref": "#/definitions/EvalNewResult"
	//   default:
	//     description: JSONRPC error
	EvalNew() (res EvalNewResult, err error)

	// swagger:operation POST /eval.spectrum EvalSpectrum
	//
	// Spectral evaluation [II].
	//
	// Starts a spectral evaluation of the item. Should be called right after `eval.new`.
	// On successful spectral evaluation the item might be returned back to customer with `eval.return`, otherwise the evaluation should be continued with `eval.hydro`.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Evaluation
	// responses:
	//   x-jsonrpc-success:
	//     description: Spectral evaluation result
	//     scheme:
	//       "$ref": "#/definitions/EvalSpectrumResult"
	//   default:
	//     description: JSONRPC error
	EvalSpectrum() (res EvalSpectrumResult, err error)

	// swagger:operation POST /eval.hydro EvalHydro
	//
	// Hydrostatic evaluation [III].
	//
	// Starts a hydrostatic evaluation of the item. Should be called right after `eval.spectrum`.
	// On successful hydrostatic evaluation the item might be returned back to customer with `eval.return`.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Evaluation
	// responses:
	//   x-jsonrpc-success:
	//     description: Hydrostatic evaluation result
	//     scheme:
	//       "$ref": "#/definitions/EvalHydroResult"
	//   default:
	//     description: JSONRPC error
	EvalHydro() (res EvalHydroResult, err error)

	// swagger:operation POST /eval.return EvalReturn
	//
	// Return item [IV].
	//
	// Starts a returning process of the item. Should be called after spectral/hydrostatic evaluation.
	// On successful returning outlet window should be closed manually: customer choice (preferred) or a timeout.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Evaluation
	// responses:
	//   x-jsonrpc-success:
	//     description: No payload
	//   default:
	//     description: JSONRPC error
	EvalReturn() (err error)

	// swagger:operation POST /eval.store EvalStore
	//
	// Store item [IV].
	//
	// Requires hardware to transfer successfully evaluated item into the internal storage.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Evaluation
	// responses:
	//   x-jsonrpc-success:
	//     description: Success
	//     scheme:
	//       "$ref": "#/definitions/EvalStoreResult"
	//   default:
	//     description: JSONRPC error
	EvalStore(req EvalStoreRequest) (res EvalStoreResult, err error)

	////// STORAGE //////

	// swagger:operation POST /storage.extract StorageExtract
	//
	// Extract item.
	//
	// Requires hardware to extract an item from the specified storage cell and bring it to the outlet window.
	// On successful extraction the outlet window should be closed manually: customer choice (preferred) or a timeout.
	// ---
	// consumes:
	// - application/json
	// produces:
	// - application/json
	// tags:
	//   - Storage
	// responses:
	//   x-jsonrpc-success:
	//     description: No payload
	//   default:
	//     description: JSONRPC error
	StorageExtract(req StorageExtractRequest) (err error)
}

// New evaluation result.
// swagger:model
type EvalNewResult struct {
	Success EvalNewResultSuccess
	Failure EvalNewResultFailure
}

type EvalNewResultSuccess struct {
	StorageCell string
}

type EvalNewResultFailure struct {
	NetworkUnavailable bool
	HardwareCheck      bool
	NoStorageRoom      bool
}

// Spectral evaluation result.
// swagger:model
type EvalSpectrumResult struct {
	Success EvalSpectrumResultSuccess
	Failure EvalSpectrumResultFailure
}

type EvalSpectrumResultSuccess struct {
	Alloy      string
	Purity     float64
	Millesimal int
	Carat      string
	Spectrum   map[string]float64
}

type EvalSpectrumResultFailure struct {
	NetworkUnavailable bool
	EvalRejected       bool
}

// Hydrostatic evaluation result.
// swagger:model
type EvalHydroResult struct {
	Success EvalHydroResultSuccess
	Failure EvalHydroResultFailure
}

type EvalHydroResultSuccess struct {
	Alloy      string
	Purity     float64
	Millesimal int
	Carat      string
	Weight     float64
	Confidence float64
	Risky      bool
}

type EvalHydroResultFailure struct {
	NetworkUnavailable bool
	EvalRejected       bool
	UnstableScale      bool
}

// Item storing request.
// swagger:model
type EvalStoreRequest struct {
	// TODO:
}

// Item storing result.
// swagger:model
type EvalStoreResult struct {
	// TODO:
}

// Item extraction request.
// swagger:model
type StorageExtractRequest struct {
	// TODO:
}
