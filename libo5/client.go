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
)

type APIConfig struct {
	BaseURL     string `env:"O5_API"`
	BearerToken string `env:"O5_BEARER"`
}

func (c APIConfig) APIClient() *API {
	client := NewAPI(c.BaseURL)
	client.BearerToken = c.BearerToken
	return client
}

type API struct {
	BaseURL     string
	HTTPClient  *http.Client
	BearerToken string
}

func NewAPI(baseURL string) *API {
	return &API{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
	}
}

func (c *API) Do(ctx context.Context, method, path string, req io.Reader) (*http.Response, error) {
	httpReq, err := http.NewRequest(method, c.BaseURL+path, req)
	if err != nil {
		return nil, err
	}
	if c.BearerToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.BearerToken)
	}
	res, err := c.HTTPClient.Do(httpReq.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *API) Request(ctx context.Context, method, path string, req, res interface{}) error {
	var reqBody io.Reader
	ctx = log.WithFields(ctx, map[string]interface{}{
		"method": method,
		"path":   path,
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

	httpRes, err := c.Do(ctx, method, path, reqBody)
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
			log.WithFields(ctx, map[string]interface{}{
				"status": httpRes.StatusCode,
				"body":   string(resBody),
			}).Error("unexpected status")
			if c.BearerToken == "" {
				return fmt.Errorf("unauthorized - try setting $O5_BEARER")
			} else {
				parts := strings.Split(c.BearerToken, ".")
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
				return fmt.Errorf("unauthorized - check $O5_BEARER")
			}
		}

		log.WithFields(ctx, map[string]interface{}{
			"status":  httpRes.StatusCode,
			"body":    string(resBody),
			"headers": httpRes.Header,
		}).Error("unexpected status")
		return fmt.Errorf("unexpected status %d", httpRes.StatusCode)
	}

	dec := json.NewDecoder(bytes.NewReader(resBody))
	dec.DisallowUnknownFields()
	if err := dec.Decode(res); err != nil {
		log.WithFields(ctx, map[string]interface{}{
			"status": httpRes.StatusCode,
			"error":  err,
		}).Error("bad JSON")
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
