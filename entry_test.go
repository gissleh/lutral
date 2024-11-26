package lutral

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseEntry(t *testing.T) {
	table := []struct {
		Input    string
		Expected *Entry
	}{
		{"2140:tìfmetok:n.", &Entry{ID: "2140", Word: "tìfmetok", PoS: []string{"n."}, InfixPositions: nil}},
		{"392:fm<0><1>et<2>ok:vtr.", &Entry{ID: "392", Word: "fmetok", PoS: []string{"vtr."}, InfixPositions: &[2]int{2, 4}}},
		{"2232:t<0><1><2>ok:vtr.", &Entry{ID: "2232", Word: "tok", PoS: []string{"vtr."}, InfixPositions: &[2]int{1, 1}}},
		{"11720:kelnì:prop.n.:loanword", &Entry{ID: "11720", Word: "kelnì", PoS: []string{"prop.n."}, InfixPositions: nil, Flags: []string{"loanword"}}},
		{"2232:t<0><1><2>ok", nil},
		{"2232:t<0><1><2>ok:", nil},
		{"2232::vtr.", nil},
		{"2232", nil},
	}

	for _, row := range table {
		t.Run(row.Input, func(t *testing.T) {
			assert.Equal(t, row.Expected, ParseEntry(row.Input))
		})
	}
}
