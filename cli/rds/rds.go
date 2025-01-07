package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/ettle/strcase"
	"github.com/pentops/golib/gl"
	"github.com/pentops/runner/commander"
)

func CommandSet() *commander.CommandSet {

	modSet := commander.NewCommandSet()
	modSet.Add("root", commander.NewCommand(runRoot))
	modSet.Add("admin", commander.NewCommand(runAdmin))
	modSet.Add("app", commander.NewCommand(runApp))
	return modSet
}

type IXEnv struct {
	EnvName     string `env:"IX_ENV"`
	ClusterName string `env:"IX_CLUSTER"`
	AWSRegion   string `env:"AWS_REGION"`

	_aws *aws.Config
}

func (ix *IXEnv) AWS(ctx context.Context) (aws.Config, error) {
	if ix._aws != nil {
		return *ix._aws, nil
	}
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load configuration: %w", err)
	}
	ix._aws = &awsConfig
	return awsConfig, nil
}

func (ix *IXEnv) MustAWS(ctx context.Context) aws.Config {
	awsConfig, err := ix.AWS(ctx)
	if err != nil {
		panic(err)
	}
	return awsConfig
}

type ServerCfg struct {
	IXEnv
	ServerName string `flag:"server" default:"default" description:"DB provider name"`
}

func (cfg *ServerCfg) Cluster(ctx context.Context) (*types.DBCluster, error) {
	aws, err := cfg.AWS(ctx)
	if err != nil {
		return nil, err
	}
	rdscl := rds.NewFromConfig(aws)

	res, err := rdscl.DescribeDBClusters(ctx, &rds.DescribeDBClustersInput{
		DBClusterIdentifier: gl.Ptr(fmt.Sprintf("%s-%s", cfg.ClusterName, cfg.ServerName)),
	})
	if err != nil {
		return nil, err
	}

	if len(res.DBClusters) == 0 {
		return nil, fmt.Errorf("DB cluster %s not found", fmt.Sprintf("%s-%s", cfg.ClusterName, "default"))
	}

	cluster := res.DBClusters[0]

	return &cluster, nil
}

func (cfg *ServerCfg) LoginIAM(ctx context.Context, cluster *types.DBCluster, username, dbname string) (string, error) {

	aws, err := cfg.AWS(ctx)
	if err != nil {
		return "", err
	}

	endpoint := gl.MustUnwrap(cluster.Endpoint)
	port := gl.MustUnwrap(cluster.Port)
	endpointAndPort := fmt.Sprintf("%s:%d", endpoint, port)
	authenticationToken, err := auth.BuildAuthToken(
		ctx, endpointAndPort, cfg.AWSRegion, username, aws.Credentials)
	if err != nil {
		return "", fmt.Errorf("failed to create authentication token: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		gl.MustUnwrap(cluster.Endpoint),
		gl.MustUnwrap(cluster.Port),
		username,
		dbname,
		authenticationToken,
	)

	return dsn, nil
}

func runRoot(ctx context.Context, cfg struct {
	ServerCfg
	PSQL bool `flag:"psql" default:"false" description:"Fork psql command"`
	IAM  bool `flag:"iam" default:"true" description:"Use IAM auth"`
}) error {

	cluster, err := cfg.Cluster(ctx)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "endpoint: %s", gl.MustUnwrap(cluster.Endpoint))

	var dsn string
	if cfg.IAM {
		username := gl.MustUnwrap(cluster.MasterUsername)
		dbname := gl.MustUnwrap(cluster.DatabaseName)
		dsn = gl.Must(cfg.LoginIAM(ctx, cluster, username, dbname))
	} else {
		secretscl := secretsmanager.NewFromConfig(cfg.MustAWS(ctx))

		secretARN := cluster.MasterUserSecret.SecretArn

		secret, err := secretscl.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: secretARN,
		})
		if err != nil {
			return err
		}

		cred := struct {
			Password string `json:"password"`
		}{}
		err = json.Unmarshal([]byte(gl.MustUnwrap(secret.SecretString)), &cred)
		if err != nil {
			return err
		}
		password := cred.Password
		dsn = fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
			gl.MustUnwrap(cluster.Endpoint),
			gl.MustUnwrap(cluster.Port),
			gl.MustUnwrap(cluster.MasterUsername),
			gl.MustUnwrap(cluster.DatabaseName),
			password)
	}

	if !cfg.PSQL {
		fmt.Println(dsn)
		return nil
	}

	return forkPSQL(ctx, dsn)
}

func forkPSQL(_ context.Context, dsn string) error {
	psqlPath, err := exec.LookPath("psql")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "path: %s", psqlPath)
	return syscall.Exec(psqlPath, []string{"psql", dsn}, os.Environ())
}

func runAdmin(ctx context.Context, cfg struct {
	ServerCfg
	PSQL bool `flag:"psql" default:"false" description:"Fork psql command"`
}) error {

	cluster, err := cfg.Cluster(ctx)
	if err != nil {
		return err
	}

	dbName := strings.Join([]string{
		strcase.ToSnake(cfg.ClusterName),
		"o5admin",
	}, "_")

	fmt.Fprintf(os.Stderr, "dbName: %s\n", dbName)
	fmt.Fprintf(os.Stderr, "endpoint: %s\n", gl.Opt(cluster.Endpoint))

	dsn, err := cfg.LoginIAM(ctx, cluster, dbName, dbName)
	if err != nil {
		return err
	}

	if !cfg.PSQL {
		fmt.Println(dsn)
		return nil
	}

	return forkPSQL(ctx, dsn)
}

func runApp(ctx context.Context, cfg struct {
	ServerCfg
	PSQL bool   `flag:"psql" default:"false" description:"Fork psql command"`
	App  string `flag:"app"  description:"App name"`
	Name string `flag:"name" default:"main" description:"DB name (within the o5 config)"`
}) error {

	cluster, err := cfg.Cluster(ctx)
	if err != nil {
		return err
	}

	dbName := strings.Join([]string{
		strcase.ToSnake(cfg.EnvName),
		strcase.ToSnake(cfg.App),
		strcase.ToSnake(cfg.Name),
	}, "_")

	fmt.Fprintf(os.Stderr, "dbName: %s\n", dbName)
	fmt.Fprintf(os.Stderr, "endpoint: %s\n", gl.Opt(cluster.Endpoint))

	dsn := gl.Must(cfg.LoginIAM(ctx, cluster, dbName, dbName))

	if !cfg.PSQL {
		fmt.Println(dsn)
		return nil
	}

	return forkPSQL(ctx, dsn)
}
