package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"tools/internal"

	"github.com/alecthomas/kong"
	"github.com/google/go-github/v55/github"
)

var helpVars = kong.Vars{
	"update_help":     `Update metadata.yaml from OpenAPI descriptions in github.com/github/rest-api-description.`,
	"format_help":     `Format metadata.yaml.`,
	"validate_help":   `Validate that metadata.yaml is consistent with source code.`,
	"unused_ops_help": `List operations in metadata.yaml that don't have any associated go methods'.`,
}

type rootCmd struct {
	WorkingDir     string            `kong:"short=C,default='.',help='Working directory. Must be within a go-github root.'"`
	Filename       string            `kong:"help='Path to metadata.yaml. Defaults to <go-github-root>/metadata.yaml.'"`
	GithubDir      string            `kong:"help='Path to the github package. Defaults to <go-github-root>/github.'"`
	UpdateMetadata updateMetadataCmd `kong:"cmd,help=${update_help}"`
	UpdateUrls     updateUrlsCmd     `kong:"cmd,help='Update documentation URLs in the Go source files in the github directory to match the urls in the metadata file.'"`
	Format         formatCmd         `kong:"cmd,help=${format_help}"`
	Validate       validateCmd       `kong:"cmd,help=${validate_help}"`
	UnusedOps      unusedOpsCmd      `kong:"cmd,help=${unused_ops_help}"`
}

func (c *rootCmd) metadata() (string, *internal.Metadata, error) {
	filename := c.Filename
	if filename == "" {
		dir, err := internal.ProjRootDir(c.WorkingDir)
		if err != nil {
			return "", nil, err
		}
		filename = filepath.Join(dir, "metadata.yaml")
	}
	var meta internal.Metadata
	err := internal.LoadMetadataFile(filename, &meta)
	if err != nil {
		return "", nil, err
	}
	return filename, &meta, nil
}

func (c *rootCmd) githubDir() (string, error) {
	dir := c.GithubDir
	if dir == "" {
		githubDir, err := internal.ProjRootDir(c.WorkingDir)
		if err != nil {
			return "", err
		}
		dir = filepath.Join(githubDir, "github")
	}
	return dir, nil
}

type updateMetadataCmd struct {
	Ref string `kong:"default='main',help='git ref to pull OpenAPI descriptions from'"`
}

func (c *updateMetadataCmd) Run(root *rootCmd) error {
	ctx := context.Background()
	filename, meta, err := root.metadata()
	if err != nil {
		return err
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable must be set to a GitHub personal access token with the public_repo scope")
	}
	client := github.NewClient(nil).WithAuthToken(token).Repositories
	err = meta.UpdateFromGithub(ctx, client, c.Ref)
	if err != nil {
		return err
	}
	return meta.SaveFile(filename)
}

type formatCmd struct{}

func (c *formatCmd) Run(root *rootCmd) error {
	filename, meta, err := root.metadata()
	if err != nil {
		return err
	}
	return meta.SaveFile(filename)
}

type validateCmd struct{}

func (c *validateCmd) Run(root *rootCmd) error {
	githubDir, err := root.githubDir()
	if err != nil {
		return err
	}
	filename, meta, err := root.metadata()
	if err != nil {
		return err
	}
	issues, err := internal.ValidateMetadata(githubDir, meta)
	if err != nil {
		return err
	}
	if len(issues) == 0 {
		return nil
	}
	for _, issue := range issues {
		fmt.Fprintln(os.Stderr, issue)
	}
	return fmt.Errorf("found %d issues in %s", len(issues), filename)
}

type updateUrlsCmd struct{}

func (c *updateUrlsCmd) Run(root *rootCmd) error {
	githubDir, err := root.githubDir()
	if err != nil {
		return err
	}
	_, meta, err := root.metadata()
	if err != nil {
		return err
	}
	err = internal.UpdateDocLinks(meta, githubDir)
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
	var unused []*internal.Operation
	for _, op := range meta.Operations {
		if len(op.GoMethods) == 0 {
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
		fmt.Printf("%s %s\n", op.Method(), op.EndpointURL())
		fmt.Printf("summary: %s\n", op.Summary())
		fmt.Printf("plans:   %s\n", strings.Join(op.Plans(), ", "))
		fmt.Printf("doc:     %s\n", op.DocumentationURL())
		fmt.Println("")
	}
	return nil
}

func main() {
	var cmd rootCmd
	k := kong.Parse(&cmd, helpVars)
	err := k.Run()
	k.FatalIfErrorf(err)
}
