// Goldex backend callbacks.
//
// Goldex can send a number of `POST` callbacks (`application/json`) to your backend.
// Some callbacks are `synchronous`, which means if the callback is defined it `must` respond with HTTP status `200`.
// Other callbacks are `asynchronous`: Goldex queues this callback and periodically tries to deliver it, until get HTTP status `200`.
//
// Callbacks (see Goldex dashboard for detailed list):
// 1. Item evaluation just started (sync);
// 2. Item evaluation finished successfully (async);
// 3. Storage cell occupation attempt (sync);
// 4. Storage cell release attempt (sync);
// 5. Custom named methods accessable for terminal UI (sync).
//
// Goldex sends HTTP headers with each callback:
// 1. `X-CBOT-PROJECT-ID` contains Goldex-world project ID;
// 2. `X-CBOT-BOT-ID` contains Goldex-world bot ID;
//
// Moreover, Goldex signs the callbacks, so you are able to verify it with a per-project public key (see Goldex dashboard).
//
// Actual callback payloads are described below in this document.
//
// swagger:meta
package callback

// Describes new evaluation that is being just started
//
// swagger:model
type EvalStarted struct {
	// Project ID
	//
	// example: 1
	ProjectID uint64 `json:"project_id"`
	// Bot ID
	//
	// example: 42
	BotID uint64 `json:"bot_id"`
	// Evaluation ID
	//
	// example: 1337
	EvalID uint64 `json:"eval_id"`
}

// Describes evaluation and decision data and is sent only on successful evaluation completion
//
// swagger:model
type EvalFinished struct {
	// Project ID
	//
	// example: 1
	ProjectID uint64 `json:"project_id"`
	// Bot ID
	//
	// example: 42
	BotID uint64 `json:"bot_id"`
	// Evaluation ID
	//
	// example: 1337
	EvalID uint64 `json:"eval_id"`

	// Detected primary metal/alloy
	//
	// example: "au", "ag"
	Alloy string `json:"alloy"`
	// Content of the metal in spectrum in percents
	//
	// example: 99.8
	AlloyContent float64 `json:"alloy_content"`
	// Spectrum
	//
	// example: {"au":99.8}
	Spectrum map[string]float64 `json:"spectrum"`
	// Item weight
	//
	// example: 3.141
	Weight float64 `json:"weight"`

	// Photos available
	Photo []EvalFinishedPhoto `json:"photo"`
	// Overall result confidence/score [0..1], where 1 - is fully confident result, and value below 0.88 is alarming.
	//
	// example: 0.889
	Confidence float64 `json:"confidence"`
	// Detected fineness purity in percents
	//
	// example: 99.9
	FinenessPurity float64 `json:"fineness_purity"`
	// Detected millesimal fineness: 585 stands for 58.5%, 999 => 99.9%, 9999 => 99.99%
	//
	// example: 9999
	FinenessMillesimal int `json:"fineness_millesimal"`
	// Detected fineness in carats
	//
	// example: "24K"
	FinenessCarat string `json:"fineness_carat"`
	// System decision about the evaluation
	//
	// example: false
	Risky bool `json:"risky"`

	// Warnings that should help with decision. For instance, there could be tungsten covered with gold.
	//
	// example: ["tungsten_in_gold"]
	Warnings []string `json:"warnings"`
}

// Contains evaluation photo details
//
// swagger:model
type EvalFinishedPhoto struct {
	// File ID
	//
	// example: 3dd321739a694bbab93d7aae360a4ab4
	PhotoID string `json:"photo_id"`
	// File ID
	//
	// example: eef30f5dc98e4c7d8d2f8df9df56c0d0
	PreviewID string `json:"preview_id"`
	// An origin photo comes from. Here "item" is item photo, "outer" is external camera
	//
	// example: "item", "outer"
	Origin string `json:"origin"`
}

// This is common data sent in a storage related callbacks
//
// swagger:model
type StorageCellEvent struct {
	// Project ID
	//
	// example: 1
	ProjectID uint64 `json:"project_id"`
	// Bot ID
	//
	// example: 42
	BotID uint64 `json:"bot_id"`
	// Cell address
	//
	// example: "A1", "J9"
	Cell string `json:"cell"`
	// Origin of the event in terms of UI flow. Here "dashboard" is on-bot system dashboard, "other" is some custom origin.
	//
	// example: "other", "dashboard", "buyout", "shop", "pawnshop"
	Domain string `json:"domain"`
}
