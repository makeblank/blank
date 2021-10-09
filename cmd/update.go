package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
	"github.com/makeblank/blank/cfg"
	"gopkg.in/yaml.v3"

	. "github.com/makeblank/blank/arg"
	. "github.com/makeblank/blank/std"
)

const UpdateCommandName = "update"

const updateExtraInfo = `
Target must be an existing config file. The data it
represents will be updated according to the specified
operations and then written to stdout.

The path argument to operation flags must be a path to a
member in target data (e.g. "/path/to/member"). If it is
ommitted, json must be an object and the operation is
applied to the entire target data.

The json argument must be a JSON string, or a path to a
config file by prefixing it with an "@" sign.

Examples:
  blank update package.json -s /dependencies/eslint '"^7"'
  blank update .eslintrc.json -a /extends '["standard"]'
  blank update config.yaml -m @base.yaml
`

var (
	updateExtraHelp  string
	updateOperations []string

	fileTypes    = []string{"json", "yaml"}
	fileTypesStr = strings.Join(fileTypes, ", ")
	fileTypesErr = fmt.Sprintf("must be: %s", fileTypesStr)

	marshallers = map[string]marshal{
		"json": jsonMarshal,
		"yaml": yaml.Marshal,
	}
)

type marshal func(interface{}) ([]byte, error)

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
	return updateExtraHelp
}

func (c *UpdateCommand) Run(args []string) error {
	var (
		target, a, t string
		sources      []*cfg.Source
		input        = "json"
		output       = "json"
		o            = 0
	)

	sources = make([]*cfg.Source, 0)

	for len(args) > 0 {
		if a, args = NextFlag(args, "-io", "--in", "--out"); Empty(a) {
			break
		}

		if t, args = NextArg(args, fileTypes...); Empty(t) {
			return FlagError(fileTypesErr, a)
		}

		if ok, _ := IsFlag(a, "-i", "--in"); ok {
			input = t
		} else {
			output = t
		}
	}

	if target, args = NextArg(args); Empty(target) {
		return ArgRequiredError("target")
	}

	for len(args) > 0 {
		var (
			b, path, src string
			data         []byte
			ops          []string
			err          error
			is           bool
		)

		if a, args = NextFlag(args); Empty(a) {
			return FlagRequiredError("operation")
		} else if is, ops = IsFlag(a, updateOperations...); !is {
			return FlagUnknownError(ops[0])
		}

		if a, args = NextArg(args); Empty(a) {
			return ArgRequiredError("path", "json")
		}

		name := fmt.Sprintf("Op#%d", o)

		if b, args = NextArg(args); Ok(b) {
			path = a
			src = b
		} else {
			path = "/"
			src = a
		}

		if src[0] == '@' {
			if data, err = ioutil.ReadFile(src[1:]); err != nil {
				return err
			}
		} else {
			data = []byte(src)
		}

		if src, err := newSource(name, path, ops, data); err != nil {
			return err
		} else {
			sources = append(sources, src)
		}

		o++
	}

	return updateFile(os.Stdout, target, input, output, sources)
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
			Name: "-i, --in",
			Desc: fmt.Sprintf("read target as `t` (%s)", fileTypesStr),
		},
		{
			Name: "-o, --out",
			Desc: fmt.Sprintf("output as `t` (%s)", fileTypesStr),
		},
	},

	ops: []*Flag{
		{Name: "-s", Desc: "set new values only"},
		{Name: "-m", Desc: "merge values"},
		{Name: "-a", Desc: `concatenate array values (implies "-m")`},
		{Name: "-u", Desc: `drop duplicate array values (implies "-a")`},
	},
}

// Create new cfg.Source based on operations from cmd line.
func newSource(n, p string, ops []string, js []byte) (s *cfg.Source, err error) {
	var (
		kind reflect.Kind
		file *cfg.File
		data interface{}
		map_ map[string]interface{}
		opts = make([]func(*mergo.Config), 0, len(ops))
	)

	p = strings.Trim(p, "/")

	if err = json.Unmarshal(js, &data); err != nil {
		return
	}

	kind = reflect.TypeOf(data).Kind()

	if p == "" && kind != reflect.Map {
		return nil, fmt.Errorf("json must be an object if path is omitted")
	} else {
		map_ = cfg.PointerToMap(p, data)
	}

	file = &cfg.File{
		Path: n,
		Data: map_,
	}

	for _, op := range ops {
		switch op {
		case "a", "u":
			opts = append(opts, mergo.WithOverride, mergo.WithAppendSlice)
			// TODO: "u" operation: drop duplicates transformer
		case "m":
			opts = append(opts, mergo.WithOverride)
		}
	}

	return &cfg.Source{
		File:    file,
		Options: opts,
	}, nil
}

// update config file from given sources and write updated
// data as given type to the writer.
func updateFile(w io.Writer, p, in, out string, s []*cfg.Source) (err error) {
	var (
		file *cfg.File
		fn   marshal
	)

	if fn = marshallers[out]; fn == nil {
		panic(fmt.Errorf("unknown config file type %s", out))
	}

	if file, err = cfg.ReadFile(p, in); err != nil {
		return err
	}

	if err = file.MergeSource(s...); err != nil {
		return err
	} else {
		var b []byte

		if b, err = fn(file.Data); err != nil {
			return
		} else {
			_, err = w.Write(b)
		}
	}

	return
}

func jsonMarshal(in interface{}) ([]byte, error) {
	return json.MarshalIndent(in, "", " ")
}

func init() {
	var b strings.Builder

	b.WriteString("\nOperation:\n")
	for _, op := range Update.ops {
		WriteFlagUsage(&b, op)
	}
	b.WriteString(updateExtraInfo)

	updateExtraHelp = b.String()
	updateOperations = make([]string, len(Update.ops))

	for i, f := range Update.ops {
		updateOperations[i] = f.Name
	}
}
