package lutral

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func miniDict() *Dictionary {
	dict := &Dictionary{}
	dict.Insert(*ParseEntry("2548:txo:conj."))
	dict.Insert(*ParseEntry("2224:to:part."))
	dict.Insert(*ParseEntry("616:irayo:intj.,n."))
	dict.Insert(*ParseEntry("2608:uniltìrantokx:n."))
	dict.Insert(*ParseEntry("604:ikran:n."))
	dict.Insert(*ParseEntry("2140:tìfmetok:n."))
	dict.Insert(*ParseEntry("2080:teri:adp."))
	dict.Insert(*ParseEntry("1108:mì:adp."))
	dict.Insert(*ParseEntry("676:ka:adp."))
	dict.Insert(*ParseEntry("-1008:l<0><1><2>ok:vtr.,adp."))
	dict.Insert(*ParseEntry("4468:kxa:n."))
	dict.Insert(*ParseEntry("812:k<0><1><2>in:vtr."))
	dict.Insert(*ParseEntry("2232:t<0><1><2>ok:vtr."))
	dict.Insert(*ParseEntry("2056:t<0><1><2>el:vtr."))
	dict.Insert(*ParseEntry("392:fm<0><1>et<2>ok:vtr."))
	dict.Insert(*ParseEntry("396:fm<0><1><2>i:vtrm."))
	dict.Insert(*ParseEntry("3700:fm<0><1><2>al:vtr."))
	dict.Insert(*ParseEntry("68:'eylan:n."))
	dict.Insert(*ParseEntry("2708:'ewll:n."))
	dict.Insert(*ParseEntry("56:'eveng:n."))
	dict.Insert(*ParseEntry("60:'evi:n."))
	dict.Insert(*ParseEntry("7772:uran:n."))
	dict.Insert(*ParseEntry("2644:uvan:n."))
	dict.Insert(*ParseEntry("13413:ukyom:n."))
	dict.Insert(*ParseEntry("6680:uk:n."))
	dict.Insert(*ParseEntry("1796:sìk:adv."))
	dict.Insert(*ParseEntry("13294:tìk:adv.,conj."))
	dict.Insert(*ParseEntry("8280:tsìk:adv."))
	dict.Insert(*ParseEntry("1056:ma:part."))
	dict.Insert(*ParseEntry("2224:to:part."))
	dict.Insert(*ParseEntry("2548:txo:conj."))
	dict.Insert(*ParseEntry("512:fu:conj."))
	dict.Insert(*ParseEntry("1792:sì:part."))
	dict.Insert(*ParseEntry("1200:ne:adp."))
	dict.Insert(*ParseEntry("13491:txawnulsrung a yur:n."))
	dict.Insert(*ParseEntry("13490:txawnulsrung:n."))
	dict.Insert(*ParseEntry("13495:säkahena:n."))
	dict.Insert(*ParseEntry("13565:säpxor:n."))
	dict.Insert(*ParseEntry("13567:säkeynven:n."))
	dict.Insert(*ParseEntry("13489:säsrung:n."))
	dict.Insert(*ParseEntry("11728:Nìyu Yorkì:prop.n.:loanword"))
	dict.Insert(*ParseEntry("9480:uvan letokx:n."))
	dict.Insert(*ParseEntry("2376:ts<0><1>e'<2>a:vtr."))
	dict.Insert(*ParseEntry("1340:n<0><1>um<2>e:vin."))
	dict.Insert(*ParseEntry("2476:ts<0><1><2>un:vim."))
	dict.Insert(*ParseEntry("13353:tsun:n."))
	dict.Insert(*ParseEntry("12985:tx<0><1><2>ap:vtr."))
	dict.Insert(*ParseEntry("3812:s<0><1><2>ar:vtr."))
	dict.Insert(*ParseEntry("5268:tsaw:pn."))
	dict.Insert(*ParseEntry("13309:tsar:pn."))
	dict.Insert(*ParseEntry("700:k<0><1>am<2>e:vtr."))
	dict.Insert(*ParseEntry("1380:oe:pn."))
	dict.Insert(*ParseEntry("1348:nga:pn."))
	dict.Insert(*ParseEntry("1548:po:pn."))
	dict.Insert(*ParseEntry("192:awnga:pn."))
	dict.Insert(*ParseEntry("6968:sno:pn."))
	dict.Insert(*ParseEntry("11440:fkxara:n."))
	dict.Insert(*ParseEntry("308:eyktan:n."))
	dict.Insert(*ParseEntry("508:ftx<0><1><2>ey:vtr."))
	dict.Insert(*ParseEntry("4456:ftxey:conj."))
	dict.Insert(*ParseEntry("2084:t<0><1>erk<2>up:vin."))
	dict.Insert(*ParseEntry("264:eltu:n."))
	dict.Insert(*ParseEntry("544:h<0><1>ah<2>aw:vin."))
	dict.Insert(*ParseEntry("13458:eltut heykahaw:ph."))
	dict.Insert(*ParseEntry("6520:eltur tìtxen s<0><1><2>i:ph."))
	dict.Insert(*ParseEntry("800:kifkey:n."))
	dict.Insert(*ParseEntry("692:kaltxì:intj."))
	dict.Insert(*ParseEntry("780:kerusey:adj."))
	dict.Insert(*ParseEntry("7752:fe':adj."))
	dict.Insert(*ParseEntry("10124:fe'lup:adj."))
	dict.Insert(*ParseEntry("12963:fe'p<0><1><2>ey:vin."))
	dict.Insert(*ParseEntry("9248:fe'ran:n."))
	dict.Insert(*ParseEntry("9256:fe'ranvi:n."))
	dict.Insert(*ParseEntry("9680:fe'<0><1><2>ul:vin."))
	dict.Insert(*ParseEntry("12962:'asap s<0><1><2>i:vin."))
	dict.Insert(*ParseEntry("68:'eylan:n."))
	dict.Insert(*ParseEntry("76:'<0><1>ì'<2>awn:vin."))
	dict.Insert(*ParseEntry("2708:'ewll:n."))
	dict.Insert(*ParseEntry("20:'awkx:n."))
	dict.Insert(*ParseEntry("8360:'ipu:adj."))
	dict.Insert(*ParseEntry("9032:'rrpxom:n."))
	dict.Insert(*ParseEntry("4368:'awlo:adv."))
	dict.Insert(*ParseEntry("8700:'llngo:n."))
	dict.Insert(*ParseEntry("2744:yerik:n."))
	dict.Insert(*ParseEntry("5312:polpxay:inter.:inter:adj."))
	dict.Insert(*ParseEntry("1524:pesu:inter.:inter:n."))
	dict.Insert(*ParseEntry("1496:pefya:inter.:inter:adv."))
	dict.Insert(*ParseEntry("1520:peseng:inter."))
	dict.Insert(*ParseEntry("2512:txe'lan:n."))
	dict.Insert(*ParseEntry("13238:wrrz<0><1>är<2>ìp:vtr."))
	dict.Insert(*ParseEntry("13239:txe'lanti wrrzärìp:ph."))
	dict.Insert(*ParseEntry("10368:tìtseri:n."))
	dict.Insert(*ParseEntry("11608:to tìtseri:ph."))

	return dict
}

