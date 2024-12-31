package messages

import (
	"fmt"

	"github.com/pentops/o5-aws-tool/gen/o5/messaging/v1/messaging"
)

func PrintMessageMeta(message *messaging.Message) {

	fmt.Printf(" Source App:   %s\n", message.SourceApp)
	fmt.Printf(" Source Env:   %s\n", message.SourceEnv)
	fmt.Printf(" GRPC Method:  %s\n", message.GrpcMethod)
	fmt.Printf(" GRPC Service: %s\n", message.GrpcService)
	fmt.Printf(" Message ID:   %s\n", message.MessageId)
	fmt.Printf(" Dest Topic:   %s\n", message.DestinationTopic)
	fmt.Printf(" Timestamp:    %s\n", message.Timestamp)
	fmt.Printf(" Headers: \n")

	for key, val := range message.Headers {
		fmt.Printf("  %s: %s\n", key, val)
	}
	if message.Request != nil {
		fmt.Printf("Request ReplyTo: %s\n", message.Request.ReplyTo)
	}
	if message.Reply != nil {
		fmt.Printf("Reply ReplyTo: %s\n", message.Reply.ReplyTo)
	}

}
