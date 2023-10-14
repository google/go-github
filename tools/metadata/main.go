// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/google/go-github/v56/github"
)

var helpVars = kong.Vars{
	"update_help":     `Update metadata.yaml from OpenAPI descriptions in github.com/github/rest-api-description.`,
	"format_help":     `Format metadata.yaml.`,
	"validate_help":   `Validate that metadata.yaml is consistent with source code.`,
	"unused_ops_help": `List operations in metadata.yaml that don't have any associated go methods'.`,
	"canonize_help":   `Update metadata.yaml to use canonical operation names.`,
}

type rootCmd struct {
	WorkingDir     string            `kong:"short=C,default=.,help='Working directory. Must be within a go-github root.'"`
	Filename       string            `kong:"help='Path to metadata.yaml. Defaults to <go-github-root>/metadata.yaml.'"`
	GithubDir      string            `kong:"help='Path to the github package. Defaults to <go-github-root>/github.'"`
	GithubURL      string            `kong:"hidden,default='https://api.github.com'"`
	UpdateMetadata updateMetadataCmd `kong:"cmd,help=${update_help}"`
	UpdateUrls     updateUrlsCmd     `kong:"cmd,help='Update documentation URLs in the Go source files in the github directory to match the urls in the metadata file.'"`
	Format         formatCmd         `kong:"cmd,help=${format_help}"`
	Validate       validateCmd       `kong:"cmd,help=${validate_help}"`
	UnusedOps      unusedOpsCmd      `kong:"cmd,help=${unused_ops_help}"`
	Canonize       canonizeCmd       `kong:"cmd,help=${canonize_help}"`
}

func (c *rootCmd) metadata() (string, *metadata, error) {
	filename := c.Filename
	if filename == "" {
		filename = filepath.Join(c.WorkingDir, "metadata.yaml")
	}
	meta, err := loadMetadataFile(filename)
	if err != nil {
		return "", nil, err
	}
	return filename, meta, nil
}

func githubClient(apiURL string) (*github.Client, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN environment variable must be set to a GitHub personal access token with the public_repo scope")
	}
	return github.NewClient(nil).WithAuthToken(token).WithEnterpriseURLs(apiURL, "")
}

type updateMetadataCmd struct {
	Ref string `kong:"default=main,help='git ref to pull OpenAPI descriptions from'"`
}

func (c *updateMetadataCmd) Run(root *rootCmd) error {
	ctx := context.Background()
	filename, meta, err := root.metadata()
	if err != nil {
		return err
	}
	client, err := githubClient(root.GithubURL)
	if err != nil {
		return err
	}

	err = meta.updateFromGithub(ctx, client, c.Ref)
	if err != nil {
		return err
	}
	return meta.saveFile(filename)
}

type formatCmd struct{}

func (c *formatCmd) Run(root *rootCmd) error {
	filename, meta, err := root.metadata()
	if err != nil {
		return err
	}
	return meta.saveFile(filename)
}

type validateCmd struct {
	CheckGithub bool `kong:"help='Check that metadata.yaml is consistent with the OpenAPI descriptions in github.com/github/rest-api-description.'"`
}

func (c *validateCmd) Run(k *kong.Context, root *rootCmd) error {
	ctx := context.Background()
	githubDir := filepath.Join(root.WorkingDir, "github")
	filename, meta, err := root.metadata()
	if err != nil {
		return err
	}
	issues, err := validateMetadata(githubDir, meta)
	if err != nil {
		return err
	}
	// don't waste time checking github if there are already issues
	if c.CheckGithub && len(issues) == 0 {
		client, err := githubClient(root.GithubURL)
		if err != nil {
			return err
		}
		msg, err := validateGitCommit(ctx, client, meta)
		if err != nil {
			return err
		}
		if msg != "" {
			issues = append(issues, msg)
		}
	}
	if len(issues) == 0 {
		return nil
	}
	for _, issue := range issues {
		fmt.Fprintln(k.Stderr, issue)
	}
	return fmt.Errorf("found %d issues in %s", len(issues), filename)
}

type updateUrlsCmd struct{}

func (c *updateUrlsCmd) Run(root *rootCmd) error {
	githubDir := filepath.Join(root.WorkingDir, "github")
	_, meta, err := root.metadata()
	if err != nil {
		return err
	}
	err = updateDocLinks(meta, githubDir)
	return err
}

type unusedOpsCmd struct {
	JSON bool `kong:"help='Output JSON.'"`
}

func (c *unusedOpsCmd) Run(root *rootCmd) error {
	_, meta, err := root.metadata()
	if err != nil {
		return err
	}
	var unused []*operation
	for _, op := range meta.operations() {
		goMethods := meta.operationMethods(op.Name)
		if len(goMethods) == 0 {
			unused = append(unused, op)
		}
	}
	if c.JSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(unused)
	}
	fmt.Printf("Found %d unused operations\n", len(unused))
	if len(unused) == 0 {
		return nil
	}
	fmt.Println("")
	for _, op := range unused {
		fmt.Println(op.Name)
		fmt.Printf("doc:     %s\n", op.DocumentationURL)
		fmt.Println("")
	}
	return nil
}

type canonizeCmd struct{}

func (c *canonizeCmd) Run(root *rootCmd) error {
	filename, meta, err := root.metadata()
	if err != nil {
		return err
	}
	err = meta.canonizeMethodOperations()
	if err != nil {
		return err
	}
	return meta.saveFile(filename)
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
