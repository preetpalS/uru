// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

// +build !windows

package command

import (
	"fmt"
	"os"
	"os/exec"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["install"] = Command{
		Name:    "install",
		Aliases: []string{"install", "in"},
		Usage:   "admin install",
		HelpMsg: "install uru",
		Eg:      `admin install`}
}

func adminInstall(ctx *env.Context) {
	_, err := exec.LookPath("uru_rt")
	if err != nil {
		fmt.Printf("[ERROR] uru_rt must be present in a directory on PATH\n")
		os.Exit(1)
	}

	fmt.Printf(env.BashWrapper)
}
