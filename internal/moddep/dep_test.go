package moddep_test

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"testing"

	. "github.com/torana-us/tfdir/internal/moddep"
)

func Test_GetDependency(t *testing.T) {
	executed_dirs := []string{
		"../../testdata/terraform/environments/test_1/",
		"../../testdata/terraform/environments/test_2/",
	}
	expected := map[string][]string{
		"../../testdata/terraform/modules/module_1": {
			"../../testdata/terraform/environments/test_2/",
			"../../testdata/terraform/environments/test_1/",
		},
		"../../testdata/terraform/modules/module_2": {
			"../../testdata/terraform/environments/test_1/",
		},
		"../../testdata/terraform/modules/module_2/modules/module_2_1": {
			"../../testdata/terraform/environments/test_1/",
		},
	}

	mod_map, d := GetDependency(executed_dirs)

	if d.HasErrors() {
		fmt.Println(expected)
		t.Fatal(d)
	}

	for k, v := range expected {
		executed, ok := mod_map.GetRaw(k)

		if !ok {
			log.Fatalf("%s is not found", k)
		}

		sort.SliceStable(executed, func(i, j int) bool {
			return executed[i] < executed[j]
		})
		sort.SliceStable(v, func(i, j int) bool {
			return v[i] < v[j]
		})

		if !reflect.DeepEqual(executed, v) {
			t.Fatalf("key: %s\n%v : %v", k, executed, v)
		}
	}
}
