// Copyright 2026 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	rawBaseURL        = "https://raw.githubusercontent.com"
	descriptionsOwner = "github"
	descriptionsRepo  = "rest-api-description"
	operationsFile    = "openapi_operations.yaml"
)

// schema is the subset of the OpenAPI 3.0 schema object this tool reads.
type schema struct {
	Ref        string             `json:"$ref"`
	Nullable   bool               `json:"nullable"`
	ReadOnly   bool               `json:"readOnly"`
	WriteOnly  bool               `json:"writeOnly"`
	Required   []string           `json:"required"`
	Properties map[string]*schema `json:"properties"`
	AllOf      []*schema          `json:"allOf"`
	OneOf      []*schema          `json:"oneOf"`
	AnyOf      []*schema          `json:"anyOf"`
}

// operation is the subset of an OpenAPI operation object this tool reads.
type operation struct {
	RequestBody *struct {
		Content map[string]struct {
			Schema *schema `json:"schema"`
		} `json:"content"`
	} `json:"requestBody"`
}

// description is one OpenAPI description file, such as descriptions/ghec/ghec.json.
type description struct {
	Paths      map[string]map[string]json.RawMessage `json:"paths"`
	Components struct {
		Schemas map[string]*schema `json:"schemas"`
	} `json:"components"`

	// The lookup caches below are filled on demand, so they are guarded: a
	// description is shared, and checking it from more than one goroutine is a
	// natural thing to do with 36 MB of schemas.
	mu      sync.Mutex
	opCache map[string]*flat
	opMiss  map[string]bool
}

func newDescription() *description {
	return &description{opCache: map[string]*flat{}, opMiss: map[string]bool{}}
}

// loadDescriptionFile reads one OpenAPI description file.
func loadDescriptionFile(path string) (*description, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d := newDescription()
	if err := json.NewDecoder(f).Decode(d); err != nil {
		return nil, fmt.Errorf("%v: %w", path, err)
	}
	return d, nil
}

// resolve follows $ref chains into components/schemas.
func (d *description) resolve(s *schema) (*schema, error) {
	seen := map[string]bool{}
	for s != nil && s.Ref != "" {
		if seen[s.Ref] {
			return nil, fmt.Errorf("circular $ref %v", s.Ref)
		}
		seen[s.Ref] = true
		name := strings.TrimPrefix(s.Ref, "#/components/schemas/")
		next, ok := d.Components.Schemas[name]
		if !ok {
			return nil, fmt.Errorf("unresolved $ref %v", s.Ref)
		}
		s = next
	}
	return s, nil
}

// propMeta is what a property says about itself, independent of any one operation.
type propMeta struct {
	nullable  bool
	readOnly  bool
	writeOnly bool
}

// flat is an object schema flattened for field-optionality comparison.
type flat struct {
	required map[string]bool
	props    map[string]propMeta
}

func newFlat() *flat {
	return &flat{required: map[string]bool{}, props: map[string]propMeta{}}
}

// merge adds the required properties and properties of other to f.
func (f *flat) merge(other *flat) {
	for name := range other.required {
		f.required[name] = true
	}
	f.mergeProps(other)
}

// mergeProps adds only the properties of other to f.
func (f *flat) mergeProps(other *flat) {
	maps.Copy(f.props, other.props)
}

