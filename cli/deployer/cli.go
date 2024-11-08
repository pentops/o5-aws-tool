package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/pentops/o5-aws-tool/cli"
	"github.com/pentops/o5-aws-tool/gen/j5/list/v1/list"
	"github.com/pentops/o5-aws-tool/gen/o5/aws/deployer/v1/deployer"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/runner/commander"
)

func O5CommandSet() *commander.CommandSet {
	remoteGroup := commander.NewCommandSet()
	remoteGroup.Add("trigger", commander.NewCommand(runTrigger))

	remoteGroup.Add("deployment", commander.NewCommand(runDeployment))
	remoteGroup.Add("deployments", commander.NewCommand(runDeployments))

	remoteGroup.Add("stack", commander.NewCommand(runStack))
	remoteGroup.Add("stacks", commander.NewCommand(runStacks))
	remoteGroup.Add("environments", commander.NewCommand(runEnvironments))

	remoteGroup.Add("cluster-config", commander.NewCommand(runClusterConfig))

	return remoteGroup
}

type StateCache struct {
	StateData string `env:"O5_CLI_STATE_DATA" default:"$HOME/.local/share/o5-cli/state.json"`
}

type stateData struct {
	Data map[string]string `json:"data"`
}

func (cfg *StateCache) SetVal(key, val string) error {
	data := &stateData{}
	if err := os.MkdirAll(filepath.Dir(cfg.StateData), 0700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	if _, err := os.Stat(cfg.StateData); err == nil {
		content, err := os.ReadFile(cfg.StateData)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}
		if err := json.Unmarshal(content, data); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
	}
	if data.Data == nil {
		data.Data = map[string]string{}
	}
	data.Data[key] = val
	content, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(cfg.StateData, content, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (cfg *StateCache) GetVal(key string) (string, error) {
	data := &stateData{}
	if _, err := os.Stat(cfg.StateData); err == nil {
		content, err := os.ReadFile(cfg.StateData)
		if err != nil {
			return "", fmt.Errorf("read file: %w", err)
		}
		if err := json.Unmarshal(content, data); err != nil {
			return "", fmt.Errorf("unmarshal: %w", err)
		}
	}
	val, ok := data.Data[key]
	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}
	return val, nil
}

func runClusterConfig(ctx context.Context, cfg struct {
	libo5.APIConfig
	Src string `flag:"src"`
}) error {
	client := cfg.APIClient()

	content, err := os.ReadFile(cfg.Src)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	req := &deployer.UpsertClusterRequest{}

	ext := filepath.Ext(cfg.Src)
	switch ext {
	case ".json":
		req.ConfigJson = content
	case ".yaml", ".yml":
		req.ConfigYaml = content
	default:
		return fmt.Errorf("unknown file type: %s", ext)
	}

	commandClient := deployer.NewDeploymentCommandService(client)
	res, err := commandClient.UpsertCluster(ctx, req)
	if err != nil {
		return fmt.Errorf("upsert cluster: %w", err)
	}
	cli.Print("Cluster Config", res.State)
	return nil
}

/*
func runClusters(ctx context.Context, cfg struct {
	api.BaseCommand
}) error {
	queryClient := deployer.NewDeploymentQueryService(cfg.Client())

	if err := libo5.Paged(ctx,
		&deployer.List{}, queryClient.ListClusters,
		func(cluster *deployer.ClusterState) error {
			fmt.Printf("%s\n", cluster.ClusterId)
			return nil
		}); err != nil {
		return fmt.Errorf("list clusters: %w", err)
	}
	return nil
}(*/

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

func runEnvironments(ctx context.Context, cfg struct {
	libo5.APIConfig
	API string `env:"O5_API" flag:"api"`
}) error {
	queryClient := deployer.NewCombinedClient(cfg.APIClient())

	if err := libo5.Paged(ctx,
		&deployer.ListEnvironmentsRequest{}, queryClient.ListEnvironments,
		func(env *deployer.EnvironmentState) error {
			fmt.Println("=========")
			fmt.Printf("Environment: %s\n", env.Data.Config.FullName)
			fmt.Printf("  ID: %s\n", env.EnvironmentId)

			for idx, jwks := range env.Data.Config.TrustJwks {
				fmt.Printf("  JWKS[%d]: %s\n", idx, jwks)
			}
			for _, cfg := range env.Data.Config.Vars {
				fmt.Printf("  %s: %s\n", cfg.Name, cfg.Value)
			}

			return nil
		}); err != nil {
		return fmt.Errorf("list stacks: %w", err)
	}
	return nil
}

func runDeployment(ctx context.Context, args struct {
	StateCache
	DeploymentID string   `flag:"id" default:"$last"`
	Flags        []string `flag:",remaining"`
}) error {

	if args.DeploymentID == "$last" {
		last, err := args.GetVal("last-deployment")
		if err != nil {
			return fmt.Errorf("get last deployment: %w", err)
		}
		args.DeploymentID = last
	}

	cs := commander.NewCommandSet()
	cs.Add("status", deploymentStatusCommand(args.DeploymentID))
	cs.Add("terminate", deploymentTerminateCommand(args.DeploymentID))

	if err := cs.Run(ctx, args.Flags); err != nil {
		return fmt.Errorf("sub task: %w", err)
	}

	return nil

}

func deploymentTerminateCommand(deploymentID string) commander.Runnable {
	return commander.NewCommand(func(ctx context.Context, cfg struct {
		libo5.APIConfig
	}) error {
		commandClient := deployer.NewDeploymentCommandService(cfg.APIClient())
		_, err := commandClient.TerminateDeployment(ctx, &deployer.TerminateDeploymentRequest{
			DeploymentId: deploymentID,
		})
		if err != nil {
			return fmt.Errorf("terminate deployment: %w", err)
		}
		return nil
	})
}

var stepStatusColor = map[string]color.Attribute{
	"DONE":    color.FgGreen,
	"ACTIVE":  color.FgBlue,
	"BLOCKED": color.FgYellow,
	"FAILED":  color.FgRed,
}

var deploymentStatusColor = map[string]color.Attribute{
	"RUNNING":    color.FgBlue,
	"TERMINATED": color.FgRed,
	"DONE":       color.FgGreen,
}

func deploymentStatusCommand(deploymentID string) commander.Runnable {
	return commander.NewCommand(func(ctx context.Context, cfg struct {
		libo5.APIConfig
		Wait bool `flag:"wait"`
	}) error {
		queryClient := deployer.NewDeploymentQueryService(cfg.APIClient())

		return deploymentStatus(ctx, queryClient, deploymentID, cfg.Wait)
	})
}

func deploymentStatus(ctx context.Context, queryClient *deployer.DeploymentQueryService, deploymentID string, wait bool) error {

	lastLastEvent := uint64(999)
	didDots := false
	for {
		res, err := queryClient.GetDeployment(ctx, &deployer.GetDeploymentRequest{
			DeploymentId: deploymentID,
		})
		if err != nil {
			return fmt.Errorf("get deployment: %w", err)
		}

		if lastLastEvent == res.State.Metadata.LastSequence {
			didDots = true
			fmt.Printf(".")
			time.Sleep(time.Second)
			continue
		}
		if didDots {
			fmt.Println()
			didDots = false
		}

		lastLastEvent = res.State.Metadata.LastSequence

		fmt.Printf("DeploymentID: %s\n", res.State.DeploymentId)
		fmt.Printf("Status: %s\n", res.State.Status)
		fmt.Printf("\n")
		stepMap := map[string]*deployer.DeploymentStep{}
		for _, step := range res.State.Data.Steps {
			stepMap[step.Id] = step
		}

		if err := listDeploymentEvents(ctx, queryClient, res.State); err != nil {
			return err
		}

		steps := res.State.Data.Steps
		fmt.Println()
		sort.Sort(StepsByStatus(steps))
		for _, step := range steps {
			fmt.Printf("Step: %s\n", color.MagentaString(step.Name))
			stepColor, ok := stepStatusColor[step.Status]
			if !ok {
				stepColor = color.FgWhite
			}
			fmt.Printf("  StepID: %s\n", step.Id)
			fmt.Printf("  Status: %s\n", color.New(stepColor).Sprint(step.Status))
			if step.Error != nil {
				fmt.Printf("  Error: %s\n", *step.Error)
			}
			for _, id := range step.DependsOn {
				dep := stepMap[id]
				if dep.Status != "DONE" {
					fmt.Printf("  BlockedBy: %s (%s)\n", dep.Name, dep.Status)
				}
			}

			fmt.Printf("\n")
		}
		if !wait {
			break
		}

		switch res.State.Status {
		case "DONE", "TERMINATED", "FAILED":
			return nil
		}
		time.Sleep(time.Second)
	}
	return nil
}

type StepsByStatus []*deployer.DeploymentStep

func (s StepsByStatus) Len() int {
	return len(s)
}

func (s StepsByStatus) Less(i, j int) bool {
	return s[i].Status < s[j].Status
}

func (s StepsByStatus) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func runDeployments(ctx context.Context, cfg struct {
	libo5.APIConfig
	StateCache
	AppName string `env:"APP_NAME" flag:"app" default:""`
	EnvName string `env:"ENV_NAME" flag:"env" default:""`
	Version string `env:"VERSION" flag:"version" default:""`
	All     bool   `flag:"all" description:"Include inactive"`

	Q    bool `flag:"q" help:"output only the IDs"`
	Wait bool `flag:"wait"`
}) error {
	if cfg.Wait && cfg.Q {
		return fmt.Errorf("wait and q are mutually exclusive")
	}

	queryClient := deployer.NewDeploymentQueryService(cfg.APIClient())

	foundDeployments := []*deployer.DeploymentState{}

	query := &list.QueryRequest{}

	if cfg.AppName != "" {
		query.Filters = append(query.Filters, &list.Filter{
			Field: &list.Field{
				Name: "data.spec.appName",
				Type: &list.FieldType{
					Value: cfg.AppName,
				},
			},
		})
	}

	if cfg.EnvName != "" {
		query.Filters = append(query.Filters, &list.Filter{
			Field: &list.Field{
				Name: "data.spec.environmentName",
				Type: &list.FieldType{
					Value: cfg.EnvName,
				},
			},
		})
	}

	if cfg.Version != "" {
		query.Filters = append(query.Filters, &list.Filter{
			Field: &list.Field{
				Name: "data.spec.version",
				Type: &list.FieldType{
					Value: cfg.Version,
				},
			},
		})
	}

	if err := libo5.Paged(ctx,
		&deployer.ListDeploymentsRequest{
			Query: query,
		},
		queryClient.ListDeployments,
		func(deployment *deployer.DeploymentState) error {
			foundDeployments = append(foundDeployments, deployment)
			return nil
		}); err != nil {
		return fmt.Errorf("list deployments: %w", err)
	}

	if len(foundDeployments) == 0 {
		if !cfg.Q {
			fmt.Println("No deployments found")
		}
		return nil
	}

	runningDeployments := []*deployer.DeploymentState{}

	for _, foundDeployment := range foundDeployments {

		if cfg.Q {
			fmt.Println(foundDeployment.DeploymentId)
			continue
		}

		statusColor, ok := deploymentStatusColor[foundDeployment.Status]
		if !ok {
			statusColor = color.FgWhite
		}

		if cfg.All || foundDeployment.Status == "RUNNING" {
			runningDeployments = append(runningDeployments, foundDeployment)
		}

		fmt.Printf("DeploymentID: %s\n", foundDeployment.DeploymentId)
		fmt.Printf("  Status: %s\n", color.New(statusColor).Sprint(foundDeployment.Status))
		fmt.Printf("  AppName: %s\n", foundDeployment.Data.Spec.AppName)
		fmt.Printf("  EnvName: %s\n", foundDeployment.Data.Spec.EnvironmentName)
		fmt.Printf("  Version: %s\n", foundDeployment.Data.Spec.Version)
		fmt.Printf("\n")

	}

	if !cfg.Wait {
		return nil
	}

	if len(runningDeployments) > 1 {
		return fmt.Errorf("more than one deployment found")
	}

	return deploymentStatus(ctx, queryClient, foundDeployments[0].DeploymentId, true)

}

func runStack(ctx context.Context, args struct {
	StackID string   `flag:",arg0"`
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
		API string `env:"O5_API" flag:"api"`
	}) error {
		client := libo5.NewAPI(cfg.API)
		queryClient := deployer.NewCombinedClient(client)
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
		API string `env:"O5_API" flag:"api"`
	}) error {
		client := libo5.NewAPI(cfg.API)
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

func listDeploymentEvents(ctx context.Context, queryClient *deployer.DeploymentQueryService, deployment *deployer.DeploymentState) error {
	stepMap := map[string]*deployer.DeploymentStep{}
	for _, step := range deployment.Data.Steps {
		stepMap[step.Id] = step
	}
	if err := libo5.Paged(ctx,
		&deployer.ListDeploymentEventsRequest{
			DeploymentId: deployment.DeploymentId,
			Query: &list.QueryRequest{
				Sorts: []*list.Sort{{
					Field:      "metadata.timestamp",
					Descending: false,
				}},
			},
		}, queryClient.ListDeploymentEvents,
		func(event *deployer.DeploymentEvent) error {

			fmt.Printf("%-20s %s\n", color.GreenString(event.Event.OneofKey()), event.Metadata.Timestamp.Format(time.RFC3339))

			switch et := event.Event.Type().(type) {
			case *deployer.DeploymentEventType_StepResult:
				step := stepMap[et.StepId]
				fmt.Printf("  Step: %s\n", step.Name)
				fmt.Printf("  Status: %s\n", et.Status)
				fmt.Printf("  ID: %s\n", et.StepId)

			case *deployer.DeploymentEventType_RunStep:
				step := stepMap[et.StepId]
				fmt.Printf("  Step: %s\n", step.Name)
				fmt.Printf("  ID: %s\n", step.Id)

			case *deployer.DeploymentEventType_Error:
				fmt.Printf("  Error: %s\n", et.Error)
			}
			return nil
		}); err != nil {
		return fmt.Errorf("list deployment events: %w", err)
	}

	return nil

}

func runTrigger(ctx context.Context, cfg struct {
	libo5.APIConfig
	StateCache

	AppName string `env:"APP_NAME" flag:"repo"`
	Org     string `env:"GITHUB_ORG" flag:"org"`
	EnvName string `env:"ENV_NAME" flag:"env"`
	Version string `env:"VERSION" flag:"version"`

	DBOnly            bool `flag:"db-only"`
	Quick             bool `flag:"quick"`
	RotateCredentials bool `flag:"rotate-creds"`
}) error {
	deploymentID := uuid.NewString()

	triggerBody := &deployer.TriggerDeploymentRequest{
		DeploymentId: deploymentID,
		Environment:  cfg.EnvName,
		Source: &deployer.TriggerSource{
			Github: &deployer.TriggerSource_GithubSource{
				Owner:  cfg.Org,
				Repo:   cfg.AppName,
				Commit: cfg.Version,
			},
		},
		Flags: &deployer.DeploymentFlags{
			QuickMode:         cfg.Quick,
			RotateCredentials: cfg.RotateCredentials,
			DbOnly:            cfg.DBOnly,
		},
	}

	client := cfg.APIClient()
	commandClient := deployer.NewDeploymentCommandService(client)

	_, err := commandClient.TriggerDeployment(ctx, triggerBody)
	if err != nil {
		return err
	}

	if err := cfg.SetVal("last-deployment", deploymentID); err != nil {
		return fmt.Errorf("set last deployment: %w", err)
	}

	fmt.Println(deploymentID)
	return nil
}
