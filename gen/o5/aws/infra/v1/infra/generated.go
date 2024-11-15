package infra

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/o5/aws/infra/v1/infra

import ()

// ECSTaskNetworkType_AWSVPC Proto: ECSTaskNetworkType_AWSVPC
type ECSTaskNetworkType_AWSVPC struct {
	SecurityGroups []string `json:"securityGroups,omitempty"`
	Subnets        []string `json:"subnets,omitempty"`
}

// ECSTaskNetworkType Proto Oneof: o5.aws.infra.v1.ECSTaskNetworkType
type ECSTaskNetworkType struct {
	J5TypeKey string                     `json:"!type,omitempty"`
	Awsvpc    *ECSTaskNetworkType_AWSVPC `json:"awsvpc,omitempty"`
}

func (s ECSTaskNetworkType) OneofKey() string {
	if s.Awsvpc != nil {
		return "awsvpc"
	}
	return ""
}

func (s ECSTaskNetworkType) Type() interface{} {
	if s.Awsvpc != nil {
		return s.Awsvpc
	}
	return nil
}

// AuroraConnection Proto: AuroraConnection
type AuroraConnection struct {
	Endpoint   string `json:"endpoint,omitempty"`
	Port       int32  `json:"port,omitempty"`
	DbUser     string `json:"dbUser,omitempty"`
	DbName     string `json:"dbName,omitempty"`
	Identifier string `json:"identifier,omitempty"`
}

// RDSHostType Proto Oneof: o5.aws.infra.v1.RDSHostType
type RDSHostType struct {
	J5TypeKey      string                      `json:"!type,omitempty"`
	Aurora         *RDSHostType_Aurora         `json:"aurora,omitempty"`
	SecretsManager *RDSHostType_SecretsManager `json:"secretsManager,omitempty"`
}

func (s RDSHostType) OneofKey() string {
	if s.Aurora != nil {
		return "aurora"
	}
	if s.SecretsManager != nil {
		return "secretsManager"
	}
	return ""
}

func (s RDSHostType) Type() interface{} {
	if s.Aurora != nil {
		return s.Aurora
	}
	if s.SecretsManager != nil {
		return s.SecretsManager
	}
	return nil
}

// ECSTaskContext Proto: ECSTaskContext
type ECSTaskContext struct {
	Cluster string              `json:"cluster,omitempty"`
	Network *ECSTaskNetworkType `json:"network,omitempty"`
}

// RDSHostType_Aurora Proto: RDSHostType_Aurora
type RDSHostType_Aurora struct {
	Conn *AuroraConnection `json:"conn,omitempty"`
}

// RDSHostType_SecretsManager Proto: RDSHostType_SecretsManager
type RDSHostType_SecretsManager struct {
	SecretName string `json:"secretName,omitempty"`
}