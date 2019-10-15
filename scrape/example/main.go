// The scrape tool demonstrates use of the github.com/google/go-github/scrape
// package to fetch data from GitHub.  The tool lists whether third-party app
// restrictions are enabled for an organization, and lists information about
// OAuth apps requested for the org.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/scrape"
)

var (
	username = flag.String("username", "", "github auth: username")
	password = flag.String("password", "", "github auth: password")
	otpseed  = flag.String("otpseed", "", "github auth: otp seed")
	org      = flag.String("org", "", "github org to get data for")
)

func main() {
	flag.Parse()

	// prompt for password and otpseed in case the user didn't want to specify as flags
	reader := bufio.NewReader(os.Stdin)
	if *password == "" {
		fmt.Print("password: ")
		*password, _ = reader.ReadString('\n')
		*password = strings.TrimSpace(*password)
	}
	if *otpseed == "" {
		fmt.Print("OTP seed: ")
		*otpseed, _ = reader.ReadString('\n')
		*otpseed = strings.TrimSpace(*otpseed)
	}

	client := scrape.NewClient(nil)

	if err := client.Authenticate(*username, *password, *otpseed); err != nil {
		log.Fatal(err)
	}

	enabled, err := client.AppRestrictionsEnabled(*org)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("App restrictions enabled for %q: %t\n", *org, enabled)

	apps, err := client.ListOAuthApps(*org)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("OAuth apps for %q: \n", *org)
	for _, app := range apps {
		fmt.Printf("\t%+v\n", app)
	}
}
