package lutral

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	Kind     NodeKind `json:"k"`
	Value    string   `json:"v,omitempty"`
	Children []Node   `json:"c,omitempty"`
}

func (node *Node) MergedWith(other Node) *Node {
	merged := node.MergeFrom(other)
	if !merged {
		panic("MergedWith does not allow failure")
	}

	return node
}

func (node *Node) MergeFrom(other Node) bool {
	if node.Value == other.Value && node.Kind == other.Kind {
		for _, otherChild := range other.Children {
			merged := false
			for i := range node.Children {
				if node.Children[i].MergeFrom(otherChild) {
					merged = true
					break
				}
			}

			if !merged {
				node.Children = append(node.Children, otherChild)
			}
		}

		return true
	}

	// Don't merge non-raw further
	if node.Kind != NKRaw || other.Kind != NKRaw {
		return false
	}

	longestCommon := ""
	if strings.HasPrefix(node.Value, other.Value) {
		longestCommon = other.Value
	} else if strings.HasPrefix(other.Value, node.Value) {
		longestCommon = node.Value
	} else {
		shortest := node.Value
		if len(other.Value) < len(node.Value) {
			shortest = other.Value
		}

		for i := range shortest {
			if node.Value[:i] == other.Value[:i] {
				longestCommon = shortest[:i]
			} else {
				break
			}
		}
	}
	if longestCommon == "" {
		return false
	}

	// Lenition safety with ts, ejectives and tìftang.
	if longestCommon == "'" {
		return false
	}
	if longestCommon == "t" && (strings.HasPrefix(node.Value, "ts") || strings.HasPrefix(node.Value, "tx") || strings.HasPrefix(other.Value, "ts") || strings.HasPrefix(other.Value, "tx")) {
		return false
	}
	if longestCommon == "p" && (strings.HasPrefix(node.Value, "px") || strings.HasPrefix(other.Value, "px")) {
		return false
	}
	if longestCommon == "k" && (strings.HasPrefix(node.Value, "kx") || strings.HasPrefix(other.Value, "kx")) {
		return false
	}
	if longestCommon == "n" && (strings.HasPrefix(node.Value, "ng") || strings.HasPrefix(other.Value, "ng")) {
		return false
	}

	if longestCommon == node.Value {
		other.Value = strings.TrimPrefix(other.Value, longestCommon)

		merged := false
		for i := range node.Children {
			if node.Children[i].MergeFrom(other) {
				merged = true
				break
			}
		}

		if !merged {
			node.Children = append(node.Children, other)
		}
	} else if longestCommon == other.Value {
		node2 := *node
		*node = other.Copy()

		node2.Value = strings.TrimPrefix(node2.Value, longestCommon)
		node.Children = append(node.Children, node2)
	} else {
		node2 := *node
		node2.Value = strings.TrimPrefix(node.Value, longestCommon)
		other.Value = strings.TrimPrefix(other.Value, longestCommon)

		node.Value = longestCommon
		node.Children = []Node{node2, other}
	}

	return true
}

func (node *Node) AndThenResult(id string) *Node {
	return node.AndThen(Node{Kind: NKResult, Value: id})
}

// AndThen runs AppendFrom and returns back the same object WITHOUT copying it.
func (node *Node) AndThen(other Node) *Node {
	node.AppendFrom(other)
	return node
}

// AppendFrom appends the tree at each leaf descendant of this node.
// This does create a copy of each node.
//
// If other is a NKRoot node, its children are added instead.
func (node *Node) AppendFrom(other Node) {
	var leavesBuffer [32]*Node

	newChildren := other.Children
	if other.Kind != NKRoot {
		newChildren = []Node{other}
	}

	leaves := node.Leaves(leavesBuffer[:0])

	var oldHooksParents []*Node
	var oldHooksIndices []int

	for _, leaf := range leaves {
		for i, child := range leaf.Children {
			if child.Kind == NKLeafHook {
				oldHooksParents = append(oldHooksParents, leaf)
				oldHooksIndices = append(oldHooksIndices, i)
			}
		}

		fakeLeaf := *leaf
		fakeLeaf.Children = newChildren
		*leaf = *leaf.MergedWith(fakeLeaf)
	}

	for i, parent := range oldHooksParents {
		index := oldHooksIndices[i]
		parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
	}
}

