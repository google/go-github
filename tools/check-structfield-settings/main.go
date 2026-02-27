// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// check-structfield-settings reads the settings for
// the custom `structfield` linter in ".golangci.yml" -
// specifically, the "allowed-tag-names" and "allowed-tag-types"
// exceptions, then scans the code repo to find all exceptions
// that are no longer needed and reports a list.
package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"github.com/google/go-github/v84/tools/structfield"
	"go.yaml.in/yaml/v3"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/checker"
	"golang.org/x/tools/go/packages"
)

func init() {
	register.Plugin("structfield", structfield.New)
}

func main() {
	log.SetFlags(0)
	configPath := flag.String("config", "", "path to .golangci.yml (defaults to searching up from cwd)")
	packagesFlag := flag.String("packages", "./...", "comma-separated list of package patterns to analyze")
	includeTests := flag.Bool("tests", false, "include test files in analysis")
	fix := flag.Bool("fix", false, "remove obsolete exceptions and sort/dedupe lists in .golangci.yml")
	flag.Parse()

	resolvedConfig, repoRoot, err := resolveConfig(*configPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	allowedNamesList, allowedTypesList, allowedNames, allowedTypes, err := readStructfieldSettings(resolvedConfig)
	if err != nil {
		log.Fatalf("parse config: %v", err)
	}
	if len(allowedNames) == 0 && len(allowedTypes) == 0 {
		log.Fatalf("no structfield settings found in %s", resolvedConfig)
	}

	duplicateNames := findDuplicates(allowedNamesList)
	duplicateTypes := findDuplicates(allowedTypesList)

	patterns := strings.Split(*packagesFlag, ",")
	for i, pattern := range patterns {
		patterns[i] = strings.TrimSpace(pattern)
	}

	usedNames, usedTypes, err := analyzeRepo(repoRoot, patterns, *includeTests, allowedNames, allowedTypes)
	if err != nil {
		log.Fatalf("analyze: %v", err)
	}

	obsoleteNames := diffKeys(allowedNames, usedNames)
	obsoleteTypes := diffKeys(allowedTypes, usedTypes)

	if len(obsoleteNames) == 0 && len(obsoleteTypes) == 0 && len(duplicateNames) == 0 && len(duplicateTypes) == 0 {
		return
	}

	if *fix {
		if err := removeObsoleteExceptions(resolvedConfig, obsoleteNames, obsoleteTypes); err != nil {
			log.Fatalf("fix: %v", err)
		}
		return
	}

	if len(obsoleteNames) > 0 {
		fmt.Println("Obsolete allowed-tag-names:")
		for _, name := range obsoleteNames {
			fmt.Printf("  - %v\n", name)
		}
	}
	if len(obsoleteTypes) > 0 {
		if len(obsoleteNames) > 0 {
			fmt.Println()
		}
		fmt.Println("Obsolete allowed-tag-types:")
		for _, name := range obsoleteTypes {
			fmt.Printf("  - %v\n", name)
		}
	}
	if len(duplicateNames) > 0 {
		if len(obsoleteNames) > 0 || len(obsoleteTypes) > 0 {
			fmt.Println()
		}
		fmt.Println("Duplicate allowed-tag-names:")
		for _, name := range sortedKeys(duplicateNames) {
			fmt.Printf("  - %v (%v)\n", name, duplicateNames[name])
		}
	}
	if len(duplicateTypes) > 0 {
		if len(obsoleteNames) > 0 || len(obsoleteTypes) > 0 || len(duplicateNames) > 0 {
			fmt.Println()
		}
		fmt.Println("Duplicate allowed-tag-types:")
		for _, name := range sortedKeys(duplicateTypes) {
			fmt.Printf("  - %v (%v)\n", name, duplicateTypes[name])
		}
	}
}

type golangciConfig struct {
	Linters struct {
		Settings struct {
			Custom struct {
				Structfield struct {
					Settings struct {
						AllowedTagNames []string `yaml:"allowed-tag-names"`
						AllowedTagTypes []string `yaml:"allowed-tag-types"`
					} `yaml:"settings"`
				} `yaml:"structfield"`
			} `yaml:"custom"`
		} `yaml:"settings"`
	} `yaml:"linters"`
}

func resolveConfig(configPath string) (string, string, error) {
	if configPath != "" {
		resolved, err := filepath.Abs(configPath)
		if err != nil {
			return "", "", err
		}
		repoRoot := filepath.Dir(resolved)
		return resolved, repoRoot, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	for dir := cwd; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, ".golangci.yml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, dir, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return "", "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
	}

	return "", "", errors.New("unable to locate .golangci.yml")
}

func readStructfieldSettings(configPath string) ([]string, []string, map[string]bool, map[string]bool, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var cfg golangciConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, nil, nil, nil, err
	}

	allowedNamesList := cfg.Linters.Settings.Custom.Structfield.Settings.AllowedTagNames
	allowedTypesList := cfg.Linters.Settings.Custom.Structfield.Settings.AllowedTagTypes

	allowedNames := make(map[string]bool)
	for _, name := range allowedNamesList {
		allowedNames[name] = true
	}

	allowedTypes := make(map[string]bool)
	for _, name := range allowedTypesList {
		allowedTypes[name] = true
	}

	return allowedNamesList, allowedTypesList, allowedNames, allowedTypes, nil
}

