package lutral

import "strings"

func PronounFromEntry(entry Entry) *Node {
	return GeneratePronoun(strings.ToLower(entry.Word)).AndThenResult(entry.ID)
}

func GeneratePronoun(word string) *Node {
	switch {
	case strings.HasSuffix(word, "yä"),
		strings.HasSuffix(word, "ri"),
		strings.HasSuffix(word, "t"),
		strings.HasSuffix(word, "r"):
		return BuildTree("$np", word)
	case strings.HasSuffix(word, "saw"):
		base := strings.TrimSuffix(word, "w")
		genBase := strings.TrimSuffix(word, "aw")
		return CombineTrees(
			BuildTree("$np", word, "$nsmod|$ncedw"),
			BuildTree("$np", base, "$nsmod|$ncev"),
			BuildTree("$np", genBase, "e", "-yä"),
			BuildTree("nì-", word),
		)
	case strings.HasSuffix(word, "o"):
		base := strings.TrimSuffix(word, "o")
		return CombineTrees(
			BuildTree("$np", word, "$nsmod|$pce_o"),
			BuildTree("$np", word, "/hook"),
			BuildTree("$np", base, "e", "-yä"),
			BuildTree("nì-", word),
		)
	case strings.HasSuffix(word, "ng"):
		return CombineTrees(
			BuildTree("$np", word, "/hook"),
			BuildTree("$np", word, "$nsmod"),
			BuildTree("$np", word, "a", "$nsmod"),
			BuildTree("$np", word, "a", "$pce_ng_a"),
			BuildTree("$np", word, "e", "-yä"),
			BuildTree("nì-", word),
		)
	case strings.HasSuffix(word, "nga"):
		base := strings.TrimSuffix(word, "a")
		return CombineTrees(
			BuildTree("$np", word, "/hook"),
			BuildTree("$np", word, "$nsmod|$pce_ng_a"),
			BuildTree("$np", base, "e", "-yä|-y"),
			BuildTree("nì-", word),
		)
	default:
		return GenerateNoun(ParseNode(word))
	}
}
