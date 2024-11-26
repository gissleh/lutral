package lutral

import (
	"strings"
)

func NounFromEntry(entry Entry) *Node {
	resultValue := entry.ID
	if len(entry.PoS) > 1 {
		for _, pos := range entry.PoS {
			if pos == "n." || pos == "prop.n." {
				resultValue = entry.ID + ":" + pos
			}
		}
	}

	word := strings.ToLower(entry.Word)
	if entry.HasFlag("loanword") && strings.HasSuffix(word, "ì") {
		return CombineTrees(
			BuildTree("$np", word).AndThenResult(resultValue),
			GenerateNoun(ParseNode(strings.TrimSuffix(word, "ì"))).AndThenResult(resultValue),
		)
	} else if strings.Contains(word, " ") {
		// This is pretty much just drawing boundaries based on known words.
		// As of November 2024, they are (* marks where suffixes shall go)
		//  toruk makto*         tsko swizaw*
		//  eltu* lefngap        pängkxoyu* lekoren
		//  tìftia* kifkeyä      uvan* letokx
		//  rel* arusikx         tìftiatu* kifkeyä
		//  swoasey* ayll        yomyo* lerìk
		//  mo* a fngä'          pamrelvul* lerìn
		//  koren* ayll          tslikxyu* latopin
		//  tslikxyu* tsawlak    renu* ngampamä
		//  txawnulsrung* a yur  txawnulsrung* a tswayon
		//  trrpxì* Sawtuteyä    mo* letrrtrr
		//  mo* a yom            mo* a hahaw

		split := strings.Split(word, " ")
		if len(split) > 2 { // txawnulsrung a tswayon, mo a hahaw, etc...
			preWord := split[0]
			postWord := word[len(split[0]):]

			return GenerateNoun(ParseNode(preWord)).AndThen(ParseNode(postWord)).AndThenResult(resultValue)
		} else {
			leftWord := split[0]
			rightWord := split[1]

			if strings.HasSuffix(leftWord, "yu") ||
				strings.HasSuffix(leftWord, "tu") ||
				strings.HasPrefix(rightWord, "a") ||
				strings.HasPrefix(rightWord, "le") ||
				strings.HasSuffix(rightWord, "yä") ||
				strings.HasSuffix(rightWord, "ä") {
				return GenerateNoun(ParseNode(leftWord)).
					AndThen(ParseNode(" ")).
					AndThen(ParseNode(rightWord)).
					AndThenResult(resultValue)
			} else {
				return BuildTree(leftWord).
					AndThen(ParseNode(" ")).
					AndThen(*GenerateNoun(ParseNode(rightWord))).
					AndThenResult(resultValue)
			}
		}
	} else {
		if entry.HasPoS("inter.") {
			return GenerateUnPrefixedNoun(ParseNode(word)).AndThenResult(resultValue)
		} else {
			return GenerateNoun(ParseNode(word)).AndThenResult(resultValue)
		}
	}
}

func GenerateNoun(core Node) *Node {
	return generateNoun(*BuildTree("$np").AndThen(core))
}

func GenerateUnPrefixedNoun(core Node) *Node {
	return generateNoun(*EmptyTree().AndThen(core))
}

func generateNoun(core Node) *Node {
	res := &core

	res.SearchReplace(func(node *Node) *Node {
		if len(node.Children) == 0 {
			if node.Kind == NKRaw {
				return generateNounTail(node.Value).MergedWith(*BuildTree(node.Value, "$nsmod"))
			} else if node.Kind == NKSuffix && node.Value != "y" {
				return generateNounTail("-" + node.Value).MergedWith(*BuildTree("-"+node.Value, "$nsmod"))
			}
		}

		return nil
	})

	return res
}

func generateNounTail(tail string) *Node {
	res := EmptyTree()

	res.MergeFrom(*BuildTree(tail, findNounSuffix(tail)))

	if !strings.HasPrefix(tail, "-") {
		// Edge case: soaia -> soaiä
		if strings.HasSuffix(tail, "ia") {
			res.MergeFrom(*BuildTree(strings.TrimSuffix(tail, "a"), "-ä"))
		}

		// Edge case: omatikaya -> omatikayaä
		if strings.HasSuffix(tail, "aya") {
			res.MergeFrom(*BuildTree(tail, "-ä"))
		}
	}

	return res
}

func findNounSuffix(tail string) string {
	suffixNode := "$ncec"
	for _, kv := range nounSuffixMappings {
		if strings.HasSuffix(tail, kv[0]) {
			suffixNode = kv[1]
			break
		}
	}

	return suffixNode
}

var nounSuffixMappings = [][2]string{
	{"ey", "$ncedy"},
	{"ay", "$ncedy"},
	{"aw", "$ncedw"},
	{"ew", "$ncedw"},
	{"o", "$ncevou"},
	{"u", "$ncevou"},
	{"ù", "$ncevou"},
	{"ia", "$nceia"},
	{"a", "$ncev|-y"},
	{"ä", "$ncev"},
	{"e", "$ncev|-y"},
	{"é", "$ncev|-y"},
	{"i", "$ncev"},
	{"ì", "$ncev"},
	{"'", "$ncec'"},
	{"t", "$ncect"},
	{"s", "$ncevìlw"},
	{"ts", "$ncevìlw"},
	{"ln", "$ncevìlw"},
	{"rk", "$ncevìlw"},

	// Fallback to $ncec
}
