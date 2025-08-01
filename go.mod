module github.com/pentops/o5-aws-tool

go 1.24.0

toolchain go1.24.2

require (
	github.com/aws/aws-sdk-go-v2 v1.36.3
	github.com/aws/aws-sdk-go-v2/config v1.29.9
	github.com/aws/aws-sdk-go-v2/feature/rds/auth v1.5.2
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.55.5
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.43.2
	github.com/aws/aws-sdk-go-v2/service/ecs v1.49.2
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.43.0
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.35.6
	github.com/aws/aws-sdk-go-v2/service/rds v1.93.2
	github.com/aws/aws-sdk-go-v2/service/s3 v1.67.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.34.8
	github.com/ettle/strcase v0.2.0
	github.com/fatih/color v1.18.0
	github.com/google/uuid v1.6.0
	github.com/pentops/golib v0.0.0-20250326060930-8c83d58ddb63
	github.com/pentops/log.go v0.0.0-20250521181902-0b84b98a60de
	github.com/pentops/o5-messaging v0.0.0-20250618204836-8a488dd0569c
	github.com/pentops/runner v0.0.0-20250619010747-2bb7a5385324
	golang.org/x/sync v0.14.0
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v2 v2.4.0
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250603165357-b52ab10f4468.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.62 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.30 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.34 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.17 // indirect
	github.com/aws/smithy-go v1.22.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pentops/j5 v0.0.0-20250617223808-91fae5a3b112 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.71.0 // indirect
)
