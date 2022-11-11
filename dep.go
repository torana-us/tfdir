package main

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-config-inspect/tfconfig"
)

type ModuleMap map[Module][]ExecutedDir

type Module string
type ExecutedDir string

func (m ModuleMap) GetRaw(key string) ([]string, bool) {
	dirs, ok := m[Module(key)]
	raws := make([]string, len(dirs))

	if !ok {
		return raws, ok
	}

	for i, v := range dirs {
		raws[i] = string(v)
	}

	return raws, ok
}

func isLocalModule(raw string) bool {
	prefix_list := []string{
		"./",
		"../",
		".\\",
		"..\\",
	}

	for _, p := range prefix_list {
		if strings.HasPrefix(raw, p) {
			return true
		}
	}

	return false
}

func GetDependency(dirs []string) (mod_map ModuleMap, diags tfconfig.Diagnostics) {
	mod_map = ModuleMap{}
	for _, dir := range dirs {
		mod, d := tfconfig.LoadModule(string(dir))
		if d.HasErrors() {
			diags = append(diags, d...)
		}

		d = getDependencyChildren(ExecutedDir(dir), mod_map, mod)
		if d.HasErrors() {
			diags = append(diags, d...)
		}
	}
	return
}

func getDependencyChildren(dir ExecutedDir, mod_map ModuleMap, mod *tfconfig.Module) (diags tfconfig.Diagnostics) {
	parent_path := mod.Path
	for _, m := range mod.ModuleCalls {
		path := m.Source
		if !isLocalModule(path) {
			continue
		}
		key := Module(filepath.Join(parent_path, path))
		mod_map[key] = append(mod_map[key], dir)

		child_mod, d := tfconfig.LoadModule(
			filepath.Join(parent_path, path),
		)

		if d.HasErrors() {
			diags = append(diags, d...)
		}

		d = getDependencyChildren(dir, mod_map, child_mod)

		if d.HasErrors() {
			diags = append(diags, d...)
		}
	}

	return
}
