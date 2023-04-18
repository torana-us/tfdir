package cmd

import (
	"reflect"
	"sort"
	"testing"
)

func Test_getDirs(t *testing.T) {
	executed_dirs := []string{
		"../testdata/terraform/environments/test_1",
		"../testdata/terraform/environments/test_2",
	}

	aa := map[string]struct {
		diff_paths []string
		expected   []string
	}{
		"関係ないファイル": {
			[]string{"../internal/moddep/dep.go"},
			[]string{},
		},
		"executed_dirsの差分": {
			[]string{"../testdata/terraform/environments/test_1/main.tf"},
			[]string{"../testdata/terraform/environments/test_1"},
		},
		"moduleの差分": {
			[]string{"../testdata/terraform/modules/module_1/main.tf"},
			[]string{"../testdata/terraform/environments/test_1", "../testdata/terraform/environments/test_2"},
		},
		"executed_dirs, moduleの混在": {
			[]string{"../testdata/terraform/environments/test_1/main.tf", "../testdata/terraform/modules/module_1/main.tf"},
			[]string{"../testdata/terraform/environments/test_1", "../testdata/terraform/environments/test_2"},
		},
	}

	for name, a := range aa {
		t.Run(name, func(t *testing.T) {
			dirs, err := getDirs(a.diff_paths, executed_dirs)

			if err != nil {
				t.Fatal(err)
			}

			sort.SliceStable(dirs, func(i, j int) bool {
				return dirs[i] < dirs[j]
			})
			sort.SliceStable(a.expected, func(i, j int) bool {
				return a.expected[i] < a.expected[j]
			})

			if !reflect.DeepEqual(dirs, a.expected) {
				t.Fatalf("\nexpected: %#v\nactual: %#v", a.expected, dirs)
			}
		})
	}
}
