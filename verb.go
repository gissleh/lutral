package lutral

import (
	"strings"
)

func VerbFromEntry(entry Entry) *Node {
	word := strings.ToLower(entry.Word)
	res := EmptyTree()

	defaultResult := entry.ID
	if len(entry.PoS) > 1 {
		for _, pos := range entry.PoS {
			if strings.HasPrefix(pos, "v") {
				defaultResult = entry.ID + ":" + pos
			}
		}
	}

	hadSiPart := false
	for _, siPart := range []string{" si", " säpi", " seyki", " säpeyki"} {
		if nounPart := strings.TrimSuffix(word, siPart); nounPart != word {
			staticInfix0 := siPart[len(" s") : len(siPart)-len("i")]

			res.MergeFrom(*GenerateSiVerb(nounPart, staticInfix0).AndThenResult(defaultResult))
			res.MergeFrom(*GenerateSiVerbParticiple(nounPart, staticInfix0).AndThenResult(entry.ID + ":adj."))
			res.MergeFrom(*GenerateSiVerbAgent(nounPart, staticInfix0).AndThenResult(entry.ID + ":n."))

			if staticInfix0 == "" {
				res.MergeFrom(*GenerateSiVerbTswo(nounPart).AndThenResult(entry.ID + ":n."))
			}

			hadSiPart = true
			break
		}
	}

	if !hadSiPart {
		res.MergeFrom(*GenerateVerb(word, *entry.InfixPositions).AndThenResult(defaultResult))
		res.MergeFrom(*GenerateNegatedVerb(word, *entry.InfixPositions).AndThenResult(defaultResult))
		res.MergeFrom(*GenerateVerbParticiple(word, *entry.InfixPositions).AndThenResult(entry.ID + ":adj."))
		res.MergeFrom(*GenerateVerbTsuk(word, *entry.InfixPositions).AndThenResult(entry.ID + ":adj."))
		res.MergeFrom(*GenerateVerbGerund(word, *entry.InfixPositions).AndThenResult(entry.ID + ":n."))
		res.MergeFrom(*GenerateVerbTswo(word, *entry.InfixPositions).AndThenResult(entry.ID + ":n."))
		res.MergeFrom(*GenerateVerbAgent(word, *entry.InfixPositions).AndThenResult(entry.ID + ":n."))
	}

	return res
}

func GenerateNegatedVerb(word string, infixes [2]int) *Node {
	return CombineTrees(
		BuildTree("[rä'ä]|[rää=rä'ä]|[ke]", " ").AndThen(*GenerateVerb(word, infixes)),
		GenerateVerb(word, infixes).AndThen(*BuildTree(" ", "[rä'ä]|[rää=rä'ä]")),
	)
}

func GenerateVerb(word string, infixes [2]int) *Node {
	split3 := splitAtInfixes(word, infixes)

	res := EmptyTree()

	// The base verb + special case: set-in-stone <0>
	if strings.HasSuffix(split3[0], "eyk") || strings.HasSuffix(split3[0], "äp") {
		res.MergeFrom(*BuildTree(split3[0], "<1>", split3[1], "<2>", split3[2]))
	} else {
		res.MergeFrom(*BuildTree(split3[0], "<0>", "<1>", split3[1], "<2>", split3[2]))
	}

	if infixes[0] != infixes[1] {
		// frrfen -> *ferfen (should only be allowed in unstressed syllables, but that info is not tracked)
		if strings.HasPrefix(split3[1], "rr") {
			res.MergeFrom(*BuildTree(split3[0], "<0>", "<er>", strings.TrimPrefix(split3[1], "rr"), "<2>", split3[2]))
		}

		// plltxe -> poltxe (should only be allowed in unstressed syllables, but that info is not tracked)
		if strings.HasPrefix(split3[1], "ll") {
			res.MergeFrom(*BuildTree(split3[0], "<0>", "<ol>", strings.TrimPrefix(split3[1], "ll"), "<2>", split3[2]))
		}

		// The ceremonial puke
		if strings.HasSuffix(split3[1], "u") {
			res.MergeFrom(*BuildTree(split3[0], "<0>", "<1>", split3[1], "<y =uy>", split3[2]))
		}
	} else {
		// The ceremonial puke
		if strings.HasSuffix(split3[0], "u") {
			res.MergeFrom(*BuildTree(split3[0], "<y =uy>", split3[2]))
		}

		// nrr -> ner
		if strings.HasPrefix(split3[2], "rr") {
			res.MergeFrom(*BuildTree(split3[0], "<0>", "<er>", strings.TrimPrefix(split3[2], "rr")))
		}

		// vll -> vol
		if strings.HasPrefix(split3[2], "ll") {
			res.MergeFrom(*BuildTree(split3[0], "<0>", "<ol>", strings.TrimPrefix(split3[2], "ll")))
		}
	}

	return res
}

