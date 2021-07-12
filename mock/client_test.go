package mock

import (
	"context"
	"testing"

	"github.com/google/go-github/v37/github"
)

func TestMockClient(t *testing.T) {
	ctx := context.Background()

	mockedHttpClient := NewMockHttpClient(
		WithRequestMatch(
			RequestMatchUsersGet,
			MustMarshall(github.User{
				Name: github.String("foobar"),
			}),
		),
		WithRequestMatch(
			RequestMatchOrganizationsList,
			MustMarshall([]github.Organization{
				{
					Name: github.String("foobar123thisorgwasmocked"),
				},
			}),
		),
	)

	c := github.NewClient(mockedHttpClient)

	user, _, userErr := c.Users.Get(ctx, "someUser")

	if user == nil || user.Name == nil || *user.Name != "foobar" {
		t.Fatalf("User name is %s, want foobar", user)
	}

	if userErr != nil {
		t.Errorf("User err is %s, want nil", userErr.Error())
	}

	orgs, _, err := c.Organizations.List(
		ctx,
		*user.Name,
		nil,
	)

	if len(orgs) != 1 {
		t.Errorf("Orgs len is %d want 1", len(orgs))
	}

	if err != nil {
		t.Errorf("Err is %s, want nil", err.Error())
	}

	if *(orgs[0].Name) != "foobar123thisorgwasmocked" {
		t.Errorf("orgs[0].Name is %s, want %s", *orgs[0].Name, "foobar123thisorgdoesnotexist")
	}
}
