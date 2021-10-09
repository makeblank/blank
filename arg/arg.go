// Provides low level utilities for command line argument parsing.
package arg

import (
	"strings"

	"github.com/makeblank/blank/std"
)

// Is the given argument a valid flag?
func IsAFlag(a string) bool {
	return len(a) > 1 && a[0] == '-'
}

// Is the given argument one of the flags?
func IsFlag(a interface{}, flags ...string) (bool, []string) {
	var splat []string

	switch t := a.(type) {
	case []string:
		splat = t
	case string:
		splat = SplitFlags(t)
	default:
		panic("First argument to IsFlag must be a string or []string")
	}

	for i, a := range splat {
		is := false
		ix := -1

		for _, f := range flags {
			for _, f := range SplitFlags(f) {
				is = is || (a == f)

				if !is && ix == -1 {
					ix = i
				}
			}
		}

		if !is {
			if ix > 0 {
				splat[0], splat[ix] = splat[ix], splat[0]
			}

			return false, splat
		}
	}

	return true, splat
}

func SplitFlags(a string) (flags []string) {
	if len(a) > 2 && a[0:2] == "--" {
		// long --flag
		flags = []string{a[2:]}
	} else if len(a) > 1 && a[0] == '-' && a[1] != '-' {
		// short -abc flags
		flags = strings.Split(a[1:], "")
	}
	return
}

// Is the given string any of the given words?
func IsWord(s string, words ...string) bool {
	for _, w := range words {
		if w == s {
			return true
		}
	}
	return false
}

// If flag is true, and first element in args is a flag,
// return that element and the rest of args. Otherwise,
// return the empty string and args unchanged.
func nextOpt(flag bool, args []string) (a string, rest []string) {
	h := std.Head(args)

	if IsAFlag(h) == flag {
		a = h
		rest = std.Tail(args)
	} else {
		rest = args
	}

	return
}

// If first element in args is a flag, return it and the
// rest of args, otherwise return the empty string and args
// unchanged.
func NextFlag(args []string, flags ...string) (f string, rest []string) {
	f, rest = nextOpt(true, args)

	if len(flags) > 0 {
		if is, _ := IsFlag(f, flags...); !is {
			f, rest = "", args
		}
	}

	return
}

// Inverse of NextFlag.
func NextArg(args []string, words ...string) (a string, rest []string) {
	a, rest = nextOpt(false, args)

	if len(words) > 0 && !IsWord(a, words...) {
		a, rest = "", args
	}

	return
}

// Is string not empty?
func Ok(s string) bool {
	return len(s) > 0
}
