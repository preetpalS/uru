// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"log"
	"io/ioutil"
	"reflect"
	"sort"
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
			Description: `ruby 2.1.1p1 (2013-12-27 revision 44443) [x86_64-linux]`,
		},
		{
			ID:          `1.7.9`,
			TagLabel:    `179`,
			Exe:         `jruby`,
			Home:        `C:\Apps\rubies\jruby\bin`,
			GemHome:     ``,
			Description: `jruby 1.7.9 (1.9.3p392) 2013-12-06 87b108a on Java HotSpot(TM) 64-Bit Server VM 1.7.0_45-b18 [Windows 8-amd64]`,
		},
		{
			ID:          `1.7.10`,
			TagLabel:    `1710`,
			Exe:         `jruby`,
			Home:        `C:\Apps\rubies\jruby_new\bin`,
			GemHome:     ``,
			Description: `jruby 1.7.10 (1.9.3p392) 2014-01-09 c4ecd6b on Java HotSpot(TM) 64-Bit Server VM 1.7.0_45-b18 [Windows 8-amd64]`,
		},
	}
	testTagLabels = []string{`211`, `179`, `1710`}
	testTags      = []string{`3577244517`, `444332046`, `3091568265`}
)

func init() {
	// silence any logging done in the package files
	log.SetOutput(ioutil.Discard)
}

func TestNewTag(t *testing.T) {
	ctx := NewContext()

	for i, rb := range testRubies {
		rv, _ := NewTag(ctx, rb)
		if rv != testTags[i] {
			t.Errorf("NewTag not returning correct value\n  want: `%v`\n  got: `%v`",
				testTags[i],
				rv)
		}
	}
}

func TestTagLabelToTag(t *testing.T) {
	ctx := NewContext()
	ctx.Registry = RubyRegistry{
		Version: RubyRegistryVersion,
		Rubies: RubyMap{
			testTags[0]: testRubies[0],
			testTags[1]: testRubies[1],
			testTags[2]: testRubies[2],
		},
	}

	// nonexistent tag label test
	tags, err := TagLabelToTag(ctx, `200`)
	if err == nil {
		t.Error("TagLabelToTag() should return error for nonexistent tag label")
	}

	// valid tag label tests
	for i, rb := range testRubies {
		tags, err = TagLabelToTag(ctx, testTagLabels[i])
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

func TestTagInfoSorter(t *testing.T) {
	ti := []tagInfo{
		{testTags[0], testTagLabels[0]},
		{testTags[1], testTagLabels[1]},
		{testTags[2], testTagLabels[2]},
	}
	tis := &tagInfoSorter{ti}

	sort.Sort(tis)
	if !sort.IsSorted(tis) {
		t.Error("Unable to sort tagInfoSorter")
	}

	expected := []string{`3091568265`, `444332046`, `3577244517`}
	actual := []string{ti[0].tag, ti[1].tag, ti[2].tag}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("tagInfoSorter incorrectly sorted\n  want: `%v`\n  got: `%v`",
			expected, actual)
	}
}

func TestSortTagsByTagLabel(t *testing.T) {
	rubyMap := &RubyMap{
		testTags[0]: testRubies[0],
		testTags[1]: testRubies[1],
		testTags[2]: testRubies[2],
	}

	expected := []string{`3091568265`, `444332046`, `3577244517`}
	actual, err := SortTagsByTagLabel(rubyMap)
	if err != nil {
		t.Error("SortTagsByTagLabel() should not return error for valid RubyMap")
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("SortTagsByTagLabel() not returning correct value\n  want: `%v`\n  got: `%v`",
			expected, actual)
	}
}
