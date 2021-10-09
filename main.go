// Provides a command line program, "blank" that can help
// generate blank dev projects using makefiles.
package main

import (
	"os"

	"github.com/makeblank/blank/cmd"
)

func main() {
	cmd.Main(os.Args[1:])
}
