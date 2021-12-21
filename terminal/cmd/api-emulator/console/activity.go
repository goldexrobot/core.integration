package console

import (
	"github.com/chzyer/readline"
)

// Activity handles activity
type Activity interface {
	// Create fires on activity creation
	Create(ctx Context) error
	// Destroy fires on activity destruction
	Destroy(ctx Context)
	// Activate fires on activity activation
	Activate(ctx Context)
	// Deactivate fires on activity deactivation
	Deactivate(ctx Context)
	// Completer returns available commands
	Completer() (pci []readline.PrefixCompleterInterface, helpPostfix []string)
	// Command consumes commands
	Command(ctx Context, args Args) (consumed bool)
}