func BenchmarkMiniDict(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dict := miniDict()
		if dict == nil {
			b.Fail()
		}
	}
}

func TestDictionary_Lookup(t *testing.T) {
	dict := miniDict()
	t.Log("Graph Size:", dict.Root.Size())

	table := []struct {
		Lookup   string
		Expected string
	}{
		{"kaltxì", "692"},
		{"kaldì", "692"},
		{"ma", "1056"},
		{"kifkey", "800"},
		{"uvan", "2644"},
		{"uran", "7772"},
		{"tìfmetok", "2140"},
		{"tìfmusetok", "392:n. tì- <us>"},
		{"aysìfmetok", "2140 ay- t→s"},
		{"täpeykìyeverkeiup", "2084 <äpeyk,ìyev,ei>"},
		{"fepesìfmusetoktsyìpoka", "392:n. pe-pxe-tì- <us> -tsyìp-o-ka p→f,px→p,t→s"},
		{"eltuti", "264 -ti"},
		{"uvanterisì letokx", "2644 -teri-sì + letokx;9480 -teri-sì"},
		{"fe'ranìl", "9248 -ìl"},
		{"fe'erul", "9680 <er>"},
		{"'asap si", "12962"},
		{"tsayuvane", "2644 tsa-ay- -ä;2644 tsa-ay- -ne"},
		{"ayerik", "2744 ay-"},
		{"apolpxay", "5312 a-"},
		{"pesul", "1524 -l"},
		{"apeseng", "1520 a-"},
		{"pesengìl", "1520 -ìl"},
		{"polpxayìl", ""},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s", row.Lookup), func(t *testing.T) {
			dict := miniDict()
			runner := dict.Runner()
			res := runner.Run(row.Lookup)

			assert.Equal(t, dict.Lookup(row.Lookup), dict.Runner().Run(row.Lookup))

			dict.Optimize()
			res2 := runner.Run(row.Lookup)
			assert.Equal(t, res, res2)

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
		})
	}
}

func TestDictionary_Extract(t *testing.T) {
	dict := miniDict()
	t.Log("Graph Size:", dict.Root.Size())

	table := []struct {
		Lookup   string
		Expected string
	}{
		{"kaltxì, ma kifkey!", "[1] 692;[2] 1056;[3] 800"},
		{"eltu herahaw", "[1] 264;[2] 544 <er>"},
		{"eltut heykahaw", "[1] 13458"},
		{"fraeltut ke heykahängaw ukìl", "[1] 13458 fra- <äng> [ke];[3] 6680 -ìl"},
		{"eltuot heykahaw", "[1] 13458 -o"},
		{"eltur tìtxen soli", "[1] 6520 <ol>"},
		{"eltur tìtxen rää sivi", "[1] 6520 <iv> [rä'ä]"},
		{"fmäpetok glurb", "[1] 392 <äp>"},
		{"te'lanti wrrzärìp", "[1] 13239 tx→t"},
		{"to tìtseri", "[1] 11608"},
		{"fmäpetok to tìtseri", "[1] 392 <äp>;[2] 11608"},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s", row.Lookup), func(t *testing.T) {
			runner := dict.Runner()
			res := runner.Extract(row.Lookup)

			assert.Equal(t, dict.Extract(row.Lookup), dict.Runner().Extract(row.Lookup))

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
		})
	}
}

func BenchmarkDictionary_Example(b *testing.B) {
	dict := miniDict()

	for _, word := range []string{"uvan", "tìfmetok", "täpeykìyeverkeiup", "fepesìfmusetoktsyìpoka"} {
		runner := dict.Runner()
		b.Run(word, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				res := runner.Run(word)
				if len(res) == 0 {
					b.Fail()
				}
			}
		})
	}
}
