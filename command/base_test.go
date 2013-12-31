// Author: Jon Maken, All Rights Reserved
// License: 3-clause BSD

package command

import (
	"testing"
)

var testGemsetNames = []struct {
	RawName string
	Ruby    string
	Gemset  string
}{
	{RawName: `200p376-x32@gemset`, Ruby: `200p376-x32`, Gemset: `gemset`},
	{RawName: `210@custom_user_gemset`, Ruby: `210`, Gemset: `custom_user_gemset`},
	{RawName: `211@fake-gemset`, Ruby: `211`, Gemset: `fake-gemset`},
}

func TestParseGemsetName(t *testing.T) {
	for _, v := range testGemsetNames {
		ruby, gemset, err := parseGemsetName(v.RawName)
		if err != nil {
			t.Error("parseGemsetName() returned error")
		}
		if ruby != v.Ruby {
			t.Errorf("parseGemsetName() returning incorrect ruby value\n  want: `%v`\n  got: `%v`",
				v.Ruby,
				ruby)
		}
		if gemset != v.Gemset {
			t.Errorf("parseGemsetName() returning incorrect gemset value\n  want: `%v`\n  got: `%v`",
				v.Gemset,
				gemset)
		}
	}
}