func GenerateVerbParticiple(word string, infixes [2]int) *Node {
	split2 := splitAtInfix(word, infixes[0])

	if strings.HasSuffix(split2[0], "eyk") {
		return GenerateAdjective(*BuildTree(split2[0], "<us,awn>", split2[1]))
	} else if strings.HasSuffix(split2[0], "äp") {
		return GenerateAdjective(*BuildTree(split2[0], "<us>", split2[1]))
	} else {
		return GenerateAdjective(*BuildTree(split2[0], "<eyk,>", "<us,awn>", split2[1])).
			MergedWith(*GenerateAdjective(*BuildTree(split2[0], "<äp>", "<us>", split2[1])))
	}
}

func GenerateVerbGerund(word string, infixes [2]int) *Node {
	split2 := splitAtInfix(word, infixes[0])

	if strings.HasSuffix(split2[0], "eyk") || strings.HasSuffix(split2[0], "äp") {
		return GenerateNoun(*BuildTree("tì-", split2[0], "<us>", split2[1]))
	} else {
		return GenerateNoun(*BuildTree("tì-", split2[0], "<0>", "<us>", split2[1]))
	}
}

func GenerateVerbTswo(word string, infixes [2]int) *Node {
	split2 := splitAtInfix(word, infixes[0])

	if strings.HasSuffix(split2[0], "eyk") || strings.HasSuffix(split2[1], "äp") {
		return GenerateNoun(*BuildTree(word, "-tswo"))
	} else {
		return GenerateNoun(*BuildTree(split2[0], "<0>", split2[1], "-tswo"))
	}
}

func GenerateVerbTsuk(word string, infixes [2]int) *Node {
	split2 := splitAtInfix(word, infixes[0])

	if strings.HasSuffix(split2[0], "eyk") || strings.HasSuffix(split2[1], "äp") {
		return GenerateAdjective(*BuildTree("tsuk-|ketsuk-", split2[0], split2[1]))
	} else {
		return GenerateAdjective(*BuildTree("tsuk-|ketsuk-", split2[0], "<0>", split2[1]))
	}
}

func GenerateVerbAgent(word string, infixes [2]int) *Node {
	split2 := splitAtInfix(word, infixes[0])

	if strings.HasSuffix(split2[0], "eyk") || strings.HasSuffix(split2[1], "äp") {
		return GenerateNoun(*BuildTree(word, "-yu"))
	} else {
		return GenerateNoun(*BuildTree(split2[0], "<0>", split2[1], "-yu"))
	}
}

func GenerateSiVerb(nounPart string, staticInfix0 string) *Node {
	if staticInfix0 != "" {
		return CombineTrees(
			BuildTree("$np", nounPart, "$nsmod_si", " s", staticInfix0, "<1>", "<2>", "i"),
			BuildTree("$np", nounPart, "$nsmod_si", " ", "[rä'ä]|[rää=rä'ä]|[ke]", " s", staticInfix0, "<1>", "<2>", "i"),
		)
	} else {
		return CombineTrees(
			BuildTree("$np", nounPart, "$nsmod_si", " s", "<eyk,äpeyk,>", "<1>", "<2>", "i"),
			BuildTree("$np", nounPart, "$nsmod_si", " ", "[rä'ä]|[rää=rä'ä]|[ke]", " s", "<eyk,äpeyk,>", "<1>", "<2>", "i"),
		)
	}
}

func GenerateSiVerbParticiple(nounPart string, staticInfix0 string) *Node {
	if staticInfix0 == "eyk" {
		return GenerateAdjective(*BuildTree(nounPart, "\\-seyk", "<us,awn>", "i"))
	} else if staticInfix0 == "äp" {
		return GenerateAdjective(*BuildTree(nounPart, "\\-säp", "<us>", "i"))
	} else {
		return GenerateAdjective(*BuildTree(nounPart, "\\-s", "<eyk>", "<us,awn>", "i")).
			MergedWith(*GenerateAdjective(*BuildTree(nounPart, "\\-s", "<us>", "i")))
	}
}

func GenerateSiVerbTswo(nounPart string) *Node {
	return GenerateNoun(*BuildTree(nounPart, "-tswo"))
}

func GenerateSiVerbAgent(nounPart string, staticInfix0 string) *Node {
	if staticInfix0 != "" {
		return GenerateNoun(*BuildTree(nounPart, "s", staticInfix0, "i", "-yu"))
	} else {
		return GenerateNoun(*BuildTree(nounPart, "si", "-yu"))
	}
}

func splitAtInfixes(s string, infixPositions [2]int) [3]string {
	return [3]string{s[:infixPositions[0]], s[infixPositions[0]:infixPositions[1]], s[infixPositions[1]:]}
}

func splitAtInfix(s string, infixPosition int) [2]string {
	return [2]string{s[:infixPosition], s[infixPosition:]}
}
