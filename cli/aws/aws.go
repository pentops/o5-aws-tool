package aws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecs_types "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/pentops/o5-aws-tool/awsinspect"
	"github.com/pentops/runner/commander"
)

func CommandSet() *commander.CommandSet {

	cmdGroup := commander.NewCommandSet()

	cmdGroup.Add("redeploy", commander.NewCommand(runRedeploy))
	cmdGroup.Add("rules", commander.NewCommand(runRules))
	cmdGroup.Add("scale", commander.NewCommand(runScale))
	cmdGroup.Add("ecs-status", commander.NewCommand(runEcsStatus))

	stacksGroup := commander.NewCommandSet()
	stacksGroup.Add("ls", commander.NewCommand(runStacksList))
	cmdGroup.Add("stacks", stacksGroup)

	return cmdGroup
}

func runStacksList(ctx context.Context, cfg struct {
	Cluster     string `flag:"cluster" required:"false"`
	Environment string `flag:"env" required:"false"`

	Status bool `flag:"status"`
}) error {
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
		if cfg.Cluster != "" && !strings.HasPrefix(*stack.StackName, cfg.Cluster) {
			continue
		}
		if cfg.Environment != "" && !strings.HasSuffix(*stack.StackName, cfg.Environment) {
			continue
		}

		fmt.Printf("Stack:   %s\n", *stack.StackName)

		if cfg.Status {
			if err := printStackStatus(ctx, ecs.NewFromConfig(awsConfig), formationClient, *stack.StackName); err != nil {
				return err
			}
			fmt.Printf("===============\n")
		}

	}

	return nil
}

func runRedeploy(ctx context.Context, cfg struct {
	ClusterName string `flag:"cluster"`
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

func runScale(ctx context.Context, cfg struct {
	StackName    string `flag:"stack"`
	DesiredCount int32  `flag:"count"`
}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	ecsClient := ecs.NewFromConfig(awsConfig)

	formationClient := cloudformation.NewFromConfig(awsConfig)

	serviceSummary, err := awsinspect.GetStackServices(ctx, formationClient, cfg.StackName)
	if err != nil {
		return err
	}

	if len(serviceSummary.ServiceARNs) != 1 {
		return fmt.Errorf("expected 1 service, got %d", len(serviceSummary.ServiceARNs))
	}
	arn := serviceSummary.ServiceARNs[0]

	_, err = ecsClient.UpdateService(ctx, &ecs.UpdateServiceInput{
		ForceNewDeployment: true,
		Service:            aws.String(arn),
		Cluster:            aws.String(serviceSummary.ClusterName),
		DesiredCount:       aws.Int32(cfg.DesiredCount),
	})
	if err != nil {
		return err
	}

	return nil
}

func runRules(ctx context.Context, cfg struct {
	StackName string `flag:"stack"`
}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	formationClient := cloudformation.NewFromConfig(awsConfig)

	eventRules := []string{}
	listenerRules := []string{}
	{
		res, err := formationClient.DescribeStackResources(ctx, &cloudformation.DescribeStackResourcesInput{
			StackName: aws.String(cfg.StackName),
		})
		if err != nil {
			return err
		}

		for _, resource := range res.StackResources {
			switch *resource.ResourceType {
			case "AWS::Events::Rule":
				eventRules = append(eventRules, *resource.PhysicalResourceId)
			case "AWS::ElasticLoadBalancingV2::ListenerRule":
				listenerRules = append(listenerRules, *resource.PhysicalResourceId)
			}
		}
	}

	{
		eventClient := eventbridge.NewFromConfig(awsConfig)
		for _, rule := range eventRules {
			rp := strings.Split(rule, "|")
			if len(rp) != 2 {
				return fmt.Errorf("unexpected rule format: %s", rule)

			}
			busName, ruleName := rp[0], rp[1]
			resp, err := eventClient.DescribeRule(ctx, &eventbridge.DescribeRuleInput{
				Name:         aws.String(ruleName),
				EventBusName: aws.String(busName),
			})

			if err != nil {
				return err
			}
			fmt.Printf("  Event Rule: %s\n", *resp.Name)
			fmt.Printf("    Pattern: %s\n", *resp.EventPattern)
		}
	}

	{
		albClient := elasticloadbalancingv2.NewFromConfig(awsConfig)
		rulesRes, err := albClient.DescribeRules(ctx, &elasticloadbalancingv2.DescribeRulesInput{
			RuleArns: listenerRules,
		})

		if err != nil {
			return err
		}

		for _, rule := range rulesRes.Rules {
			fmt.Printf("  Listener Rule: %s\n", *rule.RuleArn)
			fmt.Printf("    Priority: %s\n", *rule.Priority)
			fmt.Printf("    Conditions:\n")
			for _, condition := range rule.Conditions {
				fmt.Printf("      Field: %s\n", *condition.Field)
				fmt.Printf("      Values: %v\n", condition.Values)
			}
		}

	}
	return nil
}

func runEcsStatus(ctx context.Context, cfg struct {
	StackName string `flag:"stack"`
}) error {
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	ecsClient := ecs.NewFromConfig(awsConfig)

	formationClient := cloudformation.NewFromConfig(awsConfig)

	return printStackStatus(ctx, ecsClient, formationClient, cfg.StackName)
}

func printStackStatus(ctx context.Context, ecsClient *ecs.Client, formationClient *cloudformation.Client, stackName string) error {

	serviceSummary, err := awsinspect.GetStackServices(ctx, formationClient, stackName)
	if err != nil {
		return err
	}

	if len(serviceSummary.ServiceARNs) != 1 {
		return fmt.Errorf("expected 1 service, got %d", len(serviceSummary.ServiceARNs))
	}
	arn := serviceSummary.ServiceARNs[0]

	resp, err := ecsClient.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Services: []string{arn},
		Cluster:  aws.String(serviceSummary.ClusterName),
	})
	if err != nil {
		return err
	}

	taskDefs := map[string]*ecs_types.TaskDefinition{}

	for _, service := range resp.Services {
		fmt.Printf("Service: %s\n", *service.ServiceName)
		for _, deployment := range service.Deployments {
			fmt.Printf("  Deployment: %s %s (%s)\n", *deployment.Status, *deployment.Id, deployment.CreatedAt.In(time.Local).Format("2006-01-02 15:04:05"))
			fmt.Printf("    Running: %d, Desired: %d, Failed: %d, Pending %d\n", deployment.RunningCount, deployment.DesiredCount, deployment.FailedTasks, deployment.PendingCount)

			taskDef := *deployment.TaskDefinition
			fmt.Printf("    Task Definition: %s\n", taskDef)

			if _, ok := taskDefs[taskDef]; !ok {
				taskDef, err := ecsClient.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
					TaskDefinition: aws.String(taskDef),
				})
				if err != nil {
					return err
				}
				taskDefs[*taskDef.TaskDefinition.TaskDefinitionArn] = taskDef.TaskDefinition
			}

			detail := taskDefs[taskDef]
			for _, containerDef := range detail.ContainerDefinitions {
				fmt.Printf("      Container: %s\n", *containerDef.Name)
				fmt.Printf("        Image: %s\n", *containerDef.Image)
			}
		}
	}

	return nil
}
