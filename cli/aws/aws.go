package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/fatih/color"
	"github.com/pentops/o5-aws-tool/awsinspect"
	"github.com/pentops/runner"
	"github.com/pentops/runner/commander"
)

func CommandSet() *commander.CommandSet {

	cmdGroup := commander.NewCommandSet()

	cmdGroup.Add("logs", commander.NewCommand(runAWSLogs))
	cmdGroup.Add("events", commander.NewCommand(runEventLogs))
	cmdGroup.Add("redeploy", commander.NewCommand(runRedeploy))

	stacksGroup := commander.NewCommandSet()
	stacksGroup.Add("ls", commander.NewCommand(runStacksList))
	cmdGroup.Add("stacks", stacksGroup)

	return cmdGroup
}

func runStacksList(ctx context.Context, cfg struct{}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	formationClient := cloudformation.NewFromConfig(awsConfig)

	allStatuses := (types.StackStatus("")).Values()
	wantStatus := make([]types.StackStatus, 0, len(allStatuses))
	for _, status := range allStatuses {
		if status == types.StackStatusDeleteComplete {
			continue
		}
		wantStatus = append(wantStatus, status)
	}

	stacks, err := formationClient.ListStacks(ctx, &cloudformation.ListStacksInput{
		StackStatusFilter: wantStatus,
	})
	if err != nil {
		return err
	}

	for _, stack := range stacks.StackSummaries {
		fmt.Printf("Stack: %s\n", *stack.StackName)
	}

	return nil
}

func runRedeploy(ctx context.Context, cfg struct {
	ClusterName string `flag:"cluster-name"`
}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	ecsClient := ecs.NewFromConfig(awsConfig)

	listRes, err := ecsClient.ListServices(ctx, &ecs.ListServicesInput{
		Cluster: aws.String(cfg.ClusterName),
	})
	if err != nil {
		return err
	}

	for _, arn := range listRes.ServiceArns {
		fmt.Printf("Service: %s\n", arn)

		_, err := ecsClient.UpdateService(ctx, &ecs.UpdateServiceInput{
			ForceNewDeployment: true,
			Service:            aws.String(arn),
			Cluster:            aws.String(cfg.ClusterName),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func runEventLogs(ctx context.Context, cfg struct {
	ClusterName string `flag:"cluster-name"`
}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	_ = awsConfig
	return nil
}

func runAWSLogs(ctx context.Context, cfg struct {
	StackName string `flag:"stack-name"`
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

	rg := runner.NewGroup()
	for _, logStream := range logStreams {
		logStream := logStream
		rg.Add(logStream.Container, func(ctx context.Context) error {
			return awsinspect.TailLogStream(ctx, logClient, logStream, fromTime, prettyLog)
		})
	}

	return rg.Run(ctx)

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
