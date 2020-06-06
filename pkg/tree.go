package pkg

import (
	"path/filepath"
	"strings"
)

// Node represents a node in a repos tree
type Node struct {
	val      string
	depth    int // depth is a nesting depth used when rendering a tree, not an depth level of a node inside the tree
	parent   *Node
	children []*Node
	repo     *Repo
}

// Root creates a new root of a tree
func Root(val string) *Node {
	root := &Node{
		val: val,
	}
	return root
}

// Add adds a child node
func (n *Node) Add(val string) *Node {
	if n.children == nil {
		n.children = make([]*Node, 0)
	}

	child := &Node{
		val:    val,
		parent: n,
	}
	n.children = append(n.children, child)
	return child
}

// GetChild finds a node with val inside this node's children (only 1 level deep).
// Returns pointer to found child or nil if node doesn't have any children or doesn't have a child with sought value.
func (n *Node) GetChild(val string) *Node {
	if n.children == nil {
		return nil
	}

	for _, child := range n.children {
		if child.val == val {
			return child
		}
	}

	return nil
}

// BuildTree builds a directory tree of paths to repositories.
// Each node represents a directory in the repo path.
// Each leaf (final node) contains a pointer to the repo.
func BuildTree(root string, repos []*Repo) *Node {
	tree := Root(root)

	for _, repo := range repos {
		path := strings.TrimPrefix(repo.path, root)
		path = strings.Trim(path, string(filepath.Separator))
		subs := strings.Split(path, string(filepath.Separator))

		// For each path fragment, start at the root of the tree
		// and check if the fragment exist among the children of the node.
		// If not, add it to node's children and move to next fragment.
		// If it does, just move to the next fragment.
		node := tree
		for i, sub := range subs {
			child := node.GetChild(sub)
			if child == nil {
				node = node.Add(sub)

				// If that's the last fragment, it's a tree leaf and needs a *Repo attached.
				if i == len(subs)-1 {
					node.repo = repo
				}

				continue
			}
			node = child
		}
	}
	return tree
}

// RenderTree returns a string representation of repos tree.
// It recursively traverses the tree and prints its nodes.
// If a node contains multiple children, they are be printed in new lines and indented.
// If a node contains only a single child, it is printed in the same line using path separator.
// For better readability the first level (repos hosts) is not indented.
//
// Example:
// Following paths:
//   /repos/github.com/user/repo1
//   /repos/github.com/user/repo2
//   /repos/github.com/another/repo
//
// will render a tree:
//   /repos/
//   github.com/
//       user/
//           repo1
//           repo2
//       another/repo
//
func RenderTree(node *Node) string {
	if node.children == nil {
		// If node is a leaf, print repo name and its status and finish processing this node.
		return node.val + " " + renderWorktreeStatus(node.repo)
	}

	shift := ""
	if node.parent == nil {
		// If node is a root, print its children on a new line without indentation.
		shift = "\n"
	} else if len(node.children) == 1 {
		// If node has only a single child, print it on the same line as its parent.
		// Setting node's depth to the same as parent's ensures that its children will be indented only once even if
		// node's path has multiple levels above.
		node.depth = node.parent.depth
	} else {
		// If node has multiple children, print each of them on a new line
		// and indent them once relative to the parent
		node.depth = node.parent.depth + 1
		shift = "\n" + strings.Repeat("\t", node.depth)
	}

	val := node.val + string(filepath.Separator)
	for _, child := range node.children {
		val += shift + RenderTree(child)
	}

	return val
}