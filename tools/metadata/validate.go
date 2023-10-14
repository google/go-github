// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/go-github/v56/github"
)

// validateMetadata returns a list of issues with the metadata file. An error means
// there was an error validating the file, not that there are issues with the file.
//
// validations:
//   - Methods in the github package must exist in metadata.yaml
//   - Methods in metadata.yaml must exist in github package
//   - Methods in metadata.yaml must have unique names
//   - Methods in metadata.yaml must have at least one operation
//   - Methods in metadata.yaml may not have duplicate operations
//   - Methods in metadata.yaml must use the canonical operation name
//   - All operations mapped from a method must exist in either ManualOps or OpenAPIOps
//   - No operations are duplicated between ManualOps and OpenAPIOps
//   - All operations in OverrideOps must exist in either ManualOps or OpenAPIOps
func validateMetadata(dir string, meta *metadata) ([]string, error) {
	serviceMethods, err := getServiceMethods(dir)
	if err != nil {
		return nil, err
	}
	var result []string
	result = validateServiceMethodsExist(result, meta, serviceMethods)
	result = validateMetadataMethods(result, meta, serviceMethods)
	result = validateOperations(result, meta)
	return result, nil
}

// validateGitCommit validates that building meta.OpenapiOps from the commit at meta.GitCommit
// results in the same operations as meta.OpenapiOps.
func validateGitCommit(ctx context.Context, client *github.Client, meta *metadata) (string, error) {
	ops, err := getOpsFromGithub(ctx, client, meta.GitCommit)
	if err != nil {
		return "", err
	}
	if !operationsEqual(ops, meta.OpenapiOps) {
		msg := fmt.Sprintf("openapi_operations does not match operations from git commit %s", meta.GitCommit)
		return msg, nil
	}
	return "", nil
}

func validateMetadataMethods(result []string, meta *metadata, serviceMethods []string) []string {
	smLookup := map[string]bool{}
	for _, method := range serviceMethods {
		smLookup[method] = true
	}
	seenMethods := map[string]bool{}
	for _, method := range meta.Methods {
		if seenMethods[method.Name] {
			msg := fmt.Sprintf("Method %s is duplicated in metadata.yaml.", method.Name)
			result = append(result, msg)
			continue
		}
		seenMethods[method.Name] = true
		if !smLookup[method.Name] {
			msg := fmt.Sprintf("Method %s in metadata.yaml does not exist in github package.", method.Name)
			result = append(result, msg)
		}
		result = validateMetaMethodOperations(result, meta, method)
	}
	return result
}

func validateMetaMethodOperations(result []string, meta *metadata, method *method) []string {
	if len(method.OpNames) == 0 {
		msg := fmt.Sprintf("Method %s in metadata.yaml does not have any operations.", method.Name)
		result = append(result, msg)
	}
	seenOps := map[string]bool{}
	for _, opName := range method.OpNames {
		if seenOps[opName] {
			msg := fmt.Sprintf("Method %s in metadata.yaml has duplicate operation: %s.", method.Name, opName)
			result = append(result, msg)
		}
		seenOps[opName] = true
		if meta.getOperation(opName) != nil {
			continue
		}
		normalizedMatch := meta.getOperationsWithNormalizedName(opName)
		if len(normalizedMatch) > 0 {
			msg := fmt.Sprintf("Method %s has operation which is does not use the canonical name. You may be able to automatically fix this by running 'script/metadata.sh canonize': %s.", method.Name, opName)
			result = append(result, msg)
			continue
		}
		msg := fmt.Sprintf("Method %s has operation which is not defined in metadata.yaml: %s.", method.Name, opName)
		result = append(result, msg)
	}
	return result
}

func validateServiceMethodsExist(result []string, meta *metadata, serviceMethods []string) []string {
	for _, method := range serviceMethods {
		if meta.getMethod(method) == nil {
			msg := fmt.Sprintf("Method %s does not exist in metadata.yaml. Please add it.", method)
			result = append(result, msg)
		}
	}
	return result
}

func validateOperations(result []string, meta *metadata) []string {
	names := map[string]bool{}
	openapiNames := map[string]bool{}
	overrideNames := map[string]bool{}
	for _, op := range meta.OpenapiOps {
		if openapiNames[op.Name] {
			msg := fmt.Sprintf("Name duplicated in openapi_operations: %s", op.Name)
			result = append(result, msg)
		}
		openapiNames[op.Name] = true
	}
	for _, op := range meta.ManualOps {
		if names[op.Name] {
			msg := fmt.Sprintf("Name duplicated in operations: %s", op.Name)
			result = append(result, msg)
		}
		names[op.Name] = true
		if openapiNames[op.Name] {
			msg := fmt.Sprintf("Name exists in both operations and openapi_operations: %s", op.Name)
			result = append(result, msg)
		}
	}
	for _, op := range meta.OverrideOps {
		if overrideNames[op.Name] {
			msg := fmt.Sprintf("Name duplicated in override_operations: %s", op.Name)
			result = append(result, msg)
		}
		overrideNames[op.Name] = true
		if !names[op.Name] && !openapiNames[op.Name] {
			msg := fmt.Sprintf("Name in override_operations does not exist in operations or openapi_operations: %s", op.Name)
			result = append(result, msg)
		}
	}
	return result
}

func getServiceMethods(dir string) ([]string, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var serviceMethods []string
	for _, filename := range dirEntries {
		var sm []string
		sm, err = getServiceMethodsFromFile(filepath.Join(dir, filename.Name()))
		if err != nil {
			return nil, err
		}
		serviceMethods = append(serviceMethods, sm...)
	}
	sort.Strings(serviceMethods)
	return serviceMethods, nil
}

// getServiceMethodsFromFile returns the service methods in filename.
func getServiceMethodsFromFile(filename string) ([]string, error) {
	if !strings.HasSuffix(filename, ".go") ||
		strings.HasSuffix(filename, "_test.go") {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Only look at the github package
	if f.Name.Name != "github" {
		return nil, nil
	}
	var serviceMethods []string
	ast.Inspect(f, func(n ast.Node) bool {
		sm := serviceMethodFromNode(n)
		if sm == "" {
			return true
		}
		serviceMethods = append(serviceMethods, sm)
		return false
	})
	return serviceMethods, nil
}
