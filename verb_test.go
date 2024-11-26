package lutral

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerbFromEntry(t *testing.T) {
	table := []struct {
		Entry    string
		Test     string
		Expected string
	}{
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "fmetok",
			Expected: "392",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "fmäpeykereteiok",
			Expected: "392 <äpeyk,er,ei>",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "fmusetoka",
			Expected: "392:adj. <us> -a",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "tsukfmetok",
			Expected: "392:adj. tsuk-",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "tsukfmäpetoka",
			Expected: "392:adj. tsuk- <äp> -a",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "aketsukfmetok",
			Expected: "392:adj. a-ketsuk-",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "fmetoktswoori",
			Expected: "392:n. -tswo-o-ri",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "fmetoktswo-ori",
			Expected: "392:n. -tswo +-ori;392:n. -tswo-o-ri",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "fmetokyu",
			Expected: "392:n. -yu",
		},
		{
			Entry: "44:'<0><1>ek<2>o:n.", Test: "ekoyu",
			Expected: "44:n. -yu 'e→e",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "ayfmetokyu",
			Expected: "392:n. ay- -yu",
		},
		{
			Entry: "6428:zeyk<0><1><2>o:vtr.", Test: "zeykero",
			Expected: "6428 <er>",
		},
		{
			Entry: "6428:zeyk<0><1><2>o:vtr.", Test: "zeykoyu",
			Expected: "6428:n. -yu",
		},
		{
			Entry: "6428:zeyk<0><1><2>o:vtr.", Test: "zeykeyko",
			Expected: "",
		},
		{
			Entry: "6428:zeyk<0><1><2>o:vtr.", Test: "atsukzeyko",
			Expected: "6428:adj. a-tsuk-",
		},
		{
			Entry: "9632:ng<0><1>ä'<2>än:vin.", Test: "tìngusä'änìri",
			Expected: "9632:n. tì- <us> -ìri",
		},
		{
			Entry: "3980:srung s<0><1><2>i:vin.", Test: "asrung-susi",
			Expected: "3980:adj. a- <us>",
		},
		{
			Entry: "10232:sngum s<0><1><2>i:vin.", Test: "sngum seyki",
			Expected: "10232 <eyk>",
		},
		{
			Entry: "2648:uvan s<0><1><2>i:vin.", Test: "fneuvantswo",
			Expected: "2648:n. fne- -tswo",
		},
		{
			Entry: "13167:piak säp<0><1><2>i:vin.", Test: "piak-säpusi",
			Expected: "13167:adj. <us>",
		},
		{
			Entry: "13167:piak säp<0><1><2>i:vin.", Test: "piak-säpäpusi",
			Expected: "",
		},
		{
			Entry: "392:fm<0><1>et<2>ok:vtr.", Test: "fmäpeykìlmetängok",
			Expected: "392 <äpeyk,ìlm,äng>",
		},
		{
			Entry: "2084:t<0><1>erk<2>up:vin.", Test: "täpeykìyeverkeiup",
			Expected: "2084 <äpeyk,ìyev,ei>",
		},
		{
			Entry: "6428:zeyk<0><1><2>o:vtr.", Test: "zeykeyko",
			Expected: "",
		},
		{
			Entry: "6420:z<0><1><2>o:vin.", Test: "zeykeyko",
			Expected: "",
		},
		{
			Entry: "9568:späp<0><1><2>eng:vin.", Test: "späpäpeng",
			Expected: "",
		},
		{
			Entry: "9568:späp<0><1><2>eng:vin.", Test: "späpeykeng",
			Expected: "",
		},
		{
			Entry: "9568:späp<0><1><2>eng:vin.", Test: "späpereng",
			Expected: "9568 <er>",
		},
		{
			Entry: "-1:l<0><1><2>ok:vtr.,adp.", Test: "lerok", // They're separate entries in the fwew data.
			Expected: "-1:vtr. <er>",
		},
		{
			Entry: "3892:tìsraw seyk<0><1><2>i:vtr.", Test: "tìsraw seykeri",
			Expected: "3892 <er>",
		},
		{
			Entry: "3892:tìsraw seyk<0><1><2>i:vtr.", Test: "tìsraw seykeykeri",
			Expected: "",
		},
		{
			Entry: "10964:<0><1>u<2>e':vtr.", Test: "eruye'",
			Expected: "10964 <er,uy>",
		},
		{
			Entry: "-10964:u<0><1><2>e':vtr.", Test: "uye'", // Fake entries for coverage's sake.
			Expected: "-10964 <uy>",
		},
		{
			Entry: "-10964:u<0><1><2>e':vtr.", Test: "uerye'", // Fake entries for coverage's sake.
			Expected: "",
		},
		{
			Entry: "464:f<0><1>rrf<2>en:vtr.", Test: "ferfen", // Incorrect, syllable is stressed, but I cannot find better example
			Expected: "464 <er>",
		},
		{
			Entry: "1544:p<0><1>lltx<2>e:vtr.", Test: "poltxe",
			Expected: "1544 <ol>",
		},
		{
			Entry: "8852:v<0><1><2>ll:vtr.", Test: "vol",
			Expected: "8852 <ol>",
		},
		{
			Entry: "9116:n<0><1><2>rr:vin.", Test: "ner",
			Expected: "9116 <er>",
		},
	}

	for _, row := range table {
		t.Run(fmt.Sprintf("%s %s", row.Entry, row.Test), func(t *testing.T) {
			entry := ParseEntry(row.Entry)
			if !assert.NotNil(t, entry) {
				return
			}

			tree := VerbFromEntry(*entry)
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
