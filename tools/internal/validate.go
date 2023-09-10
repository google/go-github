// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

// ValidateMetadata returns a list of issues with the metadata file. An error means
// there was an error validating the file, not that there are issues with the file.
func ValidateMetadata(dir string, meta *Metadata) ([]string, error) {
	undocumentedMethods, err := getUndocumentedMethods(dir, meta)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, m := range undocumentedMethods {
		msg := fmt.Sprintf("Undocumented method %s. Please add it to metadata.", m.name())
		result = append(result, msg)
	}
	metaMethodMap := map[string]bool{}
	for _, m := range meta.UndocumentedMethods {
		metaMethodMap[m] = true
	}
	methodOperations := map[string][]*Operation{}
	for _, op := range meta.Operations {
		for _, m := range op.GoMethods {
			metaMethodMap[m] = true
			methodOperations[m] = append(methodOperations[m], op)
		}
	}
	var metaMethods []string
	for m := range metaMethodMap {
		metaMethods = append(metaMethods, m)
	}
	sort.Strings(metaMethods)
	missing, err := missingMethods(dir, metaMethods)
	if err != nil {
		return nil, err
	}
	for _, m := range missing {
		msg := fmt.Sprintf("Method %s in metadata does not exist in github package.", m)
		result = append(result, msg)
	}
	for _, m := range meta.UndocumentedMethods {
		if len(methodOperations[m]) == 0 {
			continue
		}
		var ops []string
		for _, op := range methodOperations[m] {
			ops = append(ops, fmt.Sprintf("%s %s", op.Method(), op.EndpointURL()))
		}
		msg := fmt.Sprintf("Method %s is listed as undocumented in metadata and also in these operations: %s", m, strings.Join(ops, ", "))
		result = append(result, msg)
	}
	return result, nil
}

// getUndocumentedMethods returns a list of methods that are not mapped to any operation in metadata.yaml
func getUndocumentedMethods(dir string, metadata *Metadata) ([]*serviceMethod, error) {
	var result []*serviceMethod
	methods, err := getServiceMethods(dir)
	if err != nil {
		return nil, err
	}
	for _, method := range methods {
		ops := metadata.operationsForMethod(method.name())
		if len(ops) > 0 {
			continue
		}
		if slices.Contains(metadata.UndocumentedMethods, method.name()) {
			continue
		}
		result = append(result, method)
	}
	return result, nil
}

// missingMethods returns the set from methods that do not exist in the github package.
func missingMethods(dir string, methods []string) ([]string, error) {
	var result []string
	existingMap := map[string]bool{}
	sm, err := getServiceMethods(dir)
	if err != nil {
		return nil, err
	}
	for _, method := range sm {
		existingMap[method.name()] = true
	}
	for _, m := range methods {
		if !existingMap[m] {
			result = append(result, m)
		}
	}
	return result, nil
}
