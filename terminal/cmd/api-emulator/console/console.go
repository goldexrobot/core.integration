package console

import (
	"container/list"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

// Console handles user input (readline) and provides colorized output
type Console struct {
	*readline.Instance
	activities *list.List
}

// New instance
func New() (c *Console, cls func(), err error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryLimit:      256,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			// block CtrlZ feature
			if r == readline.CharCtrlZ {
				return r, false
			}
			return r, true
		},
	})
	if err != nil {
		return
	}

	c = &Console{
		Instance:   rl,
		activities: list.New(),
	}
	cls = func() {
		rl.Close()
	}
	return
}

// Printf prints!
func (c *Console) Printf(f string, args ...interface{}) {
	fmt.Fprintf(c.Stderr(), f+"\n", args...)
}

// Debugf prints!
func (c *Console) Debugf(f string, args ...interface{}) {
	c.Printf(color.HiBlackString(fmt.Sprintf(f, args...)))
}

// Infof prints!
func (c *Console) Infof(f string, args ...interface{}) {
	c.Printf(color.CyanString(fmt.Sprintf(f, args...)))
}

// Warningf prints!
func (c *Console) Warningf(f string, args ...interface{}) {
	c.Printf(color.YellowString(fmt.Sprintf(f, args...)))
}

// Errorf prints!
func (c *Console) Errorf(f string, args ...interface{}) {
	c.Printf(color.RedString(fmt.Sprintf(f, args...)))
}

func (c *Console) push(ctx context.Context, a Activity, prompt string) bool {
	prompt = strings.Trim(prompt, "> \t")

	ctx, cancel := context.WithCancel(ctx)
	actx := Context{
		Context:       ctx,
		cancelContext: cancel,
		console:       c,
		prompt:        prompt,
		activity:      a,
	}

	// create activity
	if err := a.Create(actx); err != nil {
		c.Errorf("Failed to create new %v-activity for: %v", prompt, err)
		return false
	}

	// deactivate current activity
	if c.activities.Len() > 0 {
		actx := (c.activities.Front().Value.(Context))
		actx.activity.Deactivate(actx)
	}

	// push new activity
	c.activities.PushFront(actx)

	// activate new current activity
	{
		actx := (c.activities.Front().Value.(Context))
		actx.activity.Activate(actx)
	}

	return true
}

func (c *Console) pop() {
	// deactivate current activity
	{
		actx := (c.activities.Front().Value.(Context))
		actx.activity.Activate(actx)
	}

	// destroy current activity
	{
		actx := (c.activities.Front().Value.(Context))
		actx.activity.Destroy(actx)
	}

	// pop current activity
	c.activities.Remove(c.activities.Front())

	// activate previous activity
	{
		actx := (c.activities.Front().Value.(Context))
		actx.activity.Activate(actx)
	}
}

// Run console loop
func (c *Console) Run(ctx context.Context, root Activity, prompt string, reading chan<- struct{}) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c.push(ctx, root, prompt)

	go func() {
		<-ctx.Done()
		c.Close()
	}()

	for {
		// current activity
		actx := (c.activities.Front().Value.(Context))

		// currently available commands
		var completer *readline.PrefixCompleter
		var helpPostfix []string
		{
			cfg := c.Config

			pc := []readline.PrefixCompleterInterface{
				DPCItem("exit", "Close the program"),
				DPCItem("help", "This help"),
				DPCItem("shortcuts", "Editor shortcuts"),
				DPCSeparator(),
			}
			pcis, postHelp := actx.activity.Completer()
			pc = append(pc, pcis...)
			completer = readline.NewPrefixCompleter(pc...)
			cfg.AutoComplete = completer
			helpPostfix = postHelp
			c.SetConfig(cfg)
		}

		// current prompt
		c.SetPrompt(color.CyanString(actx.prompt + "> "))

		// wait input
		select {
		case reading <- struct{}{}:
		default:
		}
		line, err := c.Readline()
		switch {
		case err == readline.ErrInterrupt:
			if len(line) == 0 {
				return
			}
			continue
		case err == io.EOF:
			return
		}

		// handle input
		line = strings.TrimSpace(line)
		switch line {
		case "":
			continue
		case "exit":
			return
		case "help":
			c.Printf(completer.Tree("  "))
			if len(helpPostfix) > 0 {
				for _, v := range helpPostfix {
					c.Printf("  " + v)
				}
			}
			continue
		case "shortcuts":
			c.Printf("https://github.com/chzyer/readline/blob/master/doc/shortcut.md")
		default:
			if !actx.activity.Command(actx, Args{Line: line, Args: strings.Split(line, " ")}) {
				c.Errorf("No such command")
			}
		}
	}
}
