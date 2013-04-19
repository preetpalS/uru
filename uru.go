// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// Uru is a lightweight, minimal install, multi-platform tool that helps you use
// Ruby more productively. Uru untethers your workflow from a single Ruby.
package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/command"
)

func main() {

	if len(os.Args) == 1 || *help == true {
		command.Help(&ctx)
	}

	cmd := flag.Arg(0)
	ctx.SetCmdAndArgs(cmd, flag.Args()[1:])

	switch {
	case ctx.CmdRegex(`admin`).MatchString(cmd):
		command.Admin(&ctx)
	case ctx.CmdRegex(`gem`).MatchString(cmd):
		command.Gem(&ctx)
	case ctx.CmdRegex(`help`).MatchString(cmd):
		command.Help(&ctx)
	case ctx.CmdRegex(`ls`).MatchString(cmd):
		command.List(&ctx)
	case ctx.CmdRegex(`ruby`).MatchString(cmd):
		command.Ruby(&ctx)
	case ctx.CmdRegex(`use`).MatchString(cmd):
		command.Use(&ctx, ``)
	default:
		fmt.Printf("[ERROR] I don't understand the `%s` command\n\n", cmd)
		command.Help(&ctx)
	}

}
