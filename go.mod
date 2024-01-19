module github.com/pentops/o5-aws-tool

go 1.21.4

require (
	github.com/aws/aws-sdk-go-v2 v1.24.1
	github.com/aws/aws-sdk-go-v2/config v1.26.5
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.42.6
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.31.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.37.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.48.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.26.2
	github.com/awslabs/goformation/v7 v7.13.1
	github.com/pentops/log.go v0.0.0-20231218074934-67aedcab3fa4
	github.com/pentops/o5-deploy-aws v0.0.0-20240119203808-a59a18bbd10e
	github.com/pentops/o5-go v0.0.0-20240119203826-8268f62deb80
	github.com/pentops/o5-runtime-sidecar v0.0.7-0.20240119184739-8c7fc14b9e5b
	github.com/pentops/runner v0.0.0-20240119184422-1878cd4dc14d
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.31.0-20231115204500-e097f827e652.2 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.5.4 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.16.16 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.14.11 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.10 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.10.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.2.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.10.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.16.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.26.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sqs v1.29.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.18.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.21.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.26.7 // indirect
	github.com/aws/smithy-go v1.19.0 // indirect
	github.com/bufbuild/protovalidate-go v0.4.3 // indirect
	github.com/elgris/sqrl v0.0.0-20210727210741-7e0198b30236 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/goccy/go-yaml v1.11.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/cel-go v0.18.2 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pentops/jsonapi v0.0.0-20240119183000-75b3b4ca86bc // indirect
	github.com/pentops/jwtauth v0.0.0-20231218034817-a97d0d7fe8cc // indirect
	github.com/pentops/outbox.pg.go v0.0.0-20231222014950-493c01cfbcc7 // indirect
	github.com/pentops/protostate v0.0.0-20240119200340-2268e85845ec // indirect
	github.com/pentops/sqrlx.go v0.0.0-20240108202916-8687fdf983c0 // indirect
	github.com/pentops/sugar-go v0.0.0-20231029194349-ec12ec0132c5 // indirect
	github.com/pquerna/cachecontrol v0.2.0 // indirect
	github.com/rs/cors v1.10.1 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/tidwall/gjson v1.17.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	golang.org/x/crypto v0.16.0 // indirect
	golang.org/x/exp v0.0.0-20240103183307-be819d1f06fc // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/genproto v0.0.0-20231127180814-3a041ad873d4 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231127180814-3a041ad873d4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231127180814-3a041ad873d4 // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.1-0.20231027082548-f4a6c1f6e5c1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
