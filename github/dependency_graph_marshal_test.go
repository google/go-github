// Copyright 2023 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"testing"
	"time"
)

func TestSBOM_Marshal(t *testing.T) {
	testJSONMarshal(t, &SBOM{}, `{}`)

	u := &SBOM{
		SBOM: &SBOMInfo{
			SPDXID:       Ptr("SPDXRef-DOCUMENT"),
			SPDXVersion:  Ptr("SPDX-2.3"),
			CreationInfo: &CreationInfo{Creators: []string{"GitHub"}},
			Name:         Ptr("owner/repo"),
		},
	}

	want := `{
		"sbom": {
			"SPDXID": "SPDXRef-DOCUMENT",
			"spdxVersion": "SPDX-2.3",
			"creationInfo": {
				"creators": [
					"GitHub"
				]
			},
			"name": "owner/repo"
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestCreationInfo_Marshal(t *testing.T) {
	testJSONMarshal(t, &CreationInfo{}, `{}`)

	created := Timestamp{time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)}
	u := &CreationInfo{
		Created:  &created,
		Creators: []string{"GitHub", "go-github"},
	}

	want := `{
		"created": "2024-01-02T03:04:05Z",
		"creators": [
			"GitHub",
			"go-github"
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestRepoDependencies_Marshal(t *testing.T) {
	testJSONMarshal(t, &RepoDependencies{}, `{}`)

	u := &RepoDependencies{
		SPDXID:           Ptr("SPDXRef-Package"),
		Name:             Ptr("pkg:golang/github.com/google/go-github"),
		VersionInfo:      Ptr("v88.0.0"),
		DownloadLocation: Ptr("NOASSERTION"),
		FilesAnalyzed:    Ptr(false),
		LicenseConcluded: Ptr("Apache-2.0"),
		LicenseDeclared:  Ptr("Apache-2.0"),
		ExternalRefs: []*PackageExternalRef{
			{
				ReferenceCategory: "PACKAGE-MANAGER",
				ReferenceType:     "purl",
				ReferenceLocator:  "pkg:golang/github.com/google/go-github@v88.0.0",
			},
		},
	}

	want := `{
		"SPDXID": "SPDXRef-Package",
		"name": "pkg:golang/github.com/google/go-github",
		"versionInfo": "v88.0.0",
		"downloadLocation": "NOASSERTION",
		"filesAnalyzed": false,
		"licenseConcluded": "Apache-2.0",
		"licenseDeclared": "Apache-2.0",
		"externalRefs": [
			{
				"referenceCategory": "PACKAGE-MANAGER",
				"referenceType": "purl",
				"referenceLocator": "pkg:golang/github.com/google/go-github@v88.0.0"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestPackageExternalRef_Marshal(t *testing.T) {
	testJSONMarshal(t, &PackageExternalRef{}, `{
		"referenceCategory": "",
		"referenceType": "",
		"referenceLocator": ""
	}`)

	u := &PackageExternalRef{
		ReferenceCategory: "PACKAGE-MANAGER",
		ReferenceType:     "purl",
		ReferenceLocator:  "pkg:npm/@actions/core@1.1.9",
	}

	want := `{
		"referenceCategory": "PACKAGE-MANAGER",
		"referenceType": "purl",
		"referenceLocator": "pkg:npm/@actions/core@1.1.9"
	}`

	testJSONMarshal(t, u, want)
}

func TestSBOMRelationship_Marshal(t *testing.T) {
	testJSONMarshal(t, &SBOMRelationship{}, `{
		"spdxElementId": "",
		"relatedSpdxElement": "",
		"relationshipType": ""
	}`)

	u := &SBOMRelationship{
		SPDXElementID:      "SPDXRef-DOCUMENT",
		RelatedSPDXElement: "SPDXRef-Package",
		RelationshipType:   "DESCRIBES",
	}

	want := `{
		"spdxElementId": "SPDXRef-DOCUMENT",
		"relatedSpdxElement": "SPDXRef-Package",
		"relationshipType": "DESCRIBES"
	}`

	testJSONMarshal(t, u, want)
}

func TestSBOMInfo_Marshal(t *testing.T) {
	testJSONMarshal(t, &SBOMInfo{}, `{}`)

	created := Timestamp{time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)}
	u := &SBOMInfo{
		SPDXID:            Ptr("SPDXRef-DOCUMENT"),
		SPDXVersion:       Ptr("SPDX-2.3"),
		CreationInfo:      &CreationInfo{Created: &created, Creators: []string{"GitHub"}},
		Name:              Ptr("owner/repo"),
		DataLicense:       Ptr("CC0-1.0"),
		DocumentDescribes: []string{"SPDXRef-Package"},
		DocumentNamespace: Ptr("https://github.com/owner/repo/dependency_graph/sbom"),
		Packages:          []*RepoDependencies{{Name: Ptr("github.com/google/go-github")}},
		Relationships:     []*SBOMRelationship{{SPDXElementID: "SPDXRef-DOCUMENT", RelatedSPDXElement: "SPDXRef-Package", RelationshipType: "DESCRIBES"}},
	}

	want := `{
		"SPDXID": "SPDXRef-DOCUMENT",
		"spdxVersion": "SPDX-2.3",
		"creationInfo": {
			"created": "2024-01-02T03:04:05Z",
			"creators": [
				"GitHub"
			]
		},
		"name": "owner/repo",
		"dataLicense": "CC0-1.0",
		"documentDescribes": [
			"SPDXRef-Package"
		],
		"documentNamespace": "https://github.com/owner/repo/dependency_graph/sbom",
		"packages": [
			{
				"name": "github.com/google/go-github"
			}
		],
		"relationships": [
			{
				"spdxElementId": "SPDXRef-DOCUMENT",
				"relatedSpdxElement": "SPDXRef-Package",
				"relationshipType": "DESCRIBES"
			}
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestDependencyGraphSnapshotResolvedDependency_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependencyGraphSnapshotResolvedDependency{}, `{}`)

	u := &DependencyGraphSnapshotResolvedDependency{
		PackageURL:   Ptr("pkg:npm/@actions/core@1.1.9"),
		Metadata:     map[string]any{"license": "MIT"},
		Relationship: Ptr("direct"),
		Scope:        Ptr("runtime"),
		Dependencies: []string{"@actions/http-client"},
	}

	want := `{
		"package_url": "pkg:npm/@actions/core@1.1.9",
		"metadata": {
			"license": "MIT"
		},
		"relationship": "direct",
		"scope": "runtime",
		"dependencies": [
			"@actions/http-client"
		]
	}`

	testJSONMarshal(t, u, want)
}

func TestDependencyGraphSnapshotJob_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependencyGraphSnapshotJob{}, `{}`)

	u := &DependencyGraphSnapshotJob{
		Correlator: Ptr("workflow_action"),
		ID:         Ptr("123456"),
		HTMLURL:    Ptr("https://github.com/owner/repo/actions/runs/123456"),
	}

	want := `{
		"correlator": "workflow_action",
		"id": "123456",
		"html_url": "https://github.com/owner/repo/actions/runs/123456"
	}`

	testJSONMarshal(t, u, want)
}

func TestDependencyGraphSnapshotDetector_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependencyGraphSnapshotDetector{}, `{}`)

	u := &DependencyGraphSnapshotDetector{
		Name:    Ptr("go-github"),
		Version: Ptr("v88.0.0"),
		URL:     Ptr("https://github.com/google/go-github"),
	}

	want := `{
		"name": "go-github",
		"version": "v88.0.0",
		"url": "https://github.com/google/go-github"
	}`

	testJSONMarshal(t, u, want)
}

