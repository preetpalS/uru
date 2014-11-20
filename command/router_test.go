// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"regexp"
	"testing"

	"bitbucket.org/jonforums/uru/env"
)

func TestRouterConfig(t *testing.T) {
	r := NewRouter(func(ctx *env.Context) {})
	r.Handle([]string{`gem`}, func(ctx *env.Context) {})
	r.Handle([]string{`ls`, `list`}, func(ctx *env.Context) {})

	count := 3

	if r.defHandler == nil {
		t.Error("CommandRouter default handler is nil")
	}
	if num := len(r.handlers); num != count {
		t.Errorf("Incorrect CommandRouter handler count\n  want: `%v`\n  got: `%v`\n",
			count,
			num)
	}
}

func BenchmarkRegexCompare(b *testing.B) {
	r, _ := regexp.Compile("gem")
	for i := 0; i < b.N; i++ {
		switch {
		case r.MatchString("foo"):
			break
		default:
			break
		}
	}
}

func BenchmarkStringCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switch {
		case "gem" == "foo":
			break
		default:
			break
		}
	}
}

func BenchmarkMultiStringCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if "gem" == "foo" || "bar" == "baz" {
		}
	}
}

func BenchmarkCommandRouter(b *testing.B) {
	ctx := env.NewContext()
	cmds := []string{"admin", "gem", "help", "ls", "ruby", "version", "215"}

	r := NewRouter(func (*env.Context) {})
	r.Handle([]string{`admin`}, func(ctx *env.Context) {})
	r.Handle([]string{`gem`}, func(ctx *env.Context) {})
	r.Handle([]string{`help`}, func(ctx *env.Context) {})
	r.Handle([]string{`ls`, `list`}, func(ctx *env.Context) {})
	r.Handle([]string{`ruby`, `rb`}, func(ctx *env.Context) {})
	r.Handle([]string{`ver`, `version`}, func(ctx *env.Context) {})

	for i := 0; i < b.N; i++ {
		for _, c := range cmds {
			r.Dispatch(ctx, c)
		}
	}
}