// flatten resolves refs and merges allOf. For oneOf and anyOf only the properties
// required by EVERY variant stay required: a property required in all variants is
// unconditionally required, while one required in only some variants is conditionally
// required and must be left to author judgement. That makes "issue_field_id OR
// name+data_type" style bodies checkable rather than an error.
func (d *description) flatten(s *schema, visiting map[string]bool) (*flat, error) {
	resolved, err := d.resolve(s)
	if err != nil {
		return nil, err
	}
	if resolved == nil {
		return newFlat(), nil
	}
	out := newFlat()

	variants := resolved.OneOf
	if len(variants) == 0 {
		variants = resolved.AnyOf
	}
	if len(variants) > 0 {
		var reqSets []map[string]bool
		for _, v := range variants {
			f, err := d.flatten(v, visiting)
			if err != nil {
				return nil, err
			}
			// Each variant states its own required properties, so only the
			// intersection below may promote one to required. Unioning them here
			// would make every alternative look mandatory.
			reqSets = append(reqSets, f.required)
			out.mergeProps(f)
		}
		for name := range out.props {
			inEveryVariant := true
			for _, required := range reqSets {
				if !required[name] {
					inEveryVariant = false
					break
				}
			}
			if inEveryVariant {
				out.required[name] = true
			}
		}
	}

	for _, r := range resolved.Required {
		out.required[r] = true
	}
	for name, ps := range resolved.Properties {
		psr, err := d.resolve(ps)
		if err != nil {
			return nil, err
		}
		meta := propMeta{}
		if psr != nil {
			meta = propMeta{nullable: psr.Nullable, readOnly: psr.ReadOnly, writeOnly: psr.WriteOnly}
		}
		out.props[name] = meta
	}

	for _, a := range resolved.AllOf {
		if a == nil {
			continue
		}
		if key := a.Ref; key != "" {
			if visiting[key] {
				continue
			}
			visiting[key] = true
		}
		f, err := d.flatten(a, visiting)
		if err != nil {
			return nil, err
		}
		out.merge(f)
	}
	return out, nil
}

// requestSchema returns the application/json request body schema for an operation. ok is
// false when the operation has no JSON object request body.
func (d *description) requestSchema(op *opRef) (*flat, bool, error) {
	// The caches are written below, so the lookup holds the lock throughout.
	// requestSchema does not call itself, so a plain mutex cannot deadlock.
	d.mu.Lock()
	defer d.mu.Unlock()

	key := op.String()
	if f, ok := d.opCache[key]; ok {
		return f, true, nil
	}
	if d.opMiss[key] {
		return nil, false, nil
	}
	item, ok := d.Paths[op.path]
	if !ok {
		d.opMiss[key] = true
		return nil, false, nil
	}
	raw, ok := item[strings.ToLower(op.method)]
	if !ok {
		d.opMiss[key] = true
		return nil, false, nil
	}
	var o operation
	if err := json.Unmarshal(raw, &o); err != nil {
		return nil, false, err
	}
	if o.RequestBody == nil {
		d.opMiss[key] = true
		return nil, false, nil
	}
	mt, ok := o.RequestBody.Content["application/json"]
	if !ok || mt.Schema == nil {
		d.opMiss[key] = true
		return nil, false, nil
	}
	f, err := d.flatten(mt.Schema, map[string]bool{})
	if err != nil {
		return nil, false, err
	}
	if len(f.props) == 0 {
		d.opMiss[key] = true
		return nil, false, nil
	}
	d.opCache[key] = f
	return f, true, nil
}

// descriptions holds the OpenAPI plans in load order (api.github.com, then ghec, then
// ghes), mirroring how the metadata tool resolves operations across plans.
type descriptions struct{ files []*description }

// requestSchema returns the request body schema from the first plan that documents the
// operation, so GHEC-only and GHES-only operations resolve as well.
func (ds *descriptions) requestSchema(op *opRef) (*flat, bool, error) {
	for _, d := range ds.files {
		f, ok, err := d.requestSchema(op)
		if err != nil {
			return nil, false, err
		}
		if ok {
			return f, true, nil
		}
	}
	return nil, false, nil
}

// planRe matches the description paths stored in openapi_operations.yaml.
var planRe = regexp.MustCompile(`^descriptions/([^/]+)/([^/]+)\.json$`)

// pinnedDescriptions returns the git ref and the description paths that the repository
// pins in openapi_operations.yaml.
func pinnedDescriptions(root string) (string, []string, error) {
	path := filepath.Join(root, operationsFile)
	f, err := os.Open(path)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	var ref string
	var files []string
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "openapi_commit:"):
			ref = strings.TrimSpace(strings.TrimPrefix(line, "openapi_commit:"))
		case strings.HasPrefix(line, "- descriptions/"):
			files = append(files, strings.TrimPrefix(line, "- "))
		}
	}
	if err := scanner.Err(); err != nil {
		return "", nil, err
	}
	if ref == "" {
		return "", nil, fmt.Errorf("%v: no openapi_commit field", path)
	}
	return ref, selectPlans(files), nil
}

