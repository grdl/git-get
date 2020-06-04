package pkg

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

var paths = []string{

	// "/home/grdl/repositories/gitlab.com/grdl/testflux",
	"/home/grdl/repositories/bitbucket.org/gridarrow/istio",
	"/home/grdl/repositories/bitbucket.org/grdl/bob",
	"/home/grdl/repositories/github.com/fboender/multi-git-status",
	"/home/grdl/repositories/github.com/grdl/git-get",
	"/home/grdl/repositories/github.com/grdl/testflux",
	"/home/grdl/repositories/github.com/johanhaleby/kubetail",
	"/home/grdl/repositories/gitlab.com/grdl/git-get",
	"/home/grdl/repositories/gitlab.com/grdl/grafana-dashboard-builder",
	"/home/grdl/repositories/gitlab.com/grdl/dotfiles",
}

func TestTree(t *testing.T) {
	InitConfig()
	root := viper.GetString(KeyReposRoot)

	tree := Root(root)

	for _, path := range paths {
		p := strings.TrimPrefix(path, root)
		p = strings.Trim(p, string(filepath.Separator))
		subs := strings.Split(p, string(filepath.Separator))

		node := tree
		for _, sub := range subs {
			child := node.GetChild(sub)
			if child == nil {
				node = node.Add(sub)
				continue
			}
			node = child
		}
	}

	fmt.Println(tree)
}

func process(node *Node, val string) *Node {
	found := node.GetChild(val)
	if found == nil {
		added := node.Add(val)
		return added
	}
	return found
}

type Node struct {
	val      string
	parent   *Node
	children []*Node
}

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

	new := &Node{
		val:    val,
		parent: n,
	}
	n.children = append(n.children, new)
	return new
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
