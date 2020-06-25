package print

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xlab/treeprint"
)

// TreePrinter prints list of repos in a directory tree format.
type TreePrinter struct {
}

// NewTreePrinter creates a TreePrinter.
func NewTreePrinter() *TreePrinter {
	return &TreePrinter{}
}

// Print generates a tree view of repos and their statuses.
func (p *TreePrinter) Print(root string, repos []Printable) string {
	if len(repos) == 0 {
		return fmt.Sprintf("There are no git repos under %s", root)
	}

	tree := buildTree(root, repos)
	tp := treeprint.New()
	tp.SetValue(root)

	p.printTree(tree, tp)

	return tp.String()
}

// Node represents a path fragment in repos tree.
type Node struct {
	val      string
	parent   *Node
	children []*Node
	repo     Printable
}

// Root creates a new root of a tree.
func Root(val string) *Node {
	root := &Node{
		val: val,
	}
	return root
}

// Add adds a child node with given value to a current node.
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
func buildTree(root string, repos []Printable) *Node {
	tree := Root(root)

	for _, r := range repos {
		path := strings.TrimPrefix(r.Path(), root)
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
		r := node.repo
		current := r.BranchStatus(r.Current())
		worktree := r.WorkTreeStatus()

		if worktree != "" {
			worktree = fmt.Sprintf("[ %s ]", worktree)
		}

		if worktree == "" && current == "" {
			tp.SetValue(node.val + " " + blue(r.Current()) + " " + green("ok"))
		} else {
			tp.SetValue(node.val + " " + blue(r.Current()) + " " + strings.Join([]string{yellow(current), red(worktree)}, " "))
		}

		for _, branch := range r.Branches() {
			status := r.BranchStatus(branch)
			if status == "" {
				status = green("ok")
			}

			tp.AddNode(blue(branch) + " " + yellow(status))
		}
	}

	for _, child := range node.children {
		branch := tp.AddBranch(child.val)
		p.printTree(child, branch)
	}
}
