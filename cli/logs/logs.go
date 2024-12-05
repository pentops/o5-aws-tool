package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/fatih/color"
	"github.com/pentops/o5-aws-tool/awsinspect"
	"github.com/pentops/o5-messaging/gen/o5/messaging/v1/messaging_pb"
	"github.com/pentops/o5-messaging/gen/o5/messaging/v1/messaging_tpb"
	"github.com/pentops/runner"
	"github.com/pentops/runner/commander"
	"google.golang.org/protobuf/encoding/protojson"
)

func CommandSet() *commander.CommandSet {
	cmdGroup := commander.NewCommandSet()
	cmdGroup.Add("logs", commander.NewCommand(runAWSLogs))
	cmdGroup.Add("events", commander.NewCommand(runEventLogs))
	return cmdGroup
}

func runAWSLogs(ctx context.Context, cfg struct {
	StackName string `flag:"stack"`
	Since     string `flag:"since" default:"0"`
}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	formationClient := cloudformation.NewFromConfig(awsConfig)

	serviceSummary, err := awsinspect.GetStackServices(ctx, formationClient, cfg.StackName)
	if err != nil {
		return err
	}

	ecsClient := ecs.NewFromConfig(awsConfig)
	logStreams, err := awsinspect.GetAllLogStreams(ctx, ecsClient, serviceSummary)
	if err != nil {
		return err
	}

	logClient := cloudwatchlogs.NewFromConfig(awsConfig)

	fromTime := time.Now()
	if cfg.Since != "0" {
		dur, err := time.ParseDuration(cfg.Since)
		if err != nil {
			return err
		}
		fromTime = time.Now().Add(-dur)
	}

	rg := runner.NewGroup()
	for _, logStream := range logStreams {
		logStream := logStream
		rg.Add(logStream.Container, func(ctx context.Context) error {
			return awsinspect.TailLogStream(ctx, logClient, logStream, fromTime, prettyLog)
		})
	}

	return rg.Run(ctx)

}

func runEventLogs(ctx context.Context, cfg struct {
	EnvName string `flag:"env"`
}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logClient := cloudwatchlogs.NewFromConfig(awsConfig)

	groupName := fmt.Sprintf("/eventbus/app-events-%s", cfg.EnvName)
	res, err := logClient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(groupName),
	})
	if err != nil {
		return err
	}

	if len(res.LogGroups) != 1 {
		return fmt.Errorf("expected 1 log group, got %d", len(res.LogGroups))
	}

	groupARN := *res.LogGroups[0].Arn
	groupARN = strings.TrimSuffix(groupARN, ":*")

	tail, err := logClient.StartLiveTail(ctx, &cloudwatchlogs.StartLiveTailInput{
		LogGroupIdentifiers: []string{groupARN},
	})
	if err != nil {
		return err
	}

	stream := tail.GetStream()
	chEvents := stream.Events()

	for event := range chEvents {
		switch ev := event.(type) {
		case *types.StartLiveTailResponseStreamMemberSessionStart:
			log.Println("Received SessionStart event")
		case *types.StartLiveTailResponseStreamMemberSessionUpdate:
			for _, logEvent := range ev.Value.SessionResults {
				logEventEvent(*logEvent.Message)
			}
		default:
			// Handle on-stream exceptions
			if err := stream.Err(); err != nil {
				return err
			} else if event == nil {
				return fmt.Errorf("stream closed")
			} else {
				log.Fatalf("Unknown event type: %T", ev)
			}
		}
	}

	return nil

}

type msgWrap struct {
	ID         string          `json:"id"`
	DetailType string          `json:"detail-type"`
	Source     string          `json:"source"`
	Detail     json.RawMessage `json:"detail"`
}

func logEventEvent(message string) {
	if len(message) < 1 || message[0] != '{' {
		fmt.Println(message)
		return
	}

	wrap := msgWrap{}
	err := json.Unmarshal([]byte(message), &wrap)
	if err != nil {
		fmt.Println(message)
		return
	}

	msg := &messaging_pb.Message{}
	err = protojson.Unmarshal([]byte(wrap.Detail), msg)
	if err != nil {
		fmt.Println(message)
		return
	}

	fmt.Printf("Event: %s\n", msg.MessageId)
	fmt.Printf("  Topic: %s\n", msg.DestinationTopic)
	fmt.Printf("  Source App: %s\n", msg.SourceApp)
	fmt.Printf("  Source Env: %s\n", msg.SourceEnv)
	fmt.Printf("  Timestamp: %s\n", msg.Timestamp.AsTime().Format(time.RFC3339))

	fmt.Printf("  GRPCService: %s\n", msg.GrpcService)
	fmt.Printf("  GRPCMethod: %s\n", msg.GrpcMethod)
	fmt.Printf("  Body.TypeUrl: %s\n", msg.Body.TypeUrl)
	/*
		body, err := base64.URLEncoding.DecodeString(string(msg.Body.Value))
		if err != nil {
			fmt.Printf("  Body.ERROR: %s\n", err)
			fmt.Printf("  Body: %s\n", msg.Body.Value)
		} else {*/

	body := msg.Body.Value
	switch msg.Body.TypeUrl {
	case "type.googleapis.com/o5.messaging.v1.topic.RawMessage":
		rawBody := &messaging_tpb.RawMessage{}
		err := protojson.Unmarshal(body, rawBody)
		if err != nil {
			fmt.Printf("  Body.ERROR: %s\n", err)
			fmt.Printf("  Body: %s\n", string(body))
		} else {
			fmt.Printf("  Body.Topic: %v\n", rawBody.Topic)
			fmt.Printf("  Body.Payload: %s\n", rawBody.Payload)
		}

	default:

		fmt.Printf("  Body: %s\n", body)
	}
	//	}

}

var levelColors = map[string]color.Attribute{
	"debug": color.FgBlue,
	"info":  color.FgGreen,
	"warn":  color.FgYellow,
	"error": color.FgRed,
}

type logLine struct {
	Level   string                 `json:"level"`
	Time    string                 `json:"time"`
	Message string                 `json:"message"`
	Fields  map[string]interface{} `json:"fields"`
}

func prettyLog(ctx context.Context, logGroup awsinspect.LogStream, message string) {
	out := os.Stdout
	if len(message) < 1 {
		return
	}
	if message[0] != '{' {
		fmt.Fprintf(out, "[%s] %s\n", logGroup.Container, message)
		return
	}

	line := logLine{}
	err := json.Unmarshal([]byte(message), &line)
	if err != nil {
		fmt.Fprintf(out, "[%s] %s\n", logGroup.Container, message)
		return
	}

	whichColor, ok := levelColors[strings.ToLower(line.Level)]
	if !ok {
		whichColor = color.FgWhite
	}
	levelColor := color.New(whichColor).SprintFunc()
	fmt.Fprintf(out, "%s [%s] %s\n", levelColor(line.Level), logGroup.LogStream, levelColor(line.Message))

	for k, v := range line.Fields {

		switch v.(type) {
		case string, int, int64, int32, float64, bool:
			fmt.Fprintf(out, "  | %s: %v\n", k, v)
		default:
			nice, _ := json.MarshalIndent(v, "  |  ", "  ")
			fmt.Fprintf(out, "  | %s: %s\n", k, string(nice))
		}
	}

}
