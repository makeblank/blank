// Provides command handlers for the "blank" program.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/makeblank/blank/arg"
	. "github.com/makeblank/blank/std"
)

const (
	BlankCommandName = "blank"
	vBLANK_PATH      = "BLANK_PATH"
	vBLANK           = "BLANK"
)

const blankCommandHelp = `
If none of the above commands is provided, blank defaults to
the "make" command.
`

var blankPaths []string

// The main "blank" command type.
type BlankCommand struct {
	info     *Info
	flags    []*Flag
	commands []Command
}

func (c *BlankCommand) Name() string {
	return BlankCommandName
}

func (c *BlankCommand) Info() *Info {
	return c.info
}

// Returns a help screen containing a list of sub-commands
// and environment variables.
func (c *BlankCommand) Help() string {
	var b strings.Builder

	b.WriteString("\nCommands:\n")

	for _, cmd := range c.commands {
		desc := cmd.Info().Desc
		fmt.Fprintf(&b, FlagLineFormat, cmd.Name(), desc)
	}

	b.WriteString(blankCommandHelp)

	WriteEnvironment(&b, vBLANK, vBLANK_PATH)

	return b.String()
}

func (c *BlankCommand) Flags() []*Flag {
	return c.flags
}

func (c *BlankCommand) Run(args []string) error {
	var (
		err error
		a   string
	)

	for len(args) > 0 {
		if a, args = NextFlag(args, "-aP"); Empty(a) {
			break
		}

		f := a

		if a, args = NextArg(args); Ok(a) {
			if f == "-a" {
				addBlankPath(a)
			} else {
				blankPaths = filepath.SplitList(a)
			}
		} else {
			err = ArgRequiredError(f)
		}
	}

	exportBlankPath()

	if err == nil {
		if a, args = NextArg(args); Ok(a) {
			if cmd := c.findSubcommand(a); cmd != nil {
				RunWithHelp(cmd, args, c.Name())
				return nil
			}
		}

		// 'make' as default command
		RunWithHelp(Make, args, c.Name())
		return nil
	}

	return err
}

func (c *BlankCommand) findSubcommand(name string) Command {
	// TODO: use subcommand map
	for _, cmd := range c.commands {
		if cmd.Name() == name {
			return cmd
		}
	}

	return nil
}

// The default "blank" command instance.
var Blank = &BlankCommand{
	&Info{
		Line: "%s [options] [command [command_options]...]",
		Desc: "Program to generate blank dev projects.",
	},

	[]*Flag{
		{Name: "-a", Desc: fmt.Sprintf("append `path` to %s", vBLANK_PATH)},
		{Name: "-P", Desc: fmt.Sprintf("set `paths` as %s", vBLANK_PATH)},
	},

	[]Command{
		Make,
		Update,
		Help,
	},
}

// Run the main "blank" program.
func Main(args []string) {
	RunWithHelp(Blank, args)
}

func init() {
	// set default BLANK_PATH in user's config home directory
	if env := os.Getenv(vBLANK_PATH); env == "" {
		blankPaths = make([]string, 1)

		if cfg, err := os.UserConfigDir(); err == nil {
			blankPaths[0] = filepath.Join(cfg, "blank")
			exportBlankPath()
		}
	} else {
		blankPaths = filepath.SplitList(env)
	}

	// set BLANK to this program so that makefiles can call
	// it recursively
	os.Setenv(vBLANK, os.Args[0])
}

func addBlankPath(path string) {
	blankPaths = append(blankPaths, path)
}

func exportBlankPath() {
	sep := string(filepath.ListSeparator)
	os.Setenv(vBLANK_PATH, strings.Join(blankPaths, sep))
}