func (node *Node) RemoveLeafHooks() {
	for i, child := range node.Children {
		if child.Kind == NKLeafHook {
			node.Children = append(node.Children[:i], node.Children[i+1:]...)
			return
		}
	}

	for i := range node.Children {
		node.Children[i].RemoveLeafHooks()
	}
}

func (node *Node) Leaves(buf []*Node) []*Node {
	if len(node.Children) == 0 {
		return append(buf, node)
	}

	// Make room for more if there are leaf-hooks
	for i := range node.Children {
		child := &node.Children[i]
		if child.Kind == NKLeafHook {
			if cap(node.Children) < len(node.Children)+2 {
				node.Children = append(make([]Node, 0, len(node.Children)+4), node.Children...)
			}
			buf = append(buf, node)
			break
		}
	}

	for i := range node.Children {
		child := &node.Children[i]
		if child.Kind == NKLeafHook {
			continue
		}

		buf = child.Leaves(buf)
	}

	return buf
}

// SearchReplace goes through the entire tree to replace nodes with those returned by the callback.
// If it gets a Root node from the callback, it will incorporate all its children instead.
func (node *Node) SearchReplace(cb func(*Node) *Node) {
	for i := range node.Children {
		res := cb(&node.Children[i])
		if res != nil && res != &node.Children[i] {
			if len(res.Children) >= 1 {
				node.Children[i] = res.Children[0]
			}
			if len(res.Children) > 1 {
				node.Children = append(node.Children, res.Children[1:]...)
			}
		}

		node.Children[i].SearchReplace(cb)
	}
}

func (node *Node) Copy() Node {
	nodeCopy := *node
	if nodeCopy.Children != nil {
		nodeCopy.Children = make([]Node, 0, len(node.Children))
		for _, child := range node.Children {
			nodeCopy.Children = append(nodeCopy.Children, child.Copy())
		}
	}

	return nodeCopy
}

func (node *Node) String() string {
	switch node.Kind {
	case NKRoot:
		return "/root"
	case NKResult:
		return "=" + node.Value
	case NKRaw:
		if ParseNode(node.Value).Kind == NKRaw {
			return node.Value
		} else {
			return "\\" + node.Value
		}
	case NKPrefix:
		if strings.HasSuffix(node.Value, "+") {
			return node.Value
		} else {
			return node.Value + "-"
		}
	case NKInfix:
		return "<" + node.Value + ">"
	case NKSuffix:
		return "-" + node.Value
	case NKSubTree:
		return "$" + node.Value
	case NKReturn:
		return "/return"
	case NKLeafHook:
		return "/hook"
	case NKParticle:
		return "[" + node.Value + "]"
	default:
		return fmt.Sprintf("??? (kind: %d, value: %#v)", node.Kind, node.Value)
	}
}

func (node *Node) Size() int {
	if len(node.Children) == 0 {
		return 1
	}

	total := 1
	for _, child := range node.Children {
		total += child.Size()
	}

	return total
}

func (node *Node) SortChildren() {
	sort.Slice(node.Children, func(i, j int) bool {
		ci := &node.Children[i]
		cj := &node.Children[j]
		ri := ci.Kind == NKRaw
		rj := cj.Kind == NKRaw
		if ri && rj {
			return ci.Value < cj.Value
		}
		if ri != rj {
			return ri == false
		}

		if ci.Kind == cj.Kind {
			return ci.Value < cj.Value
		} else {
			return ci.Kind < cj.Kind
		}
	})

	for i := range node.Children {
		node.Children[i].SortChildren()
	}
}

func (node *Node) Compact() {
	if node.Kind == NKRaw && len(node.Children) == 1 && node.Children[0].Kind == NKRaw {
		node.Value = node.Value + node.Children[0].Value
		node.Children = node.Children[0].Children
		node.Compact()
	}

	for i := range node.Children {
		node.Children[i].Compact()
	}
}

