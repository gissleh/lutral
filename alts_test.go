package lutral

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAlternativeSpellings(t *testing.T) {
	table := []struct {
		Str string
		Res string
	}{
		{"skxawng", ""},
		{"kxanì", "ganì"},
		{"syawm", "shawm"},
		{"'awkx", "'awg"},
		{"futa", "fta"},
		{"tsata", ""},
		{"fra'u", "frau"},
		{"ngä'än", "ngään"},
		{"kä'ärìp", "käärìp"},
		{"tìngus<0><1>ä'<2>än", "tìngus<0><1>ä<2>än"},
		{"kxan'epe", "gan'epe"},
		{"mu'ni", ""},
		{"tìsyìmawnun'i", "tsyìmawnun'i,tìshìmawnun'i,chìmawnun'i"},
	}

	for _, row := range table {
		t.Run(row.Str, func(t *testing.T) {
			assert.Equal(t, row.Res, strings.Join(AlternativeSpellings(row.Str), ","))
		})
	}
}
