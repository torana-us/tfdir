package main

import (
	"path/filepath"
	"reflect"
	"strings"
)

type Node struct {
	value    string
	nextMap  map[string]bool
	children Nodes
}

type Nodes []*Node

type Path struct {
	Value string
}

func NewPath(value string) *Path {
	return &Path{
		Value: value,
	}
}

func (p *Path) HeadAndRest() (head string, rest string) {
	b, a, _ := strings.Cut(p.Value, "/")
	return b, a
}

func NewNode(value string) *Node {
	return &Node{
		value:    value,
		children: Nodes{},
		nextMap:  map[string]bool{},
	}
}

func NewTree(path Path) *Node {
	h, r := path.HeadAndRest()
	if h == "" {
		return nil
	}
	root := NewNode(h)
	child := NewTree(*NewPath(r))

	if child != nil {
		root.Add(child)
	}

	return root
}

func (n *Node) Eq(other *Node) bool {
	return reflect.DeepEqual(n, other)
}

func (n *Node) isLeaf() bool {
	return len(n.children) == 0
}

func (n *Node) findChild(value string) *Node {
	for _, child := range n.children {
		if child.value == value {
			return child
		}
	}

	return nil
}

func (n *Node) Add(others ...*Node) *Node {
	for _, other := range others {
		if n.nextMap[other.value] {
			child := n.findChild(other.value)
			child.Add(other.children...)
		} else {
			n.nextMap[other.value] = true
			n.children = append(n.children, other)
		}
	}
	return n
}

func MakeTreeMap(dirs []Path) map[string]*Node {
	tree_map := map[string]*Node{}

	for _, path := range dirs {
		if path.Value == "" {
			continue
		}

		h, rest := path.HeadAndRest()

		tree, ok := tree_map[h]

		if !ok {
			tree = NewTree(path)
			tree_map[h] = tree
		}

		if rest != "" {
			tree.Add(NewTree(*NewPath(rest)))
		}
	}

	return tree_map
}

func (n *Node) allPathHelper(prefix string) []Path {
	if n.isLeaf() {
		return []Path{
			*NewPath(filepath.Join(prefix, n.value)),
		}
	}

	paths := []Path{}
	prefix = filepath.Join(prefix, n.value)

	for _, child := range n.children {
		paths = append(paths, child.allPathHelper(prefix)...)
	}

	return paths
}

func (n *Node) AllPath() []Path {
	return n.allPathHelper("")
}

func (n *Node) Search(path Path) (*Node, bool) {
	h, r := path.HeadAndRest()

	if h != n.value {
		return nil, false
	}

	if r == "" {
		return n, true
	}

	r_path := NewPath(r)

	h, _ = r_path.HeadAndRest()

	if !n.nextMap[h] {
		return nil, false
	}

	return n.findChild(h).Search(*r_path)
}
