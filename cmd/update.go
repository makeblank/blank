package cmd

import (
	"fmt"
	"strings"

	. "github.com/makeblank/blank/arg"
	. "github.com/makeblank/blank/std"
)

const UpdateCommandName = "update"

const updateCommandInfo = `
Target must be an existing JSON/YAML file. The data it
represents will be updated according to the specified
operations and then written to stdout.

The path argument to operation flags must be a path to a
member in target data (e.g. "/path/to/member"). If it is
ommitted, the operation is applied to the entire data.

The json argument must be a JSON string, or a path to a
JSON/YAML file by prefixing it with an "@" sign.

Examples:
  blank update package.json -s /dependencies/eslint '"^7"'
  blank update .eslintrc.json -a /extends '"standard"'
  blank update config.yaml -m @base.yaml
`

var (
	updateCommandHelp string

	operationFlags = []string{
		"-sam", "--set", "--append", "--merge",
	}

	outputTypes = []string{
		"json", "yaml",
	}

	outputTypesStr   = strings.Join(outputTypes, ", ")
	outputTypesError = fmt.Sprintf("must be: %s", outputTypesStr)
)

type Operation struct {
	Ops  []string
	Path string
	Json string
}

// The "update" subcommand type.
type UpdateCommand struct {
	info  *Info
	flags []*Flag
	ops   []*Flag

	NoHelp
}

func (c *UpdateCommand) Name() string {
	return UpdateCommandName
}

func (c *UpdateCommand) Info() *Info {
	return c.info
}

func (c *UpdateCommand) Help() string {
	return updateCommandHelp
}

func (c *UpdateCommand) Run(args []string) error {
	var (
		target, a, b string
		ops          []string
		operations   []*Operation
		output       = "json"
	)

	operations = make([]*Operation, 0)

	if a, args = NextFlag(args, "-t", "--type"); Ok(a) {
		if output, args = NextArg(args, outputTypes...); Empty(output) {
			return FlagError(outputTypesError, "-t")
		}
	}

	if target, args = NextArg(args); Empty(target) {
		return ArgRequiredError("target")
	}

	for len(args) > 0 {
		var is bool

		if a, args = NextFlag(args); Empty(a) {
			return FlagRequiredError("operation")
		} else if is, ops = IsFlag(a, operationFlags...); !is {
			return FlagUnknownError(ops[0])
		}

		if a, args = NextArg(args); Empty(a) {
			return ArgRequiredError("path", "json")
		}

		if b, args = NextArg(args); Empty(b) {
			operations = append(operations, &Operation{ops, "/", a})
		} else {
			operations = append(operations, &Operation{ops, a, b})
		}
	}

	// TODO: update config file

	fmt.Printf("input: %s\n", target)
	fmt.Printf("output: %s\n", output)

	for _, o := range operations {
		fmt.Printf("%q: %q %v\n", o.Ops, o.Path, o.Json)
	}

	return nil
}

func (c *UpdateCommand) Flags() []*Flag {
	return c.flags
}

// The default "update" subcommand instance.
var Update = &UpdateCommand{
	info: &Info{
		Line: "%s [options] target [operation [path] json]...",
		Desc: "Update/patch config files.",
	},

	flags: []*Flag{
		{
			Name: "-t, --type",
			Desc: fmt.Sprintf("output file as `t` (%s)", outputTypesStr),
		},
	},

	ops: []*Flag{
		{Name: "-s, --set", Desc: "set value at path"},
		{Name: "-m, --merge", Desc: "merge with value at path"},
		{Name: "-a, --append", Desc: "append array values (modifier)"},
	},
}

func init() {
	var b strings.Builder

	b.WriteString("\nOperation:\n")

	for _, op := range Update.ops {
		WriteFlagUsage(&b, op)
	}

	b.WriteString(updateCommandInfo)

	updateCommandHelp = b.String()
}
