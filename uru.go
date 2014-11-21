// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// Uru is a lightweight, minimal install, multi-platform tool that helps you use
// Ruby more productively. Uru untethers your workflow from a single Ruby.
package main

import (
	"flag"
	"io/ioutil"
	"log"

	"bitbucket.org/jonforums/uru/command"
	"bitbucket.org/jonforums/uru/env"
)

func main() {
	// initialization
	debug := flag.Bool(`debug`, false, "enable debug mode")
	flag.Parse()

	if !*debug {
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("[DEBUG] initializing uru v%s\n", env.AppVersion)

	ctx := env.NewContext()
	cmdRouter := command.NewRouter(command.Use)

	initHome(ctx)
	initCommandRouter(cmdRouter)
	initRubies(ctx)

	if len(flag.Args()) == 0 {
		command.Help(ctx)
	}

	cmd := flag.Arg(0)
	ctx.SetCmdAndArgs(cmd, flag.Args()[1:])
	log.Printf("[DEBUG] cmd = %s\n", cmd)

	// dispatch command to registered handler
	cmdRouter.Dispatch(ctx, cmd)
}
