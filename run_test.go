package lutral

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func buildTestTree() *Node {
	return CombineTrees(
		BuildTree("fm", "<0>", "<1>", "et", "<2>", "ok", "=392"),
		BuildTree("fm", "<0>", "<1>", "<2>", "i", "=396"),
		BuildTree("fm", "<0>", "<1>", "<2>", "al", "=3700"),
		BuildTree("f", "<0>", "<1>", "rrf", "<2>", "en", "=464"),
		BuildTree("s", "<0>", "<1>", "<2>", "i", "=1788"),
		BuildTree("$np", "tì-", "fm", "<us>", "etok", "=392:n."),
		BuildTree("$np", "fmetok", "-yu", "=392:n."),
		BuildTree("$np", "uvan", "$nsmod|$ncec", "=2644"),
		BuildTree("$np", "uvan", " s", "<0>", "<1>", "<2>", "i", "=2648"),
		BuildTree("$np", "tìtaron", "=7336"),
		BuildTree("sìk", "=1796"),
		BuildTree("tìk", "=13294"),
		BuildTree("tsìk", "=8280"),
		BuildTree("ma", "=1056"),
		BuildTree("$np", "'eylan", "=56"),
		BuildTree("z", "<0>", "<1>", "<2>", "o", "=6420"),
		BuildTree("'", "<0>", "<1>", "rrk", "<2>", "o", "=10496"),
		BuildTree("$np", "fe'ranvi", "$nsmod|$ncec", "=9256"),
		BuildTree("$np", "fe'ran", "$nsmod|$ncec", "=9248"),
		BuildTree("$np", "'e-", "lì'u", "$nsmod|$ncec", "=-1"), // fake prefix for coverage
		BuildTree("tsuk-", "k", "anom", "=8392"),
		BuildTree("'a", "sap", "$nsmod_si", " ", "s", "<1>", "i", "=12962"),
		BuildTree("a-", "$np", "'asap-s", "<us>", "i", "=12962:adj."),
	)
}

func TestRunner_Run(t *testing.T) {
	dict := buildTestTree()

	table := []struct {
		Lookup  string
		Results []string
	}{
		{"fmetok", []string{"392"}},
		{"fmeretok", []string{"392 <er>"}},
		{"fmeyketok", []string{"392 <eyk>"}},
		{"fmeteiok", []string{"392 <ei>"}},
		{"fmäpeykìlmetängok", []string{"392 <äpeyk,ìlm,äng>"}},
		{"folrrfen", []string{"464 <ol>"}},
		{"ferrrfen", []string{}},
		{"seiyi", []string{"1788 <eiy>"}},
		{"frrfeiyen", []string{}},
		{"tìfmusetok", []string{"392:n. tì- <us>"}},
		{"sìtaron", []string{"7336 t→s"}},
		{"fraysìtaron", []string{"7336 fra-ay- t→s"}},
		{"saysìtaron", []string{"7336 tsa-ay- ts→s,t→s"}},
		{"sìfmusetok", []string{"392:n. tì- <us> t→s"}},
		{"saysìfmusetok", []string{"392:n. tsa-ay-tì- <us> ts→s,t→s"}},
		{"fepesìfmusetok", []string{"392:n. pe-pxe-tì- <us> p→f,px→p,t→s"}},
		{"fìuvan", []string{"2644 fì-"}},
		{"fneuvan", []string{"2644 fne-"}},
		{"fìfneuvan", []string{"2644 fì-fne-"}},
		{"ayuvan", []string{"2644 ay-"}},
		{"frayuvan", []string{"2644 fra-ay-"}},
		{"sayuvan", []string{"2644 tsa-ay- ts→s"}},
		{"uvanä", []string{"2644 -ä"}},
		{"uvan si", []string{"2644 + si", "2648"}},
		{"sayfneuvanti", []string{"2644 tsa-ay-fne- -ti ts→s"}},
		{"sìk", []string{"1796", "13294 t→s", "8280 ts→s"}},
		{"'eylan", []string{"56"}},
		{"meeylan", []string{"56 me- 'e→e"}},
		{"pxeeylan", []string{"56 pxe- 'e→e"}},
		{"meylan", []string{"56 me- 'e→e"}},
		{"peylan", []string{"56 pxe- px→p,'e→e", "56 pe- 'e→e"}},
		{"pxeylan", []string{"56 pxe- 'e→e"}},
		{"zeyko", []string{"6420 <eyk>"}},
		{"rrko", []string{}},
		{"eykrrko", []string{"10496 <eyk> 'e→e"}},
		{"zeykeyko", []string{}},
		{"zerero", []string{}},
		{"fe'ran", []string{"9248"}},
		{"fe'ranvi", []string{"9256"}},
		{"tsukkanom", []string{"8392 tsuk-"}},
		{"tsukanom", []string{"8392 tsuk-"}},
		{"'asap soli", []string{"12962 <ol>"}},
		{"ayasap-susi", []string{"12962:adj. a-ay- <us> 'a→a"}},
		{"melì'u", []string{"-1 me-'e- 'e→e"}},
	}

	for _, row := range table {
		t.Run(row.Lookup, func(t *testing.T) {
			runner := &Runner{Root: dict}
			before := time.Now()
			results := runner.Run(row.Lookup)
			duration := time.Now().Sub(before)
			resStr := make([]string, 0, len(results))
			for _, result := range results {
				resStr = append(resStr, result.String())
			}

			assert.Equal(t, row.Results, resStr)
			t.Log("Step Count:", runner.StepCount)
			t.Log("Comparison Count:", runner.SubStepCount)
			t.Log("Cold Time:", duration)
		})
	}
}

func TestRunner_Extract(t *testing.T) {
	dict := buildTestTree()
	table := []struct {
		Lookup  string
		Results []string
	}{
		{"fmetok fìuvanti, ma eylan", []string{
			"[1] 392",
			"[2] 2644 fì- -ti",
			"[3] 1056",
			"[4] 56 'e→e",
		}},
		{"fmetok fìkeyawralì'uti, ma eylan", []string{
			"[1] 392",
			"[3] 1056",
			"[4] 56 'e→e",
		}},
		{"tsauvan seri", []string{
			"[1] 2648 tsa- <er>",
		}},
		{"blerg?!", []string{}},
	}

	for _, row := range table {
		t.Run(row.Lookup, func(t *testing.T) {
			runner := &Runner{Root: dict}
			res := runner.Extract(row.Lookup)
			resStr := make([]string, 0, len(res))
			for _, result := range res {
				resStr = append(resStr, result.String())
			}
			assert.Equal(t, row.Results, resStr)
		})
	}
}

func TestRunner_runStep_Panic(t *testing.T) {
	assert.Panics(t, func() {
		runner := &Runner{}
		runner.runStep(&Node{Kind: NKReturn}, "glurb", 0, "", nil)
	})

	assert.Panics(t, func() {
		runner := &Runner{SubtreeMap: map[string]*Node{}}
		runner.runStep(&Node{Kind: NKSubTree}, "non_existent", 0, "", nil)
	})
}

func BenchmarkGenerateInitialSubTreeMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		gen := GenerateInitialSubTreeMap()
		if gen == nil {
			b.Fail()
		}
	}
}