func TestDependencyGraphSnapshotManifestFile_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependencyGraphSnapshotManifestFile{}, `{}`)

	u := &DependencyGraphSnapshotManifestFile{SourceLocation: Ptr("go.mod")}

	want := `{
		"source_location": "go.mod"
	}`

	testJSONMarshal(t, u, want)
}

func TestDependencyGraphSnapshotManifest_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependencyGraphSnapshotManifest{}, `{}`)

	u := &DependencyGraphSnapshotManifest{
		Name:     Ptr("go.mod"),
		File:     &DependencyGraphSnapshotManifestFile{SourceLocation: Ptr("go.mod")},
		Metadata: map[string]any{"ecosystem": "go"},
		Resolved: map[string]*DependencyGraphSnapshotResolvedDependency{
			"github.com/google/go-github": {
				PackageURL:   Ptr("pkg:golang/github.com/google/go-github@v88.0.0"),
				Relationship: Ptr("direct"),
				Scope:        Ptr("runtime"),
			},
		},
	}

	want := `{
		"name": "go.mod",
		"file": {
			"source_location": "go.mod"
		},
		"metadata": {
			"ecosystem": "go"
		},
		"resolved": {
			"github.com/google/go-github": {
				"package_url": "pkg:golang/github.com/google/go-github@v88.0.0",
				"relationship": "direct",
				"scope": "runtime"
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestDependencyGraphSnapshot_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependencyGraphSnapshot{}, `{"version": 0}`)

	scanned := Timestamp{time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)}
	u := &DependencyGraphSnapshot{
		Version:  1,
		Sha:      Ptr("ce587453ced02b1526dfb4cb910479d431683101"),
		Ref:      Ptr("refs/heads/main"),
		Job:      &DependencyGraphSnapshotJob{Correlator: Ptr("workflow_action"), ID: Ptr("123456")},
		Detector: &DependencyGraphSnapshotDetector{Name: Ptr("go-github"), Version: Ptr("v88.0.0")},
		Scanned:  &scanned,
		Metadata: map[string]any{"build": "release"},
		Manifests: map[string]*DependencyGraphSnapshotManifest{
			"go.mod": {
				Name: Ptr("go.mod"),
				File: &DependencyGraphSnapshotManifestFile{SourceLocation: Ptr("go.mod")},
			},
		},
	}

	want := `{
		"version": 1,
		"sha": "ce587453ced02b1526dfb4cb910479d431683101",
		"ref": "refs/heads/main",
		"job": {
			"correlator": "workflow_action",
			"id": "123456"
		},
		"detector": {
			"name": "go-github",
			"version": "v88.0.0"
		},
		"scanned": "2024-01-02T03:04:05Z",
		"metadata": {
			"build": "release"
		},
		"manifests": {
			"go.mod": {
				"name": "go.mod",
				"file": {
					"source_location": "go.mod"
				}
			}
		}
	}`

	testJSONMarshal(t, u, want)
}

func TestDependencyGraphSnapshotCreationData_Marshal(t *testing.T) {
	testJSONMarshal(t, &DependencyGraphSnapshotCreationData{}, `{"id": 0}`)

	createdAt := Timestamp{time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)}
	u := &DependencyGraphSnapshotCreationData{
		ID:        12345,
		CreatedAt: &createdAt,
		Message:   Ptr("Dependency results for the repo have been successfully updated."),
		Result:    Ptr("SUCCESS"),
	}

	want := `{
		"id": 12345,
		"created_at": "2024-01-02T03:04:05Z",
		"message": "Dependency results for the repo have been successfully updated.",
		"result": "SUCCESS"
	}`

	testJSONMarshal(t, u, want)
}
