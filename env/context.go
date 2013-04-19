// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"fmt"
	"os"
	"regexp"
)

type RubyRegistry map[string]Ruby

type Context struct {
	commandRegex map[string]*regexp.Regexp
	home         string
	command      string
	commandArgs  []string

	Rubies RubyRegistry
}

func (c *Context) Init() {
	c.commandRegex = make(map[string]*regexp.Regexp)
	c.Rubies = make(RubyRegistry, 4)
}

func (c *Context) CmdRegex(cmd string) *regexp.Regexp {
	if c.commandRegex == nil {
		panic("Context has not been initialized")
	}

	return c.commandRegex[cmd]
}
func (c *Context) SetCmdRegex(cmd string, r string) error {
	if c.commandRegex == nil {
		panic("Context has not been initialized")
	}

	var e error
	c.commandRegex[cmd], e = regexp.Compile(r)
	if e != nil {
		fmt.Fprintf(os.Stderr, "error: unable to compile `%s` regexp", cmd)
	}
	return e
}

func (c *Context) Home() string {
	return c.home
}
func (c *Context) SetHome(h string) {
	c.home = h
}

func (c *Context) Cmd() string {
	return c.command
}
func (c *Context) SetCmd(cmd string) {
	c.command = cmd
}

func (c *Context) CmdArgs() []string {
	return c.commandArgs
}
func (c *Context) SetCmdArgs(args []string) {
	c.commandArgs = args
}

func (c *Context) SetCmdAndArgs(cmd string, args []string) {
	c.command = cmd
	c.commandArgs = args
}
