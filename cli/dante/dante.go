package dante

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pentops/o5-aws-tool/cli"
	"github.com/pentops/o5-aws-tool/cli/api"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/o5-aws-tool/libo5/o5/dante/v1/dante"
	"github.com/pentops/runner/commander"
)

func DanteCommandSet() *commander.CommandSet {
	remoteGroup := commander.NewCommandSet()

	remoteGroup.Add("ls", commander.NewCommand(runDanteLs))
	remoteGroup.Add("reject", commander.NewCommand(runDanteReject))

	return remoteGroup
}

func runDanteLs(ctx context.Context, cfg struct {
	api.BaseCommand
	Interactive bool `flag:"i" help:"Interactive mode"`
}) error {
	client := cfg.Client()

	queryClient := dante.NewDeadMessageQueryService(client)
	commandClient := dante.NewDeadMessageCommandService(client)
	printMessage := func(state *dante.DeadMessageState) {
		fmt.Printf("DL %s\n", state.Keys.MessageId)
		msg := state.Data.CurrentVersion.Message
		fmt.Printf("  /%s/%s\n", msg.GrpcService, msg.GrpcMethod)

		buff := &bytes.Buffer{}
		json.Indent(buff, msg.Body.Value, "  |", "  ")
		body := string(buff.Bytes())[:250]
		fmt.Printf("  |%s\n", body)

		if state.Data.Notification.Problem.UnhandledError != nil {
			fmt.Printf("  Error: %s", state.Data.Notification.Problem.UnhandledError.Error)
		}
	}
	runInteraction := func(ctx context.Context, state *dante.DeadMessageState) error {
		hasMods := false
		newVersion := &dante.DeadMessageVersion{}
		for {
			options := []string{"delete", "[r]eplay", "set-url"}

			if hasMods {
				options = append(options, "[s]ave - Saves the changes made to the message")
				options = append(options, "discard - Discards the pending changes")
			} else {
				options = append(options, "[n]ext - Ignore this message and move on", "[q]uit")
			}
			for _, opt := range options {
				fmt.Printf("  %s\n", opt)
			}
			option := cli.Ask("Action")
			if option == "" {
				fmt.Printf("invalid option\n")
				continue
			}
			fmt.Printf("option: %s\n", strings.ToLower(option))
			switch strings.ToLower(option) {
			case "d", "delete":
				res, err := commandClient.RejectDeadMessage(ctx, &dante.RejectDeadMessageRequest{
					MessageId: state.Keys.MessageId,
				})
				if err != nil {
					return fmt.Errorf("reject dead messages: %w", err)
				}
				cli.Print("DL", res.Message.Status)
				return nil

			case "r", "replay":
				res, err := commandClient.ReplayDeadMessage(ctx, &dante.ReplayDeadMessageRequest{
					MessageId: state.Keys.MessageId,
				})
				if err != nil {
					return fmt.Errorf("replay dead messages: %w", err)
				}
				cli.Print("DL", res.Message.Status)
				return nil

			case "s", "save":
				if !hasMods {
					fmt.Printf("No changes to save\n")
					continue
				}
				newID := uuid.NewString()
				newVersion.VersionId = newID
				res, err := commandClient.UpdateDeadMessage(ctx, &dante.UpdateDeadMessageRequest{
					MessageId:         state.Keys.MessageId,
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

			case "print-infra":
				fmt.Printf("Notification\n")
				fmt.Printf("  Type: %s\n", state.Data.Notification.Infra.Type)
				for key, val := range state.Data.Notification.Infra.Metadata {
					fmt.Printf("    %s: %s\n", key, val)
				}

			case "print-sqs":
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
				continue

			case "set-url":
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
				continue

			case "q", "quit":
				return libo5.ErrStopPaging
			default:
				fmt.Printf("Invalid option %q\n", option)
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
	api.BaseCommand
	ID string `flag:"id" help:"ID of the dead message to shelve"`
}) error {
	client := cfg.Client()

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
