package console

import (
	"bytes"
	"strings"

	"github.com/chzyer/readline"
)

var helpSeparator = "~~~~~\n"

// DescribedPrefixCompleter is readline.PrefixCompleter with description
type DescribedPrefixCompleter struct {
	*readline.PrefixCompleter
	desc      string
	separator bool
	nohelp    bool
}

// Print impl
func (p DescribedPrefixCompleter) Print(prefix string, level int, buf *bytes.Buffer) {
	if p.nohelp {
		return
	}
	if p.separator {
		buf.WriteString(prefix)
		buf.WriteString(helpSeparator)
		return
	}
	if strings.TrimSpace(string(p.GetName())) != "" {
		buf.WriteString(prefix)
		if level > 0 {
			buf.WriteString("├")
			buf.WriteString(strings.Repeat("─", (level*4)-2))
			buf.WriteString(" ")
		}
		buf.WriteString(string(p.GetName()) + "  " + p.desc + "\n")
		level++
	}
	for _, ch := range p.GetChildren() {
		ch.Print(prefix, level, buf)
	}
}

// DPCItem makes DescribedPrefixCompleter
func DPCItem(name, desc string, pc ...readline.PrefixCompleterInterface) DescribedPrefixCompleter {
	name += " "
	return DescribedPrefixCompleter{
		PrefixCompleter: &readline.PrefixCompleter{
			Name:     []rune(name),
			Dynamic:  false,
			Children: pc,
		},
		desc: desc,
	}
}

// DPCItemNoHelp makes DescribedPrefixCompleter
func DPCItemNoHelp(name string, pc ...readline.PrefixCompleterInterface) DescribedPrefixCompleter {
	name += " "
	return DescribedPrefixCompleter{
		PrefixCompleter: &readline.PrefixCompleter{
			Name:     []rune(name),
			Dynamic:  false,
			Children: pc,
		},
		nohelp: true,
	}
}

// DPCSeparator makes DescribedPrefixCompleter
func DPCSeparator() DescribedPrefixCompleter {
	return DescribedPrefixCompleter{
		PrefixCompleter: &readline.PrefixCompleter{
			Name:     []rune(" "),
			Dynamic:  false,
			Children: nil,
		},
		separator: true,
	}
}
