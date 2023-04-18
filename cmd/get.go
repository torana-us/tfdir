/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/torana-us/tfdir/internal/moddep"
	"github.com/torana-us/tfdir/internal/tree"
)

// testCmd represents the test command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the directory where terraform is run.",
	Long: `
Get the directory where terraform is run.
example:
	git diff --name-only | tfdir get
	`,
	RunE: getCmdRunner,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func getCmdRunner(cmd *cobra.Command, args []string) error {
	var config Config
	viper.Unmarshal(&config)

	dirs, err := getDirs(read_from_stdin(), config.ExecutedDirs)

	if err != nil {
		return err
	}

	for _, dir := range dirs {
		fmt.Println(dir)
	}

	return nil
}

func sliceUnique(target []string) []string {
	m := map[string]struct{}{}
	unique := []string{}

	for _, v := range target {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			unique = append(unique, v)
		}
	}

	return unique
}

func getDirs(diff_paths []string, executed_dirs []string) ([]string, error) {
	result := []string{}

	mod_map, diags := moddep.GetDependency(executed_dirs)
	tree_map := tree.MakeTreeMap(executed_dirs)

	if diags.HasErrors() {
		return result, diags
	}

	for _, d := range diff_paths {
		diff_dir := filepath.Dir(d)
		// from module
		dirs, ok := mod_map.GetRaw(diff_dir)
		if ok {
			result = append(result, dirs...)
		}

		// frm diff
		diff_dir_path := tree.NewPath(diff_dir)
		h, _ := diff_dir_path.HeadAndRest()

		tree, ok := tree_map[h]

		if !ok {
			continue
		}

		n, ok := tree.Search(*diff_dir_path)

		if !ok {
			continue
		}

		for _, p := range n.AllPathString() {
			result = append(result, filepath.Join(filepath.Dir(diff_dir), p))
		}
	}

	return sliceUnique(result), nil
}

func read_from_stdin() (list []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	return
}
