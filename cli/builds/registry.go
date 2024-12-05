package builds

import (
	"context"
	"fmt"
	"os"

	"github.com/pentops/o5-aws-tool/gen/j5/builds/github/v1/github"
	"github.com/pentops/o5-aws-tool/libo5"

	"github.com/pentops/runner/commander"
)

func BuildsCommandSet() *commander.CommandSet {

	registryGroup := commander.NewCommandSet()

	registryGroup.Add("ls", commander.NewCommand(runLs))
	registryGroup.Add("sync-config", commander.NewCommand(runRegistryConfig))
	registryGroup.Add("trigger-j5", commander.NewCommand(runTriggerJ5))
	registryGroup.Add("trigger-o5", commander.NewCommand(runTriggerO5))
	registryGroup.Add("trigger-o5-group", commander.NewCommand(runTriggerO5Group))

	return registryGroup

}

func runLs(ctx context.Context, cfg struct {
	libo5.APIConfig
}) error {
	client := cfg.APIClient()

	registryClient := github.NewRepoQueryService(client)

	repos, err := libo5.PagedAll(ctx, &github.ListReposRequest{}, registryClient.ListRepos)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		fmt.Printf("Repo: %s/%s\n", repo.Owner, repo.Name)
		for _, branch := range repo.Data.Branches {
			fmt.Printf("  Branch: %s\n", branch.BranchName)
			for _, deployTarget := range branch.DeployTargets {
				fmt.Printf("    Deploy Target: %s\n", deployTarget.OneofKey())
			}
		}
	}
	return nil

}

func runTriggerO5Group(ctx context.Context, cfg struct {
	libo5.APIConfig
	File string `flag:"file"`
}) error {

	client := cfg.APIClient()

	regClient := github.NewCombinedClient(client)

	data, err := os.ReadFile(cfg.File)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	repos, err := parseRepos(data)
	if err != nil {
		return fmt.Errorf("parse repos: %w", err)
	}

	for _, repo := range repos {
		fmt.Printf("Config: %v\n", repo)

		fmt.Printf("Checking %s/%s\n", repo.Owner, repo.Repo)

		for _, branch := range repo.DeployBranches {
			fmt.Printf("  Deploy Branch %s to %s\n", branch.Name, branch.Env)

			res, err := regClient.Trigger(ctx, &github.TriggerRequest{
				Commit: fmt.Sprintf("refs/heads/%s", branch.Name),
				Owner:  repo.Owner,
				Repo:   repo.Repo,
				Target: &github.DeployTargetType{
					O5Build: &github.DeployTargetType_O5Build{
						Environment: branch.Env,
					},
				},
			})
			if err != nil {
				return fmt.Errorf("trigger: %w", err)
			}

			for _, target := range res.Targets {
				fmt.Printf("Target: %s\n", target)
			}
			fmt.Printf("Triggered %d targets\n", len(res.Targets))
		}

	}

	return nil
}

func runTriggerJ5(ctx context.Context, cfg struct {
	libo5.APIConfig
	Owner string `flag:"owner"`
	Repo  string `flag:"repo"`
	Ref   string `flag:"ref" default:"refs/heads/main"`
}) error {
	client := cfg.APIClient()

	registryClient := github.NewRepoCommandService(client)

	res, err := registryClient.Trigger(ctx, &github.TriggerRequest{
		Commit: cfg.Ref,
		Owner:  cfg.Owner,
		Repo:   cfg.Repo,
		Target: &github.DeployTargetType{
			J5Build: &github.DeployTargetType_J5Build{},
		},
	})
	if err != nil {
		return fmt.Errorf("trigger: %w", err)
	}

	for _, target := range res.Targets {
		fmt.Printf("Target: %s\n", target)
	}
	fmt.Printf("Triggered %d targets\n", len(res.Targets))

	return nil
}

func runTriggerO5(ctx context.Context, cfg struct {
	libo5.APIConfig
	Owner       string `flag:"owner"`
	Repo        string `flag:"repo"`
	Ref         string `flag:"ref" default:"refs/heads/main"`
	Environment string `flag:"env"`
}) error {
	client := cfg.APIClient()

	registryClient := github.NewRepoCommandService(client)

	res, err := registryClient.Trigger(ctx, &github.TriggerRequest{
		Commit: cfg.Ref,
		Owner:  cfg.Owner,
		Repo:   cfg.Repo,
		Target: &github.DeployTargetType{
			O5Build: &github.DeployTargetType_O5Build{
				Environment: cfg.Environment,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("trigger: %w", err)
	}

	for _, target := range res.Targets {
		fmt.Printf("Target: %s\n", target)
	}
	fmt.Printf("Triggered %d targets\n", len(res.Targets))

	return nil
}
