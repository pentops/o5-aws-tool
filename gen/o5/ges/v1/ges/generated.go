package ges

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/o5/ges/v1/ges

import (
	context "context"
	json "encoding/json"
	url "net/url"
	strings "strings"
	time "time"

	list "github.com/pentops/o5-aws-tool/gen/j5/list/v1/list"
	state "github.com/pentops/o5-aws-tool/gen/j5/state/v1/state"
)

type Requester interface {
	Request(ctx context.Context, method string, path string, body interface{}, response interface{}) error
}

// CommandService
type CommandService struct {
	Requester
}

func NewCommandService(requester Requester) *CommandService {
	return &CommandService{
		Requester: requester,
	}
}

func (s CommandService) ReplayEvents(ctx context.Context, req *ReplayEventsRequest) (*ReplayEventsResponse, error) {
	pathParts := make([]string, 5)
	pathParts[0] = ""
	pathParts[1] = "ges"
	pathParts[2] = "v1"
	pathParts[3] = "events"
	pathParts[4] = "replay"
	path := strings.Join(pathParts, "/")
	resp := &ReplayEventsResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s CommandService) ReplayUpserts(ctx context.Context, req *ReplayUpsertsRequest) (*ReplayUpsertsResponse, error) {
	pathParts := make([]string, 5)
	pathParts[0] = ""
	pathParts[1] = "ges"
	pathParts[2] = "v1"
	pathParts[3] = "upserts"
	pathParts[4] = "replay"
	path := strings.Join(pathParts, "/")
	resp := &ReplayUpsertsResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ReplayEventsRequest
type ReplayEventsRequest struct {
	QueueUrl    string `json:"queueUrl"`
	GrpcService string `json:"grpcService"`
	GrpcMethod  string `json:"grpcMethod"`
}

// ReplayEventsResponse
type ReplayEventsResponse struct {
}

// ReplayUpsertsRequest
type ReplayUpsertsRequest struct {
	QueueUrl    string `json:"queueUrl"`
	GrpcService string `json:"grpcService"`
	GrpcMethod  string `json:"grpcMethod"`
}

// ReplayUpsertsResponse
type ReplayUpsertsResponse struct {
}

// QueryService
type QueryService struct {
	Requester
}

func NewQueryService(requester Requester) *QueryService {
	return &QueryService{
		Requester: requester,
	}
}

func (s QueryService) EventsList(ctx context.Context, req *EventsListRequest) (*EventsListResponse, error) {
	pathParts := make([]string, 4)
	pathParts[0] = ""
	pathParts[1] = "ges"
	pathParts[2] = "v1"
	pathParts[3] = "events"
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &EventsListResponse{}
	err := s.Request(ctx, "GET", path, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s QueryService) UpsertList(ctx context.Context, req *UpsertListRequest) (*UpsertListResponse, error) {
	pathParts := make([]string, 4)
	pathParts[0] = ""
	pathParts[1] = "ges"
	pathParts[2] = "v1"
	pathParts[3] = "upsert"
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &UpsertListResponse{}
	err := s.Request(ctx, "GET", path, nil, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// EventsListRequest
type EventsListRequest struct {
	Page  *list.PageRequest  `json:"-" query:"page"`
	Query *list.QueryRequest `json:"-" query:"query"`
}

func (s *EventsListRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

func (s EventsListRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	if s.Page != nil {
		bb, err := json.Marshal(s.Page)
		if err != nil {
			return nil, err
		}
		values.Set("page", string(bb))
	}
	if s.Query != nil {
		bb, err := json.Marshal(s.Query)
		if err != nil {
			return nil, err
		}
		values.Set("query", string(bb))
	}
	return values, nil
}

// EventsListResponse
type EventsListResponse struct {
	Page   *list.PageResponse `json:"page,omitempty"`
	Events []*Event           `json:"events,omitempty"`
}

func (s EventsListResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s EventsListResponse) GetItems() []*Event {
	return s.Events
}

// UpsertListRequest
type UpsertListRequest struct {
	Page  *list.PageRequest  `json:"-" query:"page"`
	Query *list.QueryRequest `json:"-" query:"query"`
}

func (s *UpsertListRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

func (s UpsertListRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	if s.Page != nil {
		bb, err := json.Marshal(s.Page)
		if err != nil {
			return nil, err
		}
		values.Set("page", string(bb))
	}
	if s.Query != nil {
		bb, err := json.Marshal(s.Query)
		if err != nil {
			return nil, err
		}
		values.Set("query", string(bb))
	}
	return values, nil
}

// UpsertListResponse
type UpsertListResponse struct {
	Page   *list.PageResponse `json:"page,omitempty"`
	Events []*Upsert          `json:"events,omitempty"`
}

func (s UpsertListResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s UpsertListResponse) GetItems() []*Upsert {
	return s.Events
}

// Event Proto: Event
type Event struct {
	EntityName   string                      `json:"entityName"`
	Metadata     *state.EventPublishMetadata `json:"metadata"`
	GrpcMethod   string                      `json:"grpcMethod"`
	GrpcService  string                      `json:"grpcService"`
	BodyType     string                      `json:"bodyType"`
	EventType    string                      `json:"eventType"`
	EntityKeys   interface{}                 `json:"entityKeys"`
	EventData    interface{}                 `json:"eventData"`
	EntityState  interface{}                 `json:"entityState"`
	EntityStatus string                      `json:"entityStatus"`
}

// Upsert Proto: Upsert
type Upsert struct {
	EntityName         string      `json:"entityName"`
	EntityId           string      `json:"entityId"`
	GrpcMethod         string      `json:"grpcMethod"`
	GrpcService        string      `json:"grpcService"`
	LastEventId        string      `json:"lastEventId"`
	LastEventTimestamp *time.Time  `json:"lastEventTimestamp"`
	Data               interface{} `json:"data"`
}

// CombinedClient
type CombinedClient struct {
	*CommandService
	*QueryService
}

func NewCombinedClient(requester Requester) *CombinedClient {
	return &CombinedClient{
		CommandService: NewCommandService(requester),
		QueryService:   NewQueryService(requester),
	}
}
