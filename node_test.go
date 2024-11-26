package lutral

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNode_String(t *testing.T) {
	nodes := []Node{
		{Kind: NKRoot},
		{Kind: NKResult, Value: "1234"},
		{Kind: NKRaw, Value: "kaltxì"},
		{Kind: NKRaw, Value: "-ke_lu_oe_eolì'uvi"},
		{Kind: NKPrefix, Value: "pxe+"},
		{Kind: NKPrefix, Value: "fì"},
		{Kind: NKSuffix, Value: "ti"},
		{Kind: NKInfix, Value: "äpeyk"},
		{Kind: NKSubTree, Value: "np"},
		{Kind: NKReturn},
		{Kind: NKLeafHook},
		{Kind: NKParticle, Value: "rä'ä|ke"},
	}

	for _, node := range nodes {
		assert.Equal(t, node, ParseNode(node.String()))
	}

	badNode := Node{Kind: -1, Value: "hello"}
	assert.Equal(t, "??? (kind: -1, value: \"hello\")", badNode.String())
}

func TestCombineTrees(t *testing.T) {
	t.Run("'ewll/'eveng/'awlo", func(t *testing.T) {
		assert.Equal(t, &Node{
			Kind: NKRoot,
			Children: []Node{
				{Kind: NKRaw, Value: "'e", Children: []Node{
					{Kind: NKRaw, Value: "veng", Children: []Node{
						{Kind: NKResult, Value: "56"},
					}},
					{Kind: NKRaw, Value: "wll", Children: []Node{
						{Kind: NKResult, Value: "2708"},
					}},
				}},
				{Kind: NKRaw, Value: "'awlo", Children: []Node{
					{Kind: NKResult, Value: "4368"},
				}},
			},
		}, CombineTrees(
			BuildTree("'eveng", "=56"),
			BuildTree("'ewll", "=2708"),
			BuildTree("'awlo", "=4368"),
		))
	})

	t.Run("tsatseng/tsat", func(t *testing.T) {
		assert.Equal(t, &Node{
			Kind: NKRoot,
			Children: []Node{
				{Kind: NKSubTree, Value: "np", Children: []Node{
					{Kind: NKRaw, Value: "tsat", Children: []Node{
						{Kind: NKResult, Value: "5980"},
						{Kind: NKRaw, Value: "seng", Children: []Node{
							{Kind: NKSubTree, Value: "nsmod", Children: []Node{
								{Kind: NKResult, Value: "2352"},
							}},
							{Kind: NKSubTree, Value: "nscec", Children: []Node{
								{Kind: NKResult, Value: "2352"},
							}},
						}},
					}},
				}},
			},
		}, CombineTrees(
			BuildTree("$np", "tsatseng", "$nsmod", "=2352"),
			BuildTree("$np", "tsatseng", "$nscec", "=2352"),
			BuildTree("$np", "tsat", "=5980"),
		))
	})

	t.Run("tìk/sìk/tsìk", func(t *testing.T) {
		assert.Equal(t, &Node{
			Kind: NKRoot,
			Children: []Node{
				{Kind: NKRaw, Value: "sìk", Children: []Node{
					{Kind: NKResult, Value: "1796"},
				}},
				{Kind: NKRaw, Value: "t", Children: []Node{
					{Kind: NKRaw, Value: "ìk", Children: []Node{
						{Kind: NKResult, Value: "13294"},
					}},
					{Kind: NKRaw, Value: "are", Children: []Node{
						{Kind: NKResult, Value: "8872"},
					}},
				}},
				{Kind: NKRaw, Value: "ts", Children: []Node{
					{Kind: NKRaw, Value: "ìk", Children: []Node{
						{Kind: NKResult, Value: "8280"},
					}},
					{Kind: NKRaw, Value: "tal", Children: []Node{
						{Kind: NKResult, Value: "4628"},
					}},
				}},
			},
		}, CombineTrees(
			BuildTree("sìk", "=1796"),
			BuildTree("tìk", "=13294"),
			BuildTree("tsìk", "=8280"),
			BuildTree("tare", "=8872"),
			BuildTree("tstal", "=4628"),
		))
	})

	t.Run("uvan/utraltsyìp/utral", func(t *testing.T) {
		tree := CombineTrees(
			BuildTree("uvan", "=2644"),
			BuildTree("utraltsyìp", "=6100"),
			BuildTree("utral", "=2628"),
		)
		tree2 := CombineTrees(
			BuildTree("uvan", "=2644"),
			BuildTree("utral", "=2628"),
			BuildTree("utraltsyìp", "=6100"),
		)

		assert.Equal(t, &Node{
			Kind: NKRoot,
			Children: []Node{
				{Kind: NKRaw, Value: "u", Children: []Node{
					{Kind: NKRaw, Value: "van", Children: []Node{
						{Kind: NKResult, Value: "2644"},
					}},
					{Kind: NKRaw, Value: "tral", Children: []Node{
						{Kind: NKResult, Value: "2628"},
						{Kind: NKRaw, Value: "tsyìp", Children: []Node{
							{Kind: NKResult, Value: "6100"},
						}},
					}},
				}},
			},
		}, tree)

		if !t.Failed() {
			assert.Equal(t, tree, tree2)
		}

		assert.Equal(t, *tree, tree.Copy())
	})

	t.Run("big", func(t *testing.T) {
		tree := CombineTrees(
			BuildTree("tì-", "t", "<us>", "aron", "=13289:n."),
			BuildTree("tì-", "t", "<us>", "are", "=8872:n."),
			BuildTree("tì-", "t", "<us>", "arep", "=13289:n."),
			BuildTree("uvan", "=2644"),
			BuildTree("utraltsyìp", "=6100"),
			BuildTree("utral", "=2628"),
			BuildTree("tì-", "k", "<us>", "ar", "=712:n."),
			BuildTree("tì-", "t", "<us>", "ìran", "=2196:n."),
			BuildTree("tì-", "t", "<us>", "am", "=1988:n."),
			BuildTree("\\k", "<0>", "<1>", "<2>", "ä", "=680"),
		)

		assert.Equal(t, &Node{
			Kind: NKRoot,
			Children: []Node{
				{Kind: NKPrefix, Value: "tì", Children: []Node{
					{Kind: NKRaw, Value: "t", Children: []Node{
						{Kind: NKInfix, Value: "us", Children: []Node{
							{Kind: NKRaw, Value: "a", Children: []Node{
								{Kind: NKRaw, Value: "r", Children: []Node{
									{Kind: NKRaw, Value: "on", Children: []Node{
										{Kind: NKResult, Value: "13289:n."},
									}},
									{Kind: NKRaw, Value: "e", Children: []Node{
										{Kind: NKResult, Value: "8872:n."},
										{Kind: NKRaw, Value: "p", Children: []Node{
											{Kind: NKResult, Value: "13289:n."},
										}},
									}},
								}},
								{Kind: NKRaw, Value: "m", Children: []Node{
									{Kind: NKResult, Value: "1988:n."},
								}},
							}},
							{Kind: NKRaw, Value: "ìran", Children: []Node{
								{Kind: NKResult, Value: "2196:n."},
							}},
						}},
					}},
					{Kind: NKRaw, Value: "k", Children: []Node{
						{Kind: NKInfix, Value: "us", Children: []Node{
							{Kind: NKRaw, Value: "ar", Children: []Node{
								{Kind: NKResult, Value: "712:n."},
							}},
						}},
					}},
				}},
				{Kind: NKRaw, Value: "u", Children: []Node{
					{Kind: NKRaw, Value: "van", Children: []Node{
						{Kind: NKResult, Value: "2644"},
					}},
					{Kind: NKRaw, Value: "tral", Children: []Node{
						{Kind: NKResult, Value: "2628"},
						{Kind: NKRaw, Value: "tsyìp", Children: []Node{
							{Kind: NKResult, Value: "6100"},
						}},
					}},
				}},
				{Kind: NKRaw, Value: "k", Children: []Node{
					{Kind: NKInfix, Value: "0", Children: []Node{
						{Kind: NKInfix, Value: "1", Children: []Node{
							{Kind: NKInfix, Value: "2", Children: []Node{
								{Kind: NKRaw, Value: "ä", Children: []Node{
									{Kind: NKResult, Value: "680"},
								}},
							}},
						}},
					}},
				}},
			},
		}, tree)

		assert.Equal(t, *tree, tree.Copy())
	})
}

