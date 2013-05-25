// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// Uru is a lightweight, minimal install, multi-platform tool that helps you use
// Ruby more productively. Uru untethers your workflow from a single Ruby.
package main

import (
	"flag"
	"log"
	"os"

	"bitbucket.org/jonforums/uru/command"
)

func main() {
	if len(os.Args) == 1 || *help == true {
		command.Help(&ctx)
	}
	if *version == true {
		command.Version(&ctx)
		os.Exit(0)
	}

	cmd := flag.Arg(0)
	ctx.SetCmdAndArgs(cmd, flag.Args()[1:])
	log.Printf("[DEBUG] cmd = %s\n", cmd)

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
	case ctx.CmdRegex(`version`).MatchString(cmd):
		command.Version(&ctx)
	default:
		command.Use(&ctx, ``)
	}
}
