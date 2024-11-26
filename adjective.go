package lutral

import "strings"

func AdjectiveFromEntry(entry Entry) *Node {
	core := *BuildTree(strings.ToLower(entry.Word))

	return CombineTrees(
		GenerateAdjective(core).AndThenResult(entry.ID),
		GenerateAdjectiveAdverb(core).AndThenResult(entry.ID),
	)
}

func AffixedOnlyAdjectiveFromEntry(entry Entry) *Node {
	core := *BuildTree(strings.ToLower(entry.Word))
	return GenerateAffixedOnlyAdjective(core).AndThenResult(entry.ID)
}

func GenerateAdjective(core Node) *Node {
	return CombineTrees(
		CopyTree(core).AndThen(ParseNode("/hook")),
		CopyTree(core).AndThen(ParseNode("-sì")),
		GenerateAffixedOnlyAdjective(core),
	)
}

func GenerateAffixedOnlyAdjective(core Node) *Node {
	return CombineTrees(
		BuildTree("a-").AndThen(core),
		CopyTree(core).AndThen(ParseNode("-a")),
	)
}

func GenerateAdjectiveAdverb(core Node) *Node {
	return BuildTree("nì-").AndThen(core)
}
