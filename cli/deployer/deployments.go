package deployer

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/pentops/o5-aws-tool/gen/j5/drss/v1/drss"
	"github.com/pentops/o5-aws-tool/gen/j5/list/v1/list"
	"github.com/pentops/o5-aws-tool/gen/o5/aws/deployer/v1/deployer"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/runner/commander"
	"golang.org/x/sync/errgroup"
)

func runDeployment(ctx context.Context, args struct {
	StateCache
	DeploymentID string `flag:"id" default:"$last"`

	Flags []string `flag:",remaining"`
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

var stepStatusColor = map[drss.StepStatus]color.Attribute{
	"DONE":    color.FgGreen,
	"ACTIVE":  color.FgBlue,
	"BLOCKED": color.FgYellow,
	"FAILED":  color.FgRed,
}

var deploymentStatusColor = map[deployer.DeploymentStatus]color.Attribute{
	"RUNNING":    color.FgBlue,
	"TERMINATED": color.FgRed,
	"DONE":       color.FgGreen,
}

func deploymentStatusCommand(deploymentID string) commander.Runnable {
	return commander.NewCommand(func(ctx context.Context, cfg struct {
		libo5.APIConfig
		deploymentStatusConfig
	}) error {
		queryClient := deployer.NewDeploymentQueryService(cfg.APIClient())
		return deploymentStatus(ctx, queryClient, deploymentID, &deploymentStatusContext{
			cfg:       cfg.deploymentStatusConfig,
			printLock: sync.Mutex{},
		})
	})
}

type deploymentStatusConfig struct {
	Wait    bool `flag:"wait"`
	Verbose bool `flag:"verbose"`
}

type deploymentStatusContext struct {
	cfg       deploymentStatusConfig
	noDots    bool
	printLock sync.Mutex
}

func deploymentStatus(ctx context.Context, queryClient *deployer.DeploymentQueryService, id string, run *deploymentStatusContext) error {

	lastLastEvent := uint64(999)
	didDots := false
	for {
		res, err := queryClient.GetDeployment(ctx, &deployer.GetDeploymentRequest{
			DeploymentId: id,
		})
		if err != nil {
			return fmt.Errorf("get deployment: %w", err)
		}

		err = func() error {
			run.printLock.Lock()
			defer run.printLock.Unlock()

			if lastLastEvent == res.State.Metadata.LastSequence {
				didDots = true
				if !run.noDots {
					fmt.Printf(".")
				}
				time.Sleep(time.Second)
				return nil
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
				stepMap[step.Meta.StepId] = step
			}

			if err := listDeploymentEvents(ctx, queryClient, res.State); err != nil {
				return err
			}

			steps := res.State.Data.Steps
			fmt.Println()
			sort.Sort(StepsByStatus(steps))
			for _, step := range steps {
				fmt.Printf("Step: %s\n", color.MagentaString(step.Meta.Name))
				stepColor, ok := stepStatusColor[step.Meta.Status]
				if !ok {
					stepColor = color.FgWhite
				}
				fmt.Printf("  StepID: %s\n", step.Meta.StepId)
				fmt.Printf("  Status: %s\n", color.New(stepColor).Sprint(step.Meta.Status))
				if step.Meta.Error != nil {
					fmt.Printf("  Error: %s\n", *step.Meta.Error)
				}
				for _, id := range step.Meta.DependsOn {
					dep := stepMap[id]
					if dep.Meta.Status != "DONE" {
						fmt.Printf("  BlockedBy: %s (%s)\n", dep.Meta.Name, dep.Meta.Status)
					}
				}

				fmt.Printf("\n")
			}
			return nil
		}()
		if err != nil {
			return err
		}
		if !run.cfg.Wait {
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
	return s[i].Meta.Status < s[j].Meta.Status
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
	Recent  bool   `flag:"recent" description:"Include inactive < 1 hour"`

	Q bool `flag:"q" help:"output only the IDs"`
	deploymentStatusConfig
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

	if cfg.Recent {
		query.Filters = append(query.Filters, &list.Filter{
			Field: &list.Field{
				Name: "metadata.createdAt",
				Type: &list.FieldType{
					Range: &list.Range{
						Min: time.Now().Add(-5 * time.Hour).Format(time.RFC3339),
						Max: time.Now().Add(time.Minute).Format(time.RFC3339),
					},
				},
			},
		})
		cfg.All = true
	}

	for {
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
		if cfg.Q || len(foundDeployments) > 0 {
			break
		}

		time.Sleep(2 * time.Second)

		fmt.Printf("Waiting for deployments...\n")
	}

	runningDeployments := []*deployer.DeploymentState{}

	for _, foundDeployment := range foundDeployments {
		if !cfg.All && foundDeployment.Status != "RUNNING" && foundDeployment.Status != "QUEUED" {

			continue
		}

		if cfg.Q {
			fmt.Println(foundDeployment.DeploymentId)
			continue
		}
		runningDeployments = append(runningDeployments, foundDeployment)

	}

	if cfg.Q {
		return nil
	}

	for _, foundDeployment := range runningDeployments {

		statusColor, ok := deploymentStatusColor[foundDeployment.Status]
		if !ok {
			statusColor = color.FgWhite
		}

		fmt.Printf("DeploymentID: %s\n", foundDeployment.DeploymentId)
		fmt.Printf("  Status: %s\n", color.New(statusColor).Sprint(foundDeployment.Status))
		fmt.Printf("  AppName: %s\n", foundDeployment.Data.Spec.AppName)
		fmt.Printf("  EnvName: %s\n", foundDeployment.Data.Spec.EnvironmentName)
		fmt.Printf("  Version: %s\n", foundDeployment.Data.Spec.Version)
		fmt.Printf("  Started: %s\n", foundDeployment.Metadata.CreatedAt.Format(time.RFC3339))
		fmt.Printf("\n")

	}

	if !cfg.Wait {
		return nil
	}

	ccc := &deploymentStatusContext{
		cfg:    cfg.deploymentStatusConfig,
		noDots: false,
	}
	if len(runningDeployments) < 1 {
		fmt.Printf("No Deployments")
		return nil
	}

	if len(runningDeployments) > 1 {
		ccc.noDots = true

		eg := &errgroup.Group{}

		for _, deployment := range runningDeployments {
			deployment := deployment
			eg.Go(func() error {
				return deploymentStatus(ctx, queryClient, deployment.DeploymentId, ccc)
			})
		}

	}

	return deploymentStatus(ctx, queryClient, runningDeployments[0].DeploymentId, ccc)

}

func listDeploymentEvents(ctx context.Context, queryClient *deployer.DeploymentQueryService, deployment *deployer.DeploymentState) error {
	stepMap := map[string]*deployer.DeploymentStep{}
	for _, step := range deployment.Data.Steps {
		stepMap[step.Meta.StepId] = step
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
				step := stepMap[et.Result.StepId]
				fmt.Printf("  Step: %s\n", step.Meta.Name)
				fmt.Printf("  Status: %s\n", step.Meta.Status)
				fmt.Printf("  ID: %s\n", step.Meta.StepId)

			case *deployer.DeploymentEventType_RunStep:
				step := stepMap[et.StepId]
				fmt.Printf("  Step: %s\n", step.Meta.Name)
				fmt.Printf("  ID: %s\n", step.Meta.StepId)

			case *deployer.DeploymentEventType_Error:
				fmt.Printf("  Error: %s\n", et.Error)
			}
			return nil
		}); err != nil {
		return fmt.Errorf("list deployment events: %w", err)
	}

	return nil

}
