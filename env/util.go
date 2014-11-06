// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type tagInfo struct {
	tag      string // unique internal identifier for a particular ruby
	tagLabel string // modifiable, user friendly name for a particular ruby
}

// tagInfoSorter sorts slices of tagInfo structs by implementing sort.Interface by
// providing Len, Swap, and Less
type tagInfoSorter struct {
	tags []tagInfo
}

func (s *tagInfoSorter) Len() int {
	return len(s.tags)
}

func (s *tagInfoSorter) Swap(i, j int) {
	s.tags[i], s.tags[j] = s.tags[j], s.tags[i]
}

func (s *tagInfoSorter) Less(i, j int) bool {
	return s.tags[i].tagLabel < s.tags[j].tagLabel
}

// CopyFile copies a source file to a destination file.
func CopyFile(dst, src string) (written int64, err error) {
	sf, err := os.Open(src)
	if err != nil {
		return
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return
	}
	defer df.Close()

	written, err = io.Copy(df, sf)

	log.Printf("[DEBUG] copied file\n  src: %s\n  dst: %s\n  bytes copied: %d\n",
		src, dst, written)

	return
}

// StringSplitPath splits the PATH env var into a slice of strings.
func StringSplitPath() (path []string, err error) {
	rawPath := os.Getenv(`PATH`)
	if rawPath == `` {
		return nil, errors.New("unable to get PATH env var value")
	}

	path = strings.Split(rawPath, string(os.PathListSeparator))

	return
}

// NewTag generates a new tag value used to identify a specific ruby.
func NewTag(ctx *Context, rb Ruby) (tag string, err error) {
	hash := fnv.New32a()
	b := bytes.NewBufferString(fmt.Sprintf("%s%s", rb.Description, rb.Home))

	_, err = hash.Write(b.Bytes())

	return fmt.Sprintf("%d", hash.Sum32()), err
}

// TagLabelToTag returns a map of registered ruby tags whose TagLabel's match that
// of the specified tag label string.
func TagLabelToTag(ctx *Context, label string) (tags RubyMap, err error) {
	tags = make(RubyMap, 4)

	for t, ri := range ctx.Registry.Rubies {
		switch {
		// fuzzy match on TagLabel
		case strings.Contains(ri.TagLabel, label):
			tags[t] = ri
		// full match on ID
		case label == ri.ID:
			tags[t] = ri
		}
	}
	if len(tags) == 0 {
		return nil, errors.New(fmt.Sprintf("---> unable to find ruby matching `%s`\n", label))
	}
	log.Printf("[DEBUG] tags matching `%s`\n%v\n", label, tags)

	return
}

// PathListForTag returns a PATH list appropriate for a given ruby tag.
func PathListForTag(ctx *Context, tag string) (path []string, err error) {
	// get current PATH and split it on the canary separator demarcating the
	// head (current ruby path) and tail (base path):
	//
	//   C:\ruby\bin;;;C:\other;D:\more -or- /.rubies/193/bin:::/other:/more
	//              ^^^                                      ^^^
	envPath := os.Getenv(`PATH`)
	if envPath == `` {
		return nil, errors.New("unable to get PATH env var value")
	}

	paths := strings.Split(envPath, Canary)

	// create the new PATH list by prepending the new ruby PATH and a canary
	// separator to the base PATH unless the new ruby is the system ruby
	var tmp string
	switch len(paths) {
	case 1:
		tmp = paths[0] // PATH is the original PATH
	case 2:
		tmp = paths[1]
	}

	newRb := ctx.Registry.Rubies[tag]
	tail := strings.Split(tmp, string(os.PathListSeparator))

	if SysRbRegex.MatchString(newRb.TagLabel) {
		// system ruby already on base PATH so set new PATH to base PATH
		path = tail
	} else {
		// prepend base PATH with computed GEM_HOME/bin, new ruby PATH,
		// and a canary initiator
		gemBinDir := filepath.Join(newRb.GemHome, `bin`)
		head := []string{gemBinDir, newRb.Home, string(os.PathListSeparator)}

		if runtime.GOOS == `windows` {
			// assume windows users always install gems to ruby installation
			// so do not prepend a generated GEM_HOME bindir to PATH
			head = []string{newRb.Home, string(os.PathListSeparator)}
		}

		path = append(head, tail...)
	}
	log.Printf("[DEBUG] === path list ===\n  %v\n", path)

	return
}

// SortTagsByTagLabel returns a string slice of tags sorted by tag label.
func SortTagsByTagLabel(rubyMap *RubyMap) (tags []string, err error) {
	if len(*rubyMap) == 0 {
		return nil, errors.New("nothing in input RubyMap; no sorted tags to return")
	}

	tis := new(tagInfoSorter)
	tis.tags = []tagInfo{}
	for t, ri := range *rubyMap {
		tis.tags = append(tis.tags, tagInfo{tag: t, tagLabel: ri.TagLabel})
	}
	sort.Sort(tis)

	for _, ti := range tis.tags {
		tags = append(tags, ti.tag)
	}
	if len(tags) == 0 {
		return nil, errors.New("no sorted tags to return")
	}

	return
}
