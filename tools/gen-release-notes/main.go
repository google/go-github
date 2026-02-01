// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gen-release-notes calls `git` to determine what the prior release was, (e.g. "v76.0.0")
// then calls `git` again to find out what changes were made since then.
//
// Finally, it writes the release notes to stdout, summarizing the
// breaking and non-breaking changes since that release.
//
// Usage:
//
//	go run tools/gen-release-notes/main.go [--tag v76.0.0]
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	sinceTag = flag.String("tag", "", "List all changes since this tag (e.g. 'v76.0.0')")

	descriptionRE = regexp.MustCompile(`^\* (.*?\((#[^\)]+)\))`)
	releaseTagRE  = regexp.MustCompile(`[^a-zA-Z0-9.\-_]+`)
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	priorRelease := *sinceTag
	if priorRelease == "" {
		priorRelease = getPriorRelease()
		log.Printf("Prior release: %v", priorRelease)
	}

	newChanges := newChangesSinceRelease(priorRelease)

	releaseNotes := genReleaseNotes(newChanges)
	fmt.Printf("%v%v", releaseNotes, "\n")

	log.Print("Done.")
}

func genReleaseNotes(text string) string {
	allLines := splitIntoPRs(text)
	fullBreakingLines, fullNonBreakingLines := splitBreakingLines(allLines)
	refBreakingLines, refNonBreakingLines := genRefLines(fullBreakingLines, fullNonBreakingLines)

	return fmt.Sprintf(releaseNotesFmt,
		strings.Join(fullBreakingLines, "\n"),
		strings.Join(fullNonBreakingLines, "\n"),
		strings.Join(refBreakingLines, "\n"),
		strings.Join(refNonBreakingLines, "\n"))
}

func splitIntoPRs(text string) []string {
	parts := strings.Split("\n"+text, "\ncommit ")
	if len(parts) < 2 {
		log.Fatal("unable to find PRs")
	}
	prs := make([]string, 0, len(parts)-1)
	for _, part := range parts {
		if part == "" {
			continue
		}
		lines := strings.Split(part, "\n")
		if len(lines) < 5 { // commit, Author:, Date:, blank, msg
			continue
		}
		var newPR []string
		for _, line := range lines[1:] {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "Author: ") || strings.HasPrefix(line, "Date: ") {
				continue
			}
			if len(newPR) == 0 {
				newPR = append(newPR, "* "+line)
			} else {
				newPR = append(newPR, "  "+line)
			}
		}
		prs = append(prs, strings.Join(newPR, "\n"))
	}
	return prs
}

func splitBreakingLines(allLines []string) (breaking, nonBreaking []string) {
	for _, pr := range allLines {
		if strings.Contains(pr, "!: ") {
			breaking = append(breaking, pr)
		} else {
			nonBreaking = append(nonBreaking, pr)
		}
	}
	return breaking, nonBreaking
}

func genRefLines(breaking, nonBreaking []string) (ref, refNon []string) {
	for _, pr := range breaking {
		m := descriptionRE.FindStringSubmatch(pr)
		if len(m) == 3 {
			ref = append(ref, strings.Replace(pr, m[1], m[2], 1))
		}
	}
	for _, pr := range nonBreaking {
		m := descriptionRE.FindStringSubmatch(pr)
		if len(m) == 3 {
			refNon = append(refNon, strings.Replace(pr, m[1], m[2], 1))
		}
	}
	return ref, refNon
}

func runCommand(cmdArgs []string) string {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //nolint:gosec
	out := &bytes.Buffer{}
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	log.Printf("Running command: %v", strings.Join(cmdArgs, " "))
	if err := cmd.Run(); err != nil {
		log.Fatalf("command failed: %v", err)
	}

	return strings.TrimSpace(out.String())
}

func newChangesSinceRelease(priorRelease string) string {
	priorRelease = releaseTagRE.ReplaceAllString(priorRelease, "")
	cmdArgs := []string{"git", "log", priorRelease + "..", "--no-color"}
	return runCommand(cmdArgs)
}

func getPriorRelease() string {
	cmdArgs := []string{"git", "describe", "--tags", "--abbrev=0"}
	return runCommand(cmdArgs)
}

const releaseNotesFmt = `
This release contains the following breaking API changes:

%v

...and the following additional changes:

%v

&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&

This release contains the following breaking API changes:

%v

...and the following additional changes:

%v
`
