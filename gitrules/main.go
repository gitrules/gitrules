package main

import (
	"fmt"
	"os"

	"github.com/gitrules/gitrules/gitrules/cmd"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/must"
	_ "github.com/gitrules/gitrules/runtime"
)

func main() {
	err, stk := must.TryWithStack(
		func() { cmd.Execute() },
	)
	if err != nil {
		if base.IsVerbose() {
			fmt.Fprintln(os.Stderr, string(stk))
		}
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}
