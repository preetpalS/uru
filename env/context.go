// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"fmt"
	"os"
	"regexp"
)

const RubyRegistryVersion = `1.0.0`

type RubyMap map[string]Ruby

type RubyRegistry struct {
	Version string
	Rubies  RubyMap
}

type Context struct {
	commandRegex map[string]*regexp.Regexp
	home         string
	command      string
	commandArgs  []string

	marshaller   *RubyMarshaller

	Registry RubyRegistry
}

func (c *Context) CmdRegex(cmd string) *regexp.Regexp {
	if c.commandRegex == nil {
		panic("Context has not been initialized")
	}

	return c.commandRegex[cmd]
}
func (c *Context) SetCmdRegex(cmd string, r string) (err error) {
	if c.commandRegex == nil {
		panic("Context has not been initialized")
	}

	c.commandRegex[cmd], err = regexp.Compile(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to compile `%s` regexp", cmd)
	}
	return
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

func (c *Context) Marshaller() *RubyMarshaller {
	return c.marshaller
}

func (c *Context) SetMarshaller(m *RubyMarshaller) {
	c.marshaller = m
}

func NewContext() *Context {
	return &Context{
		commandRegex: make(map[string]*regexp.Regexp),
		marshaller: NewRubyMarshaller(marshalRubies),
		Registry: RubyRegistry{
			Version: RubyRegistryVersion,
			Rubies:  make(RubyMap, 4),
		},
	}
}
