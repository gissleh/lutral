package lutral

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Runner struct {
	Root       *Node
	SubtreeMap map[string]*Node
	PhraseMap  map[string][]Result

	StepCount    int64
	SubStepCount int64

	res      []Result
	isSorted bool
}

const (
	noLenition        = 0
	allowLenition     = 1
	mandatoryLenition = 2
)

func (runner *Runner) Run(text string) []Result {
	if runner.SubtreeMap == nil {
		runner.SubtreeMap = GenerateInitialSubTreeMap()
	}

	runner.res = runner.res[:0]
	runner.runStep(runner.Root, strings.ToLower(text), allowLenition, "", nil)

	return append(runner.res[:0:0], runner.res...)
}

// Extract is like run, but it works through the text from left to right, returning all the entries that
// are the longest at each step (i.e. a si-verb wins over its noun or adjective component).
func (runner *Runner) Extract(text string) []Result {
	return runner.extract(text, false)
}

func (runner *Runner) ExtractWithoutSkipping(text string) []Result {
	return runner.extract(text, true)
}

func (runner *Runner) extract(text string, doNotSkip bool) []Result {
	if runner.SubtreeMap == nil {
		runner.SubtreeMap = GenerateInitialSubTreeMap()
	}

	runner.res = runner.res[:0]
	position := 0

	text = strings.ToLower(text)
	for ; len(text) > 0; text = strings.Trim(text, punctuation) {
		// Record where we are and run.
		resOffset := len(runner.res)
		runner.runStep(runner.Root, text, allowLenition, "", nil)
		position += 1

		// Skip word if no results
		if len(runner.res) == resOffset {
			if doNotSkip {
				return nil
			}

			next := strings.IndexAny(text, punctuation)
			if next == -1 {
				text = text[len(text):]
			} else {
				text = text[next:]
			}

			continue
		}

		if len(runner.res) > resOffset+1 {
			shortestRemainder := runner.res[resOffset].Remainder
			for i := resOffset + 1; i < len(runner.res); i++ {
				if len(runner.res[i].Remainder) < len(shortestRemainder) {
					shortestRemainder = runner.res[i].Remainder
				}
			}

			next := resOffset
			for i := resOffset; i < len(runner.res); i++ {
				if len(runner.res[i].Remainder) <= len(shortestRemainder) {
					runner.res[next] = runner.res[i]
					next += 1
				}
			}
			runner.res = runner.res[:next]
			text = shortestRemainder
		} else {
			text = runner.res[resOffset].Remainder
		}

		for i := range runner.res[resOffset:] {
			runner.res[resOffset+i].Position = position
		}
	}

	for i := range runner.res {
		runner.res[i].Remainder = ""
	}

	if position > 1 {
		for phraseID, phrase := range runner.PhraseMap {
			if len(runner.res) < len(phrase) {
				continue
			}
			positionFoundMap := uint32(0)
			phraseResult := Result{ID: phraseID}

			basePosition := len(runner.res)
			for i := range runner.res {
				relativePosition := runner.res[i].Position - basePosition
				if relativePosition < 0 {
					relativePosition = 0
				}
				if relativePosition >= len(phrase) {
					break
				}

				if runner.res[i].CoveredBy(phrase[relativePosition]) {
					if basePosition > runner.res[i].Position {
						basePosition = runner.res[i].Position
						phraseResult.Position = basePosition
					}

					phraseResult.AddAffixesFrom(runner.res[i], phrase[relativePosition])
					positionFoundMap |= 1 << relativePosition
				}
			}

			if positionFoundMap == (1<<len(phrase))-1 {
				l := 0
				added := false
				redIndex := -1

				for _, res := range runner.res {
					if res.Position < basePosition || res.Position >= basePosition+len(phrase) {
						runner.res[l] = res
						l += 1
					} else if res.Position == basePosition && !added {
						runner.res[l] = res
						redIndex = l
						l += 1
						added = true
					}
				}
				runner.res[redIndex] = phraseResult
				runner.res = runner.res[:l]
				break
			}
		}
	}

	return runner.res
}

