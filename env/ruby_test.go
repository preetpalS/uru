// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package env

import (
	"testing"
)

type rubyInfo struct {
	VersionString string
	Exe string
	Version string
	PatchLevel string
}

var rubies = map[string]rubyInfo {
	`ruby-windows-187`: rubyInfo{
		`ruby 1.8.7 (2012-10-12 patchlevel 371) [i386-mingw32]`,
		`ruby`,
		`1.8.7`,
		`371`,
	},
	`ruby-darwin-187`: rubyInfo{
		`ruby 1.8.7 (2009-06-12 patchlevel 174) [universal-darwin10.0]`,
		`ruby`,
		`1.8.7`,
		`174`,
	},
	`ruby-windows-193`: rubyInfo{
		`ruby 1.9.3p430 (2013-05-15 revision 40754) [i386-mingw32]`,
		`ruby`,
		`1.9.3`,
		`p430`,
	},
	`ruby-windows-200`: rubyInfo{
		`ruby 2.0.0p197 (2013-05-20 revision 40843) [i386-mingw32]`,
		`ruby`,
		`2.0.0`,
		`p197`,
	},
	`ruby-linux-200`: rubyInfo{
		`ruby 2.0.0p197 (2013-05-20 revision 40843) [i686-linux]`,
		`ruby`,
		`2.0.0`,
		`p197`,
	},
	`ruby-darwin-200`: rubyInfo{
		`ruby 2.0.0p197 (2013-05-20 revision 40843) [i386-darwin10.8.0]`,
		`ruby`,
		`2.0.0`,
		`p197`,
	},
	`jruby-windows-174`: rubyInfo{
		`jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Client VM 1.7.0_21-b11 +indy [Windows 7-x86]`,
		`jruby`,
		`1.7.4`,
		``,
	},
	`jruby-linux-174`: rubyInfo{
		`jruby 1.7.4 (1.9.3p392) 2013-05-16 2390d3b on Java HotSpot(TM) Server VM 1.7.0_21-b11 [linux-i386]`,
		`jruby`,
		`1.7.4`,
		``,
	},
	`ruby-linux-dev`: rubyInfo{
		`ruby 2.1.0dev (2013-05-25 trunk 40932) [i686-linux]`,
		`ruby`,
		`2.1.0`,
		`dev`,
	},
}

func TestRubyRegex(t *testing.T) {

	for _, ri := range rubies {
		matches := rbRegex.FindStringSubmatch(ri.VersionString)
		if matches == nil {
			t.Error("ruby regex did not match full ruby version string")
		}

		if matches[1] != ri.Exe {
			t.Errorf("ruby regex did not match ruby executable string\n  want: `%s`\n  got: `%s`",
				ri.Exe,
				matches[1])
		}
		if matches[2] != ri.Version {
			t.Errorf("ruby regex did not match ruby version string\n  want: `%s`\n  got: `%s`",
				ri.Version,
				matches[2])
		}
		if matches[3] != ri.PatchLevel && matches[4] == ``{
			t.Errorf("ruby regex did not match ruby patchlevel string\n  want: `%s`\n  got: `%s`",
				ri.PatchLevel,
				matches[3])
		}
		if matches[4] != `` && matches[4] != ri.PatchLevel {
			t.Errorf("ruby regex did not match ruby patchlevel string\n  want: `%s`\n  got: `%s`",
				ri.PatchLevel,
				matches[4])
		}
	}

}
