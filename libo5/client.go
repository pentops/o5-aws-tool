package libo5

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/pentops/log.go/log"
)

type API struct {
	BaseURL    string
	HTTPClient *http.Client
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
	res, err := c.HTTPClient.Do(httpReq.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *API) Request(ctx context.Context, method, path string, req, res interface{}) error {
	var reqBody io.Reader
	if req != nil {
		reqBytes, err := json.Marshal(req)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(reqBytes)
	}

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
		log.WithFields(ctx, map[string]interface{}{
			"status": httpRes.StatusCode,
			"body":   string(resBody),
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
		fmt.Printf(string(resBody))
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
