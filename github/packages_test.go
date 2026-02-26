// Copyright 2021 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPackageRegistry_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageRegistry{}, "{}")

	o := &PackageRegistry{
		AboutURL: Ptr("aurl"),
		Name:     Ptr("name"),
		Type:     Ptr("type"),
		URL:      Ptr("url"),
		Vendor:   Ptr("vendor"),
	}
	want := `{
		"about_url": "aurl",
		"name": "name",
		"type": "type",
		"url": "url",
		"vendor": "vendor"
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageFile_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageFile{}, "{}")

	o := &PackageFile{
		DownloadURL: Ptr("durl"),
		ID:          Ptr(int64(1)),
		Name:        Ptr("name"),
		SHA256:      Ptr("sha256"),
		SHA1:        Ptr("sha1"),
		MD5:         Ptr("md5"),
		ContentType: Ptr("ct"),
		State:       Ptr("state"),
		Author: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Size:      Ptr(int64(1)),
		CreatedAt: &Timestamp{referenceTime},
		UpdatedAt: &Timestamp{referenceTime},
	}

	want := `{
		"download_url": "durl",
		"id": 1,
		"name": "name",
		"sha256": "sha256",
		"sha1": "sha1",
		"md5": "md5",
		"content_type": "ct",
		"state": "state",
		"author": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"size": 1,
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageRelease_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageRelease{}, "{}")

	o := &PackageRelease{
		URL:             Ptr("url"),
		HTMLURL:         Ptr("hurl"),
		ID:              Ptr(int64(1)),
		TagName:         Ptr("tn"),
		TargetCommitish: Ptr("tcs"),
		Name:            Ptr("name"),
		Draft:           Ptr(true),
		Author: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Prerelease:  Ptr(true),
		CreatedAt:   &Timestamp{referenceTime},
		PublishedAt: &Timestamp{referenceTime},
	}

	want := `{
		"url": "url",
		"html_url": "hurl",
		"id": 1,
		"tag_name": "tn",
		"target_commitish": "tcs",
		"name": "name",
		"draft": true,
		"author": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"prerelease": true,
		"created_at": ` + referenceTimeStr + `,
		"published_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageVersion_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageVersion{}, "{}")

	metadata, _ := json.Marshal([]string{"a", "b"})

	o := &PackageVersion{
		ID:             Ptr(int64(1)),
		Name:           Ptr("n"),
		URL:            Ptr("url"),
		PackageHTMLURL: Ptr("phurl"),
		License:        Ptr("l"),
		Description:    Ptr("d"),
		CreatedAt:      &Timestamp{referenceTime},
		UpdatedAt:      &Timestamp{referenceTime},
		Metadata:       metadata,
		Version:        Ptr("ver"),
		Summary:        Ptr("sum"),
		Body:           json.RawMessage(`"body"`),
		BodyHTML:       Ptr("btnhtml"),
		Release: &PackageRelease{
			URL:             Ptr("url"),
			HTMLURL:         Ptr("hurl"),
			ID:              Ptr(int64(1)),
			TagName:         Ptr("tn"),
			TargetCommitish: Ptr("tcs"),
			Name:            Ptr("name"),
			Draft:           Ptr(true),
			Author: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			Prerelease:  Ptr(true),
			CreatedAt:   &Timestamp{referenceTime},
			PublishedAt: &Timestamp{referenceTime},
		},
		Manifest:        Ptr("mani"),
		HTMLURL:         Ptr("hurl"),
		TagName:         Ptr("tn"),
		TargetCommitish: Ptr("tcs"),
		TargetOID:       Ptr("tid"),
		Draft:           Ptr(true),
		Prerelease:      Ptr(true),
		ContainerMetadata: &PackageEventContainerMetadata{
			Labels: map[string]any{"l1": true, "l2": "a"},
		},
		DockerMetadata: []any{"a", "b"},
		NPMMetadata: Ptr(PackageNPMMetadata{
			Name: Ptr("n"),
		}),
		NugetMetadata: []*PackageNugetMetadata{
			{Name: Ptr("n")},
		},
		RubyMetadata: map[string]any{"k1": "v1", "k2": "v2"},
		PackageFiles: []*PackageFile{
			{
				DownloadURL: Ptr("durl"),
				ID:          Ptr(int64(1)),
				Name:        Ptr("name"),
				SHA256:      Ptr("sha256"),
				SHA1:        Ptr("sha1"),
				MD5:         Ptr("md5"),
				ContentType: Ptr("ct"),
				State:       Ptr("state"),
				Author: &User{
					Login:           Ptr("l"),
					ID:              Ptr(int64(1)),
					URL:             Ptr("u"),
					AvatarURL:       Ptr("a"),
					GravatarID:      Ptr("g"),
					Name:            Ptr("n"),
					Company:         Ptr("c"),
					Blog:            Ptr("b"),
					Location:        Ptr("l"),
					Email:           Ptr("e"),
					Hireable:        Ptr(true),
					Bio:             Ptr("b"),
					TwitterUsername: Ptr("t"),
					PublicRepos:     Ptr(1),
					Followers:       Ptr(1),
					Following:       Ptr(1),
					CreatedAt:       &Timestamp{referenceTime},
					SuspendedAt:     &Timestamp{referenceTime},
				},
				Size:      Ptr(int64(1)),
				CreatedAt: &Timestamp{referenceTime},
				UpdatedAt: &Timestamp{referenceTime},
			},
		},
		PackageURL: Ptr("purl"),
		Author: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		SourceURL:           Ptr("surl"),
		InstallationCommand: Ptr("ic"),
		DeletedAt:           &Timestamp{referenceTime},
	}

	want := `{
		"id": 1,
		"name": "n",
		"url": "url",
		"package_html_url": "phurl",
		"license": "l",
		"description": "d",
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"metadata": ["a", "b"],
		"version": "ver",
		"summary": "sum",
		"body": "body",
		"body_html": "btnhtml",
		"release": {
			"url": "url",
			"html_url": "hurl",
			"id": 1,
			"tag_name": "tn",
			"target_commitish": "tcs",
			"name": "name",
			"draft": true,
			"author": {
				"login": "l",
				"id": 1,
				"avatar_url": "a",
				"gravatar_id": "g",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e",
				"hireable": true,
				"bio": "b",
				"twitter_username": "t",
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"suspended_at": ` + referenceTimeStr + `,
				"url": "u"
			},
			"prerelease": true,
			"created_at": ` + referenceTimeStr + `,
			"published_at": ` + referenceTimeStr + `
		},
		"manifest": "mani",
		"html_url": "hurl",
		"tag_name": "tn",
		"target_commitish": "tcs",
		"target_oid": "tid",
		"draft": true,
		"prerelease": true,
		"container_metadata": {
			"labels": {
				"l1": true,
				"l2": "a"
			}
		},
		"docker_metadata": ["a", "b"],
		"npm_metadata": {
			"name": "n"
		},
		"nuget_metadata": [{"name": "n"}],
		"ruby_metadata": {
			"k1": "v1",
			"k2": "v2"
		},
		"package_files": [
			{
				"download_url": "durl",
				"id": 1,
				"name": "name",
				"sha256": "sha256",
				"sha1": "sha1",
				"md5": "md5",
				"content_type": "ct",
				"state": "state",
				"author": {
					"login": "l",
					"id": 1,
					"avatar_url": "a",
					"gravatar_id": "g",
					"name": "n",
					"company": "c",
					"blog": "b",
					"location": "l",
					"email": "e",
					"hireable": true,
					"bio": "b",
					"twitter_username": "t",
					"public_repos": 1,
					"followers": 1,
					"following": 1,
					"created_at": ` + referenceTimeStr + `,
					"suspended_at": ` + referenceTimeStr + `,
					"url": "u"
				},
				"size": 1,
				"created_at": ` + referenceTimeStr + `,
				"updated_at": ` + referenceTimeStr + `
			}
		],
		"package_url": "purl",
		"author": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"source_url": "surl",
		"installation_command": "ic",
		"deleted_at": ` + referenceTimeStr + `
	}`

	testJSONMarshal(t, o, want, cmpJSONRawMessageComparator())
}

func TestPackageVersion_GetBody(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pv        *PackageVersion
		wantValue string
		wantOk    bool
	}{
		"pv nil": {
			pv:        nil,
			wantValue: "",
			wantOk:    false,
		},
		"body nil": {
			pv: &PackageVersion{
				Body: nil,
			},
			wantValue: "",
			wantOk:    false,
		},
		"invalid body": {
			pv: &PackageVersion{
				Body: json.RawMessage(`{
					"repository": {
						"name": "n"
					},
					"info": {
						"type": "t"
					}
				}`),
			},
			wantValue: "",
			wantOk:    false,
		},
		"valid body": {
			pv: &PackageVersion{
				Body: json.RawMessage(`"body"`),
			},
			wantValue: "body",
			wantOk:    true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resValue, resOk := test.pv.GetBody()

			if resValue != test.wantValue || resOk != test.wantOk {
				t.Errorf("PackageVersion.GetBody() - got: %v, %v; want: %v, %v", resValue, resOk, test.wantValue, test.wantOk)
			}
		})
	}
}

func TestPackageVersion_GetBodyAsPackageVersionBody(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pv        *PackageVersion
		wantValue *PackageVersionBody
		wantOk    bool
	}{
		"pv nil": {
			pv:        nil,
			wantValue: nil,
			wantOk:    false,
		},
		"body nil": {
			pv: &PackageVersion{
				Body: nil,
			},
			wantValue: nil,
			wantOk:    false,
		},
		"invalid body": {
			pv: &PackageVersion{
				Body: json.RawMessage(`"body"`),
			},
			wantValue: nil,
			wantOk:    false,
		},
		"valid body": {
			pv: &PackageVersion{
				Body: json.RawMessage(`{
					"repository": {
						"name": "n"
					},
					"info": {
						"type": "t"
					}
				}`),
			},
			wantValue: &PackageVersionBody{
				Repo: &Repository{
					Name: Ptr("n"),
				},
				Info: &PackageVersionBodyInfo{
					Type: Ptr("t"),
				},
			},
			wantOk: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resValue, resOk := test.pv.GetBodyAsPackageVersionBody()

			if !cmp.Equal(resValue, test.wantValue) || resOk != test.wantOk {
				t.Errorf("PackageVersion.GetBodyAsPackageVersionBody() - got: %v, %v; want: %v, %v", resValue, resOk, test.wantValue, test.wantOk)
			}
		})
	}
}

func TestPackageVersion_GetMetadata(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pv        *PackageVersion
		wantValue *PackageMetadata
		wantOk    bool
	}{
		"pv nil": {
			pv:        nil,
			wantValue: nil,
			wantOk:    false,
		},
		"metadata nil": {
			pv: &PackageVersion{
				Metadata: nil,
			},
			wantValue: nil,
			wantOk:    false,
		},
		"invalid metadata": {
			pv: &PackageVersion{
				Metadata: json.RawMessage(`[]`),
			},
			wantValue: nil,
			wantOk:    false,
		},
		"valid metadata": {
			pv: &PackageVersion{
				Metadata: json.RawMessage(`{
					"package_type": "container",
					"container": {
						"tags": ["a"]
					}
				}`),
			},
			wantValue: &PackageMetadata{
				PackageType: Ptr("container"),
				Container: &PackageContainerMetadata{
					Tags: []string{"a"},
				},
			},
			wantOk: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resValue, resOk := test.pv.GetMetadata()

			if !cmp.Equal(resValue, test.wantValue) || resOk != test.wantOk {
				t.Errorf("PackageVersion.GetMetadata() - got: %v, %v; want: %v, %v", resValue, resOk, test.wantValue, test.wantOk)
			}
		})
	}
}

func TestPackageVersion_GetRawMetadata(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pv   *PackageVersion
		want json.RawMessage
	}{
		"pv nil": {
			pv:   nil,
			want: nil,
		},
		"metadata nil": {
			pv: &PackageVersion{
				Metadata: nil,
			},
			want: json.RawMessage{},
		},
		"valid metadata": {
			pv: &PackageVersion{
				Metadata: json.RawMessage(`"a"`),
			},
			want: json.RawMessage(`"a"`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			res := test.pv.GetRawMetadata()

			if string(res) != string(test.want) {
				t.Errorf("PackageVersion.GetRawMetadata() - got: %v; want: %v", res, test.want)
			}
		})
	}
}

func TestPackage_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &Package{}, "{}")

	o := &Package{
		ID:          Ptr(int64(1)),
		Name:        Ptr("name"),
		PackageType: Ptr("pt"),
		HTMLURL:     Ptr("hurl"),
		Visibility:  Ptr("private"),
		Owner: &User{
			Login:           Ptr("l"),
			ID:              Ptr(int64(1)),
			URL:             Ptr("u"),
			AvatarURL:       Ptr("a"),
			GravatarID:      Ptr("g"),
			Name:            Ptr("n"),
			Company:         Ptr("c"),
			Blog:            Ptr("b"),
			Location:        Ptr("l"),
			Email:           Ptr("e"),
			Hireable:        Ptr(true),
			Bio:             Ptr("b"),
			TwitterUsername: Ptr("t"),
			PublicRepos:     Ptr(1),
			Followers:       Ptr(1),
			Following:       Ptr(1),
			CreatedAt:       &Timestamp{referenceTime},
			SuspendedAt:     &Timestamp{referenceTime},
		},
		Repository:  &Repository{ID: Ptr(int64(1))},
		CreatedAt:   &Timestamp{referenceTime},
		UpdatedAt:   &Timestamp{referenceTime},
		Namespace:   Ptr("n"),
		Description: Ptr("d"),
		Ecosystem:   Ptr("e"),
		PackageVersion: &PackageVersion{
			ID:       Ptr(int64(1)),
			Version:  Ptr("ver"),
			Summary:  Ptr("sum"),
			Body:     json.RawMessage(`"body"`),
			BodyHTML: Ptr("btnhtml"),
			Release: &PackageRelease{
				URL:             Ptr("url"),
				HTMLURL:         Ptr("hurl"),
				ID:              Ptr(int64(1)),
				TagName:         Ptr("tn"),
				TargetCommitish: Ptr("tcs"),
				Name:            Ptr("name"),
				Draft:           Ptr(true),
				Author: &User{
					Login:           Ptr("l"),
					ID:              Ptr(int64(1)),
					URL:             Ptr("u"),
					AvatarURL:       Ptr("a"),
					GravatarID:      Ptr("g"),
					Name:            Ptr("n"),
					Company:         Ptr("c"),
					Blog:            Ptr("b"),
					Location:        Ptr("l"),
					Email:           Ptr("e"),
					Hireable:        Ptr(true),
					Bio:             Ptr("b"),
					TwitterUsername: Ptr("t"),
					PublicRepos:     Ptr(1),
					Followers:       Ptr(1),
					Following:       Ptr(1),
					CreatedAt:       &Timestamp{referenceTime},
					SuspendedAt:     &Timestamp{referenceTime},
				},
				Prerelease:  Ptr(true),
				CreatedAt:   &Timestamp{referenceTime},
				PublishedAt: &Timestamp{referenceTime},
			},
			Manifest:        Ptr("mani"),
			HTMLURL:         Ptr("hurl"),
			TagName:         Ptr("tn"),
			TargetCommitish: Ptr("tcs"),
			TargetOID:       Ptr("tid"),
			Draft:           Ptr(true),
			Prerelease:      Ptr(true),
			CreatedAt:       &Timestamp{referenceTime},
			UpdatedAt:       &Timestamp{referenceTime},
			PackageFiles: []*PackageFile{
				{
					DownloadURL: Ptr("durl"),
					ID:          Ptr(int64(1)),
					Name:        Ptr("name"),
					SHA256:      Ptr("sha256"),
					SHA1:        Ptr("sha1"),
					MD5:         Ptr("md5"),
					ContentType: Ptr("ct"),
					State:       Ptr("state"),
					Author: &User{
						Login:           Ptr("l"),
						ID:              Ptr(int64(1)),
						URL:             Ptr("u"),
						AvatarURL:       Ptr("a"),
						GravatarID:      Ptr("g"),
						Name:            Ptr("n"),
						Company:         Ptr("c"),
						Blog:            Ptr("b"),
						Location:        Ptr("l"),
						Email:           Ptr("e"),
						Hireable:        Ptr(true),
						Bio:             Ptr("b"),
						TwitterUsername: Ptr("t"),
						PublicRepos:     Ptr(1),
						Followers:       Ptr(1),
						Following:       Ptr(1),
						CreatedAt:       &Timestamp{referenceTime},
						SuspendedAt:     &Timestamp{referenceTime},
					},
					Size:      Ptr(int64(1)),
					CreatedAt: &Timestamp{referenceTime},
					UpdatedAt: &Timestamp{referenceTime},
				},
			},
			Author: &User{
				Login:           Ptr("l"),
				ID:              Ptr(int64(1)),
				URL:             Ptr("u"),
				AvatarURL:       Ptr("a"),
				GravatarID:      Ptr("g"),
				Name:            Ptr("n"),
				Company:         Ptr("c"),
				Blog:            Ptr("b"),
				Location:        Ptr("l"),
				Email:           Ptr("e"),
				Hireable:        Ptr(true),
				Bio:             Ptr("b"),
				TwitterUsername: Ptr("t"),
				PublicRepos:     Ptr(1),
				Followers:       Ptr(1),
				Following:       Ptr(1),
				CreatedAt:       &Timestamp{referenceTime},
				SuspendedAt:     &Timestamp{referenceTime},
			},
			InstallationCommand: Ptr("ic"),
		},
		Registry: &PackageRegistry{
			AboutURL: Ptr("aurl"),
			Name:     Ptr("name"),
			Type:     Ptr("type"),
			URL:      Ptr("url"),
			Vendor:   Ptr("vendor"),
		},
		URL:          Ptr("url"),
		VersionCount: Ptr(int64(1)),
	}

	want := `{
		"id": 1,
		"name": "name",
		"package_type": "pt",
		"html_url": "hurl",
		"visibility": "private",
		"owner": {
			"login": "l",
			"id": 1,
			"avatar_url": "a",
			"gravatar_id": "g",
			"name": "n",
			"company": "c",
			"blog": "b",
			"location": "l",
			"email": "e",
			"hireable": true,
			"bio": "b",
			"twitter_username": "t",
			"public_repos": 1,
			"followers": 1,
			"following": 1,
			"created_at": ` + referenceTimeStr + `,
			"suspended_at": ` + referenceTimeStr + `,
			"url": "u"
		},
		"repository": {
			"id": 1
		},
		"created_at": ` + referenceTimeStr + `,
		"updated_at": ` + referenceTimeStr + `,
		"namespace": "n",
		"description": "d",
		"ecosystem": "e",
		"package_version": {
			"id": 1,
			"version": "ver",
			"summary": "sum",
			"body": "body",
			"body_html": "btnhtml",
			"release": {
				"url": "url",
				"html_url": "hurl",
				"id": 1,
				"tag_name": "tn",
				"target_commitish": "tcs",
				"name": "name",
				"draft": true,
				"author": {
					"login": "l",
					"id": 1,
					"avatar_url": "a",
					"gravatar_id": "g",
					"name": "n",
					"company": "c",
					"blog": "b",
					"location": "l",
					"email": "e",
					"hireable": true,
					"bio": "b",
					"twitter_username": "t",
					"public_repos": 1,
					"followers": 1,
					"following": 1,
					"created_at": ` + referenceTimeStr + `,
					"suspended_at": ` + referenceTimeStr + `,
					"url": "u"
				},
				"prerelease": true,
				"created_at": ` + referenceTimeStr + `,
				"published_at": ` + referenceTimeStr + `
			},
			"manifest": "mani",
			"html_url": "hurl",
			"tag_name": "tn",
			"target_commitish": "tcs",
			"target_oid": "tid",
			"draft": true,
			"prerelease": true,
			"created_at": ` + referenceTimeStr + `,
			"updated_at": ` + referenceTimeStr + `,
			"package_files": [
				{
					"download_url": "durl",
					"id": 1,
					"name": "name",
					"sha256": "sha256",
					"sha1": "sha1",
					"md5": "md5",
					"content_type": "ct",
					"state": "state",
					"author": {
						"login": "l",
						"id": 1,
						"avatar_url": "a",
						"gravatar_id": "g",
						"name": "n",
						"company": "c",
						"blog": "b",
						"location": "l",
						"email": "e",
						"hireable": true,
						"bio": "b",
						"twitter_username": "t",
						"public_repos": 1,
						"followers": 1,
						"following": 1,
						"created_at": ` + referenceTimeStr + `,
						"suspended_at": ` + referenceTimeStr + `,
						"url": "u"
					},
					"size": 1,
					"created_at": ` + referenceTimeStr + `,
					"updated_at": ` + referenceTimeStr + `
				}
			],
			"author": {
				"login": "l",
				"id": 1,
				"avatar_url": "a",
				"gravatar_id": "g",
				"name": "n",
				"company": "c",
				"blog": "b",
				"location": "l",
				"email": "e",
				"hireable": true,
				"bio": "b",
				"twitter_username": "t",
				"public_repos": 1,
				"followers": 1,
				"following": 1,
				"created_at": ` + referenceTimeStr + `,
				"suspended_at": ` + referenceTimeStr + `,
				"url": "u"
			},
			"installation_command": "ic"
		},
		"registry": {
			"about_url": "aurl",
			"name": "name",
			"type": "type",
			"url": "url",
			"vendor": "vendor"
		},
		"url": "url",
		"version_count": 1
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageMetadata_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageMetadata{}, "{}")

	o := &PackageMetadata{
		PackageType: Ptr("pt"),
		Container: Ptr(PackageContainerMetadata{
			Tags: []string{"a"},
		}),
	}

	want := `{
		"package_type": "pt",
		"container": {
			"tags": ["a"]
		}
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageContainerMetadata_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageContainerMetadata{}, "{}")

	o := &PackageContainerMetadata{
		Tags: []string{"a", "b"},
	}

	want := `{
		"tags": ["a","b"]
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageVersionBody_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageVersionBody{}, "{}")

	o := &PackageVersionBody{
		Repo: Ptr(Repository{
			ID: Ptr(int64(1)),
		}),
		Info: Ptr(PackageVersionBodyInfo{
			Type: Ptr("t"),
		}),
	}

	want := `{
		"repository": {
			"id": 1
		},
		"info": {
			"type": "t"
		}
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageVersionBodyInfo_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageVersionBodyInfo{}, "{}")

	o := &PackageVersionBodyInfo{
		Type:       Ptr("t"),
		OID:        Ptr("o"),
		Mode:       Ptr(int64(1)),
		Name:       Ptr("n"),
		Path:       Ptr("p"),
		Size:       Ptr(int64(1)),
		Collection: Ptr(true),
	}

	want := `{
		"type": "t",
		"oid": "o",
		"mode": 1,
		"name": "n",
		"path": "p",
		"size": 1,
		"collection": true
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageEventContainerMetadata_Marshal(t *testing.T) {
	t.Parallel()
	testJSONMarshal(t, &PackageEventContainerMetadata{}, "{}")

	o := &PackageEventContainerMetadata{
		Labels: map[string]any{
			"k": "v",
		},
		Manifest: map[string]any{
			"k": float64(1),
		},
		Tag: &PackageEventContainerMetadataTag{
			Name:   Ptr("n"),
			Digest: Ptr("d"),
		},
	}

	want := `{
		"labels": {
			"k": "v"
		},
		"manifest": {
			"k": 1
		},
		"tag": {
			"name": "n",
			"digest": "d"
		}
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageEventContainerMetadataTag_Marshal(t *testing.T) {
	t.Parallel()

	testJSONMarshal(t, &PackageEventContainerMetadataTag{}, "{}")

	o := &PackageEventContainerMetadataTag{
		Name:   Ptr("n"),
		Digest: Ptr("d"),
	}

	want := `{
		"name": "n",
		"digest": "d"
	}`

	testJSONMarshal(t, o, want)
}

func TestPackageNugetMetadata_Marshal(t *testing.T) {
	t.Parallel()

	o, _ := json.Marshal(map[string]string{
		"k1": "v1",
		"k2": "v2",
	})

	tests := map[string]struct {
		input *PackageNugetMetadata
		want  string
	}{
		"zero": {
			input: &PackageNugetMetadata{},
			want:  "{}",
		},
		"string": {
			input: &PackageNugetMetadata{
				ID:    json.RawMessage(`1`),
				Name:  Ptr("n"),
				Value: json.RawMessage(`"s"`),
			},
			want: `{
				"id": 1,
				"name": "n",
				"value": "s"
			}`,
		},
		"int": {
			input: &PackageNugetMetadata{
				ID:    json.RawMessage(`1`),
				Name:  Ptr("n"),
				Value: json.RawMessage(`1`),
			},
			want: `{
				"id": 1,
				"name": "n",
				"value": 1
			}`,
		},
		"object": {
			input: &PackageNugetMetadata{
				ID:    json.RawMessage(`1`),
				Name:  Ptr("n"),
				Value: o,
			},
			want: `{
				"id": 1,
				"name": "n",
				"value": {
					"k1": "v1",
					"k2": "v2"
				}
			}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testJSONMarshal(t, test.input, test.want, cmpJSONRawMessageComparator())
		})
	}
}

