package messaging

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/libo5/o5/messaging/v1/messaging

import (
	time "time"
)

// Message Proto: o5.messaging.v1.Message
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

// Any Proto: o5.messaging.v1.Any
type Any struct {
	TypeUrl  string `json:"typeUrl,omitempty"`
	Value    []byte `json:"value,omitempty"`
	Encoding string `json:"encoding,omitempty"`
}

// Message_Request Proto: o5.messaging.v1.Message.Request
type Message_Request struct {
	ReplyTo string `json:"replyTo,omitempty"`
}

// Message_Reply Proto: o5.messaging.v1.Message.Reply
type Message_Reply struct {
	ReplyTo string `json:"replyTo,omitempty"`
}

// Problem Proto: o5.messaging.v1.topic.Problem
type Problem struct {
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

// Problem_UnhandledError Proto: o5.messaging.v1.topic.Problem.UnhandledError
type Problem_UnhandledError struct {
	Error string `json:"error,omitempty"`
}

// Infra Proto: o5.messaging.v1.topic.Infra
type Infra struct {
	Type     string            `json:"type,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DeadMessage Proto: o5.messaging.v1.topic.DeadMessage
type DeadMessage struct {
	DeathId    string   `json:"deathId,omitempty"`
	HandlerApp string   `json:"handlerApp,omitempty"`
	HandlerEnv string   `json:"handlerEnv,omitempty"`
	Message    *Message `json:"message,omitempty"`
	Problem    *Problem `json:"problem,omitempty"`
	Infra      *Infra   `json:"infra,omitempty"`
}
