// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gen-release-notes first reads the web page https://github.com/google/go-github
// to determine what the prior release was, (e.g. "v76.0.0")
// then reads https://github.com/google/go-github/compare/${PRIOR_RELEASE}...master
// to find out what changes were made since then.
//
// Finally, it writes the release notes to stdout, summarizing the
// breaking and non-breaking changes since that release.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

const (
	baseWebURL = "https://github.com/google/go-github"

	// fragile, but works for now.
	detailsDiv = `<div class="flex-auto min-width-0 js-details-container Details">`
)

var (
	releaseRE     = regexp.MustCompile(`<span [^>]+>([^<]+)</span>\s*<span title="Label: Latest"[^>]*>`)
	linkPrimaryRE = regexp.MustCompile(`(?ms)(.*?<a class="Link--primary[^>]+>)`)
	startAvatarRE = regexp.MustCompile(`(?ms)(<div class="AvatarStack)`)
	bracketsRE    = regexp.MustCompile(`(?ms)(<[^>]+>)`)
	newlinesRE    = regexp.MustCompile(`(?m)(\n+)`)
	descriptionRE = regexp.MustCompile(`^\* (.*?\((#[^\)]+)\))`)
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	priorRelease := getPriorRelease()

	newChanges := newChangesSinceRelease(priorRelease)

	releaseNotes := genReleaseNotes(newChanges)
	fmt.Printf("%v%v", releaseNotes, "\n")

	log.Print("Done.")
}

func genReleaseNotes(text string) string {
	// strip everything before first detailsDiv:
	idx := strings.Index(text, detailsDiv)
	if idx < 0 {
		log.Fatal("Could not find detailsDiv")
	}
	text = text[idx:]

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
	parts := strings.Split(text, detailsDiv)
	if len(parts) < 2 {
		log.Fatal("unable to find PRs")
	}
	prs := make([]string, 0, len(parts)-1)
	for _, part := range parts {
		if part == "" {
			continue
		}
		newDiv := matchDivs(part)
		for {
			oldDiv := newDiv
			newDiv = stripPRHTML(oldDiv)
			if newDiv == oldDiv {
				break
			}
		}
		prs = append(prs, newDiv)
	}
	return prs
}

func stripPRHTML(text string) string {
	_, innerText := getTagSequence(text)
	if i := strings.Index(text, "</a>"); i > 0 {
		newText := text[:i] + strings.Join(innerText, "")
		newText = strings.ReplaceAll(newText, "…", "")
		newText = newlinesRE.ReplaceAllString(newText, "\n  ")
		return newText
	}
	return text
}

func getTagSequence(text string) (tagSeq, innerText []string) {
	m := bracketsRE.FindAllStringIndex(text, -1)
	var lastEnd int
	for _, pair := range m {
		start := pair[0]
		end := pair[1] - 1
		if lastEnd > 0 && start > lastEnd+1 {
			rawText := text[lastEnd+1 : start]
			s := strings.TrimSpace(rawText)
			switch s {
			case "", "&hellip;": // skip
			default:
				// Special case:
				if strings.HasPrefix(rawText, "BREAKING") {
					rawText = "\n\n" + rawText
				}
				innerText = append(innerText, rawText)
			}
		}
		lastEnd = end
		s := text[start+1 : end]
		if s == "code" {
			innerText = append(innerText, " `")
			continue
		}
		if s == "/code" {
			innerText = append(innerText, "` ")
			continue
		}
		if s[0] == '/' {
			tagSeq = append(tagSeq, s)
			continue
		}
		if i := strings.Index(s, " "); i > 0 {
			tagSeq = append(tagSeq, s[0:i])
		} else {
			tagSeq = append(tagSeq, s)
		}
	}
	return tagSeq, innerText
}

func matchDivs(text string) string {
	chunks := strings.Split(text, `</div>`)
	var divCount int
	var lastChunk int
	for i, chunk := range chunks {
		chunks[i] = strings.TrimSpace(chunks[i])
		divsInChunk := strings.Count(chunk, `<div `)
		divCount += divsInChunk
		lastChunk++
		if lastChunk == divCount {
			newDivs := strings.Join(chunks[:lastChunk], "")
			return stripLinkPrimary(newDivs)
		}
	}
	return ""
}

func stripLinkPrimary(text string) string {
	m := linkPrimaryRE.FindStringSubmatch(text)
	if len(m) != 2 {
		log.Fatalf("unable to find link primary in: '%v'", text)
	}
	newText := strings.TrimSpace(text[len(m[0]):])
	// As a special case, trim off all the Avatar stuff
	m2 := startAvatarRE.FindStringIndex(newText)
	if len(m2) > 0 {
		newText = newText[:m2[0]]
	}
	return "* " + newText
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

func newChangesSinceRelease(priorRelease string) string {
	url := fmt.Sprintf("%v/compare/%v...master", baseWebURL, priorRelease)
	resp, err := http.Get(url) //nolint:gosec
	must(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	must(err)

	return string(body)
}

func getPriorRelease() string {
	resp, err := http.Get(baseWebURL)
	must(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	must(err)

	matches := releaseRE.FindStringSubmatch(string(body))
	if len(matches) != 2 {
		log.Fatal("could not find release info")
	}

	priorRelease := strings.TrimSpace(matches[1])
	if priorRelease == "" {
		log.Fatal("found empty prior release version")
	}

	return priorRelease
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const releaseNotesFmt = `This release contains the following breaking API changes:

%v

...and the following additional changes:

%v

&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&

This release contains the following breaking API changes:

%v

...and the following additional changes:

%v
`
