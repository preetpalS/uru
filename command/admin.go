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
	subCmd := cmdArgs[0]
	ctx.SetCmdArgs(cmdArgs[1:])

	// initialize the admin subcommand router
	adminRouter := NewRouter(func (ctx *env.Context) {
		fmt.Printf("[ERROR] I don't understand the `%s` admin sub-command\n\n", subCmd)
		ctx.SetCmdAndArgs(``, nil)
		Help(ctx)
	})
	adminRouter.Handle([]string{`add`}, adminAdd)
	adminRouter.Handle([]string{`gemset`, `gs`}, adminGemset)
	adminRouter.Handle([]string{`install`, `in`}, adminInstall)
	adminRouter.Handle([]string{`refresh`}, adminRefresh)
	adminRouter.Handle([]string{`retag`, `tag`}, adminRetag)
	adminRouter.Handle([]string{`del`, `rm`}, adminRemove)

	// dispatch subcommands to registered handlers
	adminRouter.Dispatch(ctx, subCmd)
}
