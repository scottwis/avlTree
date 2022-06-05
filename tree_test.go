package avlTree_test

import (
	"fmt"
	"math"
	"strings"

	"github.com/scottwis/avlTree"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestEmptyTree(t *testing.T) {
        t.Parallel()
	tree := avlTree.EmptyTree()

	n := tree.Find(avlTree.NumericKey(4))
	require.Nil(t, n)

	n = tree.Find(avlTree.StringKey("findme"))
	require.Nil(t, n)
}

func TestSingleNodeTreeFind(t *testing.T) {
        t.Parallel()
	tree := avlTree.EmptyTree().Update(avlTree.NumericKey(10), "10")

	n := tree.Find(avlTree.NumericKey(4))
	require.Nil(t, n)

	n = tree.Find(avlTree.NumericKey(10))
	require.Equal(t, "10", n.Data())
	require.Equal(t, avlTree.NumericKey(10), n.Key())
}

func TestSingleNodeTreeSize(t *testing.T) {
        t.Parallel()
	tree := avlTree.EmptyTree().Update(avlTree.NumericKey(10), "10")

	require.False(t, tree.IsEmpty())

	require.Equal(t, int64(1), tree.Height())
	require.Nil(t, tree.Left())
	require.Nil(t, tree.Right())
}

func TestSingleNodeTreeLowerBound(t *testing.T) {
        t.Parallel()
	tree := avlTree.EmptyTree().Update(avlTree.NumericKey(10), "10")

	lower := tree.GreatestLowerBound(avlTree.NumericKey(9))
	require.Nil(t, lower)
	lower = tree.GreatestLowerBound(avlTree.NumericKey(10))
	require.Equal(t, "10", lower.Data())

	require.Equal(t, "10:10", dumpTree(tree))
}

func TestSize(t *testing.T) {
        t.Parallel()
	testCases := []struct {
		items      []int
		sizeOutput int64
	}{
		{[]int{}, 0},
		{[]int{1}, 1},
		{[]int{1, 2}, 2},
		{[]int{1, 3, 1, 3, 2}, 3},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%v", tc.items)
		sizeOutput := tc.sizeOutput
		items := tc.items
		t.Run(name, func(t *testing.T) {
                        t.Parallel()
			tree := makeTree(items...)
			require.Equal(t, sizeOutput, tree.Size())

			if len(items) == 0 {
				require.Equal(t, int64(0), tree.Height())
				return
			}
			height := 1 + int64(math.Log2(float64(len(unique(items)))))
			require.Equal(t, height, tree.Height())
		})
	}
}

func TestTreeIteration(t *testing.T) {
        t.Parallel()
	testCases := []struct {
		items           []int
		iterationOutput string
	}{
		{[]int{}, ""},
		{[]int{1}, "1:1"},
		{[]int{1, 2}, "1:1 2:2"},
		{[]int{2, 1}, "1:1 2:2"},
		{[]int{2, 1, 3}, "1:1 2:2 3:3"},
		{[]int{2, 3, 1}, "1:1 2:2 3:3"},
		{[]int{1, 3, 1, 3, 2}, "1:1 2:2 3:3"},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%v", tc.items)
		output := tc.iterationOutput
		items := tc.items
		t.Run(name, func(t *testing.T) {
                        t.Parallel()
			tree := makeTree(items...)
			require.Equal(t, output, dumpTree(tree))
		})
	}
}

func unique(items []int) []int {
	m := map[int]struct{}{}
	for _, item := range items {
		m[item] = struct{}{}
	}
	u := make([]int, 0, len(m))
	for item := range m {
		u = append(u, item)
	}
	return u
}

func dumpTree(n avlTree.Node) string {
	s := ""
	for i := n.Iter(); i.MoveNext(); {
		node := i.Current().(avlTree.Node)
		s += fmt.Sprintf("%v:%v ", node.Key(), node.Data())
	}
	return strings.TrimSpace(s)
}

func makeTree(items ...int) avlTree.Node {
	tree := avlTree.EmptyTree()
	for _, item := range items {
		tree = tree.Update(avlTree.NumericKey(item), fmt.Sprintf("%d", item))
	}
	return tree
}
