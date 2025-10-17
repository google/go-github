// Copyright 2025 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsonfieldname is a custom linter to be used by
// golangci-lint to find instances where the Go field name
// of a struct does not match the JSON tag name.
// It honors idiomatic Go initialisms and handles the
// special case of `Github` vs `GitHub` as agreed upon
// by the original author of the repo.
package jsonfieldname

import (
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"strings"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("jsonfieldname", New)
}

// JSONFieldNamePlugin is a custom linter plugin for golangci-lint.
type JSONFieldNamePlugin struct{}

// New returns an analysis.Analyzer to use with golangci-lint.
func New(_ any) (register.LinterPlugin, error) {
	return &JSONFieldNamePlugin{}, nil
}

// BuildAnalyzers builds the analyzers for the JSONFieldNamePlugin.
func (f *JSONFieldNamePlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		{
			Name: "jsonfieldname",
			Doc:  "Reports mismatches between Go field and JSON tag names. Note that the JSON tag name is the source-of-truth and the Go field name needs to match it.",
			Run:  run,
		},
	}, nil
}

// GetLoadMode returns the load mode for the JSONFieldNamePlugin.
func (f *JSONFieldNamePlugin) GetLoadMode() string {
	return register.LoadModeSyntax
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil {
				return false
			}

			switch t := n.(type) {
			case *ast.TypeSpec:
				structType, ok := t.Type.(*ast.StructType)
				if !ok {
					return true
				}
				structName := t.Name.Name

				// Only check exported structs.
				if !ast.IsExported(structName) {
					return true
				}

				for _, field := range structType.Fields.List {
					if field.Tag == nil || len(field.Names) == 0 {
						continue
					}

					goField := field.Names[0]
					tagValue := strings.Trim(field.Tag.Value, "`")
					structTag := reflect.StructTag(tagValue)
					jsonTagName, ok := structTag.Lookup("json")
					if !ok || jsonTagName == "-" {
						continue
					}
					jsonTagName = strings.TrimSuffix(jsonTagName, ",omitempty")

					checkGoFieldName(structName, goField.Name, jsonTagName, goField.Pos(), pass)
				}
			}

			return true
		})
	}
	return nil, nil
}

func checkGoFieldName(structName, goFieldName, jsonTagName string, tokenPos token.Pos, pass *analysis.Pass) {
	fullName := structName + "." + goFieldName
	if allowedExceptions[fullName] {
		return
	}

	want, alternate := jsonTagToPascal(jsonTagName)
	if goFieldName != want && goFieldName != alternate {
		const msg = "change Go field name %q to %q for JSON tag %q in struct %q"
		pass.Reportf(tokenPos, msg, goFieldName, want, jsonTagName, structName)
	}
}

func splitJSONTag(jsonTagName string) []string {
	if strings.Contains(jsonTagName, "_") {
		return strings.Split(jsonTagName, "_")
	}

	if strings.Contains(jsonTagName, "-") {
		return strings.Split(jsonTagName, "-")
	}

	if strings.ToLower(jsonTagName) == jsonTagName { // single word
		return []string{jsonTagName}
	}

	s := camelCaseRE.ReplaceAllString(jsonTagName, "$1 $2")
	parts := strings.Fields(s)
	for i, part := range parts {
		parts[i] = strings.ToLower(part)
	}

	return parts
}

var camelCaseRE = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func jsonTagToPascal(jsonTagName string) (want, alternate string) {
	parts := splitJSONTag(jsonTagName)
	alt := make([]string, len(parts))
	for i, part := range parts {
		alt[i] = part
		if part == "" {
			continue
		}
		upper := strings.ToUpper(part)
		if initialisms[upper] {
			parts[i] = upper
			alt[i] = upper
		} else if specialCase, ok := specialCases[upper]; ok {
			parts[i] = specialCase
			alt[i] = specialCase
		} else if possibleAlternate, ok := possibleAlternates[upper]; ok {
			parts[i] = possibleAlternate
			alt[i] = strings.ToUpper(part[:1]) + part[1:]
		} else {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
			alt[i] = parts[i]
		}
	}
	return strings.Join(parts, ""), strings.Join(alt, "")
}

