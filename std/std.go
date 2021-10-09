// Provides 'standard' utilities for the "blank" program.
package std

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

// Run an exec.Cmd with additional mechanics.
//
// (1) Forward all signals recieved from this program to the
// command's process, and (2) automatically exit with the
// command's exit status.
func ExecCmd(cmd *exec.Cmd) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals)

	go func() {
		err := cmd.Run()
		code := cmd.ProcessState.ExitCode()

		// is this condition possible?
		// got error, but process exit code is 0
		if err != nil && code == 0 {
			code = 1
		}

		os.Exit(code)
	}()

	for {
		// keep looping if channel closed or signal failed,
		// because above goroutine will ultimately os.Exit
		if s, ok := <-signals; !ok {
			continue
		} else {
			cmd.Process.Signal(s) //nolint:errcheck
		}
	}
}

// Print environment variables.
func WriteEnvironment(out io.Writer, names ...string) {
	if len(names) == 0 {
		return
	}

	fmt.Fprintln(out, "\nEnvironment:")

	for _, name := range names {
		n := name

		for _, path := range strings.Split(os.Getenv(name), ":") {
			fmt.Fprintf(out, "  %-10s  %s\n", n, path)
			n = ""
		}
	}
}

// Write error message to stderr.
func WriteError(v ...interface{}) {
	var (
		h   interface{}
		err error
	)

	if len(v) == 0 {
		return
	}

	h, v = v[0], v[1:]

	switch t := h.(type) {
	case error:
		err = t
	case string:
		err = fmt.Errorf(t, v...)
	default:
		err = fmt.Errorf("%v", t)
	}

	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

func Head(slice []string) (h string) {
	if len(slice) > 0 {
		h = slice[0]
	}
	return
}

func Tail(slice []string) (t []string) {
	if len(slice) > 1 {
		t = slice[1:]
	}
	return
}

func Empty(s string) bool {
	return len(s) == 0
}