func analyzeRepo(repoRoot string, patterns []string, includeTests bool, allowedNames, allowedTypes map[string]bool) (map[string]bool, map[string]bool, error) {
	plugin, err := structfield.New(nil)
	if err != nil {
		return nil, nil, err
	}
	created, err := plugin.BuildAnalyzers()
	if err != nil {
		return nil, nil, err
	}
	if len(created) == 0 {
		return nil, nil, errors.New("no analyzers returned by structfield")
	}
	analyzer := created[0]

	cfg := &packages.Config{
		Mode:  packages.LoadAllSyntax,
		Dir:   repoRoot,
		Tests: includeTests,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, nil, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return nil, nil, errors.New("package load errors")
	}

	graph, err := checker.Analyze([]*analysis.Analyzer{analyzer}, pkgs, &checker.Options{Sequential: true})
	if err != nil {
		return nil, nil, err
	}

	usedNames := make(map[string]bool)
	usedTypes := make(map[string]bool)

	for act := range graph.All() {
		if !act.IsRoot || act.Analyzer != analyzer {
			continue
		}
		for _, diag := range act.Diagnostics {
			markUsedException(diag.Message, allowedNames, allowedTypes, usedNames, usedTypes)
		}
	}

	return usedNames, usedTypes, nil
}

var (
	nameMismatchRE  = regexp.MustCompile(`^change Go field name "([^"]+)" to ".*" for .* tag ".*" in struct "([^"]+)"$`)
	typeChangeRE    = regexp.MustCompile(`^change the "([^"]+)" field type to ".*" in the struct "([^"]+)"`)
	fieldInStructRE = regexp.MustCompile(`^the "([^"]+)" field in struct "([^"]+)" .*`)
)

func markUsedException(msg string, allowedNames, allowedTypes, usedNames, usedTypes map[string]bool) {
	if match := nameMismatchRE.FindStringSubmatch(msg); match != nil {
		key := match[2] + "." + match[1]
		if allowedNames[key] {
			usedNames[key] = true
		}
		return
	}

	if match := typeChangeRE.FindStringSubmatch(msg); match != nil {
		key := match[2] + "." + match[1]
		if allowedTypes[key] {
			usedTypes[key] = true
		}
		return
	}

	if match := fieldInStructRE.FindStringSubmatch(msg); match != nil {
		key := match[2] + "." + match[1]
		if allowedTypes[key] {
			usedTypes[key] = true
		}
	}
}

func diffKeys(all, used map[string]bool) []string {
	var obsolete []string
	for key := range all {
		if !used[key] {
			obsolete = append(obsolete, key)
		}
	}
	return slices.Sorted(slices.Values(obsolete))
}

func findDuplicates(values []string) map[string]int {
	counts := make(map[string]int)
	duplicates := make(map[string]int)
	for _, value := range values {
		counts[value]++
		if counts[value] > 1 {
			duplicates[value] = counts[value]
		}
	}
	return duplicates
}

func sortedKeys(values map[string]int) []string {
	return slices.Sorted(maps.Keys(values))
}

func removeObsoleteExceptions(configPath string, obsoleteNames, obsoleteTypes []string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	fileInfo, statErr := os.Stat(configPath)
	fileMode := os.FileMode(0o644)
	if statErr == nil {
		fileMode = fileInfo.Mode()
	}

	hasTrailingNewline := strings.HasSuffix(string(data), "\n")
	lines := strings.Split(string(data), "\n")

	obsoleteNamesSet := make(map[string]bool, len(obsoleteNames))
	for _, name := range obsoleteNames {
		obsoleteNamesSet[name] = true
	}
	obsoleteTypesSet := make(map[string]bool, len(obsoleteTypes))
	for _, name := range obsoleteTypes {
		obsoleteTypesSet[name] = true
	}
	seenNames := make(map[string]bool)
	seenTypes := make(map[string]bool)
	var items []*listItem

	updated := make([]string, 0, len(lines))
	section := ""

	for i := 0; i < len(lines); {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if section == "" {
			switch {
			case strings.HasPrefix(trimmed, "allowed-tag-names:"):
				section = "names"
				items = items[:0]
				updated = append(updated, line)
				i++
				continue
			case strings.HasPrefix(trimmed, "allowed-tag-types:"):
				section = "types"
				items = items[:0]
				updated = append(updated, line)
				i++
				continue
			default:
				updated = append(updated, line)
				i++
				continue
			}
		}

		if strings.HasPrefix(trimmed, "- ") {
			value := strings.TrimSpace(trimmed[2:])
			if hash := strings.Index(value, "#"); hash >= 0 {
				value = strings.TrimSpace(value[:hash])
			}
			if section == "names" {
				if obsoleteNamesSet[value] || seenNames[value] {
					i++
					continue
				}
				seenNames[value] = true
			}
			if section == "types" {
				if obsoleteTypesSet[value] || seenTypes[value] {
					i++
					continue
				}
				seenTypes[value] = true
			}
			items = append(items, &listItem{value: value, line: line})
			i++
			continue
		}

		updated = appendSortedItems(updated, items)
		section = ""
		items = items[:0]
		continue
	}

	if section != "" {
		updated = appendSortedItems(updated, items)
	}

	content := strings.Join(updated, "\n")
	if hasTrailingNewline && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	return os.WriteFile(configPath, []byte(content), fileMode)
}

type listItem struct {
	value string
	line  string
}

func appendSortedItems(lines []string, items []*listItem) []string {
	slices.SortFunc(items, func(a, b *listItem) int {
		return cmp.Compare(a.value, b.value)
	})
	for _, item := range items {
		lines = append(lines, item.line)
	}
	return lines
}
