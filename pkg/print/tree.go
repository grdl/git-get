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

	return tp.String() + Errors(repos)
}

// Node represents a path fragment in repos tree.
type Node struct {
	val      string
	parent   *Node
	children []*Node
	repo     Printable
	depth    int
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
		depth:  n.depth + 1,
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

// printTree renders the repo tree by recursively traversing the tree nodes.
// If a node doesn't have any children, it's a leaf node containing the repo status.
func (p *TreePrinter) printTree(node *Node, tp treeprint.Tree) {
	if node.children == nil {
		tp.SetValue(printLeaf(node))
	}

	for _, child := range node.children {
		branch := tp.AddBranch(child.val)
		p.printTree(child, branch)
	}
}

func printLeaf(node *Node) string {
	r := node.repo

	// If any errors happened during status loading, don't print the status but "error" instead.
	// Actual error messages are printed in bulk below the tree.
	if len(r.Errors()) > 0 {
		return fmt.Sprintf("%s %s", node.val, red("error"))
	}

	current := r.BranchStatus(r.Current())
	worktree := r.WorkTreeStatus()

	if worktree != "" {
		worktree = fmt.Sprintf("[ %s ]", worktree)
	}

	var str strings.Builder

	if worktree == "" && current == "" {
		str.WriteString(fmt.Sprintf("%s %s %s", node.val, blue(r.Current()), green("ok")))
	} else {
		str.WriteString(fmt.Sprintf("%s %s %s", node.val, blue(r.Current()), strings.Join([]string{yellow(current), red(worktree)}, " ")))
	}

	for _, branch := range r.Branches() {
		status := r.BranchStatus(branch)
		if status == "" {
			status = green("ok")
		}

		str.WriteString(fmt.Sprintf("\n%s%s %s", indentation(node), blue(branch), yellow(status)))
	}

	return str.String()
}

// indentation generates a correct indentation for the branches row to match the links to lower rows.
// It traverses the tree "upwards" and checks if a parent node is the youngest one (ie, there are no more sibling at the same level).
// If it is, it means that level should be indented with empty spaces because there is nothing to link to anymore.
// If it isn't the youngest, that level needs to be indented using a "|" link.
func indentation(node *Node) string {
	// Slice of levels. Slice index is node depth, true value means the node is the youngest.
	levels := make([]bool, node.depth)

	// Traverse until node has no parents (ie, we reached the root).
	n := node
	for n.parent != nil {
		levels[n.depth-1] = n.isYoungest()
		n = n.parent
	}

	var indent strings.Builder

	const space = "    "
	const link = "│   "
	for _, y := range levels {
		if y {
			indent.WriteString(space)
		} else {
			indent.WriteString(link)
		}
	}

	// Finally, indent by the size of node name (to match the rest of the branches)
	indent.WriteString(strings.Repeat(" ", len(node.val)+1))

	return indent.String()
}

// isYoungest checks if the node is the last one in the slice of children
func (n *Node) isYoungest() bool {
	if n.parent == nil {
		return true
	}

	sisters := n.parent.children
	var myIndex int
	for i, sis := range sisters {
		if sis.val == n.val {
			myIndex = i
			break
		}
	}
	return myIndex == len(sisters)-1
}