type NodeKind int

const (
	// NKRoot indicate that it's a root node. It should not be used further into the tree.
	NKRoot = iota
	// NKResult is the final node in a tree that indicates that a dictionary entry has been found.
	// It should only return if there is not more remaining.
	NKResult
	// NKRaw matches the actual text.
	NKRaw
	// NKPrefix matches the text and adds it as a prefix if equal.
	NKPrefix
	// NKInfix matches the text and adds it as an infix if equal.
	// It a list with the name exists, it will check all of them.
	// It can target multiple if comma separated. If you have an infix
	// that is a beginning of another, it must come after the longer one
	// (e.g. äpeyk before äp) as it stops exploring after accepting one.
	NKInfix
	// NKSuffix matches the text and adds it as a suffix. You can add =B
	// to get it to add another suffix name when matched, e.g. "e=ä"
	NKSuffix
	// NKSubTree executes a subtree until it finds. If inside a subtree, it will replace the current.
	NKSubTree
	// NKReturn returns out of a subtree.
	NKReturn
	// NKLeafHook makes the parent return as a leaf.
	NKLeafHook
	// NKParticle allows a sub-result with the given criteria
	NKParticle
)

func ParseNode(s string) Node {
	switch {
	case strings.HasPrefix(s, "="):
		return Node{Kind: NKResult, Value: strings.TrimPrefix(s, "=")}
	case strings.HasSuffix(s, "+"):
		return Node{Kind: NKPrefix, Value: s}
	case strings.HasSuffix(s, "-"):
		return Node{Kind: NKPrefix, Value: strings.TrimSuffix(s, "-")}
	case strings.HasPrefix(s, "<") && strings.HasSuffix(s, ">"):
		return Node{Kind: NKInfix, Value: strings.TrimSuffix(strings.TrimPrefix(s, "<"), ">")}
	case strings.HasPrefix(s, "-"):
		return Node{Kind: NKSuffix, Value: strings.TrimPrefix(s, "-")}
	case strings.HasPrefix(s, "$"):
		return Node{Kind: NKSubTree, Value: strings.TrimPrefix(s, "$")}
	case strings.HasPrefix(s, "\\"):
		return Node{Kind: NKRaw, Value: strings.TrimPrefix(s, "\\")}
	case strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"):
		return Node{Kind: NKParticle, Value: strings.TrimLeft(strings.TrimRight(s, "]"), "[")}
	case s == "/return":
		return Node{Kind: NKReturn}
	case s == "/hook":
		return Node{Kind: NKLeafHook}
	case s == "/root":
		return Node{Kind: NKRoot}
	default:
		return Node{Kind: NKRaw, Value: s}
	}
}

func CombineTrees(trees ...*Node) *Node {
	current := &Node{Kind: NKRoot, Children: make([]Node, 0, len(trees)*2)}

	for _, tree := range trees {
		if !current.MergeFrom(*tree) {
			panic("CombineTree failed to merge one or more of the trees")
		}
	}

	return current
}

func EmptyTree() *Node {
	return &Node{Kind: NKRoot}
}

func CopyTree(node Node) *Node {
	nodeCopy := node.Copy()
	return &nodeCopy
}

func BuildTree(input ...string) *Node {
	var leavesBuf [128]*Node
	var newNodeBuf [16]Node

	root := &Node{Kind: NKRoot}
	leaves := leavesBuf[:0]
	newNodes := newNodeBuf[:0]

	for _, value := range input {
		if value == "" {
			continue
		}

		newNodes = newNodes[:0]
		for _, s := range strings.Split(value, "|") {
			newNodes = append(newNodes, ParseNode(s))
		}

		leaves = root.Leaves(leaves[:0])
		for _, leaf := range leaves {
			for _, newNode := range newNodes {
				leaf.Children = append(leaf.Children, newNode.Copy())
			}
		}
	}

	return root
}