var allowedExceptions = map[string]bool{
	"ActionsCacheUsageList.RepoCacheUsage": true, // TODO: RepoCacheUsages ?
	"AuditEntry.ExternalIdentityNameID":    true,
	"AuditEntry.Timestamp":                 true,
	"CheckSuite.AfterSHA":                  true,
	"CheckSuite.BeforeSHA":                 true,
	"CodeSearchResult.CodeResults":         true,
	"CodeSearchResult.Total":               true,
	"CommitAuthor.Login":                   true,
	"CommitsSearchResult.Commits":          true,
	"CommitsSearchResult.Total":            true,
	"CreateOrgInvitationOptions.TeamID":    true, // TODO: TeamIDs
	"DependencyGraphSnapshot.Sha":          true, // TODO: SHA
	"Discussion.DiscussionCategory":        true, // TODO: Category ?
	"EditOwner.OwnerInfo":                  true,
	"Event.RawPayload":                     true,
	"HookRequest.RawPayload":               true,
	"HookResponse.RawPayload":              true,
	"Issue.PullRequestLinks":               true, // TODO: PullRequest
	"IssueImportRequest.IssueImport":       true, // TODO: Issue
	"IssuesSearchResult.Issues":            true, // TODO: Items
	"IssuesSearchResult.Total":             true,
	"LabelsSearchResult.Labels":            true, // TODO: Items
	"LabelsSearchResult.Total":             true,
	"ListCheckRunsResults.Total":           true,
	"ListCheckSuiteResults.Total":          true,
	"ListCustomDeploymentRuleIntegrationsResponse.AvailableIntegrations":      true,
	"ListDeploymentProtectionRuleResponse.ProtectionRules":                    true,
	"OrganizationCustomRepoRoles.CustomRepoRoles":                             true, // TODO: CustomRoles
	"OrganizationCustomRoles.CustomRepoRoles":                                 true, // TODO: Roles
	"PreReceiveHook.ConfigURL":                                                true,
	"ProjectV2ItemEvent.ProjectV2Item":                                        true, // TODO: ProjectsV2Item
	"Protection.RequireLinearHistory":                                         true, // TODO: RequiredLinearHistory
	"ProtectionRequest.RequireLinearHistory":                                  true, // TODO: RequiredLinearHistory
	"PullRequestComment.InReplyTo":                                            true, // TODO: InReplyToID
	"PullRequestReviewsEnforcementRequest.BypassPullRequestAllowancesRequest": true, // TODO: BypassPullRequestAllowances
	"PullRequestReviewsEnforcementRequest.DismissalRestrictionsRequest":       true, // TODO: DismissalRestrictions
	"PullRequestReviewsEnforcementUpdate.BypassPullRequestAllowancesRequest":  true, // TODO: BypassPullRequestAllowances
	"PullRequestReviewsEnforcementUpdate.DismissalRestrictionsRequest":        true, // TODO: DismissalRestrictions
	"Reactions.MinusOne":                                                      true,
	"Reactions.PlusOne":                                                       true,
	"RepositoriesSearchResult.Repositories":                                   true,
	"RepositoriesSearchResult.Total":                                          true,
	"RepositoryVulnerabilityAlert.GitHubSecurityAdvisoryID":                   true,
	"SCIMDisplayReference.Ref":                                                true,
	"SecretScanningAlertLocationDetails.Startline":                            true, // TODO: StartLine
	"SecretScanningPatternOverride.Bypassrate":                                true, // TODO: BypassRate
	"StarredRepository.Repository":                                            true, // TODO: Repo
	"Timeline.Requester":                                                      true, // TODO: ReviewRequester
	"Timeline.Reviewer":                                                       true, // TODO: RequestedReviewer
	"TopicsSearchResult.Topics":                                               true, // TODO: Items
	"TopicsSearchResult.Total":                                                true,
	"TotalCacheUsage.TotalActiveCachesUsageSizeInBytes":                       true, // TODO: TotalActiveCachesSizeInBytes
	"TransferRequest.TeamID":                                                  true, // TODO: TeamIDs
	"Tree.Entries":                                                            true,
	"User.LdapDn":                                                             true, // TODO: LDAPDN
	"UsersSearchResult.Total":                                                 true,
	"UsersSearchResult.Users":                                                 true,
	"WeeklyStats.Additions":                                                   true,
	"WeeklyStats.Commits":                                                     true,
	"WeeklyStats.Deletions":                                                   true,
	"WeeklyStats.Week":                                                        true,
}

// Common Go initialisms that should be all caps.
var initialisms = map[string]bool{
	"API": true, "ASCII": true,
	"CAA": true, "CAS": true, "CNAME": true, "CPU": true,
	"CSS": true, "CWE": true, "CVE": true, "CVSS": true,
	"DN": true, "DNS": true,
	"EOF": true, "EPSS": true,
	"GB": true, "GHSA": true, "GPG": true, "GUID": true,
	"HTML": true, "HTTP": true, "HTTPS": true,
	"ID": true, "IDE": true, "IDP": true, "IP": true, "JIT": true,
	"JSON": true,
	"LDAP": true, "LFS": true, "LHS": true,
	"MD5": true, "MS": true, "MX": true,
	"NPM": true, "NTP": true, "NVD": true,
	"OID": true, "OS": true,
	"PEM": true, "PR": true, "QPS": true,
	"RAM": true, "RHS": true, "RPC": true,
	"SAML": true, "SBOM": true, "SCIM": true,
	"SHA": true, "SHA1": true, "SHA256": true,
	"SKU": true, "SLA": true, "SMTP": true, "SNMP": true,
	"SPDX": true, "SPDXID": true, "SQL": true, "SSH": true,
	"SSL": true, "SSO": true, "SVN": true,
	"TCP": true, "TFVC": true, "TLS": true, "TTL": true,
	"UDP": true, "UI": true, "UID": true, "UUID": true,
	"URI": true, "URL": true, "UTF8": true,
	"VCF": true, "VCS": true, "VM": true,
	"XML": true, "XMPP": true, "XSRF": true, "XSS": true,
}

var specialCases = map[string]string{
	"CPUS":    "CPUs",
	"CWES":    "CWEs",
	"GRAPHQL": "GraphQL",
	"HREF":    "HRef",
	"IDS":     "IDs",
	"IPS":     "IPs",
	"OAUTH":   "OAuth",
	"OPENAPI": "OpenAPI",
	"URLS":    "URLs",
}

var possibleAlternates = map[string]string{
	"ORGANIZATION":  "Org",
	"ORGANIZATIONS": "Orgs",
	"REPOSITORY":    "Repo",
	"REPOSITORIES":  "Repos",
}
