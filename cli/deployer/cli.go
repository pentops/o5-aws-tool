package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/pentops/o5-aws-tool/gen/o5/aws/deployer/v1/deployer"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/runner/commander"
)

func O5CommandSet() *commander.CommandSet {
	remoteGroup := commander.NewCommandSet()
	remoteGroup.Add("trigger", commander.NewCommand(runTrigger))
	remoteGroup.Add("local", commander.NewCommand(runLocal))

	remoteGroup.Add("deployment", commander.NewCommand(runDeployment))
	remoteGroup.Add("deployments", commander.NewCommand(runDeployments))

	remoteGroup.Add("stack", commander.NewCommand(runStack))
	remoteGroup.Add("stacks", commander.NewCommand(runStacks))

	remoteGroup.Add("environment", commander.NewCommand(runEnvironment))
	remoteGroup.Add("environments", commander.NewCommand(runEnvironments))
	remoteGroup.Add("environment-config", commander.NewCommand(runEnvironmentConfig))

	remoteGroup.Add("cluster", commander.NewCommand(runCluster))
	remoteGroup.Add("cluster-config", commander.NewCommand(runClusterConfig))
	remoteGroup.Add("clusters", commander.NewCommand(runClusters))
	remoteGroup.Add("cluster-override", commander.NewCommand(runClusterOverride))

	return remoteGroup
}

type StateCache struct {
	StateData string `env:"O5_CLI_STATE_DATA" default:"$HOME/.local/share/o5-cli/state.json"`
}

type stateData struct {
	Data map[string]string `json:"data"`
}

func (cfg *StateCache) SetVal(key, val string) error {
	dataPath := os.ExpandEnv(cfg.StateData)
	data := &stateData{}
	if err := os.MkdirAll(filepath.Dir(dataPath), 0700); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	if _, err := os.Stat(dataPath); err == nil {
		content, err := os.ReadFile(dataPath)
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
	if err := os.WriteFile(dataPath, content, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (cfg *StateCache) GetVal(key string) (string, error) {
	dataPath := os.ExpandEnv(cfg.StateData)

	data := &stateData{}
	if _, err := os.Stat(dataPath); err == nil {
		content, err := os.ReadFile(dataPath)
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

func runEnvironment(ctx context.Context, cfg struct {
	libo5.APIConfig

	Id string `flag:"id"`
}) error {
	queryClient := deployer.NewCombinedClient(cfg.APIClient())

	req := &deployer.GetEnvironmentRequest{
		EnvironmentId: cfg.Id,
	}

	res, err := queryClient.GetEnvironment(ctx, req)
	if err != nil {
		return fmt.Errorf("get environment: %w", err)
	}

	fmt.Println("Environment:")
	fmt.Println("\tName:", res.State.Data.Config.FullName)
	fmt.Println("\tID:", res.State.EnvironmentId)
	fmt.Println("\tClusterID:", res.State.ClusterId)
	fmt.Println("\tStatus: ", res.State.Status)
	fmt.Println("=========")

	b, err := json.MarshalIndent(res.State.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	fmt.Printf("Data: %v\n", string(b))

	return nil
}

func runEnvironments(ctx context.Context, cfg struct {
	libo5.APIConfig
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

func runEnvironmentConfig(ctx context.Context, cfg struct {
	libo5.APIConfig

	Id  string `flag:"id"`
	Src string `flag:"src"`
}) error {
	client := cfg.APIClient()

	content, err := readFile(ctx, cfg.Src)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	req := &deployer.UpsertEnvironmentRequest{
		EnvironmentId: cfg.Id,
	}

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
	res, err := commandClient.UpsertEnvironment(ctx, req)
	if err != nil {
		return fmt.Errorf("upsert environment: %w", err)
	}

	fmt.Println("Environment:")
	fmt.Println("\tID:", res.State.EnvironmentId)
	fmt.Println("\tName:", res.State.Data.Config.FullName)

	return nil
}

func runTrigger(ctx context.Context, cfg struct {
	libo5.APIConfig
	StateCache

	AppName string `env:"APP_NAME" flag:"repo"`
	Org     string `env:"GITHUB_ORG" flag:"org"`
	EnvName string `env:"ENV_NAME" flag:"env"`

	Version string `env:"VERSION" flag:"version"`
	Import  bool   `flag:"import"`

	DBOnly            bool `flag:"db-only"`
	InfraOnly         bool `flag:"infra-only"`
	DBReset           bool `flag:"db-reset"`
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
			InfraOnly:         cfg.InfraOnly,
			ImportResources:   cfg.Import,
			RecreateDatabase:  cfg.DBReset,
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

func runLocal(ctx context.Context, cfg struct {
	libo5.APIConfig
	StateCache

	AppName string `env:"APP_NAME" flag:"repo"`
	Org     string `env:"GITHUB_ORG" flag:"org"`
	EnvName string `env:"ENV_NAME" flag:"env"`
	Version string `env:"VERSION" flag:"version"`
	Import  bool   `flag:"import"`

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
			ImportResources:   cfg.Import,
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
