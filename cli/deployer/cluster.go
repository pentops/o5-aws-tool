package deployer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pentops/o5-aws-tool/gen/o5/aws/deployer/v1/deployer"
	"github.com/pentops/o5-aws-tool/libo5"
	"github.com/pentops/runner/commander"
)

type ClusterBaseConfig struct {
	ClusterID string `flag:"id"`
}

func runCluster(ctx context.Context, cfg struct {
	ClusterBaseConfig
	Flags []string `flag:",remaining"`
}) error {

	cs := commander.NewCommandSet()
	cs.Add("override", clusterOverride(cfg.ClusterBaseConfig))
	cs.Add("pull", clusterPull(cfg.ClusterBaseConfig))
	return cs.Run(ctx, cfg.Flags)
}

func clusterOverride(base ClusterBaseConfig) commander.Runnable {
	return commander.NewCommand(func(ctx context.Context, cfg struct {
		libo5.APIConfig
		Key string `flag:"key"`
		Val string `flag:"val"`
	}) error {
		client := cfg.APIClient()

		var val *string
		if cfg.Val != "" {
			val = &cfg.Val
		}

		cl := deployer.NewDeploymentCommandService(client)
		res, err := cl.SetClusterOverride(ctx, &deployer.SetClusterOverrideRequest{
			ClusterId: base.ClusterID,
			Overrides: []*deployer.ParameterOverride{{
				Key:   cfg.Key,
				Value: val,
			}},
		})
		if err != nil {
			return fmt.Errorf("set cluster override: %w", err)
		}
		return printCluster(res.State)
	})
}

func clusterPull(base ClusterBaseConfig) commander.Runnable {
	return commander.NewCommand(func(ctx context.Context, cfg struct {
		libo5.APIConfig
	}) error {
		queryClient := deployer.NewClusterQueryService(cfg.APIClient())

		res, err := queryClient.GetCluster(ctx, &deployer.GetClusterRequest{
			ClusterId: base.ClusterID,
		})
		if err != nil {
			return fmt.Errorf("get cluster: %w", err)
		}
		return printCluster(res.State)
	})
}

type s3API interface {
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

var s3Client s3API

func getS3Client(ctx context.Context) (s3API, error) {
	if s3Client != nil {
		return s3Client, nil
	}
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	s3Client = s3.NewFromConfig(awsConfig)
	return s3Client, nil
}

func readFile(ctx context.Context, path string) ([]byte, error) {
	if strings.HasPrefix(path, "s3://") {
		client, err := getS3Client(ctx)
		if err != nil {
			return nil, err
		}
		bucket := strings.TrimPrefix(path, "s3://")
		parts := strings.SplitN(bucket, "/", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid s3 path: %s", path)
		}
		res, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &parts[0],
			Key:    &parts[1],
		})
		if err != nil {
			return nil, fmt.Errorf("get object: %w", err)
		}

		return io.ReadAll(res.Body)
	}
	return os.ReadFile(path)
}

func runClusterConfig(ctx context.Context, cfg struct {
	libo5.APIConfig
	Src string `flag:"src"`
}) error {
	client := cfg.APIClient()

	content, err := readFile(ctx, cfg.Src)
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
	return printCluster(res.State)
}

func runClusterOverride(ctx context.Context, cfg struct {
	libo5.APIConfig
	ClusterBaseConfig
	Key string `flag:"key"`
	Val string `flag:"val"`
}) error {
	client := cfg.APIClient()

	var val *string
	if cfg.Val != "" {
		val = &cfg.Val
	}

	cl := deployer.NewDeploymentCommandService(client)
	res, err := cl.SetClusterOverride(ctx, &deployer.SetClusterOverrideRequest{
		ClusterId: cfg.ClusterID,
		Overrides: []*deployer.ParameterOverride{{
			Key:   cfg.Key,
			Value: val,
		}},
	})
	if err != nil {
		return fmt.Errorf("set cluster override: %w", err)
	}
	return printCluster(res.State)
}

func runClusters(ctx context.Context, cfg struct {
	libo5.APIConfig
}) error {
	queryClient := deployer.NewClusterQueryService(cfg.APIClient())

	if err := libo5.Paged(ctx,
		&deployer.ListClustersRequest{}, queryClient.ListClusters,
		func(cluster *deployer.ClusterState) error {
			fmt.Printf("Cluster %s\n", cluster.Data.Config.Name)
			fmt.Printf("   ID: %s\n", cluster.ClusterId)
			fmt.Printf("   Status: %s\n", cluster.Status)
			fmt.Printf("   Created: %s\n", cluster.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("   Updated: %s\n", cluster.Metadata.UpdatedAt.Format("2006-01-02 15:04:05"))
			return nil
		}); err != nil {
		return fmt.Errorf("list clusters: %w", err)
	}
	return nil
}

func printCluster(cluster *deployer.ClusterState) error {

	fmt.Printf("Cluster %s\n", cluster.Data.Config.Name)
	fmt.Printf("   ID: %s\n", cluster.ClusterId)
	fmt.Printf("   Status: %s\n", cluster.Status)
	fmt.Printf("   Created: %s\n", cluster.Metadata.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("   Updated: %s\n", cluster.Metadata.UpdatedAt.Format("2006-01-02 15:04:05"))
	jd, err := json.MarshalIndent(cluster.Data.Config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	fmt.Println(string(jd))
	return nil
}
