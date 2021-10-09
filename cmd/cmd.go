package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "github.com/makeblank/blank/std"
)

// Encapsulates only "informational" details about a flag.
// Used by WriteFlagUsage and WriteCommandUsage to display
// the help screen.
type Flag struct {
	Name string // The flag names
	Desc string // The flag description.
}

// A uniquely named command.
type Named interface {
	Name() string
}

// A command's usage and short description.
type Info struct {
	Line string
	Desc string
}

// Provides the Info and Help methods, for commands to
// provide detailed info in the help screen.
type Informative interface {
	Info() *Info
	Help() string // Returns longer help message
}

// Provides the Run method, for commands to implement their
// run function.
type Runnable interface {
	Run([]string) error
}

// Provides the Flags method, for commands to provide a list
// of command-line flags.
type Flagged interface {
	Flags() []*Flag
}

// A command.
type Command interface {
	Named
	Informative
	Runnable
	Flagged
}

// Embedable struct that implements an empty "Flags" method.
type NoFlags struct{}

func (n *NoFlags) Flags() []*Flag {
	return nil
}

// Embeddable struct that implements an empty "Help" method.
type NoHelp struct{}

func (n *NoHelp) Help() string {
	return ""
}

// The fmt string used by WriteFlagUsage.
const FlagLineFormat = "  %-13s  %s\n"

// Write program help screen.
//
// The p argument is this command's parents' names, if any.
func WriteCommandUsage(w io.Writer, c Command, p ...string) {
	p = append(p, c.Name())

	info := c.Info()
	line := fmt.Sprintf(info.Line, strings.Join(p, " "))

	fmt.Fprint(w, "Usage: ", line, "\n\n", info.Desc, "\n")

	flags := c.Flags()

	fmt.Fprintln(w, "\nOptions:")

	if len(flags) > 0 {
		for _, f := range flags {
			WriteFlagUsage(w, f)
		}
	}

	fmt.Fprintf(
		w,
		FlagLineFormat,
		"-h, --help",
		"show help screen",
	)

	help := c.Help()

	if help != "" {
		fmt.Fprintln(w, help)
	}
}

// Write a flag's usage line.
func WriteFlagUsage(w io.Writer, f *Flag) {
	word, usage := UnquoteUsage(f.Desc)
	flag := fmt.Sprintf("%s %s", f.Name, word)
	fmt.Fprintf(w, FlagLineFormat, flag, usage)
}

// Run command c, intercepting first -h or --help argument.
func RunWithHelp(c Command, args []string, p ...string) {
	for _, a := range args {
		if a == "-h" || a == "--help" {
			WriteCommandUsage(os.Stdout, c, p...)
			return
		} else {
			break
		}
	}

	if err := c.Run(args); err != nil {
		WriteError(err)
		os.Stdout.WriteString("\n")
		WriteCommandUsage(os.Stdout, c)
		os.Exit(1)
	}
}

// Extracts a back-quoted name from the usage string for a
// flag and returns it and the un-quoted usage.
//
// Given "a `name` to show" it returns ("name", "a name to
// show"). If there are no back quotes, the name is an
// educated guess of the type of the flag's value, or the
// empty string if the flag is boolean.
func UnquoteUsage(desc string) (name string, usage string) {
	usage = desc

	for i := 0; i < len(usage); i++ {
		if usage[i] == '`' {
			for j := i + 1; j < len(usage); j++ {
				if usage[j] == '`' {
					name = usage[i+1 : j]
					usage = usage[:i] + name + usage[j+1:]
					return name, usage
				}
			}
			break // Only one back quote; use empty name.
		}
	}

	return
}
