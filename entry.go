package lutral

import "strings"

// Entry is the minimum information needed to build a tree for the word.
type Entry struct {
	ID             string
	Word           string
	PoS            []string
	InfixPositions *[2]int
	// Supported Flags: "loanword"
	Flags []string
}

func (e *Entry) WordWithInfixBrackets() string {
	if e.InfixPositions == nil {
		return e.Word
	}

	res := e.Word[:e.InfixPositions[1]] + "<2>" + e.Word[e.InfixPositions[1]:]
	res = res[:e.InfixPositions[0]] + "<0><1>" + res[e.InfixPositions[0]:]
	return res
}

func (e *Entry) SetWordAndInfixes(wordWithBrackets string) {
	e.Word = wordWithBrackets

	if infix0Pos := strings.Index(wordWithBrackets, "<0><1>"); infix0Pos >= 0 {
		e.InfixPositions = &[2]int{
			infix0Pos,
			strings.Index(e.Word, "<2>") - len("<0><1>"),
		}

		e.Word = infixBracketReplacer.Replace(e.Word)
	}
}

func (e *Entry) HasFlag(pred string) bool {
	for _, flag := range e.Flags {
		if pred == flag {
			return true
		}
	}

	return false
}

func (e *Entry) HasPoS(pred string) bool {
	for _, pos := range e.PoS {
		if pred == pos {
			return true
		}
	}

	return false
}

func ParseEntry(s string) *Entry {
	res := &Entry{}
	split := strings.SplitN(s, ":", 4)
	if len(split) < 3 || split[2] == "" || split[1] == "" {
		return nil
	}

	res.ID = split[0]
	res.SetWordAndInfixes(split[1])
	res.PoS = strings.Split(split[2], ",")
	for i := range res.PoS {
		res.PoS[i] = strings.TrimSpace(res.PoS[i])
	}

	if len(split) >= 4 && split[3] != "" {
		res.Flags = strings.Split(split[3], ",")
	}

	return res
}

var infixBracketReplacer = strings.NewReplacer("<0>", "", "<1>", "", "<2>", "")
