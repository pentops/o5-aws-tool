package messaging

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/o5/messaging/v1/messaging

import (
	time "time"
)

// WireEncoding Proto Enum: o5.messaging.v1.WireEncoding
type WireEncoding string

const (
	WireEncoding_UNSPECIFIED WireEncoding = "UNSPECIFIED"
	WireEncoding_PROTOJSON   WireEncoding = "PROTOJSON"
	WireEncoding_RAW         WireEncoding = "RAW"
)

// Message_Request Proto: Message_Request
type Message_Request struct {
	ReplyTo string `json:"replyTo,omitempty"`
}

// Message Proto: Message
type Message struct {
	MessageId        string            `json:"messageId,omitempty"`
	GrpcService      string            `json:"grpcService,omitempty"`
	GrpcMethod       string            `json:"grpcMethod,omitempty"`
	Body             *Any              `json:"body,omitempty"`
	SourceApp        string            `json:"sourceApp,omitempty"`
	SourceEnv        string            `json:"sourceEnv,omitempty"`
	DestinationTopic string            `json:"destinationTopic,omitempty"`
	Timestamp        *time.Time        `json:"timestamp,omitempty"`
	Headers          map[string]string `json:"headers,omitempty"`
	Request          *Message_Request  `json:"request,omitempty"`
	Reply            *Message_Reply    `json:"reply,omitempty"`
}

// Problem_UnhandledError Proto: Problem_UnhandledError
type Problem_UnhandledError struct {
	Error string `json:"error,omitempty"`
}

// DeadMessage Proto: DeadMessage
type DeadMessage struct {
	DeathId    string   `json:"deathId,omitempty"`
	HandlerApp string   `json:"handlerApp,omitempty"`
	HandlerEnv string   `json:"handlerEnv,omitempty"`
	Message    *Message `json:"message,omitempty"`
	Problem    *Problem `json:"problem,omitempty"`
	Infra      *Infra   `json:"infra,omitempty"`
}

// Infra Proto: Infra
type Infra struct {
	Type     string            `json:"type,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Problem Proto Oneof: o5.messaging.v1.Problem
type Problem struct {
	J5TypeKey      string                  `json:"!type,omitempty"`
	UnhandledError *Problem_UnhandledError `json:"unhandledError,omitempty"`
}

func (s Problem) OneofKey() string {
	if s.UnhandledError != nil {
		return "unhandledError"
	}
	return ""
}

func (s Problem) Type() interface{} {
	if s.UnhandledError != nil {
		return s.UnhandledError
	}
	return nil
}

// Any Proto: Any
type Any struct {
	TypeUrl  string `json:"typeUrl,omitempty"`
	Value    []byte `json:"value,omitempty"`
	Encoding string `json:"encoding,omitempty"`
}

// Message_Reply Proto: Message_Reply
type Message_Reply struct {
	ReplyTo string `json:"replyTo,omitempty"`
}
