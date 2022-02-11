// Goldex backend callbacks.
//
// swagger:meta
package callback

// A new evaluation is just started
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

// New photo is received from the terminal during the evaluation
//
// swagger:model
type EvalPhoto struct {
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
	// File ID
	//
	// example: 3dd321739a694bbab93d7aae360a4ab4
	PhotoID string `json:"photo_id"`
	// File ID
	//
	// example: eef30f5dc98e4c7d8d2f8df9df56c0d0
	PreviewID string `json:"preview_id"`
	// An origin the photo comes from. Here "item" is the item photo, "outer" is the external camera
	//
	// example: "item", "outer"
	Origin string `json:"origin"`
}

// Evaluation is cancelled by customer or failed
//
// swagger:model
type EvalCancelled struct {
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
	// Origin of the event in terms of UI flow. Here "dashboard" is on-bot system dashboard, "other" is some custom origin, "collection" is storage loading/unloading dashboard.
	//
	// example: "other", "dashboard", "buyout", "shop", "pawnshop", "collection"
	Domain string `json:"domain"`
	// Transaction ID as unique cell changing operation
	//
	// example: "240702e1e8be42d39137c1d802dd6941"
	Transaction string `json:"transaction"`
}
