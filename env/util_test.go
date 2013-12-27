// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestStringSplitPath(t *testing.T) {
	prevPath := os.Getenv(`PATH`)
	newPath := []string{`A`, `B`, `C`, `D`}

	os.Setenv(`PATH`, strings.Join(newPath, string(os.PathListSeparator)))
	rv, _ := StringSplitPath()
	if !reflect.DeepEqual(rv, newPath) {
		t.Errorf("StringSplitPath not returning correct value\n  want: `%v`\n  got: `%v`",
			newPath,
			rv)
	}

	os.Setenv(`PATH`, prevPath)
}

func TestNewTag(t *testing.T) {
	ctx := NewContext()
	testRuby := Ruby{
		ID:          `2.1.1-p1`,
		TagLabel:    `211p1`,
		Exe:         `ruby`,
		Home:        `/home/fake/.rubies/ruby-2.1.0/bin`,
		GemHome:     `/home/fake/.gem/ruby/2.1.0`,
		Description: `ruby 2.1.1p1 (2013-12-27 revision 44443) [i386-mingw32]`,
	}
	tag := `2292638747`

	rv, _ := NewTag(ctx, testRuby)
	if rv != tag {
		t.Errorf("NewTag not returning correct value\n  want: `%v`\n  got: `%v`",
			tag,
			rv)
	}
}
