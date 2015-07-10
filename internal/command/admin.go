// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"fmt"
	"os"

	"bitbucket.org/jonforums/uru/internal/env"
)

var adminRouter *Router = NewRouter(defAdminHandler)

var adminCmd *Command = &Command{
	Name:    "admin",
	Aliases: []string{"admin"},
	Usage:   "admin SUBCMD ARGS",
	Eg:      `admin add C:\Apps\rubies\ruby-2.1\bin`,
	Short:   "administer uru installation",
	Run:     admin,
}

func init() {
	CmdRouter.Handle(adminCmd.Aliases, adminCmd)
}

func admin(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	if len(cmdArgs) == 0 {
		return
	}
	subCmd := cmdArgs[0]
	ctx.SetCmdArgs(cmdArgs[1:])

	adminRouter.Dispatch(ctx, subCmd)
}

func defAdminHandler(ctx *env.Context) {
	fmt.Printf("[ERROR] I don't understand the `%s` admin sub-command\n\n", ctx.CmdArgs()[0])
	os.Exit(1)
}
