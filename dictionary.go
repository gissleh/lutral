package lutral

import (
	"strings"
)

// Dictionary collects all the functionality in one place where you can add Entry objects rather than
// building from them.
type Dictionary struct {
	Root       Node                `json:"root"`
	IsSorted   bool                `json:"isSorted"`
	SubTreeMap map[string]*Node    `json:"subtreeMap"`
	Phrases    map[string][]Result `json:"phrases"`
}

func (dictionary *Dictionary) Runner() *Runner {
	return &Runner{Root: &dictionary.Root, SubtreeMap: dictionary.SubTreeMap, PhraseMap: dictionary.Phrases, res: make([]Result, 0, 8), isSorted: dictionary.IsSorted}
}

func (dictionary *Dictionary) Lookup(word string) []Result {
	return dictionary.Runner().Run(word)
}

func (dictionary *Dictionary) Extract(words string) []Result {
	return dictionary.Runner().Extract(words)
}

func (dictionary *Dictionary) Optimize() {
	dictionary.Root.Compact()
	dictionary.Root.SortChildren()
	dictionary.IsSorted = true
}

func (dictionary *Dictionary) Insert(entry Entry) {
	root := &dictionary.Root
	dictionary.IsSorted = false

	if dictionary.SubTreeMap == nil {
		dictionary.SubTreeMap = GenerateInitialSubTreeMap()
	}
	if dictionary.Phrases == nil {
		dictionary.Phrases = make(map[string][]Result)
	}

	for _, spelling := range WithAlternativeSpellings(strings.ToLower(entry.WordWithInfixBrackets())) {
		entry := entry
		entry.SetWordAndInfixes(spelling)

		uninflectables := EmptyTree()
		uninflectableCount := 0

		for _, pos := range entry.PoS {
			switch pos {
			case "adj.", "num.":
				root.MergeFrom(*AdjectiveFromEntry(entry))
			case "n.", "prop.n.":
				root.MergeFrom(*NounFromEntry(entry))
			case "pn.":
				root.MergeFrom(*PronounFromEntry(entry))
			case "vin.", "vim.", "vtr.", "vtrm.":
				root.MergeFrom(*VerbFromEntry(entry))
			case "inter.":
				hasFlag := false
				if entry.HasFlag("inter:adj.") {
					root.MergeFrom(*AdjectiveFromEntry(entry))
					hasFlag = true
				}
				if entry.HasFlag("inter:n.") {
					root.MergeFrom(*NounFromEntry(entry))
					hasFlag = true
				}
				if entry.HasFlag("inter:adv.") && !entry.HasFlag("inter:n.") {
					root.MergeFrom(*UninflectableWordFromEntry(entry, ""))
					hasFlag = true
				}

				if !hasFlag {
					root.MergeFrom(*AffixedOnlyAdjectiveFromEntry(entry))
					root.MergeFrom(*NounFromEntry(entry))
				}
			case "ph.":
				runner := dictionary.Runner()
				runner.PhraseMap = nil

				entries := runner.ExtractWithoutSkipping(entry.Word)
				if entries != nil {
					dictionary.Phrases[entry.ID] = simplestResultSet(entries)
				} else if entry.InfixPositions != nil {
					root.MergeFrom(*VerbFromEntry(entry))
				} else {
					uninflectables.MergeFrom(*UninflectableWordFromEntry(entry, pos))
					uninflectableCount++
				}
			case "adp.":
				adposition, suffix := AdpositionFromEntry(entry)
				root.MergeFrom(*adposition)
				dictionary.SubTreeMap["nsadp"].MergeFrom(*suffix)
			default:
				uninflectables.MergeFrom(*UninflectableWordFromEntry(entry, pos))
				uninflectableCount++
			}
		}

		if uninflectableCount != 0 {
			if uninflectableCount != len(entry.PoS) {
				root.MergeFrom(*uninflectables)
			} else {
				root.MergeFrom(*UninflectableWordFromEntry(entry, ""))
			}
		}
	}
}

// UninflectableWordFromEntry generates a plain word. It will use the `pos` argument if there are multiple
// for the entry. While it says uninflectable, it will still support lenition as any initial raw-node.
func UninflectableWordFromEntry(entry Entry, pos string) *Node {
	resultValue := entry.ID
	if len(entry.PoS) > 1 && pos != "" {
		resultValue = entry.ID + ":" + pos
	}

	return CombineTrees(
		BuildTree(strings.ToLower(entry.Word)).AndThenResult(resultValue),
		BuildTree(strings.ToLower(entry.Word), "-sì").AndThenResult(resultValue),
	)
}
