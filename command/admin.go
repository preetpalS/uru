// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	CommandRegistry["admin"] = Command{
		Name:    "admin",
		Aliases: nil,
		Usage:   "admin SUBCMD ARGS",
		HelpMsg: "administer uru installation",
		Eg:      `admin add C:\ruby200\bin`}
}

func Admin(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	if len(cmdArgs) == 0 {
		return
	}
	ctx.SetCmdArgs(cmdArgs[1:])

	switch subCmd := cmdArgs[0]; {
	case ctx.CmdRegex(`add`).MatchString(subCmd):
		adminAdd(ctx)
	case ctx.CmdRegex(`install`).MatchString(subCmd):
		adminInstall(ctx)
	case ctx.CmdRegex(`rm`).MatchString(subCmd):
		adminRemove(ctx)
	default:
		fmt.Printf("[ERROR] I don't understand the `%s` admin sub-command\n\n", subCmd)
		ctx.SetCmdAndArgs(``, nil)
		Help(ctx)
	}
}
