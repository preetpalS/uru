// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/jonforums/uru/env"
)

var helpCmd *Command = &Command{
	Name:    "help",
	Aliases: []string{"help"},
	Usage:   "help",
	Eg:      "help",
	Short:   "help",
	Run:     help,
}

func init() {
	CmdRouter.Handle(helpCmd.Aliases, helpCmd)
}

func help(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	if len(cmdArgs) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] CMD ARG ...\n", env.AppName)
		fmt.Fprintln(os.Stderr, "\nwhere CMD is one of:")
		printCommandSummary()
		fmt.Fprintf(os.Stderr, "\nfor help on a particular command, type `%s help CMD`\n",
			env.AppName)
	} else {
		commandHelp(cmdArgs[0])
	}

	os.Exit(0)
}

func printCommandSummary() {
	for k, v := range *CmdRouter.Commands() {
		fmt.Fprintf(os.Stderr, "%6.6s   %s\n", k, v.Short)
	}
}

func printAdminCommandSummary() {
	fmt.Fprintln(os.Stderr, "\nwhere SUBCMD is one of:")
	for k, v := range *adminRouter.Commands() {
		fmt.Fprintf(os.Stderr, "%8.8s   %s\n", k, v.Short)
		if v.Aliases != nil {
			fmt.Fprintf(os.Stderr, "%8.8s   aliases: %s\n", "", strings.Join(v.Aliases, ", "))
		}
		fmt.Fprintf(os.Stderr, "%8.8s   usage: %s\n", "", v.Usage)
		fmt.Fprintf(os.Stderr, "%8.8s   eg: %s %s\n\n", "", env.AppName, v.Eg)
	}
}

func commandHelp(cmd string) {
	command, err := CmdRouter.Handler(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "---> No help available on `%s`\n", cmd)
		return
	}

	buf := bytes.NewBufferString("  Description: %s\n")
	if command.Aliases != nil {
		buf.WriteString(fmt.Sprintf("  Aliases: %s\n", strings.Join(command.Aliases, ", ")))
	}
	buf.WriteString("  Usage: %s\n  Example: %s %s\n")

	fmt.Fprintf(os.Stderr,
		buf.String(),
		command.Short,
		env.AppName, command.Usage,
		env.AppName, command.Eg)

	if cmd == `admin` {
		printAdminCommandSummary()
	}
}