func TestPackageNPMMetadata_Marshal(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input *PackageNPMMetadata
		want  string
	}{
		"zero": {
			input: &PackageNPMMetadata{},
			want:  "{}",
		},
		"string": {
			input: &PackageNPMMetadata{
				Name:                 Ptr("n"),
				Version:              Ptr("v"),
				NPMUser:              Ptr("nu"),
				Description:          Ptr("d"),
				GitHead:              Ptr("gh"),
				Homepage:             Ptr("h"),
				License:              Ptr("l"),
				Main:                 Ptr("m"),
				ID:                   Ptr("id"),
				NodeVersion:          Ptr("nv"),
				NPMVersion:           Ptr("npmv"),
				Readme:               Ptr("r"),
				InstallationCommand:  Ptr("ic"),
				CommitOID:            Ptr("coid"),
				HasShrinkwrap:        Ptr(true),
				PublishedViaActions:  Ptr(true),
				ReleaseID:            Ptr(int64(1)),
				DeletedByID:          Ptr(int64(1)),
				Author:               map[string]string{"k1": "v1"},
				Bugs:                 map[string]string{"k1": "v1"},
				Dependencies:         map[string]string{"k1": "v1"},
				DevDependencies:      map[string]string{"k1": "v1"},
				PeerDependencies:     map[string]string{"k1": "v1"},
				OptionalDependencies: map[string]string{"k1": "v1"},
				Dist:                 map[string]string{"k1": "v1"},
				Repository:           map[string]string{"k1": "v1"},
				Engines:              map[string]string{"k1": "v1"},
				Directories:          map[string]string{"k1": "v1"},
				Scripts:              map[string]any{"k1": float64(1)},
				Bin:                  map[string]any{"k1": true},
				Man:                  map[string]any{"k1": "v1"},
				Keywords:             []string{"kw1", "kw2"},
				Files:                []string{"f1", "f2"},
				OS:                   []string{"os1", "os2"},
				CPU:                  []string{"cpu1", "cpu2"},
				Maintainers:          []any{"m1", "m2"},
				Contributors:         []any{"c1", "c2"},
			},
			want: `{
				"name": "n",
				"version": "v",
				"npm_user": "nu",
				"description": "d",
				"git_head": "gh",
				"homepage": "h",
				"license": "l",
				"main": "m",
				"id": "id",
				"node_version": "nv",
				"npm_version": "npmv",
				"readme": "r",
				"installation_command": "ic",
				"commit_oid": "coid",
				"has_shrinkwrap": true,
				"published_via_actions": true,
				"release_id": 1,
				"deleted_by_id": 1,
				"author": {
					"k1": "v1"
				},
				"bugs": {
					"k1": "v1"
				},
				"dependencies": {
					"k1": "v1"
				},
				"dev_dependencies": {
					"k1": "v1"
				},
				"peer_dependencies": {
					"k1": "v1"
				},
				"optional_dependencies": {
					"k1": "v1"
				},
				"dist": {
					"k1": "v1"
				},
				"repository": {
					"k1": "v1"
				},
				"engines": {
					"k1": "v1"
				},
				"directories": {
					"k1": "v1"
				},
				"scripts": {
					"k1": 1
				},
				"bin": {
					"k1": true
				},
				"man": {
					"k1": "v1"
				},
				"keywords": ["kw1", "kw2"],
				"files": ["f1", "f2"],
				"os": ["os1", "os2"],
				"cpu": ["cpu1", "cpu2"],
				"maintainers": ["m1", "m2"],
				"contributors": ["c1", "c2"]
			}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testJSONMarshal(t, test.input, test.want)
		})
	}
}
