package lutral

import "strings"

func WithAlternativeSpellings(str string) []string {
	return append([]string{str}, AlternativeSpellings(str)...)
}

func AlternativeSpellings(str string) []string {
	var res []string
	if strings.HasPrefix(str, "sä") && !strings.HasPrefix(str, "säts") {
		afterPrefix := strings.TrimPrefix(str, "sä")
		for _, clusterable := range clusterables {
			if strings.HasPrefix(afterPrefix, clusterable) {
				res = append(res, "s"+afterPrefix)
				break
			}
		}
	}

	if strings.HasPrefix(str, "tìs") {
		afterPrefix := strings.TrimPrefix(str, "tìs")
		for _, clusterable := range clusterables {
			if strings.HasPrefix(afterPrefix, clusterable) {
				res = append(res, "ts"+afterPrefix)
				break
			}
		}
	}

	beforeReef := len(res)
	if reefAlt := generateReefSpelling(str); reefAlt != str {
		res = append(res, reefAlt)
	}

	for _, alt := range res[:beforeReef] {
		reefAlt := generateReefSpelling(alt)
		if reefAlt != alt {
			res = append(res, reefAlt)
		}
	}

	if strings.HasPrefix(str, "fut") && !strings.HasPrefix(str, "futs") {
		res = append(res, "ft"+str[len("fut"):])
	}

	return res
}

func generateReefSpelling(s string) string {
	res := reefReplacer.Replace(s)
	if strings.ContainsRune(res, '\'') {
		sb := strings.Builder{}
		prev := rune(0)
		prev2 := rune(0)
		for _, ch := range res {
			if prev == '\'' && strings.ContainsRune(vowels, prev2) && strings.ContainsRune(vowels, ch) {
				sb.WriteRune(prev2)
				sb.WriteRune(ch)

				prev2 = 0
				prev = 0
			} else {
				if prev2 != 0 {
					sb.WriteRune(prev2)
				}

				prev2 = prev
				prev = ch
			}
		}

		if prev2 != 0 {
			sb.WriteRune(prev2)
		}
		if prev != 0 {
			sb.WriteRune(prev)
		}

		res = sb.String()
	}

	return res
}

var reefReplacer = strings.NewReplacer(
	"tsy", "ch",
	"sy", "sh",
	"nkx", "n-g",
	"tskx", "tskx",
	"skx", "skx",
	"fkx", "fkx",
	"kx", "g",
	"tskx", "tspx",
	"skx", "spx",
	"fkx", "fpx",
	"tx", "d",
	"tskx", "tspx",
	"skx", "spx",
	"fkx", "fpx",
	"px", "b",
)

var clusterables = []string{
	"px", "tx", "kx", "ng", "p", "t", "k", "m", "n", "r", "l", "w", "y",
}

var vowels = "aäeéiìouù<"