// selectPlans picks the plans to check against: api.github.com, then ghec, then the newest
// GHES release. It is the same preference order the metadata tool uses when it resolves an
// operation, and it keeps the download to three files instead of all ten.
func selectPlans(files []string) []string {
	var api, ghec, ghes string
	ghesMajor, ghesMinor := -1, -1
	for _, path := range slices.Sorted(slices.Values(files)) { // Also removes duplicates.
		m := planRe.FindStringSubmatch(path)
		if m == nil {
			continue
		}
		switch dir := m[1]; {
		case dir == "api.github.com":
			api = path
		case dir == "ghec":
			ghec = path
		case strings.HasPrefix(dir, "ghes-"):
			major, minor, err := parseVersion(strings.TrimPrefix(dir, "ghes-"))
			if err != nil || major < 3 {
				continue
			}
			if major > ghesMajor || (major == ghesMajor && minor > ghesMinor) {
				ghesMajor, ghesMinor, ghes = major, minor, path
			}
		}
	}
	var out []string
	for _, plan := range []string{api, ghec, ghes} {
		if plan != "" {
			out = append(out, plan)
		}
	}
	return out
}

// parseVersion parses a "major.minor" release version.
func parseVersion(s string) (major, minor int, err error) {
	rawMajor, rawMinor, ok := strings.Cut(s, ".")
	if !ok {
		return 0, 0, fmt.Errorf("malformed version %q", s)
	}
	major, err = strconv.Atoi(rawMajor)
	if err != nil {
		return 0, 0, err
	}
	minor, err = strconv.Atoi(rawMinor)
	return major, minor, err
}

// fetchDescriptions downloads the given description paths at the given ref, caching them
// in cacheDir so that repeated runs only download each file once. It returns the local
// path of each file, in the order given.
func fetchDescriptions(ctx context.Context, ref string, paths []string, cacheDir string) ([]string, error) {
	dir := filepath.Join(cacheDir, ref)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 5 * time.Minute}
	local := make([]string, 0, len(paths))
	for _, path := range paths {
		dest := filepath.Join(dir, strings.ReplaceAll(path, "/", "_"))
		local = append(local, dest)
		if info, err := os.Stat(dest); err == nil && info.Size() > 0 {
			continue
		}
		if err := download(ctx, client, rawBaseURL, ref, path, dest); err != nil {
			return nil, err
		}
	}
	return local, nil
}

// download fetches one description file from baseURL and writes it to dest.
func download(ctx context.Context, client *http.Client, baseURL, ref, path, dest string) error {
	url := baseURL + "/" + descriptionsOwner + "/" + descriptionsRepo + "/" + ref + "/" + path
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v: %v", url, resp.Status)
	}
	f, err := os.CreateTemp(filepath.Dir(dest), filepath.Base(dest)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), dest)
}

// defaultCacheDir returns the directory used to cache downloaded descriptions.
func defaultCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return filepath.Join(os.TempDir(), "go-github-schemafields")
	}
	return filepath.Join(dir, "go-github-schemafields")
}

// loadDescriptions loads the given description files, or the pinned ones when no files are
// named.
func loadDescriptions(ctx context.Context, root string, files []string, cacheDir string) (*descriptions, error) {
	if len(files) == 0 {
		ref, paths, err := pinnedDescriptions(root)
		if err != nil {
			return nil, err
		}
		if len(paths) == 0 {
			return nil, errors.New("no OpenAPI description files found in " + operationsFile)
		}
		files, err = fetchDescriptions(ctx, ref, paths, cacheDir)
		if err != nil {
			return nil, err
		}
	}
	ds := &descriptions{}
	for _, path := range files {
		d, err := loadDescriptionFile(path)
		if err != nil {
			return nil, err
		}
		ds.files = append(ds.files, d)
	}
	return ds, nil
}
