// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// CopyFile copies a source file to a destination file.
func CopyFile(dst, src string) (written int64, err error) {
	sf, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer df.Close()

	return io.Copy(df, sf)
}

// StringSplitPath splits the PATH env var into a slice of strings.
func StringSplitPath() (path []string, err error) {
	rawPath := os.Getenv(`PATH`)
	if rawPath == `` {
		return nil, errors.New("error: unable to get PATH env var value")
	}

	path = strings.Split(rawPath, string(os.PathListSeparator))

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

	newRb := ctx.Rubies[tag]
	tail := strings.Split(tmp, string(os.PathListSeparator))

	if SysRbRegex.MatchString(tag) {
		// system ruby already on base PATH so set new PATH to base PATH
		path = tail
	} else {
		gemBinDir := filepath.Join(newRb.GemHome, `bin`)
		if runtime.GOOS == `windows` {
			// assume windows users always install gems to ruby installation
			// so do not prepend a generated GEM_HOME bindir to PATH
			gemBinDir = ``
		}

		// prepend base PATH with computed GEM_HOME/bin, new ruby PATH,
		// and a canary initiator
		head := []string{gemBinDir, newRb.Home, string(os.PathListSeparator)}
		path = append(head, tail...)
	}
	log.Printf("[DEBUG] %v\n", path)

	return
}
