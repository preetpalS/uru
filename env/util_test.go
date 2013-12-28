// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	testRubies = []Ruby{
		{
			ID:          `2.1.1-p1`,
			TagLabel:    `211p1`,
			Exe:         `ruby`,
			Home:        `/home/fake/.rubies/ruby-2.1.0/bin`,
			GemHome:     `/home/fake/.gem/ruby/2.1.0`,
			Description: `ruby 2.1.1p1 (2013-12-27 revision 44443) [i386-mingw32]`,
		},
		{
			ID:          `1.7.9`,
			TagLabel:    `179`,
			Exe:         `jruby`,
			Home:        `C:\Apps\rubies\jruby\bin`,
			GemHome:     ``,
			Description: `jruby 1.7.9 (1.9.3p392) 2013-12-06 87b108a on Java HotSpot(TM) 64-Bit Server VM 1.7.0_45-b18 [Windows 8-amd64]`,
		},
	}
	testLabels = []string{`211`, `179`}
	testTags   = []string{`2292638747`, `444332046`}
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

	rv, _ := NewTag(ctx, testRubies[0])
	if rv != testTags[0] {
		t.Errorf("NewTag not returning correct value\n  want: `%v`\n  got: `%v`",
			testTags[0],
			rv)
	}
}

func TestTagLabelToTag(t *testing.T) {
	ctx := NewContext()
	ctx.Registry = RubyRegistry{
		Version: RubyRegistryVersion,
		Rubies: RubyMap{
			testTags[0]: testRubies[0],
			testTags[1]: testRubies[1],
		},
	}

	// nonexistent tag label test
	tags, err := TagLabelToTag(ctx, `200`)
	if err == nil {
		t.Error("TagLabelToTag() should return error for nonexistent tag label")
	}

	// valid tag label tests
	for i, rb := range testRubies {
		tags, err = TagLabelToTag(ctx, testLabels[i])
		if err != nil {
			t.Error("TagLabelToTag() should not return error for valid tag label")
		}
		if tags[testTags[i]].ID != rb.ID {
			t.Errorf("TagLabelToTag() not returning correct value\n  want: `%v`\n  got: `%v`",
				rb.ID,
				tags[testTags[i]].ID)
		}
	}
}
