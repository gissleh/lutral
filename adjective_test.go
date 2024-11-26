package lutral

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGenerateAdjective(t *testing.T) {
	var table = []struct {
		ID        string
		Adjective string
		Test      string
		Result    string
	}{
		{
			ID: "9116:adj.", Adjective: "n;<us,awn>;rr",
			Test: "anusrr", Result: "9116:adj. a- <us>",
		},
		{
			ID: "9116:adj.", Adjective: "n;<us,awn>;rr",
			Test: "nawnrra", Result: "9116:adj. <awn> -a",
		},
		{
			ID: "392:adj.", Adjective: "fm;<us,awn>;etok",
			Test: "fmawnetok", Result: "392:adj. <awn>",
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s %s", row.Adjective, row.Test), func(t *testing.T) {
			tree := GenerateAdjective(*BuildTree(strings.Split(row.Adjective, ";")...)).AndThen(ParseNode("=" + row.ID))
			r := Runner{Root: tree}
			res := r.Run(row.Test)

			resStr := ""
			for _, res := range res {
				if len(resStr) > 0 {
					resStr += ";"
				}
				resStr += res.String()
			}

			assert.Equal(t, row.Result, resStr)

			if t.Failed() {
				j, _ := json.MarshalIndent(tree, "", "  ")
				t.Log(string(j))
			}
		})
	}
}

func TestGenerateAdjectiveAdverb(t *testing.T) {
	var table = []struct {
		ID        string
		Adjective string
		Test      string
		Result    string
	}{
		{
			ID: "500:adv.", Adjective: "ftue",
			Test: "nìftue", Result: "500:adv. nì-",
		},
		{
			ID: "9144:adv.", Adjective: "fyin",
			Test: "nìfyin", Result: "9144:adv. nì-",
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s %s", row.Adjective, row.Test), func(t *testing.T) {
			tree := GenerateAdjectiveAdverb(*BuildTree(strings.Split(row.Adjective, ";")...)).AndThen(ParseNode("=" + row.ID))
			r := Runner{Root: tree}
			res := r.Run(row.Test)

			resStr := ""
			for _, res := range res {
				if len(resStr) > 0 {
					resStr += ";"
				}
				resStr += res.String()
			}

			assert.Equal(t, row.Result, resStr)

			if t.Failed() {
				j, _ := json.MarshalIndent(tree, "", "  ")
				t.Log(string(j))
			}
		})
	}
}
