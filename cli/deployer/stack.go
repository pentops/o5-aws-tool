package deployer

import (
	"context"
	"fmt"

	"github.com/pentops/o5-aws-tool/gen/o5/aws/deployer/v1/deployer"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/runner/commander"
)

func runStacks(ctx context.Context, cfg struct {
	libo5.APIConfig
}) error {
	queryClient := deployer.NewCombinedClient(cfg.APIClient())

	if err := libo5.Paged(ctx,
		&deployer.ListStacksRequest{}, queryClient.ListStacks,
		func(stack *deployer.StackState) error {
			fmt.Printf("%s %s\n", stack.Data.EnvironmentName, stack.Data.ApplicationName)
			fmt.Printf("   Status: %s\n", stack.Status)
			fmt.Printf("   ID: %s\n", stack.StackId)

			return nil
		}); err != nil {
		return fmt.Errorf("list stacks: %w", err)
	}
	return nil
}

func runStack(ctx context.Context, args struct {
	StackID string   `flag:"id"`
	Flags   []string `flag:",remaining"`
}) error {

	cs := commander.NewCommandSet()
	cs.Add("status", stackStatusCommand(args.StackID))
	cs.Add("terminate-all", stackTerminateCommand(args.StackID))

	if err := cs.Run(ctx, args.Flags); err != nil {
		return fmt.Errorf("sub task: %w", err)
	}
	return nil

}

func stackStatusCommand(stackID string) commander.Runnable {
	return commander.NewCommand(func(ctx context.Context, cfg struct {
		libo5.APIConfig
	}) error {
		queryClient := deployer.NewCombinedClient(cfg.APIClient())
		res, err := queryClient.GetStack(ctx, &deployer.GetStackRequest{
			StackId: stackID,
		})
		if err != nil {
			return fmt.Errorf("get stack: %w", err)
		}
		fmt.Printf("StackID: %s\n", res.State.StackId)
		fmt.Printf("Status: %s\n", res.State.Status)
		if res.State.Data.CurrentDeployment != nil {
			fmt.Printf("  Current Deployment: %s\n", res.State.Data.CurrentDeployment.DeploymentId)
		}
		for _, queued := range res.State.Data.QueuedDeployments {
			fmt.Printf("  Queued Deployment: %s\n", queued.DeploymentId)
		}
		return nil
	})
}
func stackTerminateCommand(stackID string) commander.Runnable {
	return commander.NewCommand(func(ctx context.Context, cfg struct {
		libo5.APIConfig
	}) error {
		client := cfg.APIClient()
		commandClient := deployer.NewCombinedClient(client)
		queryClient := deployer.NewCombinedClient(client)
		stack, err := queryClient.GetStack(ctx, &deployer.GetStackRequest{
			StackId: stackID,
		})
		if err != nil {
			return fmt.Errorf("get stack: %w", err)
		}

		deployments := stack.State.Data.QueuedDeployments
		deployments = append(deployments, stack.State.Data.CurrentDeployment)
		for _, queued := range deployments {
			_, err := commandClient.TerminateDeployment(ctx, &deployer.TerminateDeploymentRequest{
				DeploymentId: queued.DeploymentId,
			})
			if err != nil {
				return fmt.Errorf("terminate deployment: %w", err)
			}
		}
		return nil
	})
}
