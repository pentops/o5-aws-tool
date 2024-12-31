package dante

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pentops/o5-aws-tool/cli"
	"github.com/pentops/o5-aws-tool/gen/o5/dante/v1/dante"
	"github.com/pentops/o5-aws-tool/gen/o5/messaging/v1/messaging"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/runner/commander"
)

func DanteCommandSet() *commander.CommandSet {
	remoteGroup := commander.NewCommandSet()

	remoteGroup.Add("ls", commander.NewCommand(runDanteLs))
	remoteGroup.Add("reject", commander.NewCommand(runDanteReject))

	return remoteGroup
}

var errExitMessage = fmt.Errorf("exit message")
var errExitLoop = fmt.Errorf("exit loop")

func prettyPrintJSON(data []byte, truncate bool) {

	var body string

	buff := &bytes.Buffer{}
	if err := json.Indent(buff, data, "  | ", "  "); err != nil {
		fmt.Printf("Error indenting JSON body : %s\n", err)
		body = string(data)
	} else {
		body = buff.String()
	}

	if truncate && len(body) > 250 {
		body = body[:250] + "..."
	}
	fmt.Printf("  | %s\n", body)
}

func printMetadata(message *messaging.Message) {

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

func runDanteLs(ctx context.Context, cfg struct {
	libo5.APIConfig
	Interactive bool `flag:"i" help:"Interactive mode"`
}) error {
	client := cfg.APIClient()

	queryClient := dante.NewDeadMessageQueryService(client)
	commandClient := dante.NewDeadMessageCommandService(client)
	printMessage := func(state *dante.DeadMessageState) {
		fmt.Printf("DL: %s\n", state.MessageId)
		msg := state.Data.CurrentVersion.Message
		fmt.Printf("  Method: /%s/%s\n", msg.GrpcService, msg.GrpcMethod)
		fmt.Printf("  Handler Env: %s\n", state.Data.Notification.HandlerEnv)
		fmt.Printf("  Handler App: %s\n", state.Data.Notification.HandlerApp)
		prettyPrintJSON(msg.Body.Value, true)
		if state.Data.Notification.Problem.UnhandledError != nil {
			fmt.Printf("  Error: %s", state.Data.Notification.Problem.UnhandledError.Error)
		}
	}
	runInteraction := func(ctx context.Context, state *dante.DeadMessageState) error {
		hasMods := false
		newVersion := &dante.DeadMessageVersion{}

		deleteCommand := &cli.Command{
			Name:    "delete",
			Short:   "d",
			Summary: "Mark the message as Rejected (deleted)",
			Run: func() error {
				res, err := commandClient.RejectDeadMessage(ctx, &dante.RejectDeadMessageRequest{
					MessageId: state.MessageId,
				})
				if err != nil {
					return fmt.Errorf("reject dead messages: %w", err)
				}
				cli.Print("DL", res.Message.Status)
				return errExitMessage
			},
		}

		replayCommand := &cli.Command{
			Name:    "replay",
			Short:   "r",
			Summary: "Replay the message",
			Run: func() error {
				res, err := commandClient.ReplayDeadMessage(ctx, &dante.ReplayDeadMessageRequest{
					MessageId: state.MessageId,
				})
				if err != nil {
					return fmt.Errorf("replay dead messages: %w", err)
				}
				cli.Print("DL", res.Message.Status)
				return errExitMessage
			},
		}

		nextCommand := &cli.Command{
			Name:    "next",
			Short:   "n",
			Summary: "Ignore this message and move on",
			Run: func() error {
				return errExitMessage
			},
		}

		saveCommand := &cli.Command{
			Name:    "save",
			Short:   "s",
			Summary: "Saves the changes made to the message",
			Run: func() error {
				if !hasMods {
					fmt.Printf("No changes to save\n")

					return nil
				}
				newID := uuid.NewString()
				newVersion.VersionId = newID

				res, err := commandClient.UpdateDeadMessage(ctx, &dante.UpdateDeadMessageRequest{
					MessageId:         state.MessageId,
					ReplacesVersionId: state.Data.CurrentVersion.VersionId,
					VersionId:         &newID,
					Message:           newVersion,
				})
				if err != nil {
					return fmt.Errorf("update dead messages: %w", err)
				}
				fmt.Printf("Saved, new state: %s\n", res.Message.Status)
				printMessage(res.Message)
				hasMods = false
				return nil
			},
		}

		discardCommand := &cli.Command{
			Name:    "discard",
			Summary: "Discards the changes made to the message",
			Run: func() error {
				if !hasMods {
					fmt.Printf("No changes to discard\n")
					return nil
				}
				hasMods = false
				newVersion = &dante.DeadMessageVersion{}
				fmt.Printf("Discarded\n")
				return nil
			},
		}

		printFullCommand := &cli.Command{
			Name:    "body",
			Short:   "b",
			Summary: "Prints the full message body",
			Run: func() error {
				fmt.Printf("Full Message\n")
				prettyPrintJSON(state.Data.CurrentVersion.Message.Body.Value, false)
				return nil
			},
		}

		printMetadataCommand := &cli.Command{
			Name:    "metadata",
			Short:   "m",
			Summary: "Prints the metadata",

			Run: func() error {
				fmt.Printf("Metadata\n")
				fmt.Printf(" Handler Env:  %s\n", state.Data.Notification.HandlerEnv)
				fmt.Printf(" Handler App:  %s\n", state.Data.Notification.HandlerApp)
				fmt.Printf(" Death ID:     %s\n", state.Data.Notification.DeathId)
				printMetadata(state.Data.Notification.Message)
				/*

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
				*/
				return nil
			},
		}

		printInfraCommand := &cli.Command{
			Name:    "print-infra",
			Summary: "Prints the infra notification",
			Run: func() error {
				fmt.Printf("Notification\n")
				fmt.Printf("  Type: %s\n", state.Data.Notification.Infra.Type)
				for key, val := range state.Data.Notification.Infra.Metadata {
					fmt.Printf("    %s: %s\n", key, val)
				}
				return nil
			},
		}

		printSQSCommand := &cli.Command{
			Name:    "print-sqs",
			Summary: "Prints the SQS message",
			Run: func() error {
				if state.Data.CurrentVersion.SqsMessage != nil {
					fmt.Printf("Current SQS message\n")
					fmt.Printf("  URL: %s\n", state.Data.CurrentVersion.SqsMessage.QueueUrl)
					for key, val := range state.Data.CurrentVersion.SqsMessage.Attributes {
						fmt.Printf("    %s: %s\n", key, val)
					}
				}

				if newVersion.SqsMessage != nil {
					fmt.Printf("Modified SQS message\n")
					fmt.Printf("URL: %s\n", newVersion.SqsMessage.QueueUrl)
					for key, val := range newVersion.SqsMessage.Attributes {
						fmt.Printf("    %s: %s\n", key, val)
					}
				}
				return nil
			},
		}

		setURLCommand := &cli.Command{
			Name:    "set-url",
			Summary: "Set the URL of the SQS message",
			Run: func() error {
				url := cli.Ask("URL: ")
				if newVersion.SqsMessage == nil {
					newVersion.SqsMessage = &dante.DeadMessageVersion_SQSMessage{}
				}
				newVersion.SqsMessage.QueueUrl = url
				hasMods = true

				fmt.Printf("Queue URL: %s\n", newVersion.SqsMessage.QueueUrl)
				for key, val := range newVersion.SqsMessage.Attributes {
					fmt.Printf("  %s: %s\n", key, val)
				}
				return nil
			},
		}

		editBodyCommand := &cli.Command{
			Name:    "edit-body",
			Summary: "Edit the body of the message using $EDITOR",
			Run: func() error {
				buff := &bytes.Buffer{}
				if err := json.Indent(buff, state.Data.CurrentVersion.Message.Body.Value, "", "  "); err != nil {
					return fmt.Errorf("indent body: %w", err)
				}
				newBody, err := cli.Edit(ctx, buff.String())
				if err != nil {
					return fmt.Errorf("edit body: %w", err)
				}
				buff = &bytes.Buffer{}
				if err := json.Indent(buff, []byte(newBody), "  ", "  "); err != nil {
					return fmt.Errorf("indent body: %w", err)
				}

				newVersion.Message = &messaging.Message{
					MessageId:   state.Data.CurrentVersion.Message.MessageId,
					SourceApp:   state.Data.CurrentVersion.Message.SourceApp,
					SourceEnv:   state.Data.CurrentVersion.Message.SourceEnv,
					GrpcMethod:  state.Data.CurrentVersion.Message.GrpcMethod,
					GrpcService: state.Data.CurrentVersion.Message.GrpcService,
					Headers:     state.Data.CurrentVersion.Message.Headers,

					Body: &messaging.Any{
						Value:    []byte(newBody),
						TypeUrl:  state.Data.CurrentVersion.Message.Body.TypeUrl,
						Encoding: state.Data.CurrentVersion.Message.Body.Encoding,
					},
					MessageId:   state.Data.CurrentVersion.Message.MessageId,
					GrpcService: state.Data.CurrentVersion.Message.GrpcService,
					GrpcMethod:  state.Data.CurrentVersion.Message.GrpcMethod,
				}
				hasMods = true
				return nil
			},
		}

		quitCommand := &cli.Command{
			Name:  "quit",
			Short: "q",

			Summary: "Quit the interaction",
			Run: func() error {
				return errExitLoop
			},
		}

		for {
			options := []*cli.Command{
				printMetadataCommand,
				printFullCommand,
				printInfraCommand,
				printSQSCommand,
				deleteCommand,
				setURLCommand,
				editBodyCommand,
				quitCommand,
			}

			if hasMods {
				options = append(options, saveCommand, discardCommand)
			} else {
				options = append(options, nextCommand, replayCommand)
			}

			if err := cli.RunOneCommand(ctx, options); err != nil {
				if err == errExitMessage {
					return nil
				}
				if err == errExitLoop {
					return libo5.ErrStopPaging
				}
				return err
			}

		}
	}

	if err := libo5.Paged(ctx,
		&dante.ListDeadMessagesRequest{}, queryClient.ListDeadMessages,
		func(state *dante.DeadMessageState) error {

			printMessage(state)

			fmt.Printf("\n\n")

			if cfg.Interactive {
				err := runInteraction(ctx, state)
				if err != nil {
					return err
				}
				fmt.Print("============\nNextMessage\n============\n\n")
			}
			return nil
		}); err != nil {
		return fmt.Errorf("list stacks: %w", err)
	}

	fmt.Printf("All Done\n")
	return nil
}

func runDanteReject(ctx context.Context, cfg struct {
	libo5.APIConfig
	ID string `flag:"id" help:"ID of the dead message to shelve"`
}) error {
	client := cfg.APIClient()

	queryClient := dante.NewDeadMessageCommandService(client)

	res, err := queryClient.RejectDeadMessage(ctx, &dante.RejectDeadMessageRequest{
		MessageId: cfg.ID,
	})
	if err != nil {
		return fmt.Errorf("shelve dead messages: %w", err)
	}
	cli.Print("DL", res)

	return nil
}
