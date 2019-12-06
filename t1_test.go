package dblib

import (
	"testing"

	"github.com/alrusov/misc"
)

//--------------------------------------------------------------------------------------------------------------------------------------------------------------//

func TestSecureString(t *testing.T) {
	data := make(misc.StringMap)
	data[`aaa`] = `aaa`
	data[`\`] = `\\`
	data[`'`] = `''`
	data[`'\'`] = `''\\''`
	data[`'12\34\\56'78''9'`] = `''12\\34\\\\56''78''''9''`

	for src, goal := range data {
		res := SecureString(src)
		if res != goal {
			t.Errorf(`SecureString(%q): expect %q, got %q`, src, goal, res)
		}
	}
}

//--------------------------------------------------------------------------------------------------------------------------------------------------------------//
