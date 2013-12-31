// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	AdminCmdRegistry["gemset"] = Command{
		Name:    "gemset",
		Aliases: []string{"gemset", "gs"},
		Usage:   "admin gemset init | ls | rm",
		HelpMsg: "administer gemset installations",
		Eg:      `admin gemset init`}
}

func adminGemset(ctx *env.Context) {
	if len(ctx.CmdArgs()) == 0 {
		fmt.Println("[ERROR] must specify a gemset operation")
		os.Exit(1)
	}

	switch subCmd := ctx.CmdArgs()[0]; subCmd {
	case `init`:
		if err:= gemsetInit(ctx); err != nil {
			fmt.Println("[ERROR] unable to initialize gemset")
			os.Exit(1)
		}
	case `ls`:
		if err := gemsetList(ctx); err != nil {
			fmt.Println("[ERROR] unable to list the gemset gems")
			os.Exit(1)
		}
	case `rm`:
		if err := gemsetRemove(ctx); err != nil {
			fmt.Println("[ERROR] unable to remove the gemset")
			os.Exit(1)
		}
	default:
		fmt.Printf("[ERROR] I don't understand the `%s` gemset sub-command\n\n", subCmd)
		ctx.SetCmdAndArgs(``, nil)
		Help(ctx)
	}
}

func gemsetInit(ctx *env.Context) (err error) {
	fmt.Println("---> performing gemset init")

	return
}

func gemsetList(ctx *env.Context) (err error) {
	fmt.Println("---> performing gemset ls")

	return
}

func gemsetRemove(ctx *env.Context) (err error) {
	fmt.Println("---> performing gemset rm")

	return
}