func (runner *Runner) runStep(node *Node, remainder string, lenitionState int, skippableLetter string, returnTo *Node) bool {
	var strSliceBuf [4]string
	var intSliceBuf [4]int
	var didProceed bool

	runner.StepCount += 1

	switch node.Kind {
	case NKRoot:
		for i := range node.Children {
			runner.runStep(&node.Children[i], remainder, lenitionState, skippableLetter, returnTo)
		}
		didProceed = true

	case NKResult:
		if remainder == "" || strings.IndexAny(remainder, punctuation) == 0 {
			runner.SubStepCount += 1
			split := strings.Split(node.Value, ":")
			res := Result{
				ID:        split[0],
				Remainder: remainder,
			}
			if len(split) > 1 {
				res.PoS = split[1]
			}

			runner.res = append(runner.res, res)
			didProceed = true
		}

	case NKRaw:
		_, lastLetterLen := utf8.DecodeLastRuneInString(node.Value)
		nextSkippable := node.Value[len(node.Value)-lastLetterLen:]
		if nextSkippable == "s" {
			nextSkippable = ""
		}

		hadLenition := false
		if lenitionState&allowLenition != 0 {
			if node.Value == "'" {
				hadLenition = true
				if !strings.HasPrefix(remainder, "'") {
					for i := range node.Children {
						resOffset := len(runner.res)
						runner.runStep(&node.Children[i], remainder, noLenition, skippableLetter, returnTo)

						deleteList := intSliceBuf[:0]
						for j, res := range runner.res[resOffset:] {
							if !strings.HasPrefix(remainder, "rr") && !strings.HasPrefix(remainder, "ll") {
								firstCh, _ := utf8.DecodeRuneInString(remainder)
								runner.res[j+resOffset].Lenitions = append(res.Lenitions, fmt.Sprintf("'%c→%c", firstCh, firstCh))
							} else {
								deleteList = append(deleteList, j+resOffset-len(deleteList))
							}
						}

						for _, deleteIndex := range deleteList {
							runner.res = append(runner.res[:deleteIndex], runner.res[deleteIndex+1:]...)
						}
					}

					didProceed = true
				}
			} else {
				lenition, afterLenition := ApplyLenition(node.Value)
				hadLenition = lenition != ""
				if lenition != "" && (lenitionState&allowLenition != 0) {
					matchTexts := append(strSliceBuf[:0], afterLenition)
					if skippableLetter != "" && strings.HasPrefix(afterLenition, skippableLetter) {
						matchTexts = append(matchTexts, strings.TrimPrefix(afterLenition, skippableLetter))
					}

					for _, matchText := range matchTexts {
						runner.SubStepCount += 1

						if trimmedRemainder := strings.TrimPrefix(remainder, matchText); trimmedRemainder != remainder || matchText == "" {
							resOffset := len(runner.res)
							for i, child := range node.Children {
								nextSkippable := nextSkippable
								if child.Kind == NKRaw {
									nextSkippable = ""
								}

								runner.runStep(&node.Children[i], trimmedRemainder, noLenition, nextSkippable, returnTo)
							}
							for i := range runner.res[resOffset:] {
								// This is always the last lenition.
								runner.res[i+resOffset].Lenitions = []string{lenition}
							}

							didProceed = true
						}
					}
				}
			}
		}

		if !hadLenition || lenitionState&mandatoryLenition == 0 {
			matchTexts := append(strSliceBuf[:0], node.Value)
			if skippableLetter != "" && strings.HasPrefix(node.Value, skippableLetter) {
				matchTexts = append(matchTexts, strings.TrimPrefix(node.Value, skippableLetter))
			}

			for _, matchText := range matchTexts {
				runner.SubStepCount += 1

				if trimmedRemainder := strings.TrimPrefix(remainder, matchText); trimmedRemainder != remainder || matchText == "" {
					for i, child := range node.Children {
						nextSkippable := nextSkippable
						if child.Kind == NKRaw {
							nextSkippable = ""
						}

						runner.runStep(&node.Children[i], trimmedRemainder, noLenition, nextSkippable, returnTo)
					}

					didProceed = true
				}
			}
		}

	case NKPrefix:
		prefix := strings.TrimSuffix(node.Value, "+")
		hasLenition := prefix != node.Value
		nextLenition := noLenition
		if hasLenition {
			nextLenition = mandatoryLenition | allowLenition
		}

		lenition, afterLenition := ApplyLenition(prefix)
		resOffset := len(runner.res)

		_, lastLetterLen := utf8.DecodeLastRuneInString(prefix)
		nextSkippable := prefix[len(prefix)-lastLetterLen:]

		if lenition != "" && (lenitionState&allowLenition != 0) {
			matchTexts := append(strSliceBuf[:0], afterLenition)
			if skippableLetter != "" && strings.HasPrefix(afterLenition, skippableLetter) {
				matchTexts = append(matchTexts, strings.TrimPrefix(afterLenition, skippableLetter))
			}

			for _, matchText := range matchTexts {
				runner.SubStepCount += 1

				if trimmedRemainder := strings.TrimPrefix(remainder, matchText); trimmedRemainder != remainder || matchText == "" {
					for i := range node.Children {
						runner.runStep(&node.Children[i], trimmedRemainder, nextLenition, nextSkippable, returnTo)
					}
					for i, res := range runner.res[resOffset:] {
						runner.res[i+resOffset].Lenitions = prependToSlice(res.Lenitions, lenition)
					}

					didProceed = true
				}
			}
		}

		if lenition == "" || lenitionState&mandatoryLenition == 0 {
			matchTexts := append(strSliceBuf[:0], prefix)
			if skippableLetter != "" && strings.HasPrefix(prefix, skippableLetter) {
				matchTexts = append(matchTexts, strings.TrimPrefix(prefix, skippableLetter))
			}

			for _, matchText := range matchTexts {
				runner.SubStepCount += 1

				if trimmedRemainder := strings.TrimPrefix(remainder, matchText); trimmedRemainder != remainder || matchText == "" {
					for i := range node.Children {
						runner.runStep(&node.Children[i], trimmedRemainder, nextLenition, nextSkippable, returnTo)
					}
				}

				didProceed = true
			}
		}

		for i, res := range runner.res[resOffset:] {
			runner.res[i+resOffset].Prefixes = prependToSlice(res.Prefixes, prefix)
		}

	case NKInfix:
		prevFit := false
		sorted := node.Value == "0" || node.Value == "1" || node.Value == "2"

	infixLoop:
		for _, infix := range infixes(infixMap, strings.Split(node.Value, ",")...) {
			runner.SubStepCount += 1

			if afterInfix := strings.TrimPrefix(remainder, infix.Match); afterInfix != remainder || infix.Match == "" {
				prevFit = infix.Match != ""

				for _, notBefore := range infix.NotBefore {
					if strings.HasPrefix(afterInfix, notBefore) {
						continue infixLoop
					}
				}

				if len(infix.OnlyBefore) > 0 {
					found := false
					for _, onlyBefore := range infix.OnlyBefore {
						runner.SubStepCount += 1
						if strings.HasPrefix(afterInfix, onlyBefore) {
							found = true
							break
						}
					}

					if !found {
						continue infixLoop
					}
				}

				resOffset := len(runner.res)

				for i := range node.Children {
					runner.runStep(&node.Children[i], afterInfix, noLenition, "", returnTo)
				}

				if infix.Name != "" {
					for i, res := range runner.res[resOffset:] {
						runner.res[i+resOffset].Infixes = prependToSlice(res.Infixes, infix.Name)
					}
				}

				didProceed = true
			} else if sorted && prevFit {
				break
			}
		}

	case NKSuffix:
		suffix, suffixName, hasAlias := strings.Cut(node.Value, "=")
		if !hasAlias {
			suffixName = suffix
		}

		remainder = strings.TrimPrefix(remainder, "-")
		runner.SubStepCount += 1

		_, lastLetterLen := utf8.DecodeLastRuneInString(suffix)
		nextSkippable := suffix[len(suffix)-lastLetterLen:]
		if lastLetterLen == len(suffix) {
			nextSkippable = ""
		}

		matchTexts := append(strSliceBuf[:0], suffix)
		if skippableLetter != "" && strings.HasPrefix(suffix, skippableLetter) {
			matchTexts = append(matchTexts, strings.TrimPrefix(suffix, skippableLetter))
		}

		for _, matchText := range matchTexts {
			if afterSuffix := strings.TrimPrefix(remainder, matchText); afterSuffix != remainder {
				resOffset := len(runner.res)
				for i := range node.Children {
					runner.runStep(&node.Children[i], afterSuffix, noLenition, nextSkippable, returnTo)
				}
				for i, res := range runner.res[resOffset:] {
					runner.res[i+resOffset].Suffixes = prependToSlice(res.Suffixes, suffixName)
				}

				didProceed = true
				break
			}
		}

	case NKSubTree:
		subTree := runner.SubtreeMap[node.Value]
		if subTree == nil {
			panic("unknown subtree " + node.Value)
		}

		nextReturnTo := returnTo
		if nextReturnTo == nil {
			nextReturnTo = node
		}

		didProceed = runner.runStep(subTree, remainder, lenitionState, skippableLetter, nextReturnTo)

	case NKReturn:
		if returnTo == nil {
			panic("nowhere to /return to")
		}

		for i := range returnTo.Children {
			childProceeded := runner.runStep(&returnTo.Children[i], remainder, lenitionState, skippableLetter, nil)
			if childProceeded {
				didProceed = true
			}
		}

	case NKParticle:
		particleMatch, particleName, hasOverride := strings.Cut(node.Value, "=")
		if !hasOverride {
			particleName = particleMatch
		}

		if afterParticle := strings.TrimPrefix(remainder, particleMatch); afterParticle != remainder {
			resOffset := len(runner.res)
			for i := range node.Children {
				runner.runStep(&node.Children[i], afterParticle, noLenition, "", returnTo)
			}
			for i, res := range runner.res[resOffset:] {
				runner.res[i+resOffset].Particles = prependToSlice(res.Particles, particleName)
			}

			didProceed = true
		}

	case NKLeafHook:
		// Do nothing, this one is just for helping tree generation.
	}

	return didProceed
}

