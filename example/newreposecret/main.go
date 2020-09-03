package main

import (
	"context"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	sodium "github.com/GoKillers/libsodium-go/cryptobox"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

var (
	repo  = flag.String("repo", "", "The repo that the secret should be added to, ex. go-github")
	owner = flag.String("owner", "", "The owner of there repo this should be added to, ex. google")
)

/*
main creates a new secret in github for a given owner/repo based on a secretName passed in as an arg
and a github token, and secret value provided via the environment

usage:
	export GITHUB_AUTH_TOKEN=<auth token from github that has secret create rights>
	export SECRET_VARIABLE=<secret value of the secret variable>
	go run main.go -owner <owner name> -repo <repository name> SECRET_VARIABLE

ex:
	export GITHUB_AUTH_TOKEN=0000000000000000
	export SECRET_VARIABLE="my-secret"
	go run main.go -owner google -repo go-github SECRET_VARIABLE

Once it runs go to the github repository > settings > left side options bar > Secrets
And you should see the new secret appear there.
*/
func main() {
	flag.Parse()

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		log.Fatal(`Unauthorized: No token present in environment variables. 
Add your github api token as an environment variable named GITHUB_AUTH_TOKEN

ex:
	export GITHUB_AUTH_TOKEN=<github auth token>`)
	}

	if *repo == "" {
		log.Fatal(`required flag repo was not passed. 
please pass the repo you want to add this flag to as part of the command

ex:
	go run main.go --repo <repo name> --owner <repo owner name> 
`)
	}

	if *owner == "" {
		log.Fatal(`required flag owner was not passed. 
please pass the owner you want to add this flag to as part of the command

ex:
	go run main.go --repo <repo name> --owner <repo owner name> 
`)
	}

	ctx, client, err := GithubAuth()
	if err != nil {
		log.Fatal("unable to authorize client with github using token found in env var GITHUB_AUTH_TOKEN, error was passed: ", err.Error())
	}
	secretName, err := getSecretName()
	if err != nil {
		log.Fatal(err.Error())
	}

	secretValue, err := getSecretValue(secretName)
	if err != nil {
		log.Fatal(err.Error())
	}

	addedSecretName, err := AddRepoSecret(ctx, client, *owner, *repo, secretName, secretValue)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(fmt.Sprintf("Added Secret %s to the repo %s/%s", addedSecretName, *owner, *repo))

}

func getSecretName() (string, error) {
	if len(os.Args) < 2 {
		err := errors.New("no secrets passed in the go run command. pass in a secret name and then ensure to have the secret variable as an environment variable")
		return "", err
	}
	args := os.Args
	secret := args[3:]
	return secret[0], nil
}

func getSecretValue(secretName string) (string, error) {
	secretValue := os.Getenv(secretName)
	if secretValue == "" {
		return "", errors.New(fmt.Sprintf(`secret with name: %s not defined as an environment variable, 
please export the variable as an enviornment variable so it can be read in

ex:
	export %s="secret value"
`, secretName, secretName))
	}
	return secretValue, nil
}

// GithubAuth reads an api token from the environment
// expecting API_GITHUB_TOKEN and authenticates with github api
// and returns an authenticated client.
func GithubAuth() (context.Context, *github.Client, error) {
	token := os.Getenv("API_GITHUB_TOKEN")
	if token == "" {
		return nil, nil, errors.New(`no API_GITHUB_TOKEN was found in the environment variables.
This is needed to authenticate with github. Please generate a github token that has access to do what you are trying to do
and export it as an environment variable.

ex:
	export API_GITHUB_TOKEN="<token contents string>"
`)
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client, nil
}

// AddRepoSecret will add a secret value to a given github repo for a given owner
// It encrypts the secret using sodium before sending the secret to github api
// therefore requires libsodium to be installed on the machine running this code
// https://formulae.brew.sh/formula/libsodium
// Github is very picky over formats things need to be to be sent, and not very descriptive in how
// they will be sent to the user
// To get a secret uploaded you need to get the public key of the repo that will be receiving the
// secret. This key is used to encrypt the secret before transport on the sending side, and then decrypted by github.
// Once you have the public key you need to encrypt the secret with sodiumlib before sending it.
// The public key comes base64 encoded, and sodiumlib expects it to not be base64 encoded so you need
// to decode the public key before using it.
// Once the public key is decoded you need to convert the string secret into bytes.
// once you have the public key decoded, and the secret string in bytes you can encrypt it using
// sodium.CryptoBoxSeal
// That will produce the correctly encrypted secret as bytes, but you need to then convert it to
// a base64 encoded string. After doing that you can use that base64 encoded string as the encrypted value
// to be part of the github.EncodedSecret type.
// The name is the string (no encoding, or base64 needed) of the secret name that will appear in github secrets
// then the KeyID will be the public key of the repo's ID, which is gettable from the public key's GetKeyID method.
// Finally you can pass that object in and have it be created or updated in github.
func AddRepoSecret(ctx context.Context, client *github.Client, owner string, repo string, secretName string, secretValue string) (string, error) {
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return "", err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return "", err
	}

	_, err = client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, encryptedSecret)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Actions.CreateOrUpdateRepoSecret returned error: %v", err))
	}

	return secretName, nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err))
	}

	secretBytes := []byte(secretValue)
	encryptedBytes, exit := sodium.CryptoBoxSeal(secretBytes, decodedPublicKey)
	if exit != 0 {
		return nil, errors.New("sodium.CryptoBoxSeal exited with non zero exit code")
	}

	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}
	return encryptedSecret, nil
}
