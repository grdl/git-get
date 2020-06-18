package print

import (
	"git-get/pkg/repo"
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
)

// TreePrinter implements Printer interface and provides methods for printing repos and their statuses.
type TreePrinter struct{}

// Print generates a tree view of repos and their statuses.
func (p *TreePrinter) Print(root string, repos []*repo.Repo) string {
	tree := buildTree(root, repos)

	tp := treeprint.New()
	tp.SetValue(root)

	p.printTree(tree, tp)

	return tp.String()
}

// Node represents a node (ie. path fragment) in a repos tree.
type Node struct {
	val      string
	depth    int // depth is a nesting depth used when rendering a smart tree, not a depth level of a tree node.
	parent   *Node
	children []*Node
	repo     *repo.Repo
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

// buildTree builds a directory tree of paths to repositories.
// Each node represents a directory in the repo path.
// Each leaf (final node) contains a pointer to the repo.
func buildTree(root string, repos []*repo.Repo) *Node {
	tree := Root(root)

	for _, r := range repos {
		path := strings.TrimPrefix(r.Path, root)
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
					node.repo = r
				}

				continue
			}
			node = child
		}
	}
	return tree
}

func (p *TreePrinter) printTree(node *Node, tp treeprint.Tree) {
	if node.children == nil {
		tp.SetValue(node.val + " " + printWorktreeStatus(node.repo))

		for _, branch := range node.repo.Status.Branches {
			// Don't print the status of the current branch. It was already printed above.
			if branch.Name == node.repo.Status.CurrentBranch {
				continue
			}
			tp.AddNode(printBranchStatus(branch))
		}

	}

	for _, child := range node.children {
		branch := tp.AddBranch(child.val)
		p.printTree(child, branch)
	}
}
