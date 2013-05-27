// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/jonforums/uru/env"
)

func Help(ctx *env.Context) {
	cmdArgs := ctx.CmdArgs()
	if len(cmdArgs) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] CMD ARG ...\n", env.AppName)
		fmt.Fprintln(os.Stderr, "\nwhere [options] are:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nwhere CMD is one of:")
		PrintCommandSummary()
		fmt.Fprintf(os.Stderr, "\nfor help on a particular command, type `%s help CMD`\n",
			env.AppName)
	} else {
		commandHelp(cmdArgs[0])
	}

	os.Exit(0)
}

func PrintCommandSummary() {
	for k, v := range CommandRegistry {
		fmt.Fprintf(os.Stderr, "%6.6s   %s\n", k, v.HelpMsg)
	}

	fmt.Fprintf(os.Stderr, "%6.6s   %s\n",
		"TAG", "switch to use ruby version TAG")
}

func printAdminCommandSummary() {
	fmt.Fprintln(os.Stderr, "\nwhere SUBCMD is one of:")
	for k, v := range AdminCmdRegistry {
		fmt.Fprintf(os.Stderr, "%8.8s   %s\n", k, v.HelpMsg)
		if v.Aliases != nil {
			fmt.Fprintf(os.Stderr, "%8.8s   aliases: %s\n", "", strings.Join(v.Aliases, ", "))
		}
		fmt.Fprintf(os.Stderr, "%8.8s   usage: %s\n", "", v.Usage)
		fmt.Fprintf(os.Stderr, "%8.8s   eg: %s %s\n\n", "", env.AppName, v.Eg)
	}
}

func commandHelp(cmd string) {
	command, ok := CommandRegistry[cmd]
	if !ok {
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
		command.HelpMsg,
		command.Usage,
		env.AppName, command.Eg)

	if cmd == `admin` {
		printAdminCommandSummary()
	}
}
