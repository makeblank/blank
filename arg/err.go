package arg

import (
	"fmt"
	"strings"
)

const (
	isRequired = "is required"
	isUnknown  = "is unknown"

	FlagOpt = "flag"
	ArgOpt  = "argument"
	CmdOpt  = "command"
)

type OptionError struct {
	Type  string
	Names []string
	Desc  string
}

func (e *OptionError) Error() string {
	var (
		n string
		c = len(e.Names)
	)

	if c == 0 {
		n = "A"
	} else {
		n = fmt.Sprintf("The %s", strings.Join(e.Names, ", "))
	}

	return fmt.Sprintf("%s %s %s", n, e.Type, e.Desc)
}

func NewOptionError(t string, m string, a ...string) *OptionError {
	return &OptionError{t, a, m}
}

func FlagError(m string, a ...string) *OptionError {
	return NewOptionError(FlagOpt, m, a...)
}

func FlagRequiredError(a ...string) *OptionError {
	return NewOptionError(FlagOpt, isRequired, a...)
}

func FlagUnknownError(a ...string) *OptionError {
	return NewOptionError(FlagOpt, isUnknown, a...)
}

func ArgError(m string, a ...string) *OptionError {
	return NewOptionError(ArgOpt, m, a...)
}

func ArgRequiredError(a ...string) *OptionError {
	return NewOptionError(ArgOpt, isRequired, a...)
}

func CmdError(m string, a ...string) *OptionError {
	return NewOptionError(CmdOpt, m, a...)
}

func CmdUnknownError(a ...string) *OptionError {
	return NewOptionError(CmdOpt, isUnknown, a...)
}

func CmdRequiredError() *OptionError {
	return NewOptionError(CmdOpt, isRequired)
}
