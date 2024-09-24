package registry

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/j5/registry/registry/v1/registry

import (
	context "context"
	errors "errors"
	client "github.com/pentops/o5-aws-tool/gen/j5/client/v1/client"
	url "net/url"
	strings "strings"
)

type Requester interface {
	Request(ctx context.Context, method string, path string, body interface{}, response interface{}) error
}

// DownloadService
type DownloadService struct {
	Requester
}

func NewDownloadService(requester Requester) *DownloadService {
	return &DownloadService{
		Requester: requester,
	}
}

func (s DownloadService) DownloadImage(ctx context.Context, req *DownloadImageRequest) (*DownloadImageResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "registry"
	pathParts[2] = "v1"
	if req.Owner == "" {
		return nil, errors.New("required field \"Owner\" not set")
	}
	pathParts[3] = req.Owner
	if req.Name == "" {
		return nil, errors.New("required field \"Name\" not set")
	}
	pathParts[4] = req.Name
	if req.Version == "" {
		return nil, errors.New("required field \"Version\" not set")
	}
	pathParts[5] = req.Version
	pathParts[6] = "image.bin"
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &DownloadImageResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DownloadService) DownloadSwagger(ctx context.Context, req *DownloadSwaggerRequest) (*DownloadSwaggerResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "registry"
	pathParts[2] = "v1"
	if req.Owner == "" {
		return nil, errors.New("required field \"Owner\" not set")
	}
	pathParts[3] = req.Owner
	if req.Name == "" {
		return nil, errors.New("required field \"Name\" not set")
	}
	pathParts[4] = req.Name
	if req.Version == "" {
		return nil, errors.New("required field \"Version\" not set")
	}
	pathParts[5] = req.Version
	pathParts[6] = "swagger.json"
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &DownloadSwaggerResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DownloadService) DownloadClientAPI(ctx context.Context, req *DownloadClientAPIRequest) (*DownloadClientAPIResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "registry"
	pathParts[2] = "v1"
	if req.Owner == "" {
		return nil, errors.New("required field \"Owner\" not set")
	}
	pathParts[3] = req.Owner
	if req.Name == "" {
		return nil, errors.New("required field \"Name\" not set")
	}
	pathParts[4] = req.Name
	if req.Version == "" {
		return nil, errors.New("required field \"Version\" not set")
	}
	pathParts[5] = req.Version
	pathParts[6] = "api.json"
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &DownloadClientAPIResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DownloadImageRequest
type DownloadImageRequest struct {
	Owner   string `json:"-" path:"owner"`
	Name    string `json:"-" path:"name"`
	Version string `json:"-" path:"version"`
}

func (s DownloadImageRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	return values, nil
}

// DownloadImageResponse
type DownloadImageResponse struct {
}

// DownloadSwaggerRequest
type DownloadSwaggerRequest struct {
	Owner   string `json:"-" path:"owner"`
	Name    string `json:"-" path:"name"`
	Version string `path:"version" json:"-"`
}

func (s DownloadSwaggerRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	return values, nil
}

// DownloadSwaggerResponse
type DownloadSwaggerResponse struct {
}

// DownloadClientAPIRequest
type DownloadClientAPIRequest struct {
	Owner   string `json:"-" path:"owner"`
	Name    string `path:"name" json:"-"`
	Version string `json:"-" path:"version"`
}

func (s DownloadClientAPIRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	return values, nil
}

// DownloadClientAPIResponse
type DownloadClientAPIResponse struct {
	Version string      `json:"version,omitempty"`
	Api     *client.API `json:"api,omitempty"`
}

// CombinedClient
type CombinedClient struct {
	*DownloadService
}

func NewCombinedClient(requester Requester) *CombinedClient {
	return &CombinedClient{
		DownloadService: NewDownloadService(requester),
	}
}