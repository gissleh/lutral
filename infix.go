package lutral

import (
	"sort"
	"strings"
)

type infix struct {
	Match      string
	Name       string
	NotBefore  []string
	OnlyBefore []string
}

func infixes(m map[string][]infix, input ...string) []infix {
	var res []infix

	if len(input) >= 1 {
		if existing, ok := m[input[0]]; ok {
			res = existing[:len(existing):len(existing)]
			input = input[1:]
		}
	}

	for _, s := range input {
		if s == "" {
			res = append(res, infix{})
			continue
		}

		tokens := strings.Split(strings.Trim(s, " "), " ")
		in := infix{}
		for i, token := range tokens {
			if i == 0 {
				in.Match = token
				in.Name = token
			}

			switch {
			case strings.HasPrefix(token, "="):
				in.Name = strings.TrimPrefix(token, "=")
			case strings.HasPrefix(token, "-"):
				in.NotBefore = append(in.NotBefore, strings.TrimPrefix(token, "-"))
			case strings.HasPrefix(token, ">"):
				in.OnlyBefore = append(in.OnlyBefore, strings.TrimPrefix(token, ">"))
			}
		}

		res = append(res, in)
	}

	return res
}

func sortedInfixes(infixes []infix) []infix {
	sort.Slice(infixes, func(i, j int) bool {
		return infixes[i].Match < infixes[j].Match
	})

	return infixes
}

var infixMap = map[string][]infix{
	"0": sortedInfixes(infixes(nil,
		"",
		"äpeyk", "epeyk =äpeyk",
		"äp", "ep =äp",
		"eyk",
	)),
	"1": sortedInfixes(infixes(nil,
		"",
		"iv", "irv", "ilv", "imv", "iyev", "ìyev",
		"am", "ìm", "ìy", "ay",
		"ìsy", "ìsh=ìsy", "asy", "ash=asy",
		"er -rr", "arm", "ìrm", "ìry", "ary",
		"ol -ll", "alm", "ìlm", "ìly", "aly",
	)),
	"2": sortedInfixes(infixes(nil,
		"",
		"eiy >i >ì >rr >ll", "ei",
		"äng", "eng =äng",
		"ats",
		"uy",
	)),
}
