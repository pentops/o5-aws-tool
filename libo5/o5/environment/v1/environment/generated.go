package environment

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/libo5/o5/environment/v1/environment

import ()

// CustomVariable_Join Proto: CustomVariable_Join
type CustomVariable_Join struct {
	Delimiter string   `json:"delimiter,omitempty"`
	Values    []string `json:"values,omitempty"`
}

// CombinedConfig Proto: CombinedConfig
type CombinedConfig struct {
	Name         string         `json:"name,omitempty"`
	EcsCluster   *ECSCluster    `json:"ecsCluster,omitempty"`
	Environments []*Environment `json:"environments,omitempty"`
}

// AWSEnvironment Proto: AWSEnvironment
type AWSEnvironment struct {
	HostHeader       *string           `json:"hostHeader,omitempty"`
	EnvironmentLinks []*AWSLink        `json:"environmentLinks,omitempty"`
	IamPolicies      []*NamedIAMPolicy `json:"iamPolicies,omitempty"`
}

// Environment Proto: Environment
type Environment struct {
	FullName    string            `json:"fullName,omitempty"`
	ClusterName string            `json:"clusterName,omitempty"`
	TrustJwks   []string          `json:"trustJwks,omitempty"`
	Vars        []*CustomVariable `json:"vars,omitempty"`
	CorsOrigins []string          `json:"corsOrigins,omitempty"`
	Aws         *AWSEnvironment   `json:"aws,omitempty"`
}

// Cluster Proto: Cluster
type Cluster struct {
	Name       string      `json:"name,omitempty"`
	EcsCluster *ECSCluster `json:"ecsCluster,omitempty"`
}

// CustomVariable Proto: CustomVariable
type CustomVariable struct {
	Name  string               `json:"name,omitempty"`
	Value string               `json:"value,omitempty"`
	Join  *CustomVariable_Join `json:"join,omitempty"`
}

// ECSCluster Proto: ECSCluster
type ECSCluster struct {
	ListenerArn          string     `json:"listenerArn,omitempty"`
	EcsClusterName       string     `json:"ecsClusterName,omitempty"`
	EcsRepo              string     `json:"ecsRepo,omitempty"`
	EcsTaskExecutionRole string     `json:"ecsTaskExecutionRole,omitempty"`
	VpcId                string     `json:"vpcId,omitempty"`
	AwsAccount           string     `json:"awsAccount,omitempty"`
	AwsRegion            string     `json:"awsRegion,omitempty"`
	EventBusArn          string     `json:"eventBusArn,omitempty"`
	O5DeployerAssumeRole string     `json:"o5DeployerAssumeRole,omitempty"`
	O5DeployerGrantRoles []string   `json:"o5DeployerGrantRoles,omitempty"`
	RdsHosts             []*RDSHost `json:"rdsHosts,omitempty"`
	SidecarImageVersion  *string    `json:"sidecarImageVersion,omitempty"`
	SidecarImageRepo     *string    `json:"sidecarImageRepo,omitempty"`
	GlobalNamespace      string     `json:"globalNamespace,omitempty"`
}

// NamedIAMPolicy Proto: NamedIAMPolicy
type NamedIAMPolicy struct {
	Name      string `json:"name,omitempty"`
	PolicyArn string `json:"policyArn,omitempty"`
}

// RDSHost Proto: RDSHost
type RDSHost struct {
	ServerGroup string `json:"serverGroup,omitempty"`
	SecretName  string `json:"secretName,omitempty"`
}

// AWSLink Proto: AWSLink
type AWSLink struct {
	LookupName string `json:"lookupName,omitempty"`
	FullName   string `json:"fullName,omitempty"`
	SnsPrefix  string `json:"snsPrefix,omitempty"`
}
