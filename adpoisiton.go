package lutral

import "strings"

func AdpositionFromEntry(entry Entry) (res *Node, suffix *Node) {
	word := strings.TrimSuffix(strings.ToLower(entry.Word), "+")

	resultValue := entry.ID
	if len(entry.PoS) > 1 {
		resultValue = entry.ID + ":adp."
	}

	return BuildTree(word).AndThenResult(resultValue),
		BuildTree("-"+word, "/return")
}
