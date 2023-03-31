package github

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v45/github"
	"github.com/sirupsen/logrus"
	"github.com/tsimer123/pet-infra-yandex-tools/internal/env"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/oauth2"
)

// GitHubActionsService type
type GitHubActionsService interface {
	GetRepoPublicKey(context.Context, string, string) (*github.PublicKey, *github.Response, error)
	CreateOrUpdateRepoSecret(context.Context, string, string, *github.EncryptedSecret) (*github.Response, error)
}

// GitHub type
type GitHub struct {
	Token, Owner, Repo, SecretName string
}

func NewGithub(o *env.Options) *GitHub {
	return &GitHub{
		Token:      o.GithubToken,
		Owner:      o.GithubOwner,
		Repo:       o.GithubRepo,
		SecretName: o.GithubSecretName,
	}
}

func (t *GitHub) UpdateSecret(secretValue string) {
	ctx, client, err := githubAuth(t.Token)
	if err != nil {
		logrus.Fatalf("Unable to authenticate to GitHub: %v", err)
	}

	actionsService := client.Actions
	if err := addRepoSecret(ctx, actionsService, t.Owner, t.Repo, t.SecretName, secretValue); err != nil {
		logrus.Fatalf("Unable to add secret to GitHub: %v", err)
	}
}

// githubAuth returns a GitHub client and context.
func githubAuth(token string) (context.Context, *github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return ctx, client, nil
}

// addRepoSecret will add a secret to a GitHub repo for use in GitHub Actions.
//
// Finally, the secretName and secretValue will determine the name of the secret added and it's corresponding value.
//
// The actual transmission of the secret value to GitHub using the api requires that the secret value is encrypted
// using the public key of the target repo. This encryption is done using x/crypto/nacl/box.
//
// First, the public key of the repo is retrieved. The public key comes base64
// encoded, so it must be decoded prior to use.
//
// Second, the decode key is converted into a fixed size byte array.
//
// Third, the secret value is converted into a slice of bytes.
//
// Fourth, the secret is encrypted with box.SealAnonymous using the repo's decoded public key.
//
// Fifth, the encrypted secret is encoded as a base64 string to be used in a github.EncodedSecret type.
//
// Sixt, The other two properties of the github.EncodedSecret type are determined. The name of the secret to be added
// (string not base64), and the KeyID of the public key used to encrypt the secret.
// This can be retrieved via the public key's GetKeyID method.
//
// Finally, the github.EncodedSecret is passed into the GitHub client.Actions.CreateOrUpdateRepoSecret method to
// populate the secret in the GitHub repo.
func addRepoSecret(ctx context.Context, actionsService GitHubActionsService, owner string, repo, secretName string, secretValue string) error {
	publicKey, _, err := actionsService.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return err
	}

	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretName, secretValue)
	if err != nil {
		return err
	}

	if _, err := actionsService.CreateOrUpdateRepoSecret(ctx, owner, repo, encryptedSecret); err != nil {
		return fmt.Errorf("Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	logrus.Infof("Added GitHub secret: %s to %s/%s", secretName, owner, repo)

	return nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	var boxKey [32]byte
	copy(boxKey[:], decodedPublicKey)
	secretBytes := []byte(secretValue)
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, crypto_rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("box.SealAnonymous failed with error %w", err)
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
