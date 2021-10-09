package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/makeblank/blank/std"
)

//go:embed make.mk
var makefile []byte

const MakeCommandName = "make"

const makeCommandHelpFmt = `
Target is searched for in directories in %[1]s. The
first file named "[target].mk" found is used as the makefile
(passed as the "-f" make option.)

%[1]s is a list of paths separated by %[2]q. Each
directory in %[1]s is used in two additional ways:
1. as an "--include-dir" make option, and 2. as the VPATH
make variable.

All other options are passed directly to the 'make' program.
Run 'make --help' for additional options.
`

var makeCommandHelp = fmt.Sprintf(
	makeCommandHelpFmt,
	vBLANK_PATH,
	filepath.ListSeparator,
)

// The "make" subcommand type.
type MakeCommand struct {
	info *Info

	NoFlags
}

func (c *MakeCommand) Name() string {
	return MakeCommandName
}

func (c *MakeCommand) Info() *Info {
	return c.info
}

func (c *MakeCommand) Help() string {
	return makeCommandHelp
}

func (c *MakeCommand) Run(args []string) error {
	// args passed to make:
	//   -s Silent operation
	//   -E Eval implicit makefile
	args = append([]string{"-s", "-E", string(makefile)}, args...)

	cmd := exec.Command("make", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	ExecCmd(cmd)

	// NOTE: ExecCmd must block and exit itself
	panic("ExecCmd did not exit")
}

// The default "make" subcommand instance.
var Make = &MakeCommand{
	info: &Info{
		Line: "%s [options] [target] ...",
		Desc: "Generate blank dev projects using makefiles.",
	},
}
