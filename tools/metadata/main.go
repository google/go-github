// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// metadata is a command-line tool used to check and update this repo.
// See CONTRIBUTING.md for details.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/google/go-github/v84/github"
)

var helpVars = kong.Vars{
	"update_openapi_help": `
Update openapi_operations.yaml from OpenAPI descriptions in github.com/github/rest-api-description at the given git ref.
`,

	"update_go_help": `
Update go source code to be consistent with openapi_operations.yaml.
 - Adds and updates "// GitHub API docs:" comments for service methods.
 - Updates "//meta:operation" comments to use canonical operation names.
 - Updates formatting of "//meta:operation" comments to make sure there isn't a space between the "//" and the "meta".
 - Formats modified files with the equivalent of "go fmt".
`,

	"format_help": `Format white space in openapi_operations.yaml and sort its operations.`,
	"unused_help": `List operations in openapi_operations.yaml that aren't used by any service methods.`,

	"working_dir_help": `Working directory. Should be the root of the go-github repository.`,
	"openapi_ref_help": `Git ref to pull OpenAPI descriptions from.`,

	"openapi_validate_help": `
Instead of updating, make sure that the operations in openapi_operations.yaml's "openapi_operations" field are
consistent with the SHA listed in "openapi_commit". This is run in CI as a convenience so that reviewers can trust
changes to openapi_operations.yaml.
`,

	"output_json_help": `Output JSON.`,
}

type rootCmd struct {
	UpdateOpenAPI updateOpenAPICmd `kong:"cmd,name=update-openapi,help=${update_openapi_help}"`
	UpdateGo      updateGoCmd      `kong:"cmd,help=${update_go_help}"`
	Format        formatCmd        `kong:"cmd,help=${format_help}"`
	Unused        unusedCmd        `kong:"cmd,help=${unused_help}"`

	WorkingDir string `kong:"short=C,default=.,help=${working_dir_help}"`

	// for testing
	GithubURL string `kong:"hidden,default='https://api.github.com'"`
}

func (c *rootCmd) opsFile() (string, *operationsFile, error) {
	filename := filepath.Join(c.WorkingDir, "openapi_operations.yaml")
	opsFile, err := loadOperationsFile(filename)
	if err != nil {
		return "", nil, err
	}
	return filename, opsFile, nil
}

func githubClient(apiURL string) (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, errors.New("GITHUB_TOKEN environment variable must be set to a GitHub personal access token with the public_repo scope")
	}
	return github.NewClient(nil).WithAuthToken(token).WithEnterpriseURLs(apiURL, "")
}

type updateOpenAPICmd struct {
	Ref            string `kong:"default=main,help=${openapi_ref_help}"`
	ValidateGithub bool   `kong:"name=validate,help=${openapi_validate_help}"`
}

func (c *updateOpenAPICmd) Run(root *rootCmd) error {
	ctx := context.Background()
	if c.ValidateGithub && c.Ref != "main" {
		return errors.New("--validate and --ref are mutually exclusive")
	}
	filename, opsFile, err := root.opsFile()
	if err != nil {
		return err
	}
	origOps := make([]*operation, len(opsFile.OpenapiOps))
	copy(origOps, opsFile.OpenapiOps)
	for i := range origOps {
		origOps[i] = origOps[i].clone()
	}
	client, err := githubClient(root.GithubURL)
	if err != nil {
		return err
	}
	ref := c.Ref
	if c.ValidateGithub {
		ref = opsFile.GitCommit
		if ref == "" {
			return errors.New("openapi_operations.yaml does not have an openapi_commit field")
		}
	}
	err = opsFile.updateFromGithub(ctx, client, ref)
	if err != nil {
		return err
	}
	if !c.ValidateGithub {
		return opsFile.saveFile(filename)
	}
	if !operationsEqual(origOps, opsFile.OpenapiOps) {
		return errors.New("openapi_operations.yaml does not match the OpenAPI descriptions in github.com/github/rest-api-description")
	}
	return nil
}

type formatCmd struct{}

func (c *formatCmd) Run(root *rootCmd) error {
	filename, opsFile, err := root.opsFile()
	if err != nil {
		return err
	}
	return opsFile.saveFile(filename)
}

type updateGoCmd struct{}

func (c *updateGoCmd) Run(root *rootCmd) error {
	_, opsFile, err := root.opsFile()
	if err != nil {
		return err
	}
	err = updateDocs(opsFile, filepath.Join(root.WorkingDir, "github"))
	return err
}

type unusedCmd struct {
	JSON bool `kong:"help=${output_json_help}"`
}

func (c *unusedCmd) Run(root *rootCmd, k *kong.Context) error {
	_, opsFile, err := root.opsFile()
	if err != nil {
		return err
	}
	unused, err := unusedOps(opsFile, filepath.Join(root.WorkingDir, "github"))
	if err != nil {
		return err
	}
	if c.JSON {
		enc := json.NewEncoder(k.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(unused)
	}
	fmt.Fprintf(k.Stdout, "Found %v unused operations\n", len(unused))
	if len(unused) == 0 {
		return nil
	}
	fmt.Fprintln(k.Stdout, "")
	for _, op := range unused {
		fmt.Fprintln(k.Stdout, op.Name)
		if op.DocumentationURL != "" {
			fmt.Fprintf(k.Stdout, "doc:     %v\n", op.DocumentationURL)
		}
		fmt.Fprintln(k.Stdout, "")
	}
	return nil
}

func main() {
	err := run(os.Args[1:], nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, opts []kong.Option) error {
	var cmd rootCmd
	parser, err := kong.New(&cmd, append(opts, helpVars)...)
	if err != nil {
		return err
	}
	k, err := parser.Parse(args)
	if err != nil {
		return err
	}
	return k.Run()
}
