package libo5

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pentops/log.go/log"
	"github.com/pentops/o5-aws-tool/gen/j5/realm/v1/realm"
)

type APIConfig struct {
	BaseURL       string `env:"O5_API"`
	BearerToken   string `env:"O5_BEARER" required:"false"`
	ClientID      string `env:"O5_CLIENT_ID" required:"false"`
	ClientSecret  string `env:"O5_CLIENT_SECRET" required:"false"`
	TokenEndpoint string `env:"O5_TOKEN_ENDPOINT" required:"false"`
}

func (c APIConfig) APIClient() *API {
	client := NewAPI(c.BaseURL)
	tokenEndpoint := c.TokenEndpoint
	if tokenEndpoint == "" {
		tokenEndpoint = c.BaseURL + "/realm-auth/v1/token"
	}

	if c.BearerToken != "" {
		client.Auth = BearerToken(c.BearerToken)
	} else if c.ClientID != "" || c.ClientSecret != "" {
		client.Auth = &APIKeyAuth{
			Key:      c.ClientID,
			Secret:   c.ClientSecret,
			Client:   client.HTTPClient,
			Endpoint: tokenEndpoint,
		}
	}
	return client
}

type API struct {
	BaseURL    string
	HTTPClient *http.Client
	Auth       AuthProvider
}

func NewAPI(baseURL string) *API {
	return &API{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
}

type AuthProvider interface {
	Authenticate(req *http.Request) error
}

type BearerToken string

func (t BearerToken) Authenticate(req *http.Request) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
	return nil
}

type APIKeyAuth struct {
	Key      string
	Secret   string
	Endpoint string

	Client *http.Client

	token string
}

func (a *APIKeyAuth) Authenticate(req *http.Request) error {
	fmt.Printf("Authenticating with APIKeyAuth: %s\n", a.Key)
	if a.token == "" {
		ctx := req.Context()
		token, err := a.getToken(ctx)
		if err != nil {
			return fmt.Errorf("getting token: %w", err)
		}
		a.token = token
	}
	fmt.Printf("Using token: %s\n", a.token)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))
	return nil
}

func (a *APIKeyAuth) getToken(ctx context.Context) (string, error) {
	ctx = log.WithFields(ctx, "auth", "APIKeyAuth")
	ctx = log.WithFields(ctx, map[string]any{
		"method": "POST",
		"path":   a.Endpoint,
	})
	reqBytes, err := json.Marshal(&realm.GetTokenRequest{
		ClientId:     a.Key,
		ClientSecret: a.Secret,
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}
	reqBody := bytes.NewReader(reqBytes)
	httpReq, err := http.NewRequest("POST", a.Endpoint, reqBody)
	if err != nil {
		return "", err
	}

	resp := &realm.GetTokenResponse{}

	err = httpRoundTrip(ctx, a.Client, httpReq, resp)
	if err != nil {
		return "", fmt.Errorf("round trip: %w", err)
	}
	return resp.Jwt, nil
}

func debugToken(authHeader string) {
	split := strings.SplitN(authHeader, " ", 2)
	if len(split) != 2 || split[0] != "Bearer" {
		fmt.Printf("Invalid Authorization header format. %q\n", authHeader)
		return
	}
	token := split[1]
	parts := strings.Split(token, ".")
	if len(parts) == 3 {
		dec, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err == nil {
			fmt.Printf("\n====\n%s\n====\n", string(dec))
		}
		dec, err = base64.RawURLEncoding.DecodeString(parts[0])
		if err == nil {
			fmt.Printf("\n====\n%s\n====\n", string(dec))
		}
	}
}

func (c *API) Request(ctx context.Context, method, path string, req, res any) error {
	var reqBody io.Reader
	ctx = log.WithFields(ctx, map[string]any{
		"method":  method,
		"path":    path,
		"baseURL": c.BaseURL,
	})

	if req != nil {
		reqBytes, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(reqBytes)
		ctx = log.WithField(ctx, "body", string(reqBytes))
	}

	log.Debug(ctx, "request")

	httpReq, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return err
	}

	if c.Auth != nil {
		err = c.Auth.Authenticate(httpReq)
		if err != nil {
			return fmt.Errorf("authenticating: %w", err)
		}
	} else {
		log.Warn(ctx, "no auth provider set, request will not be authenticated")
	}

	return httpRoundTrip(ctx, c.HTTPClient, httpReq, res)
}

func httpRoundTrip(ctx context.Context, client *http.Client, httpReq *http.Request, res any) error {

	httpRes, err := client.Do(httpReq.WithContext(ctx))
	if err != nil {
		return err
	}

	defer httpRes.Body.Close()

	resBody, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	if httpRes.StatusCode != http.StatusOK {
		if httpRes.StatusCode == http.StatusUnauthorized {
			log.WithFields(ctx,
				"status", httpRes.StatusCode,
				"body", string(resBody),
			).Error("unexpected status")
			debugToken(httpReq.Header.Get("Authorization"))
			return fmt.Errorf("unauthorized - check $O5_BEARER")
		}

		log.WithFields(ctx,
			"status", httpRes.StatusCode,
			"body", string(resBody),
			"headers", httpRes.Header,
		).Error("unexpected status")
		return fmt.Errorf("unexpected status %d", httpRes.StatusCode)
	}

	dec := json.NewDecoder(bytes.NewReader(resBody))
	dec.DisallowUnknownFields()
	if err := dec.Decode(res); err != nil {
		log.WithFields(ctx, map[string]any{
			"status": httpRes.StatusCode,
			"error":  err,
		}).Error("bad JSON response")
		fmt.Print("\n====\n" + string(resBody) + "\n====\n")
		return fmt.Errorf("marshal response: %w", err)
	}

	return nil
}

var ErrStopPaging = fmt.Errorf("stop paging")

type PageRequest interface {
	SetPageToken(string)
}
type PageResponse[Item any] interface {
	GetPageToken() *string
	GetItems() []Item
}

func Paged[
	Req PageRequest,
	Res PageResponse[Item],
	Item any,
](ctx context.Context, baseReq Req, call func(context.Context, Req) (Res, error), callback func(Item) error) error {

	for {
		res, err := call(ctx, baseReq)
		if err != nil {
			return err
		}

		for _, item := range res.GetItems() {
			if err := callback(item); err != nil {
				if errors.Is(err, ErrStopPaging) {
					return nil
				}
				return err
			}
		}

		resToken := res.GetPageToken()
		if resToken == nil {
			return nil
		}

		baseReq.SetPageToken(*resToken)
	}
}

func PagedAll[
	Req PageRequest,
	Res PageResponse[Item],
	Item any,
](ctx context.Context, baseReq Req, call func(context.Context, Req) (Res, error)) ([]Item, error) {
	all := []Item{}
	err := Paged(ctx, baseReq, call, func(item Item) error {
		all = append(all, item)
		return nil
	})
	return all, err
}
