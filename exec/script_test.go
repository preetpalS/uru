// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package exec

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"bitbucket.org/jonforums/uru/env"
)

func init() {
	// silence any logging done in the package files
	log.SetOutput(ioutil.Discard)
}

func TestWinPathList2Nix(t *testing.T) {
	pth := []string{`C:\Apps\rubies\ruby-2.1.0\bin`, env.Canary, `C:\some\fake\path`}
	result := strings.Join(winPathToNix(&pth), `:`)

	if strings.Contains(result, `C:`) {
		t.Errorf("Generated *nix path contains Windows volume designator")
	}
	if strings.ContainsRune(result, '\\') {
		t.Errorf(`Generated *nix path contains '\' char`)
	}
	if !strings.Contains(result, env.Canary) {
		t.Errorf("Generated *nix path missing `%s` canary", env.Canary)
	}
}