func TestBuildTree(t *testing.T) {
	tree1 := BuildTree("tì-|sä-", "t", "<0>|<us>", "<1>", "ar", "<2>", "on", "-ti|-ìri|-ur", "=7336:n.")
	assert.Equal(t, &Node{
		Kind: NKRoot, Children: []Node{
			{Kind: NKPrefix, Value: "tì", Children: []Node{
				{Kind: NKRaw, Value: "t", Children: []Node{
					{Kind: NKInfix, Value: "0", Children: []Node{
						{Kind: NKInfix, Value: "1", Children: []Node{
							{Kind: NKRaw, Value: "ar", Children: []Node{
								{Kind: NKInfix, Value: "2", Children: []Node{
									{Kind: NKRaw, Value: "on", Children: []Node{
										{Kind: NKSuffix, Value: "ti", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ìri", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ur", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
									}},
								}},
							}},
						}},
					}},
					{Kind: NKInfix, Value: "us", Children: []Node{
						{Kind: NKInfix, Value: "1", Children: []Node{
							{Kind: NKRaw, Value: "ar", Children: []Node{
								{Kind: NKInfix, Value: "2", Children: []Node{
									{Kind: NKRaw, Value: "on", Children: []Node{
										{Kind: NKSuffix, Value: "ti", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ìri", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ur", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
									}},
								}},
							}},
						}},
					}},
				}},
			}},
			{Kind: NKPrefix, Value: "sä", Children: []Node{
				{Kind: NKRaw, Value: "t", Children: []Node{
					{Kind: NKInfix, Value: "0", Children: []Node{
						{Kind: NKInfix, Value: "1", Children: []Node{
							{Kind: NKRaw, Value: "ar", Children: []Node{
								{Kind: NKInfix, Value: "2", Children: []Node{
									{Kind: NKRaw, Value: "on", Children: []Node{
										{Kind: NKSuffix, Value: "ti", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ìri", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ur", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
									}},
								}},
							}},
						}},
					}},
					{Kind: NKInfix, Value: "us", Children: []Node{
						{Kind: NKInfix, Value: "1", Children: []Node{
							{Kind: NKRaw, Value: "ar", Children: []Node{
								{Kind: NKInfix, Value: "2", Children: []Node{
									{Kind: NKRaw, Value: "on", Children: []Node{
										{Kind: NKSuffix, Value: "ti", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ìri", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
										{Kind: NKSuffix, Value: "ur", Children: []Node{
											{Kind: NKResult, Value: "7336:n."},
										}},
									}},
								}},
							}},
						}},
					}},
				}},
			}},
		},
	}, tree1)
}

func TestCombineTrees_panics(t *testing.T) {
	assert.Panics(t, func() {
		CombineTrees(&Node{Kind: NKRaw, Value: "leykopx", Children: []Node{}})
	})
}
