// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package main

import (
	"testing"

	"bitbucket.org/jonforums/uru/env"
)

var testCommands = []struct {
	cmd   string
	input string
	match bool
}{
	{"admin", "admin", true},
	{"admin", "administer", false},
	{"add", "add", true},
	{"add", "addme", false},
	{"gem", "gem", true},
	{"gem", "gem2", false},
	{"gemset", "gemset", true},
	{"gemset", "gs", true},
	{"gemset", "gemsets", false},
	{"help", "help", true},
	{"help", "h", false},
	{"install", "install", true},
	{"install", "in", true},
	{"install", "inme", false},
	{"ls", "list", true},
	{"ls", "ls", true},
	{"ls", "listme", false},
	{"refresh", "refresh", true},
	{"refresh", "ref", false},
	{"retag", "retag", true},
	{"retag", "tag", true},
	{"retag", "tag_", false},
	{"rm", "rm", true},
	{"rm", "del", true},
	{"rm", "delete", false},
	{"ruby", "ruby", true},
	{"ruby", "rb", true},
	{"ruby", "ruby2", false},
	{"version", "version", true},
	{"version", "ver", true},
	{"version", "vers", false},
}

func TestInitCommandParser(t *testing.T) {
	ctx := env.NewContext()
	initCommandParser(ctx)

	for _, v := range testCommands {
		actual := ctx.CmdRegex(v.cmd).MatchString(v.input)
		if actual != v.match {
			t.Errorf("`%s` input did not match regex `%v`",
				v.input, ctx.CmdRegex(v.cmd))
		}
	}
}
