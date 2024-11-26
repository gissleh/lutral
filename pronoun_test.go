package lutral

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGeneratePronoun(t *testing.T) {
	table := []struct {
		Pronoun string
		Test    string
		Results string
	}{
		{"po", "po", "123"},
		{"po", "peyä", "123 -yä"},
		{"po", "pol", "123 -l"},
		{"po", "poyä", ""},
		{"po", "poä", ""},
		{"oe", "oe", "123"},
		{"oe", "oeyä", "123 -yä"},
		{"oe", "oey", "123 -y"},
		{"nga", "nga", "123"},
		{"nga", "ngey", "123 -y"},
		{"nga", "ngal", "123 -l"},
		{"nga", "ngaru", "123 -ru"},
		{"nga", "menga", "123 me-"},
		{"nga", "mengal", "123 me- -l"},
		{"nga", "ayngeyä", "123 ay- -yä"},
		{"ayoeng", "ayoeng", "123"},
		{"ayoeng", "ayoengal", "123 -l"},
		{"ayoeng", "ayoengìl", ""},
		{"ayoeng", "ayoengati", "123 -ti"},
		{"tsaw", "tsaw", "123"},
		{"tsaw", "tsal", "123 -l"},
		{"tsaw", "tsat", "123 -t"},
		{"tsaw", "tsati", "123 -ti"},
		{"tsaw", "tsari", "123 -ri"},
		{"tsaw", "tsafa", "123 -fa"},
		{"tsari", "tsari", "123"},
		{"tsar", "tsar", "123"},
		{"sno", "sneyä", "123 -yä"},
		{"sneyä", "sneyä", "123"},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s_%s", row.Pronoun, row.Test), func(t *testing.T) {
			tree := GeneratePronoun(row.Pronoun).AndThenResult("123")
			runner := Runner{Root: tree, SubtreeMap: GenerateInitialSubTreeMap()}
			runner.SubtreeMap["nsadp"].MergeFrom(*BuildTree("-fa", "/return"))

			res := runner.Run(row.Test)
			resStr := ""
			for _, res := range res {
				if len(resStr) != 0 {
					resStr += ";"
				}
				resStr += res.String()
			}

			assert.Equal(t, row.Results, resStr)
			if t.Failed() {
				j, _ := json.MarshalIndent(tree, "", "  ")
				t.Log(string(j))
			}
		})
	}
}
