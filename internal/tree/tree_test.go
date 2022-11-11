package tree_test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/torana-us/tfdir/internal/tree"
)

func Test_Eq_Simple(t *testing.T) {
	const Equal, NotEqual = true, false

	aa := map[string]struct {
		node_1   string
		node_2   string
		expected bool
	}{
		"等しい":   {"test", "test", Equal},
		"等しくない": {"test", "test2", NotEqual},
	}

	for name, a := range aa {
		t.Run(name, func(t *testing.T) {
			node_1 := NewNode(a.node_1)
			node_2 := NewNode(a.node_2)

			actual := node_1.Eq(node_2)

			if actual != a.expected {
				t.Fatalf("actual: %v, expected: %v", actual, a.expected)
			}
		})
	}
}

func Test_Eq_HadChildren(t *testing.T) {
	const Equal, NotEqual = true, false

	aa := map[string]struct {
		node_1_root     string
		node_1_children []string
		node_2_root     string
		node_2_children []string
		expected        bool
	}{
		"等しい":            {"test", []string{"child1"}, "test", []string{"child1"}, Equal},
		"childrenが等しくない": {"test", []string{"child1"}, "test", []string{"child2"}, NotEqual},
		"rootが等しくない":     {"test", []string{"child1"}, "test2", []string{"child1"}, NotEqual},
	}

	for name, a := range aa {
		t.Run(name, func(t *testing.T) {
			root_1 := NewNode(a.node_1_root)
			for _, c := range a.node_1_children {
				root_1.Add(NewNode(c))
			}
			root_2 := NewNode(a.node_2_root)
			for _, c := range a.node_2_children {
				root_2.Add(NewNode(c))
			}

			actual := root_1.Eq(root_2)

			if actual != a.expected {
				t.Fatalf("actual: %v, expected: %v", actual, a.expected)
			}
		})
	}
}

func Test_Add(t *testing.T) {
	root := NewNode("/").Add(NewNode("child"))
	added := NewNode("child").Add((NewNode("grandchild")))
	expected := NewNode("/").Add(NewNode("child").Add(NewNode("grandchild")))

	root.Add(added)

	result := root.Eq(expected)

	fmt.Println(result)
}

func Test_PathHeadAndRest(t *testing.T) {
	aa := []struct {
		path_value    string
		expected_head string
		expected_rest string
	}{
		{"", "", ""},
		{"terraform", "terraform", ""},
		{"terraform/ansible", "terraform", "ansible"},
		{"terraform/ansible/php", "terraform", "ansible/php"},
	}

	for _, a := range aa {
		t.Run(a.path_value, func(t *testing.T) {
			path := NewPath(a.path_value)
			head, rest := path.HeadAndRest()

			if head != a.expected_head {
				t.Fatalf("actual: %s, expected: %s", head, a.expected_head)
			}

			if rest != a.expected_rest {
				t.Fatalf("actual: %s, expected: %s", rest, a.expected_rest)
			}
		})
	}
}

func Test_NewTree(t *testing.T) {
	tree := NewTree(*NewPath("terraform/php/golang"))
	expected := NewNode("terraform").Add(NewNode("php").Add(NewNode("golang")))

	if !tree.Eq(expected) {
		t.Fatalf("Something wrong")
	}
}

func Test_MakeTreeMap(t *testing.T) {
	terraform_path := "terraform/environments/"
	module_path := "modules/"
	dirs := []string{
		terraform_path,
		module_path,
	}
	expected_terraform_tree := NewTree(*NewPath(terraform_path))
	expected_module_tree := NewTree(*NewPath(module_path))

	tree_map := MakeTreeMap(dirs)

	if !tree_map["terraform"].Eq(expected_terraform_tree) {
		t.Fatalf("terraform tree something wrong")
	}
	if !tree_map["modules"].Eq(expected_module_tree) {
		t.Fatalf("modules tree something wrong")
	}
}

func Test_AllPath_Simple(t *testing.T) {
	aa := []struct {
		path_value string
	}{
		{"terraform"},
		{"terraform/ansible"},
		{"terraform/ansible/php"},
	}

	for _, a := range aa {
		t.Run(a.path_value, func(t *testing.T) {
			path := *NewPath(a.path_value)
			tree := NewTree(path)

			actual := tree.AllPath()

			if !reflect.DeepEqual([]Path{path}, actual) {
				t.Fatalf("actual: %v\nexpected: %v", actual, []Path{path})
			}
		})
	}
}

func Test_AllPath(t *testing.T) {
	dirs := []string{
		"a/b/c",
		"a/b/d",
		"a/g",
	}

	tree := MakeTreeMap(dirs)["a"]

	actual := tree.AllPathString()

	if !reflect.DeepEqual(dirs, actual) {
		t.Fatalf("actual: %v\nexpected: %v", actual, dirs)
	}
}

func Test_Search(t *testing.T) {
	const Found, NotFound = true, false
	aa := []struct {
		path        string
		search_path string
		expected    bool
	}{
		{"a/b/c", "a/b/c", Found},
		{"a/b/c", "a/b", Found},
		{"a/b/c", "a", Found},
		{"a/b/c", "c", NotFound},
		{"a/b/c", "a/c", NotFound},
		{"a/b/c", "a/b/d", NotFound},
		{"a/b/c", "a/b/c/d", NotFound},
	}
	for _, a := range aa {
		t.Run(fmt.Sprintf("(%s)<---(%s)", a.path, a.search_path), func(t *testing.T) {
			tree := NewTree(*NewPath(a.path))
			_, ok := tree.Search(*NewPath(a.search_path))

			if ok != a.expected {
				t.Fatalf("Not Expected")
			}
		})
	}

}
