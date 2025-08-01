package libo5

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type IAMAuth struct {
	Endpoint string
	Realm    string
	token    string
	aws      aws.Config
}

func NewIAMAuth(endpoint string, realm string, awsConfig aws.Config) (*IAMAuth, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("IAM endpoint cannot be empty")
	}

	auth := &IAMAuth{
		Endpoint: endpoint,
		Realm:    realm,
		aws:      awsConfig,
	}

	return auth, nil
}

func (a *IAMAuth) fetchNewToken(ctx context.Context) error {

	signer := v4.NewSigner()

	bodyBytes := []byte{}

	endpoint := a.Endpoint
	qs := url.Values{}
	qs.Set("realm", a.Realm)
	endpoint += "?" + qs.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	bodyHash := sha256.New()
	bodyHash.Write(bodyBytes)
	bodyHashSum := bodyHash.Sum(nil)

	cred, err := a.aws.Credentials.Retrieve(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	service := "lambda"
	err = signer.SignHTTP(ctx, cred, req, hex.EncodeToString(bodyHashSum), service, a.aws.Region, time.Now())
	if err != nil {
		return fmt.Errorf("failed to sign request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, string(body))
		return fmt.Errorf("request failed with status code %d", res.StatusCode)
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}
	var tokenRes tokenResponse
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	a.token = tokenRes.Token

	return nil
}

func (a *IAMAuth) Token(ctx context.Context) (string, error) {

	if a.token == "" {
		if err := a.fetchNewToken(ctx); err != nil {
			return "", fmt.Errorf("failed to fetch new token: %w", err)

		}
	}
	if a.token == "" {
		return "", fmt.Errorf("token is empty after fetching")
	}
	return a.token, nil
}

func (a *IAMAuth) Authenticate(req *http.Request) error {
	token, err := a.Token(req.Context())
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}
