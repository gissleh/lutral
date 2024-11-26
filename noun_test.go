package lutral

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNounFromEntry(t *testing.T) {
	table := []struct {
		Entry    string
		Test     string
		Expected string
	}{
		{
			Entry: "604:ikran:n.", Test: "ikran",
			Expected: "604",
		},
		{
			Entry: "604:ikran:n.", Test: "ayfneikranur",
			Expected: "604 ay-fne- -ur",
		},
		{
			Entry: "604:ikran:n.", Test: "ikranti",
			Expected: "604 -ti",
		},
		{
			Entry: "616:irayo:n.,intj.", Test: "irayoru",
			Expected: "616:n. -ru",
		},
		{
			Entry: "5476:hametsì:n.:loanword", Test: "hametsì",
			Expected: "5476",
		},
		{
			Entry: "5476:hametsì:n.:loanword", Test: "hametsur",
			Expected: "5476 -ur",
		},
		{
			Entry: "5476:hametsì:n.:loanword", Test: "hametsit",
			Expected: "5476 -it",
		},
		{
			Entry: "5476:hametsì:n.:loanword", Test: "hametìti",
			Expected: "",
		},
		{
			Entry: "13492:txawnulsrung a tswayon:n.", Test: "txawnulsrung a tswayon",
			Expected: "13492",
		},
		{
			Entry: "13492:txawnulsrung a tswayon:n.", Test: "txawnulsrungit a tswayon",
			Expected: "13492 -it",
		},
		{
			Entry: "13492:txawnulsrung a tswayon:n.", Test: "petawnulsrungit a tswayon",
			Expected: "13492 pxe- -it px→p,tx→t;13492 pe- -it tx→t",
		},
		{
			Entry: "13537:mo letrrtrr:n.", Test: "mori letrrtrr",
			Expected: "13537 -ri",
		},
		{
			Entry: "13537:mo letrrtrr:n.", Test: "fnemo letrrtrr",
			Expected: "13537 fne-",
		},
		{
			Entry: "1928:tsko swizaw:n.", Test: "tsko swizaw",
			Expected: "1928",
		},
		{
			Entry: "1928:tsko swizaw:n.", Test: "tsko swizawti",
			Expected: "1928 -ti",
		},
		{
			Entry: "1928:tsko swizaw:n.", Test: "tskoti swizaw",
			Expected: "",
		},
		{
			Entry: "10352:tìftiatu kifkeyä:n.", Test: "tìftiatu kifkeyä",
			Expected: "10352",
		},
		{
			Entry: "13188:koren ayll:n.", Test: "korenìri ayll",
			Expected: "13188 -ìri",
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s %s", row.Entry, row.Test), func(t *testing.T) {
			entry := ParseEntry(row.Entry)
			if !assert.NotNil(t, entry) {
				return
			}

			tree := NounFromEntry(*entry)
			runner := Runner{Root: tree}
			res := runner.Run(row.Test)

			t.Log("Steps taken:", runner.StepCount)
			t.Log("Comparisons done:", runner.SubStepCount)

			resStr := ""
			for _, res := range res {
				if len(resStr) != 0 {
					resStr += ";"
				}
				resStr += res.String()
			}

			assert.Equal(t, row.Expected, resStr)
			if t.Failed() {
				j, _ := json.MarshalIndent(tree, "", "  ")
				t.Log(string(j))
			}
		})
	}
}

func TestGenerateNoun(t *testing.T) {
	var table = []struct {
		ID     string
		Noun   string
		Test   string
		Result string
	}{
		{
			ID: "2140", Noun: "tìfmetok",
			Test: "tìfmetok", Result: "2140",
		},
		{
			ID: "2140", Noun: "tìfmetok",
			Test: "tìfmetokti", Result: "2140 -ti",
		},
		{
			ID: "2140", Noun: "tìfmetok",
			Test: "faysìfmetokti", Result: "2140 fay- -ti t→s;2140 pay- -ti p→f,t→s",
		},
		{
			ID: "13495", Noun: "säkahena",
			Test: "fnesäkahenatsyìpti", Result: "13495 fne- -tsyìp-ti",
		},
		{
			ID: "13565", Noun: "säpxor",
			Test: "fayfnesäpxorìri", Result: "13565 fay-fne- -ìri;13565 pay-fne- -ìri p→f",
		},
		{
			ID: "4804", Noun: "soaia",
			Test: "soaiä", Result: "4804 -ä",
		},
		{
			ID: "4804", Noun: "soaia",
			Test: "soaiayä", Result: "4804 -yä",
		},
		{
			ID: "4804", Noun: "soaia",
			Test: "soaiaä", Result: "",
		},
		{
			ID: "4804", Noun: "soaia",
			Test: "soaiaru", Result: "4804 -ru",
		},
		{
			ID: "392:n.", Noun: "tì-;fm;<us>;etok",
			Test: "ayfnetìfmusetokur", Result: "392:n. ay-fne-tì- <us> -ur",
		},
		{
			ID: "9116:n.", Noun: "tì-;n;<us>;rr",
			Test: "fepesìnusrrtsyìpperi", Result: "9116:n. pe-pxe-tì- <us> -tsyìp-pe-ri p→f,px→p,t→s",
		},
		{
			ID: "2648:n.", Noun: "uvan;-tswo",
			Test: "peuvantswoä", Result: "2648:n. pxe- -tswo-ä px→p;2648:n. pe- -tswo-ä",
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s %s", row.Noun, row.Test), func(t *testing.T) {
			tree := GenerateNoun(*BuildTree(strings.Split(row.Noun, ";")...)).AndThen(ParseNode("=" + row.ID))
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
