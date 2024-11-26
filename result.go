package lutral

import (
	"fmt"
	"slices"
	"strings"
)

type Result struct {
	ID        string   `json:"id"`
	Position  int      `json:"index,omitempty"`
	PoS       string   `json:"pos,omitempty"`
	Remainder string   `json:"remainder,omitempty"`
	Prefixes  []string `json:"prefixes,omitempty"`
	Infixes   []string `json:"infixes,omitempty"`
	Suffixes  []string `json:"suffixes,omitempty"`
	Lenitions []string `json:"lenitions,omitempty"`
	Particles []string `json:"particles,omitempty"`
}

func (result *Result) String() string {
	sb := strings.Builder{}
	if result.Position > 0 {
		sb.WriteRune('[')
		_, _ = fmt.Fprintf(&sb, "%d", result.Position)
		sb.WriteString("] ")
	}

	sb.WriteString(result.ID)
	if result.PoS != "" {
		sb.WriteRune(':')
		sb.WriteString(result.PoS)
	}

	if len(result.Prefixes) > 0 {
		sb.WriteRune(' ')
		for _, prefix := range result.Prefixes {
			sb.WriteString(prefix)
			sb.WriteByte('-')
		}
	}

	if len(result.Infixes) > 0 {
		sb.WriteRune(' ')
		sb.WriteRune('<')
		for i, infix := range result.Infixes {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(infix)
		}
		sb.WriteRune('>')
	}

	if len(result.Suffixes) > 0 {
		sb.WriteRune(' ')
		for _, suffix := range result.Suffixes {
			sb.WriteByte('-')
			sb.WriteString(suffix)
		}
	}

	if len(result.Lenitions) > 0 {
		sb.WriteRune(' ')
		for i, lenition := range result.Lenitions {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(lenition)
		}
	}

	if len(result.Particles) > 0 {
		sb.WriteString(" [")
		for i, lenition := range result.Particles {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(lenition)
		}
		sb.WriteRune(']')
	}

	if len(result.Remainder) > 0 {
		sb.WriteString(" +")
		sb.WriteString(result.Remainder)
	}

	return sb.String()
}

func (result *Result) CoveredBy(template Result) bool {
	return result.ID == template.ID && result.PoS == template.PoS &&
		sliceCovered(template.Prefixes, result.Prefixes) &&
		sliceCovered(template.Suffixes, result.Suffixes) &&
		sliceCovered(template.Infixes, result.Infixes) &&
		sliceCovered(template.Lenitions, result.Lenitions) &&
		sliceCovered(template.Particles, result.Particles)
}

func (result *Result) AddAffixesFrom(other, template Result) {
	for _, prefix := range other.Prefixes {
		if !slices.Contains(template.Prefixes, prefix) {
			result.Prefixes = append(result.Prefixes, prefix)
		}
	}
	for _, infix := range other.Infixes {
		if !slices.Contains(template.Infixes, infix) {
			result.Infixes = append(result.Infixes, infix)
		}
	}
	for _, suffix := range other.Suffixes {
		if !slices.Contains(template.Suffixes, suffix) {
			result.Suffixes = append(result.Suffixes, suffix)
		}
	}
	for _, lenition := range other.Lenitions {
		if !slices.Contains(template.Lenitions, lenition) {
			result.Lenitions = append(result.Lenitions, lenition)
		}
	}
	for _, particle := range other.Particles {
		if !slices.Contains(template.Particles, particle) {
			result.Particles = append(result.Particles, particle)
		}
	}
}

func simplestResultSet(results []Result) []Result {
	if len(results) == 0 {
		return nil
	}

	res := make([]Result, results[len(results)-1].Position)
	for _, result := range results {
		p := result.Position - 1
		if res[p].ID == "" {
			res[p] = result
		}

		a := len(result.Lenitions) + len(result.Prefixes) + len(result.Infixes) + len(result.Suffixes) + len(result.Particles)
		b := len(res[p].Lenitions) + len(res[p].Prefixes) + len(res[p].Infixes) + len(res[p].Suffixes) + len(res[p].Particles)
		if a < b {
			res[p] = result
		}
	}

	return res
}

func sliceCovered(template, actual []string) bool {
	if len(template) == 0 {
		return true
	}

	for _, template := range template {
		found := false
		for _, actual := range actual {
			if actual == template {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}
