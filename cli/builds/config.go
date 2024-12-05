package builds

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pentops/o5-aws-tool/cli"
	"github.com/pentops/o5-aws-tool/gen/j5/builds/github/v1/github"
	"github.com/pentops/o5-aws-tool/libo5"
	"gopkg.in/yaml.v2"
)

type DeployBranch struct {
	Name string `json:"name"`
	Env  string `json:"env"`
}

type RepoConfig struct {
	Owner          string          `yaml:"owner"`
	Repo           string          `yaml:"repo"`
	Checks         bool            `yaml:"checks"`
	J5             bool            `yaml:"j5"`
	DeployBranches []*DeployBranch `yaml:"deployBranches"`
}

func parseRepos(data []byte) ([]RepoConfig, error) {
	rawType := struct {
		Repos []RepoConfig
	}{}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.SetStrict(true)
	if err := dec.Decode(&rawType); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return rawType.Repos, nil

}

func runRegistryConfig(ctx context.Context, cfg struct {
	libo5.APIConfig
	File   string   `flag:"file"`
	Filter []string `flag:"name" required:"false" help:"Filter by name"`
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

	allRepos, err := libo5.PagedAll(ctx, &github.ListReposRequest{}, regClient.ListRepos)
	if err != nil {
		return fmt.Errorf("list repos: %w", err)
	}

	byID := sliceToMap(allRepos, func(repo *github.RepoState) string {
		return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	})

	for _, repo := range repos {
		if len(cfg.Filter) > 0 {
			found := false
			for _, filter := range cfg.Filter {
				if filter == repo.Repo {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		fmt.Printf("Config: %v\n", repo)

		fmt.Printf("Checking %s/%s\n", repo.Owner, repo.Repo)

		branches := []*github.Branch{}

		if repo.J5 {
			fmt.Printf("  Build Branch * to j5\n")
			branches = append(branches, &github.Branch{
				BranchName: "*",
				DeployTargets: []*github.DeployTargetType{{
					J5Build: &github.DeployTargetType_J5Build{},
				}},
			})
		}

		for _, branch := range repo.DeployBranches {
			fmt.Printf("  Deploy Branch %s to %s\n", branch.Name, branch.Env)
			branches = append(branches, &github.Branch{
				BranchName: branch.Name,
				DeployTargets: []*github.DeployTargetType{{
					O5Build: &github.DeployTargetType_O5Build{
						Environment: branch.Env,
					},
				}},
			})
		}

		wantConfig := &github.RepoStateData{
			ChecksEnabled: repo.Checks,
			Branches:      branches,
		}

		if len(branches) == 0 {
			fmt.Printf("No Branches\n")
			continue
		}

		name := fmt.Sprintf("%s/%s", repo.Owner, repo.Repo)
		existing, ok := byID[name]
		if ok {
			delete(byID, name)
			cli.Print("Existing", existing)
			if configMatches(existing.Data, wantConfig) {
				fmt.Printf("No Changes\n")
				continue
			}
		}

		req := &github.ConfigureRepoRequest{
			Owner: repo.Owner,
			Name:  repo.Repo,
			Config: &github.RepoEventType_Configure{
				ChecksEnabled: wantConfig.ChecksEnabled,
				Branches:      wantConfig.Branches,
				Merge:         false,
			},
		}

		cli.Print("Configure", req)
		res, err := regClient.ConfigureRepo(ctx, req)
		if err != nil {
			log.Printf("configure repo: %v", err)
			continue

			//return fmt.Errorf("configure repo: %w", err)
		}
		cli.Print("Configured", res)
	}

	for id := range byID {
		fmt.Printf("Exists on remote: %s\n", id)
	}

	return nil
}

func configMatches(existing *github.RepoStateData, want *github.RepoStateData) bool {
	if existing.ChecksEnabled != want.ChecksEnabled {
		fmt.Printf("Checks Enabled: %t != %t\n", existing.ChecksEnabled, want.ChecksEnabled)
		return false
	}

	existingBranches := sliceToMap(existing.Branches, func(branch *github.Branch) string {
		return branch.BranchName
	})

	wantBranches := sliceToMap(want.Branches, func(branch *github.Branch) string {
		return branch.BranchName
	})

	for name, want := range wantBranches {
		existing, ok := existingBranches[name]
		if !ok {
			fmt.Printf("Missing Branch: %s\n", name)
			return false
		}
		if !branchMatches(existing, want) {
			fmt.Printf("Branch Mismatch: %s\n", name)
			return false
		}
		delete(existingBranches, name)
	}

	if len(existingBranches) > 0 {
		for name := range existingBranches {
			fmt.Printf("Extra Branch: %s\n", name)
		}
		fmt.Printf("Extra Branches\n")
		return false
	}
	return true
}

func branchMatches(existing *github.Branch, want *github.Branch) bool {

	toMap := func(target *github.DeployTargetType) string {
		if target.O5Build != nil {
			return "o5/" + target.O5Build.Environment
		}
		if target.J5Build != nil {
			return "j5"
		}
		return ""
	}
	namedExisting := sliceToMap(existing.DeployTargets, toMap)
	namedWant := sliceToMap(want.DeployTargets, toMap)

	for name := range namedWant {
		_, ok := namedExisting[name]
		if !ok {
			return false
		}
		delete(namedExisting, name)
	}

	return len(namedExisting) <= 0

}
func sliceToMap[T any](slice []T, key func(T) string) map[string]T {
	out := map[string]T{}
	for _, item := range slice {
		out[key(item)] = item
	}
	return out
}
