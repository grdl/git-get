package print

import (
	"git-get/pkg/repo"
	"path/filepath"
	"strings"
)

// SmartPrinter implements Printer interface and provides methods for printing repos and their statuses.
// It's "smart" because it automatically folds branches which only have a single child and indents branches with many children.
type SmartPrinter struct {
	// length is the size (number of chars) of the currently processed line.
	// It's used to correctly indent the lines with branches status.
	length int
}

// Print generates a list of repositories and their statuses.
func (p *SmartPrinter) Print(root string, repos []*repo.Repo) string {
	tree := buildTree(root, repos)

	return p.printSmartTree(tree)
}

// printSmartTree recursively traverses the tree and prints its nodes.
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
func (p *SmartPrinter) printSmartTree(node *Node) string {
	if node.children == nil {
		// If node is a leaf, print repo name and its status and finish processing this node.
		value := node.val

		// TODO: Ugly
		// If this is called from tests the repo will be nil and we should return just the name without the status.
		if node.repo.Repository == nil {
			return value
		}

		value += " " + printWorktreeStatus(node.repo)

		// Print the status of each branch on a new line, indented to match the position of the current branch name.
		indent := "\n" + strings.Repeat(" ", p.length+len(node.val))
		for _, branch := range node.repo.Status.Branches {
			// Don't print the status of the current branch. It was already printed above.
			if branch.Name == node.repo.Status.CurrentBranch {
				continue
			}

			value += indent + printBranchStatus(branch)
		}

		return value
	}

	val := node.val + string(filepath.Separator)

	shift := ""
	if node.parent == nil {
		// If node is a root, print its children on a new line without indentation.
		shift = "\n"
	} else if len(node.children) == 1 {
		// If node has only a single child, print it on the same line as its parent.
		// Setting node's depth to the same as parent's ensures that its children will be indented only once even if
		// node's path has multiple levels above.
		node.depth = node.parent.depth

		p.length += len(val)
	} else {
		// If node has multiple children, print each of them on a new line
		// and indent them once relative to the parent
		node.depth = node.parent.depth + 1
		shift = "\n" + strings.Repeat("    ", node.depth)
		p.length = 0
	}

	for _, child := range node.children {
		p.length += len(shift)
		val += shift + p.printSmartTree(child)
		p.length = 0
	}

	return val
}