func prependToSlice[T any, S ~[]T](slice S, value T) S {
	if slice == nil {
		return S{value}
	}

	slice = append(slice[:1], slice...)
	slice[0] = value
	return slice
}

func GenerateInitialSubTreeMap() map[string]*Node {
	return map[string]*Node{
		// Noun prefixes
		"np": CombineTrees(
			BuildTree("$np2"),
			BuildTree("me+|pxe+|ay+", "$np2"),
			BuildTree("fì-|tsa-|pe+", "$np2"),
			BuildTree("fì-|tsa-|pe+", "me+|pxe+|ay+", "$np2"),
			BuildTree("sna-|munsna-", "$np2"),
			BuildTree("fra-", "$np2"),
			BuildTree("fra-", "ay+", "$np2"),
			BuildTree("fay+|pay+", "$np2"),
		),
		// Noun prefixes for modifying the noun (called by np or np_numbers only)
		"np2": CombineTrees(
			BuildTree("/return"),
			BuildTree("fne-", "/return"),
		),
		// Noun suffixes that modify
		"nsmod": CombineTrees(
			BuildTree("$nsadp"),
			BuildTree("-sì", "/return"),
			BuildTree("-fkeyk", "$ncec|$nsmod_fkeyk|$nsadp"),
			BuildTree("-tsyìp", "$ncec|$nsmod_tsyìp|$nsadp"),
			BuildTree("-o", "$ncevou|$nsadp"),
			BuildTree("-pe", "$ncev|$nsadp"),
		),
		// Noun suffixed that modify and can follow -fkeyk
		"nsmod_fkeyk": CombineTrees(
			BuildTree("-tsyìp", "$ncec|$nsmod_tsyìp|$nsadp"),
			BuildTree("-o", "$ncevou|$nsadp"),
			BuildTree("-pe", "$ncev|$nsadp"),
		),
		// Noun suffixes that modify and can follow -tsyìp
		"nsmod_tsyìp": CombineTrees(
			BuildTree("-o", "$ncevou|$nsadp"),
			BuildTree("-pe", "$ncev|$nsadp"),
		),
		// Modify noun-part of si-verbs
		"nsmod_si": CombineTrees(
			BuildTree("/return"),
			BuildTree("-o|-pe", "/return"),
			BuildTree("-tsyìp", "/return"),
			BuildTree("-tsyìp", "-o|-pe", "/return"),
		),
		// Noun suffixes from adpositions. To be filled by Dictionary
		"nsadp": EmptyTree(),
		// Noun case endings: vowels
		"ncev": CombineTrees(
			BuildTree("/return"),
			BuildTree("-l|-t|-ti|-r|-ru|-ri|-yä|-ye=yä", "/return"),
		),
		// Noun case endings: after "ia"
		"nceia": CombineTrees(
			BuildTree("/return"),
			BuildTree("-l|-t|-ti|-r|-ru|-ri", "/return"),
			BuildTree("-yä", "/return"), // todo: error node
		),
		// Noun case endings: after "o"/"u"
		"ncevou": CombineTrees(
			BuildTree("/return"),
			BuildTree("-l|-t|-ti|-r|-ru|-ri|-ä|-e=ä", "/return"),
		),
		// Noun case endings: consonants
		"ncec": CombineTrees(
			BuildTree("/return"),
			BuildTree("-ìl|-ti|-it|-ur|-ìri|-ä|-e=ä", "/return"),
		),
		// Noun case endings: loan words (replacing ì)
		"ncevìlw": CombineTrees(
			BuildTree("/return"),
			BuildTree("-ìl|-it|-ur|-ìri|-ä|-e=ä", "/return"),
		),
		// Noun case endings: consonant "t"
		"ncect": CombineTrees(
			BuildTree("/return"),
			BuildTree("-ìl|-it|-ur|-ìri|-ä|-e=ä", "/return"),
		),
		// Noun case endings: consonant "'" (tìftang)
		"ncec'": CombineTrees(
			BuildTree("/return"),
			BuildTree("-ìl|-ti|-it|-ur|-ìri|-ä|-e=ä", "/return"),
		),
		// Noun case endings: diphthongs "ay"/"ey"
		"ncedy": CombineTrees(
			BuildTree("/return"),
			BuildTree("-l|-ìl|-t|-ti|-ur|-ru|-ri|-ä|-e=ä", "/return"),
			BuildTree("-it|-ìri", "/return"), // todo: error node
		),
		// Noun case endings: diphthongs "aw"/"ew"
		"ncedw": CombineTrees(
			BuildTree("/return"),
			BuildTree("-l|-ìl|-ti|-it|-r|-ur|-ri|-ä|-e=ä", "/return"),
		),
		// Pronoun-specific case endings
		"pce_o": CombineTrees(
			BuildTree("-l|-t|-ti|-r|-ru|-ri", "/return"),
		),
		"pce_ng_a": CombineTrees(
			BuildTree("-l|-t|-ti|-r|-ru|-ri", "/return"),
		),
	}
}

const punctuation = " ,;.…—–-?!"
