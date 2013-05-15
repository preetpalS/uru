// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package main

import (
	"flag"
	"io/ioutil"
	"log"

	"bitbucket.org/jonforums/uru/env"
)

var (
	debug   = flag.Bool(`debug`, false, "enable debug mode")
	help    = flag.Bool(`help`, false, "this help summary")
	version = flag.Bool(`version`, false, "show version info")

	ctx env.Context
)

func init() {
	flag.Parse()

	if !*debug {
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("[DEBUG] initializing uru v%s\n", env.AppVersion)

	initHome()
	initContext()
	initCommandParser()
	initRubies()
}
