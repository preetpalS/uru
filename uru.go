// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// Uru is a lightweight, minimal install, multi-platform tool that helps you use
// Ruby more productively. Uru untethers your workflow from a single Ruby.
package main

import (
	"io/ioutil"
	"log"
	"os"

	"bitbucket.org/jonforums/uru/command"
	"bitbucket.org/jonforums/uru/env"
)

func main() {
	var debug, needHelp bool
	var cmd string

	if len(os.Args) == 1 {
		needHelp = true
	}
	for _, a := range os.Args {
		switch a {
		case "-h", "--help":
			needHelp = true
		case "--debug":
			debug = true
		}
	}

	if !debug {
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("[DEBUG] initializing uru v%s\n", env.AppVersion)

	ctx := env.NewContext()
	initHome(ctx)
	initRubies(ctx)

	if needHelp {
		cmd = "help"
	} else {
		cmd = os.Args[1]
		if len(os.Args) > 2 {
			ctx.SetCmdArgs(os.Args[2:])
		}
	}
	ctx.SetCmd(cmd)
	log.Printf("[DEBUG] cmd = %s\n", cmd)

	command.CmdRouter.Dispatch(ctx, cmd)
}
