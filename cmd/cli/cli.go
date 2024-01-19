package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/pentops/o5-aws-tool/awsinspect"
	"github.com/pentops/runner"
	"github.com/pentops/runner/commander"
)

var Version string

func main() {

	cmdGroup := commander.NewCommandSet()

	cmdGroup.Add("logs", commander.NewCommand(runAWSLogs))
	cmdGroup.Add("redeploy", commander.NewCommand(runRedeploy))

	cmdGroup.RunMain("o5-aws-tool", Version)

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
			return awsinspect.TailLogStream(ctx, logClient, logStream, fromTime)
		})
	}

	return rg.Run(ctx)

}
