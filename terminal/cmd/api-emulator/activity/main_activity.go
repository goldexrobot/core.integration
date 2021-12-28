package activity

import (
	"github.com/chzyer/readline"
	"github.com/goldexrobot/core.integration/terminal/cmd/api-emulator/console"
	"github.com/sirupsen/logrus"
)

type Main struct {
	Logger *logrus.Entry
	Ctl    Controller
}

type Controller interface {
	ResetAPI() bool
	HealAPI()
	SetHardwareBusinessMult(m float64)
	ToggleFailOnHardwareAccess() bool
	ToggleFailOnNetworkAccess() bool
	ToggleFailOnStorageRoomCheck() bool
	ToggleFailOnStorageAccess() bool
	ToggleEvalRejection() bool
	ToggleWeighingScaleUnstable() bool
}

func (a *Main) Create(ctx console.Context) error {
	return nil
}

func (a *Main) Destroy(ctx console.Context) {
}

func (a *Main) Activate(ctx console.Context) {
}

func (a *Main) Deactivate(ctx console.Context) {
}

func (a *Main) Completer() ([]readline.PrefixCompleterInterface, []string) {
	return []readline.PrefixCompleterInterface{
		console.DPCItem(
			"break", "Toggle terminal emulation features:",
			console.DPCItem("hardware", "Emulate critical hardware failures"),
			console.DPCItem("network", "Emulate network outage"),
			console.DPCItem(
				"storage", "Storage features:",
				console.DPCItem("access", "Emulate storage access prohibition"),
				console.DPCItem("room", "Emulate lack of storage room"),
			),
			console.DPCItem(
				"eval", "Evaluation features:",
				console.DPCItem("acceptance", "Emulate an evaluating item rejection"),
				console.DPCItem("scale", "Emulate a mechanical vibration affecting weighing process"),
			),
		),
		console.DPCItem("delay", "Set hardware business delay multiplier [0.0,10.0]"),
		console.DPCItem("heal", "Set API fully healthy"),
		console.DPCItem("reset", "Reset API to initial state"),
	}, nil
}

// Command ...
func (a *Main) Command(ctx console.Context, args console.Args) bool {
	switch args.Args[0] {

	case "break":
		what, _ := args.String(1)
		switch what {
		case "hardware":
			if a.Ctl.ToggleFailOnHardwareAccess() {
				ctx.Warningf("Hardware failures now are enabled")
			} else {
				ctx.Infof("Hardware now is healthy")
			}
			return true
		case "network":
			if a.Ctl.ToggleFailOnNetworkAccess() {
				ctx.Warningf("Network outage now is enabled")
			} else {
				ctx.Infof("Network now is available")
			}
			return true
		case "storage":
			how, _ := args.String(2)
			switch how {
			case "access":
				if a.Ctl.ToggleFailOnStorageAccess() {
					ctx.Warningf("Storage access now is forbidden")
				} else {
					ctx.Infof("Storage access now is allowed")
				}
				return true
			case "room":
				if a.Ctl.ToggleFailOnStorageRoomCheck() {
					ctx.Warningf("Storage now has no more room")
				} else {
					ctx.Infof("Storage now has enough room")
				}
				return true
			}
		case "eval":
			how, _ := args.String(2)
			switch how {
			case "acceptance":
				if a.Ctl.ToggleEvalRejection() {
					ctx.Warningf("Evaluated item now will be rejected")
				} else {
					ctx.Infof("Evaluated item now will be accepted")
				}
				return true
			case "scale":
				if a.Ctl.ToggleWeighingScaleUnstable() {
					ctx.Warningf("Emulating unstable weighing scale result (hydrostatic evaluation)")
				} else {
					ctx.Infof("Weighing scale now is stable (hydrostatic evaluation)")
				}
				return true
			}
		}
	case "heal":
		a.Ctl.HealAPI()
		ctx.Infof("API now is healthy and reset")
		return true
	case "delay":
		v, ok := args.Float64(1)
		if !ok {
			a.Ctl.SetHardwareBusinessMult(1)
			ctx.Infof("Hardware business delay multiplier set to 1")
		} else {
			a.Ctl.SetHardwareBusinessMult(v)
			ctx.Infof("Hardware business delay multiplier set to %v", v)
		}
		return true
	case "reset":
		if a.Ctl.ResetAPI() {
			ctx.Infof("API state has been reset")
		} else {
			ctx.Errorf("API has pending requests")
		}
		return true
	}
	return false
}
