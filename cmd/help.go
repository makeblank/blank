package cmd

import (
	"os"

	. "github.com/makeblank/blank/std"
)

const HelpCommandName = "help"

// The "help" subcommand type.
type HelpCommand struct {
	info *Info
	NoFlags
	NoHelp
}

func (c *HelpCommand) Name() string {
	return HelpCommandName
}

func (c *HelpCommand) Info() *Info {
	return c.info
}

func (c *HelpCommand) Run(args []string) error {
	a := Head(args)

	if !Empty(a) {
		if cmd := Blank.findSubcommand(a); cmd != nil {
			WriteCommandUsage(os.Stdout, cmd, Blank.Name())
			return nil
		}
	}

	WriteCommandUsage(os.Stdout, Blank)

	return nil
}

// The default "help" subcommand instance.
var Help = &HelpCommand{
	info: &Info{
		Line: "%s [command]",
		Desc: "Show help screen.",
	},
}
