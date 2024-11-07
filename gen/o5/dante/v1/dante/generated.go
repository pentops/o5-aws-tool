package dante

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/o5/dante/v1/dante

import (
	context "context"
	json "encoding/json"
	errors "errors"
	list "github.com/pentops/o5-aws-tool/gen/j5/list/v1/list"
	state "github.com/pentops/o5-aws-tool/gen/j5/state/v1/state"
	messaging "github.com/pentops/o5-aws-tool/gen/o5/messaging/v1/messaging"
	url "net/url"
	strings "strings"
)

type Requester interface {
	Request(ctx context.Context, method string, path string, body interface{}, response interface{}) error
}

// DeadMessageCommandService
type DeadMessageCommandService struct {
	Requester
}

func NewDeadMessageCommandService(requester Requester) *DeadMessageCommandService {
	return &DeadMessageCommandService{
		Requester: requester,
	}
}

func (s DeadMessageCommandService) UpdateDeadMessage(ctx context.Context, req *UpdateDeadMessageRequest) (*UpdateDeadMessageResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "dante"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "messages"
	if req.MessageId == "" {
		return nil, errors.New("required field \"MessageId\" not set")
	}
	pathParts[5] = req.MessageId
	pathParts[6] = "update"
	path := strings.Join(pathParts, "/")
	resp := &UpdateDeadMessageResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeadMessageCommandService) ReplayDeadMessage(ctx context.Context, req *ReplayDeadMessageRequest) (*ReplayDeadMessageResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "dante"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "messages"
	if req.MessageId == "" {
		return nil, errors.New("required field \"MessageId\" not set")
	}
	pathParts[5] = req.MessageId
	pathParts[6] = "replay"
	path := strings.Join(pathParts, "/")
	resp := &ReplayDeadMessageResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeadMessageCommandService) RejectDeadMessage(ctx context.Context, req *RejectDeadMessageRequest) (*RejectDeadMessageResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "dante"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "messages"
	if req.MessageId == "" {
		return nil, errors.New("required field \"MessageId\" not set")
	}
	pathParts[5] = req.MessageId
	pathParts[6] = "shelve"
	path := strings.Join(pathParts, "/")
	resp := &RejectDeadMessageResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateDeadMessageRequest
type UpdateDeadMessageRequest struct {
	ReplacesVersionId string              `json:"replacesVersionId,omitempty"`
	VersionId         *string             `json:"versionId,omitempty"`
	Message           *DeadMessageVersion `json:"message"`
	MessageId         string              `json:"-" path:"messageId"`
}

// UpdateDeadMessageResponse
type UpdateDeadMessageResponse struct {
	Message *DeadMessageState `json:"message"`
}

// ReplayDeadMessageRequest
type ReplayDeadMessageRequest struct {
	MessageId string `json:"-" path:"messageId"`
}

// ReplayDeadMessageResponse
type ReplayDeadMessageResponse struct {
	Message *DeadMessageState `json:"message"`
}

// RejectDeadMessageRequest
type RejectDeadMessageRequest struct {
	Reason    string `json:"reason,omitempty"`
	MessageId string `json:"-" path:"messageId"`
}

// RejectDeadMessageResponse
type RejectDeadMessageResponse struct {
	Message *DeadMessageState `json:"message"`
}

// DeadMessageQueryService
type DeadMessageQueryService struct {
	Requester
}

func NewDeadMessageQueryService(requester Requester) *DeadMessageQueryService {
	return &DeadMessageQueryService{
		Requester: requester,
	}
}

func (s DeadMessageQueryService) GetDeadMessage(ctx context.Context, req *GetDeadMessageRequest) (*GetDeadMessageResponse, error) {
	pathParts := make([]string, 6)
	pathParts[0] = ""
	pathParts[1] = "dante"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "message"
	if req.MessageId == nil || *req.MessageId == "" {
		return nil, errors.New("required field \"MessageId\" not set")
	}
	pathParts[5] = *req.MessageId
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &GetDeadMessageResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeadMessageQueryService) ListDeadMessages(ctx context.Context, req *ListDeadMessagesRequest) (*ListDeadMessagesResponse, error) {
	pathParts := make([]string, 5)
	pathParts[0] = ""
	pathParts[1] = "dante"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "messages"
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &ListDeadMessagesResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeadMessageQueryService) ListDeadMessageEvents(ctx context.Context, req *ListDeadMessageEventsRequest) (*ListDeadMessageEventsResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "dante"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "message"
	if req.MessageId == "" {
		return nil, errors.New("required field \"MessageId\" not set")
	}
	pathParts[5] = req.MessageId
	pathParts[6] = "events"
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &ListDeadMessageEventsResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetDeadMessageRequest
type GetDeadMessageRequest struct {
	MessageId *string `json:"-" path:"messageId"`
}

func (s GetDeadMessageRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	return values, nil
}

// GetDeadMessageResponse
type GetDeadMessageResponse struct {
	Message *DeadMessageState   `json:"message"`
	Events  []*DeadMessageEvent `json:"events,omitempty"`
}

// ListDeadMessagesRequest
type ListDeadMessagesRequest struct {
	Page  *list.PageRequest  `json:"-" query:"page"`
	Query *list.QueryRequest `json:"-" query:"query"`
}

func (s *ListDeadMessagesRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

func (s ListDeadMessagesRequest) QueryParameters() (url.Values, error) {
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

// ListDeadMessagesResponse
type ListDeadMessagesResponse struct {
	Messages []*DeadMessageState `json:"messages,omitempty"`
	Page     *list.PageResponse  `json:"page,omitempty"`
}

func (s ListDeadMessagesResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListDeadMessagesResponse) GetItems() []*DeadMessageState {
	return s.Messages
}

// ListDeadMessageEventsRequest
type ListDeadMessageEventsRequest struct {
	MessageId string             `json:"-" path:"messageId"`
	Page      *list.PageRequest  `json:"-" query:"page"`
	Query     *list.QueryRequest `query:"query" json:"-"`
}

func (s *ListDeadMessageEventsRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

func (s ListDeadMessageEventsRequest) QueryParameters() (url.Values, error) {
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

// ListDeadMessageEventsResponse
type ListDeadMessageEventsResponse struct {
	Events []*DeadMessageEvent `json:"events,omitempty"`
	Page   *list.PageResponse  `json:"page,omitempty"`
}

func (s ListDeadMessageEventsResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListDeadMessageEventsResponse) GetItems() []*DeadMessageEvent {
	return s.Events
}

// DeadMessageEventType_Updated Proto: DeadMessageEventType_Updated
type DeadMessageEventType_Updated struct {
	Spec *DeadMessageVersion `json:"spec,omitempty"`
}

// DeadMessageVersion Proto: DeadMessageVersion
type DeadMessageVersion struct {
	VersionId  string                         `json:"versionId,omitempty"`
	Message    *messaging.Message             `json:"message,omitempty"`
	SqsMessage *DeadMessageVersion_SQSMessage `json:"sqsMessage,omitempty"`
}

// DeadMessageState Proto: DeadMessageState
type DeadMessageState struct {
	Metadata  *state.StateMetadata `json:"metadata"`
	MessageId string               `json:"messageId,omitempty"`
	Status    string               `json:"status,omitempty"`
	Data      *DeadMessageData     `json:"data,omitempty"`
}

// DeadMessageKeys Proto: DeadMessageKeys
type DeadMessageKeys struct {
	MessageId string `json:"messageId,omitempty"`
}

// MessageStatus Proto Enum: o5.dante.v1.MessageStatus
type MessageStatus string

const (
	MessageStatus_UNSPECIFIED MessageStatus = "UNSPECIFIED"
	MessageStatus_CREATED     MessageStatus = "CREATED"
	MessageStatus_UPDATED     MessageStatus = "UPDATED"
	MessageStatus_REPLAYED    MessageStatus = "REPLAYED"
	MessageStatus_REJECTED    MessageStatus = "REJECTED"
)

// DeadMessageEventType_Replayed Proto: DeadMessageEventType_Replayed
type DeadMessageEventType_Replayed struct {
}

// DeadMessageEventType Proto Oneof: o5.dante.v1.DeadMessageEventType
type DeadMessageEventType struct {
	J5TypeKey string                         `json:"!type,omitempty"`
	Notified  *DeadMessageEventType_Notified `json:"notified,omitempty"`
	Updated   *DeadMessageEventType_Updated  `json:"updated,omitempty"`
	Replayed  *DeadMessageEventType_Replayed `json:"replayed,omitempty"`
	Rejected  *DeadMessageEventType_Rejected `json:"rejected,omitempty"`
}

func (s DeadMessageEventType) OneofKey() string {
	if s.Notified != nil {
		return "notified"
	}
	if s.Updated != nil {
		return "updated"
	}
	if s.Replayed != nil {
		return "replayed"
	}
	if s.Rejected != nil {
		return "rejected"
	}
	return ""
}

func (s DeadMessageEventType) Type() interface{} {
	if s.Notified != nil {
		return s.Notified
	}
	if s.Updated != nil {
		return s.Updated
	}
	if s.Replayed != nil {
		return s.Replayed
	}
	if s.Rejected != nil {
		return s.Rejected
	}
	return nil
}

// DeadMessageEventType_Notified Proto: DeadMessageEventType_Notified
type DeadMessageEventType_Notified struct {
	Notification *messaging.DeadMessage `json:"notification,omitempty"`
}

// DeadMessageVersion_SQSMessage Proto: DeadMessageVersion_SQSMessage
type DeadMessageVersion_SQSMessage struct {
	QueueUrl   string            `json:"queueUrl,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

// DeadMessageEvent Proto: DeadMessageEvent
type DeadMessageEvent struct {
	Metadata  *state.EventMetadata  `json:"metadata"`
	MessageId string                `json:"messageId,omitempty"`
	Event     *DeadMessageEventType `json:"event"`
}

// DeadMessageEventType_Rejected Proto: DeadMessageEventType_Rejected
type DeadMessageEventType_Rejected struct {
	Reason string `json:"reason,omitempty"`
}

// DeadMessageData Proto: DeadMessageData
type DeadMessageData struct {
	Notification   *messaging.DeadMessage `json:"notification,omitempty"`
	CurrentVersion *DeadMessageVersion    `json:"currentVersion,omitempty"`
}

// CombinedClient
type CombinedClient struct {
	*DeadMessageCommandService
	*DeadMessageQueryService
}

func NewCombinedClient(requester Requester) *CombinedClient {
	return &CombinedClient{
		DeadMessageCommandService: NewDeadMessageCommandService(requester),
		DeadMessageQueryService:   NewDeadMessageQueryService(requester),
	}
}
