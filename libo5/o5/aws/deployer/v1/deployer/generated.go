package deployer

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/libo5/o5/aws/deployer/v1/deployer

import (
	context "context"
	errors "errors"
	list "github.com/pentops/o5-aws-tool/libo5/j5/list/v1/list"
	application "github.com/pentops/o5-aws-tool/libo5/o5/application/v1/application"
	environment "github.com/pentops/o5-aws-tool/libo5/o5/environment/v1/environment"
	messaging "github.com/pentops/o5-aws-tool/libo5/o5/messaging/v1/messaging"
	state "github.com/pentops/o5-aws-tool/libo5/psm/state/v1/state"
	url "net/url"
	strings "strings"
)

type Requester interface {
	Request(ctx context.Context, method string, path string, body interface{}, response interface{}) error
}

// DeploymentCommandService
type DeploymentCommandService struct {
	Requester
}

func NewDeploymentCommandService(requester Requester) *DeploymentCommandService {
	return &DeploymentCommandService{
		Requester: requester,
	}
}

func (s DeploymentCommandService) TriggerDeployment(ctx context.Context, req *TriggerDeploymentRequest) (*TriggerDeploymentResponse, error) {
	pathParts := make([]string, 6)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "deployments"
	if req.DeploymentId == "" {
		return nil, errors.New("required field \"DeploymentId\" not set")
	}
	pathParts[5] = req.DeploymentId
	path := strings.Join(pathParts, "/")
	resp := &TriggerDeploymentResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeploymentCommandService) TerminateDeployment(ctx context.Context, req *TerminateDeploymentRequest) (*TerminateDeploymentResponse, error) {
	pathParts := make([]string, 6)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "deployments"
	if req.DeploymentId == "" {
		return nil, errors.New("required field \"DeploymentId\" not set")
	}
	pathParts[5] = req.DeploymentId
	path := strings.Join(pathParts, "/")
	resp := &TerminateDeploymentResponse{}
	err := s.Request(ctx, "DELETE", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeploymentCommandService) UpsertCluster(ctx context.Context, req *UpsertClusterRequest) (*UpsertClusterResponse, error) {
	pathParts := make([]string, 6)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "clusters"
	pathParts[5] = "config"
	path := strings.Join(pathParts, "/")
	resp := &UpsertClusterResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeploymentCommandService) UpsertEnvironment(ctx context.Context, req *UpsertEnvironmentRequest) (*UpsertEnvironmentResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "environments"
	if req.EnvironmentId == "" {
		return nil, errors.New("required field \"EnvironmentId\" not set")
	}
	pathParts[5] = req.EnvironmentId
	pathParts[6] = "config"
	path := strings.Join(pathParts, "/")
	resp := &UpsertEnvironmentResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeploymentCommandService) UpsertStack(ctx context.Context, req *UpsertStackRequest) (*UpsertStackResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "c"
	pathParts[4] = "stacks"
	if req.StackId == "" {
		return nil, errors.New("required field \"StackId\" not set")
	}
	pathParts[5] = req.StackId
	pathParts[6] = "config"
	path := strings.Join(pathParts, "/")
	resp := &UpsertStackResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// EnvironmentQueryService
type EnvironmentQueryService struct {
	Requester
}

func NewEnvironmentQueryService(requester Requester) *EnvironmentQueryService {
	return &EnvironmentQueryService{
		Requester: requester,
	}
}

func (s EnvironmentQueryService) ListEnvironments(ctx context.Context, req *ListEnvironmentsRequest) (*ListEnvironmentsResponse, error) {
	pathParts := make([]string, 5)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "environments"
	path := strings.Join(pathParts, "/")
	resp := &ListEnvironmentsResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s EnvironmentQueryService) GetEnvironment(ctx context.Context, req *GetEnvironmentRequest) (*GetEnvironmentResponse, error) {
	pathParts := make([]string, 6)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "environment"
	if req.EnvironmentId == "" {
		return nil, errors.New("required field \"EnvironmentId\" not set")
	}
	pathParts[5] = req.EnvironmentId
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &GetEnvironmentResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s EnvironmentQueryService) ListEnvironmentEvents(ctx context.Context, req *ListEnvironmentEventsRequest) (*ListEnvironmentEventsResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "environment"
	if req.EnvironmentId == "" {
		return nil, errors.New("required field \"EnvironmentId\" not set")
	}
	pathParts[5] = req.EnvironmentId
	pathParts[6] = "events"
	path := strings.Join(pathParts, "/")
	resp := &ListEnvironmentEventsResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// StackQueryService
type StackQueryService struct {
	Requester
}

func NewStackQueryService(requester Requester) *StackQueryService {
	return &StackQueryService{
		Requester: requester,
	}
}

func (s StackQueryService) GetStack(ctx context.Context, req *GetStackRequest) (*GetStackResponse, error) {
	pathParts := make([]string, 6)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "stack"
	if req.StackId == "" {
		return nil, errors.New("required field \"StackId\" not set")
	}
	pathParts[5] = req.StackId
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &GetStackResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s StackQueryService) ListStacks(ctx context.Context, req *ListStacksRequest) (*ListStacksResponse, error) {
	pathParts := make([]string, 5)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "stacks"
	path := strings.Join(pathParts, "/")
	resp := &ListStacksResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s StackQueryService) ListStackEvents(ctx context.Context, req *ListStackEventsRequest) (*ListStackEventsResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "stack"
	if req.StackId == "" {
		return nil, errors.New("required field \"StackId\" not set")
	}
	pathParts[5] = req.StackId
	pathParts[6] = "events"
	path := strings.Join(pathParts, "/")
	resp := &ListStackEventsResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeploymentQueryService
type DeploymentQueryService struct {
	Requester
}

func NewDeploymentQueryService(requester Requester) *DeploymentQueryService {
	return &DeploymentQueryService{
		Requester: requester,
	}
}

func (s DeploymentQueryService) GetDeployment(ctx context.Context, req *GetDeploymentRequest) (*GetDeploymentResponse, error) {
	pathParts := make([]string, 6)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "deployment"
	if req.DeploymentId == "" {
		return nil, errors.New("required field \"DeploymentId\" not set")
	}
	pathParts[5] = req.DeploymentId
	path := strings.Join(pathParts, "/")
	if query, err := req.QueryParameters(); err != nil {
		return nil, err
	} else if len(query) > 0 {
		path += "?" + query.Encode()
	}
	resp := &GetDeploymentResponse{}
	err := s.Request(ctx, "GET", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeploymentQueryService) ListDeploymentEvents(ctx context.Context, req *ListDeploymentEventsRequest) (*ListDeploymentEventsResponse, error) {
	pathParts := make([]string, 7)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "deployment"
	if req.DeploymentId == "" {
		return nil, errors.New("required field \"DeploymentId\" not set")
	}
	pathParts[5] = req.DeploymentId
	pathParts[6] = "events"
	path := strings.Join(pathParts, "/")
	resp := &ListDeploymentEventsResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s DeploymentQueryService) ListDeployments(ctx context.Context, req *ListDeploymentsRequest) (*ListDeploymentsResponse, error) {
	pathParts := make([]string, 5)
	pathParts[0] = ""
	pathParts[1] = "deployer"
	pathParts[2] = "v1"
	pathParts[3] = "q"
	pathParts[4] = "deployments"
	path := strings.Join(pathParts, "/")
	resp := &ListDeploymentsResponse{}
	err := s.Request(ctx, "POST", path, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpsertStackResponse
type UpsertStackResponse struct {
	State *StackState `json:"state"`
}

// StepRequestType_PGCleanup Proto: StepRequestType_PGCleanup
type StepRequestType_PGCleanup struct {
	Spec *PostgresSpec `json:"spec,omitempty"`
}

// PostgresSpec Proto: PostgresSpec
type PostgresSpec struct {
	DbName                  string   `json:"dbName,omitempty"`
	DbExtensions            []string `json:"dbExtensions,omitempty"`
	RootSecretName          string   `json:"rootSecretName,omitempty"`
	MigrationTaskOutputName *string  `json:"migrationTaskOutputName,omitempty"`
	SecretOutputName        string   `json:"secretOutputName,omitempty"`
}

// ClusterEventType_Configured Proto: ClusterEventType_Configured
type ClusterEventType_Configured struct {
	Config *environment.Cluster `json:"config,omitempty"`
}

// DeploymentEventType_Done Proto: DeploymentEventType_Done
type DeploymentEventType_Done struct {
}

// TerminateDeploymentRequest
type TerminateDeploymentRequest struct {
	DeploymentId string `json:"-" path:"deploymentId"`
}

// UpsertEnvironmentRequest
type UpsertEnvironmentRequest struct {
	ClusterId     string                   `json:"clusterId,omitempty"`
	Config        *environment.Environment `json:"config,omitempty"`
	ConfigYaml    []byte                   `json:"configYaml,omitempty"`
	ConfigJson    []byte                   `json:"configJson,omitempty"`
	EnvironmentId string                   `json:"-" path:"environmentId"`
}

// ListDeploymentsResponse
type ListDeploymentsResponse struct {
	Deployments []*DeploymentState `json:"deployments,omitempty"`
	Page        *list.PageResponse `json:"page,omitempty"`
}

func (s ListDeploymentsResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListDeploymentsResponse) GetItems() []*DeploymentState {
	return s.Deployments
}

// DeploymentEventType_StackWait Proto: DeploymentEventType_StackWait
type DeploymentEventType_StackWait struct {
}

// StackEventType_RunDeployment Proto: StackEventType_RunDeployment
type StackEventType_RunDeployment struct {
	DeploymentId string `json:"deploymentId,omitempty"`
}

// StepRequestType Proto Message: StepRequestType
type StepRequestType struct {
	J5TypeKey  string                      `json:"!type,omitempty"`
	EvalJoin   *StepRequestType_EvalJoin   `json:"evalJoin,omitempty"`
	CfCreate   *StepRequestType_CFCreate   `json:"cfCreate,omitempty"`
	CfPlan     *StepRequestType_CFPlan     `json:"cfPlan,omitempty"`
	CfUpdate   *StepRequestType_CFUpdate   `json:"cfUpdate,omitempty"`
	CfScale    *StepRequestType_CFScale    `json:"cfScale,omitempty"`
	PgUpsert   *StepRequestType_PGUpsert   `json:"pgUpsert,omitempty"`
	PgEvaluate *StepRequestType_PGEvaluate `json:"pgEvaluate,omitempty"`
	PgCleanup  *StepRequestType_PGCleanup  `json:"pgCleanup,omitempty"`
	PgMigrate  *StepRequestType_PGMigrate  `json:"pgMigrate,omitempty"`
}

func (s StepRequestType) OneofKey() string {
	if s.EvalJoin != nil {
		return "evalJoin"
	}
	if s.CfCreate != nil {
		return "cfCreate"
	}
	if s.CfPlan != nil {
		return "cfPlan"
	}
	if s.CfUpdate != nil {
		return "cfUpdate"
	}
	if s.CfScale != nil {
		return "cfScale"
	}
	if s.PgUpsert != nil {
		return "pgUpsert"
	}
	if s.PgEvaluate != nil {
		return "pgEvaluate"
	}
	if s.PgCleanup != nil {
		return "pgCleanup"
	}
	if s.PgMigrate != nil {
		return "pgMigrate"
	}
	return ""
}

func (s StepRequestType) Type() interface{} {
	if s.EvalJoin != nil {
		return s.EvalJoin
	}
	if s.CfCreate != nil {
		return s.CfCreate
	}
	if s.CfPlan != nil {
		return s.CfPlan
	}
	if s.CfUpdate != nil {
		return s.CfUpdate
	}
	if s.CfScale != nil {
		return s.CfScale
	}
	if s.PgUpsert != nil {
		return s.PgUpsert
	}
	if s.PgEvaluate != nil {
		return s.PgEvaluate
	}
	if s.PgCleanup != nil {
		return s.PgCleanup
	}
	if s.PgMigrate != nil {
		return s.PgMigrate
	}
	return nil
}

// DeploymentEventType_StackWaitFailure Proto: DeploymentEventType_StackWaitFailure
type DeploymentEventType_StackWaitFailure struct {
	Error string `json:"error,omitempty"`
}

// ListEnvironmentsRequest
type ListEnvironmentsRequest struct {
	Page  *list.PageRequest  `json:"page,omitempty"`
	Query *list.QueryRequest `json:"query,omitempty"`
}

func (s *ListEnvironmentsRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

// ListEnvironmentEventsRequest
type ListEnvironmentEventsRequest struct {
	Page          *list.PageRequest  `json:"page,omitempty"`
	Query         *list.QueryRequest `json:"query,omitempty"`
	EnvironmentId string             `json:"-" path:"environmentId"`
}

func (s *ListEnvironmentEventsRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

// DeploymentEventType_StackAvailable Proto: DeploymentEventType_StackAvailable
type DeploymentEventType_StackAvailable struct {
	StackOutput *CFStackOutput `json:"stackOutput,omitempty"`
}

// StackDeployment Proto: StackDeployment
type StackDeployment struct {
	DeploymentId string `json:"deploymentId,omitempty"`
	Version      string `json:"version,omitempty"`
}

// GetStackResponse
type GetStackResponse struct {
	State  *StackState   `json:"state,omitempty"`
	Events []*StackEvent `json:"events,omitempty"`
}

// StepRequestType_PGUpsert Proto: StepRequestType_PGUpsert
type StepRequestType_PGUpsert struct {
	Spec              *PostgresSpec `json:"spec,omitempty"`
	InfraOutputStepId string        `json:"infraOutputStepId,omitempty"`
	RotateCredentials bool          `json:"rotateCredentials"`
}

// CFChangesetOutput Proto: CFChangesetOutput
type CFChangesetOutput struct {
	Lifecycle string `json:"lifecycle,omitempty"`
}

// UpsertEnvironmentResponse
type UpsertEnvironmentResponse struct {
	State *EnvironmentState `json:"state"`
}

// GetStackRequest
type GetStackRequest struct {
	StackId string `json:"-" path:"stackId"`
}

func (s GetStackRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	return values, nil
}

// DeploymentStep Proto: DeploymentStep
type DeploymentStep struct {
	Id        string           `json:"id,omitempty"`
	Name      string           `json:"name,omitempty"`
	Status    string           `json:"status,omitempty"`
	Request   *StepRequestType `json:"request,omitempty"`
	Output    *StepOutputType  `json:"output,omitempty"`
	Error     *string          `json:"error,omitempty"`
	DependsOn []string         `json:"dependsOn,omitempty"`
}

// StackEvent Proto: StackEvent
type StackEvent struct {
	Metadata *state.EventMetadata `json:"metadata"`
	*StackKeys
	Event *StackEventType `json:"event"`
}

// ClusterEventType Proto Message: ClusterEventType
type ClusterEventType struct {
	J5TypeKey  string                       `json:"!type,omitempty"`
	Configured *ClusterEventType_Configured `json:"configured,omitempty"`
}

func (s ClusterEventType) OneofKey() string {
	if s.Configured != nil {
		return "configured"
	}
	return ""
}

func (s ClusterEventType) Type() interface{} {
	if s.Configured != nil {
		return s.Configured
	}
	return nil
}

// DeploymentEventType_Triggered Proto: DeploymentEventType_Triggered
type DeploymentEventType_Triggered struct {
}

// UpsertStackRequest
type UpsertStackRequest struct {
	StackId string `json:"-" path:"stackId"`
}

// ListDeploymentsRequest
type ListDeploymentsRequest struct {
	Page  *list.PageRequest  `json:"page,omitempty"`
	Query *list.QueryRequest `json:"query,omitempty"`
}

func (s *ListDeploymentsRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

// StepOutputType_CFStatus Proto: StepOutputType_CFStatus
type StepOutputType_CFStatus struct {
	Output *CFStackOutput `json:"output,omitempty"`
}

// EnvironmentKeys Proto: EnvironmentKeys
type EnvironmentKeys struct {
	EnvironmentId string `json:"environmentId,omitempty"`
	ClusterId     string `json:"clusterId,omitempty"`
}

// KeyValue Proto: KeyValue
type KeyValue struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// DeploymentEventType_RunSteps Proto: DeploymentEventType_RunSteps
type DeploymentEventType_RunSteps struct {
	Steps []*DeploymentStep `json:"steps,omitempty"`
}

// UpsertClusterRequest
type UpsertClusterRequest struct {
	ClusterId  string                      `json:"clusterId,omitempty"`
	Config     *environment.CombinedConfig `json:"config,omitempty"`
	ConfigYaml []byte                      `json:"configYaml,omitempty"`
	ConfigJson []byte                      `json:"configJson,omitempty"`
}

// StepRequestType_EvalJoin Proto: StepRequestType_EvalJoin
type StepRequestType_EvalJoin struct {
}

// CFStackInput Proto: CFStackInput
type CFStackInput struct {
	StackName    string                          `json:"stackName,omitempty"`
	S3Template   *S3Template                     `json:"s3Template,omitempty"`
	TemplateBody string                          `json:"templateBody,omitempty"`
	EmptyStack   bool                            `json:"emptyStack"`
	DesiredCount int32                           `json:"desiredCount,omitempty"`
	Parameters   []*CloudFormationStackParameter `json:"parameters,omitempty"`
	SnsTopics    []string                        `json:"snsTopics,omitempty"`
}

// DeploymentState Proto: DeploymentState
type DeploymentState struct {
	Metadata *state.StateMetadata `json:"metadata"`
	*DeploymentKeys
	Status string               `json:"status,omitempty"`
	Data   *DeploymentStateData `json:"data,omitempty"`
}

// DeploymentEventType_Error Proto: DeploymentEventType_Error
type DeploymentEventType_Error struct {
	Error string `json:"error,omitempty"`
}

// StackKeys Proto: StackKeys
type StackKeys struct {
	StackId       string `json:"stackId,omitempty"`
	EnvironmentId string `json:"environmentId,omitempty"`
	ClusterId     string `json:"clusterId,omitempty"`
}

// ListStackEventsResponse
type ListStackEventsResponse struct {
	Events []*StackEvent      `json:"events,omitempty"`
	Page   *list.PageResponse `json:"page,omitempty"`
}

func (s ListStackEventsResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListStackEventsResponse) GetItems() []*StackEvent {
	return s.Events
}

// ClusterStateData Proto: ClusterStateData
type ClusterStateData struct {
	Config *environment.Cluster `json:"config,omitempty"`
}

// StackState Proto: StackState
type StackState struct {
	Metadata *state.StateMetadata `json:"metadata"`
	*StackKeys
	Status string          `json:"status,omitempty"`
	Data   *StackStateData `json:"data,omitempty"`
}

// ListStackEventsRequest
type ListStackEventsRequest struct {
	Page    *list.PageRequest  `json:"page,omitempty"`
	Query   *list.QueryRequest `json:"query,omitempty"`
	StackId string             `json:"-" path:"stackId"`
}

func (s *ListStackEventsRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

// GetDeploymentResponse
type GetDeploymentResponse struct {
	State  *DeploymentState   `json:"state,omitempty"`
	Events []*DeploymentEvent `json:"events,omitempty"`
}

// DeploymentEventType_Created Proto: DeploymentEventType_Created
type DeploymentEventType_Created struct {
	Request *messaging.RequestMetadata `json:"request,omitempty"`
	Spec    *DeploymentSpec            `json:"spec,omitempty"`
}

// S3Template Proto: S3Template
type S3Template struct {
	Bucket string `json:"bucket,omitempty"`
	Key    string `json:"key,omitempty"`
	Region string `json:"region,omitempty"`
}

// EnvironmentEventType_Configured Proto: EnvironmentEventType_Configured
type EnvironmentEventType_Configured struct {
	Config *environment.Environment `json:"config,omitempty"`
}

// CFStackOutput Proto: CFStackOutput
type CFStackOutput struct {
	Lifecycle string      `json:"lifecycle,omitempty"`
	Outputs   []*KeyValue `json:"outputs,omitempty"`
}

// CloudFormationStackParameterType_DesiredCount Proto: CloudFormationStackParameterType_DesiredCount
type CloudFormationStackParameterType_DesiredCount struct {
}

// CombinedClient
type CombinedClient struct {
	*DeploymentCommandService
	*EnvironmentQueryService
	*StackQueryService
	*DeploymentQueryService
}

func NewCombinedClient(requester Requester) *CombinedClient {
	return &CombinedClient{
		DeploymentCommandService: NewDeploymentCommandService(requester),
		EnvironmentQueryService:  NewEnvironmentQueryService(requester),
		StackQueryService:        NewStackQueryService(requester),
		DeploymentQueryService:   NewDeploymentQueryService(requester),
	}
}

// GetEnvironmentRequest
type GetEnvironmentRequest struct {
	EnvironmentId string `json:"-" path:"environmentId"`
}

func (s GetEnvironmentRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	return values, nil
}

// ListDeploymentEventsRequest
type ListDeploymentEventsRequest struct {
	Page         *list.PageRequest  `json:"page,omitempty"`
	Query        *list.QueryRequest `json:"query,omitempty"`
	DeploymentId string             `json:"-" path:"deploymentId"`
}

func (s *ListDeploymentEventsRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

// StepRequestType_CFScale Proto: StepRequestType_CFScale
type StepRequestType_CFScale struct {
	StackName    string `json:"stackName,omitempty"`
	DesiredCount int32  `json:"desiredCount,omitempty"`
}

// StepOutputType_CFPlanStatus Proto: StepOutputType_CFPlanStatus
type StepOutputType_CFPlanStatus struct {
	Output *CFChangesetOutput `json:"output,omitempty"`
}

// StepOutputType Proto Message: StepOutputType
type StepOutputType struct {
	J5TypeKey    string                       `json:"!type,omitempty"`
	CfStatus     *StepOutputType_CFStatus     `json:"cfStatus,omitempty"`
	CfPlanStatus *StepOutputType_CFPlanStatus `json:"cfPlanStatus,omitempty"`
}

func (s StepOutputType) OneofKey() string {
	if s.CfStatus != nil {
		return "cfStatus"
	}
	if s.CfPlanStatus != nil {
		return "cfPlanStatus"
	}
	return ""
}

func (s StepOutputType) Type() interface{} {
	if s.CfStatus != nil {
		return s.CfStatus
	}
	if s.CfPlanStatus != nil {
		return s.CfPlanStatus
	}
	return nil
}

// StackStateData Proto: StackStateData
type StackStateData struct {
	CurrentDeployment *StackDeployment   `json:"currentDeployment,omitempty"`
	StackName         string             `json:"stackName,omitempty"`
	ApplicationName   string             `json:"applicationName,omitempty"`
	EnvironmentName   string             `json:"environmentName,omitempty"`
	EnvironmentId     string             `json:"environmentId,omitempty"`
	QueuedDeployments []*StackDeployment `json:"queuedDeployments,omitempty"`
}

// DeploymentEventType_Terminated Proto: DeploymentEventType_Terminated
type DeploymentEventType_Terminated struct {
}

// EnvironmentStateData Proto: EnvironmentStateData
type EnvironmentStateData struct {
	Config *environment.Environment `json:"config,omitempty"`
}

// ListStacksRequest
type ListStacksRequest struct {
	Page  *list.PageRequest  `json:"page,omitempty"`
	Query *list.QueryRequest `json:"query,omitempty"`
}

func (s *ListStacksRequest) SetPageToken(pageToken string) {
	if s.Page == nil {
		s.Page = &list.PageRequest{}
	}
	s.Page.Token = &pageToken
}

// GetDeploymentRequest
type GetDeploymentRequest struct {
	DeploymentId string `json:"-" path:"deploymentId"`
}

func (s GetDeploymentRequest) QueryParameters() (url.Values, error) {
	values := url.Values{}
	return values, nil
}

// EnvironmentState Proto: EnvironmentState
type EnvironmentState struct {
	Metadata *state.StateMetadata `json:"metadata"`
	*EnvironmentKeys
	Status string                `json:"status,omitempty"`
	Data   *EnvironmentStateData `json:"data,omitempty"`
}

// DeploymentEvent Proto: DeploymentEvent
type DeploymentEvent struct {
	Metadata *state.EventMetadata `json:"metadata"`
	Keys     *DeploymentKeys      `json:"keys"`
	Event    *DeploymentEventType `json:"event"`
}

// CloudFormationStackParameterType Proto Message: CloudFormationStackParameterType
type CloudFormationStackParameterType struct {
	J5TypeKey    string                                         `json:"!type,omitempty"`
	RulePriority *CloudFormationStackParameterType_RulePriority `json:"rulePriority,omitempty"`
	DesiredCount *CloudFormationStackParameterType_DesiredCount `json:"desiredCount,omitempty"`
}

func (s CloudFormationStackParameterType) OneofKey() string {
	if s.RulePriority != nil {
		return "rulePriority"
	}
	if s.DesiredCount != nil {
		return "desiredCount"
	}
	return ""
}

func (s CloudFormationStackParameterType) Type() interface{} {
	if s.RulePriority != nil {
		return s.RulePriority
	}
	if s.DesiredCount != nil {
		return s.DesiredCount
	}
	return nil
}

// StackEventType_DeploymentCompleted Proto: StackEventType_DeploymentCompleted
type StackEventType_DeploymentCompleted struct {
	Deployment *StackDeployment `json:"deployment,omitempty"`
}

// ClusterKeys Proto: ClusterKeys
type ClusterKeys struct {
	ClusterId string `json:"clusterId,omitempty"`
}

// StackEventType Proto Message: StackEventType
type StackEventType struct {
	J5TypeKey           string                              `json:"!type,omitempty"`
	Configured          *StackEventType_Configured          `json:"configured,omitempty"`
	DeploymentRequested *StackEventType_DeploymentRequested `json:"deploymentRequested,omitempty"`
	DeploymentCompleted *StackEventType_DeploymentCompleted `json:"deploymentCompleted,omitempty"`
	DeploymentFailed    *StackEventType_DeploymentFailed    `json:"deploymentFailed,omitempty"`
	RunDeployment       *StackEventType_RunDeployment       `json:"runDeployment,omitempty"`
}

func (s StackEventType) OneofKey() string {
	if s.Configured != nil {
		return "configured"
	}
	if s.DeploymentRequested != nil {
		return "deploymentRequested"
	}
	if s.DeploymentCompleted != nil {
		return "deploymentCompleted"
	}
	if s.DeploymentFailed != nil {
		return "deploymentFailed"
	}
	if s.RunDeployment != nil {
		return "runDeployment"
	}
	return ""
}

func (s StackEventType) Type() interface{} {
	if s.Configured != nil {
		return s.Configured
	}
	if s.DeploymentRequested != nil {
		return s.DeploymentRequested
	}
	if s.DeploymentCompleted != nil {
		return s.DeploymentCompleted
	}
	if s.DeploymentFailed != nil {
		return s.DeploymentFailed
	}
	if s.RunDeployment != nil {
		return s.RunDeployment
	}
	return nil
}

// TerminateDeploymentResponse
type TerminateDeploymentResponse struct {
}

// DeploymentEventType_StepResult Proto: DeploymentEventType_StepResult
type DeploymentEventType_StepResult struct {
	StepId string          `json:"stepId,omitempty"`
	Status string          `json:"status,omitempty"`
	Output *StepOutputType `json:"output,omitempty"`
	Error  *string         `json:"error,omitempty"`
}

// EnvironmentEvent Proto: EnvironmentEvent
type EnvironmentEvent struct {
	Metadata *state.EventMetadata `json:"metadata"`
	*EnvironmentKeys
	Event *EnvironmentEventType `json:"event"`
}

// DeploymentSpec Proto: DeploymentSpec
type DeploymentSpec struct {
	AppName         string                          `json:"appName,omitempty"`
	Version         string                          `json:"version,omitempty"`
	EnvironmentName string                          `json:"environmentName,omitempty"`
	EnvironmentId   string                          `json:"environmentId,omitempty"`
	Template        *S3Template                     `json:"template,omitempty"`
	EcsCluster      string                          `json:"ecsCluster,omitempty"`
	CfStackName     string                          `json:"cfStackName,omitempty"`
	Flags           *DeploymentFlags                `json:"flags,omitempty"`
	Databases       []*PostgresSpec                 `json:"databases,omitempty"`
	Parameters      []*CloudFormationStackParameter `json:"parameters,omitempty"`
	SnsTopics       []string                        `json:"snsTopics,omitempty"`
}

// DeploymentFlags Proto: DeploymentFlags
type DeploymentFlags struct {
	QuickMode         bool `json:"quickMode"`
	RotateCredentials bool `json:"rotateCredentials"`
	CancelUpdates     bool `json:"cancelUpdates"`
	DbOnly            bool `json:"dbOnly"`
	InfraOnly         bool `json:"infraOnly"`
	ImportResources   bool `json:"importResources"`
}

// EnvironmentEventType Proto Message: EnvironmentEventType
type EnvironmentEventType struct {
	J5TypeKey  string                           `json:"!type,omitempty"`
	Configured *EnvironmentEventType_Configured `json:"configured,omitempty"`
}

func (s EnvironmentEventType) OneofKey() string {
	if s.Configured != nil {
		return "configured"
	}
	return ""
}

func (s EnvironmentEventType) Type() interface{} {
	if s.Configured != nil {
		return s.Configured
	}
	return nil
}

// ClusterState Proto: ClusterState
type ClusterState struct {
	Metadata *state.StateMetadata `json:"metadata"`
	*ClusterKeys
	Status string            `json:"status,omitempty"`
	Data   *ClusterStateData `json:"data,omitempty"`
}

// StepRequestType_CFCreate Proto: StepRequestType_CFCreate
type StepRequestType_CFCreate struct {
	Spec       *CFStackInput  `json:"spec,omitempty"`
	Output     *CFStackOutput `json:"output,omitempty"`
	EmptyStack bool           `json:"emptyStack"`
}

// TriggerDeploymentRequest
type TriggerDeploymentRequest struct {
	Environment  string           `json:"environment,omitempty"`
	Source       *TriggerSource   `json:"source,omitempty"`
	Flags        *DeploymentFlags `json:"flags,omitempty"`
	DeploymentId string           `json:"-" path:"deploymentId"`
}

// StackEventType_DeploymentFailed Proto: StackEventType_DeploymentFailed
type StackEventType_DeploymentFailed struct {
	Deployment *StackDeployment `json:"deployment,omitempty"`
	Error      string           `json:"error,omitempty"`
}

// GetEnvironmentResponse
type GetEnvironmentResponse struct {
	State  *EnvironmentState   `json:"state,omitempty"`
	Events []*EnvironmentEvent `json:"events,omitempty"`
}

// StackEventType_DeploymentRequested Proto: StackEventType_DeploymentRequested
type StackEventType_DeploymentRequested struct {
	Deployment      *StackDeployment `json:"deployment,omitempty"`
	ApplicationName string           `json:"applicationName,omitempty"`
	EnvironmentName string           `json:"environmentName,omitempty"`
	EnvironmentId   string           `json:"environmentId,omitempty"`
}

// DeploymentKeys Proto: DeploymentKeys
type DeploymentKeys struct {
	DeploymentId  string `json:"deploymentId,omitempty"`
	StackId       string `json:"stackId,omitempty"`
	EnvironmentId string `json:"environmentId,omitempty"`
	ClusterId     string `json:"clusterId,omitempty"`
}

// TriggerSource Proto Message: TriggerSource
type TriggerSource struct {
	J5TypeKey string                      `json:"!type,omitempty"`
	Github    *TriggerSource_GithubSource `json:"github,omitempty"`
	Inline    *TriggerSource_InlineSource `json:"inline,omitempty"`
}

func (s TriggerSource) OneofKey() string {
	if s.Github != nil {
		return "github"
	}
	if s.Inline != nil {
		return "inline"
	}
	return ""
}

func (s TriggerSource) Type() interface{} {
	if s.Github != nil {
		return s.Github
	}
	if s.Inline != nil {
		return s.Inline
	}
	return nil
}

// CloudFormationStackParameterType_RulePriority Proto: CloudFormationStackParameterType_RulePriority
type CloudFormationStackParameterType_RulePriority struct {
	RouteGroup string `json:"routeGroup,omitempty"`
}

// DeploymentEventType Proto Message: DeploymentEventType
type DeploymentEventType struct {
	J5TypeKey        string                                `json:"!type,omitempty"`
	Created          *DeploymentEventType_Created          `json:"created,omitempty"`
	Triggered        *DeploymentEventType_Triggered        `json:"triggered,omitempty"`
	StackWait        *DeploymentEventType_StackWait        `json:"stackWait,omitempty"`
	StackWaitFailure *DeploymentEventType_StackWaitFailure `json:"stackWaitFailure,omitempty"`
	StackAvailable   *DeploymentEventType_StackAvailable   `json:"stackAvailable,omitempty"`
	RunSteps         *DeploymentEventType_RunSteps         `json:"runSteps,omitempty"`
	StepResult       *DeploymentEventType_StepResult       `json:"stepResult,omitempty"`
	RunStep          *DeploymentEventType_RunStep          `json:"runStep,omitempty"`
	Error            *DeploymentEventType_Error            `json:"error,omitempty"`
	Done             *DeploymentEventType_Done             `json:"done,omitempty"`
	Terminated       *DeploymentEventType_Terminated       `json:"terminated,omitempty"`
}

func (s DeploymentEventType) OneofKey() string {
	if s.Created != nil {
		return "created"
	}
	if s.Triggered != nil {
		return "triggered"
	}
	if s.StackWait != nil {
		return "stackWait"
	}
	if s.StackWaitFailure != nil {
		return "stackWaitFailure"
	}
	if s.StackAvailable != nil {
		return "stackAvailable"
	}
	if s.RunSteps != nil {
		return "runSteps"
	}
	if s.StepResult != nil {
		return "stepResult"
	}
	if s.RunStep != nil {
		return "runStep"
	}
	if s.Error != nil {
		return "error"
	}
	if s.Done != nil {
		return "done"
	}
	if s.Terminated != nil {
		return "terminated"
	}
	return ""
}

func (s DeploymentEventType) Type() interface{} {
	if s.Created != nil {
		return s.Created
	}
	if s.Triggered != nil {
		return s.Triggered
	}
	if s.StackWait != nil {
		return s.StackWait
	}
	if s.StackWaitFailure != nil {
		return s.StackWaitFailure
	}
	if s.StackAvailable != nil {
		return s.StackAvailable
	}
	if s.RunSteps != nil {
		return s.RunSteps
	}
	if s.StepResult != nil {
		return s.StepResult
	}
	if s.RunStep != nil {
		return s.RunStep
	}
	if s.Error != nil {
		return s.Error
	}
	if s.Done != nil {
		return s.Done
	}
	if s.Terminated != nil {
		return s.Terminated
	}
	return nil
}

// UpsertClusterResponse
type UpsertClusterResponse struct {
	State *ClusterState `json:"state"`
}

// ListEnvironmentsResponse
type ListEnvironmentsResponse struct {
	Environments []*EnvironmentState `json:"environments,omitempty"`
	Page         *list.PageResponse  `json:"page,omitempty"`
}

func (s ListEnvironmentsResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListEnvironmentsResponse) GetItems() []*EnvironmentState {
	return s.Environments
}

// DeploymentStateData Proto: DeploymentStateData
type DeploymentStateData struct {
	Request *messaging.RequestMetadata `json:"request,omitempty"`
	Spec    *DeploymentSpec            `json:"spec,omitempty"`
	Steps   []*DeploymentStep          `json:"steps,omitempty"`
}

// StepRequestType_CFPlan Proto: StepRequestType_CFPlan
type StepRequestType_CFPlan struct {
	Spec            *CFStackInput `json:"spec,omitempty"`
	ImportResources bool          `json:"importResources"`
}

// CloudFormationStackParameter Proto: CloudFormationStackParameter
type CloudFormationStackParameter struct {
	Name    string                            `json:"name,omitempty"`
	Value   string                            `json:"value,omitempty"`
	Resolve *CloudFormationStackParameterType `json:"resolve,omitempty"`
}

// StepRequestType_PGEvaluate Proto: StepRequestType_PGEvaluate
type StepRequestType_PGEvaluate struct {
	DbName string `json:"dbName,omitempty"`
}

// DeploymentEventType_RunStep Proto: DeploymentEventType_RunStep
type DeploymentEventType_RunStep struct {
	StepId string `json:"stepId,omitempty"`
}

// ListStacksResponse
type ListStacksResponse struct {
	Stacks []*StackState      `json:"stacks,omitempty"`
	Page   *list.PageResponse `json:"page,omitempty"`
}

func (s ListStacksResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListStacksResponse) GetItems() []*StackState {
	return s.Stacks
}

// StepRequestType_PGMigrate Proto: StepRequestType_PGMigrate
type StepRequestType_PGMigrate struct {
	Spec              *PostgresSpec `json:"spec,omitempty"`
	InfraOutputStepId string        `json:"infraOutputStepId,omitempty"`
	EcsClusterName    string        `json:"ecsClusterName,omitempty"`
}

// ListDeploymentEventsResponse
type ListDeploymentEventsResponse struct {
	Events []*DeploymentEvent `json:"events,omitempty"`
	Page   *list.PageResponse `json:"page,omitempty"`
}

func (s ListDeploymentEventsResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListDeploymentEventsResponse) GetItems() []*DeploymentEvent {
	return s.Events
}

// StepRequestType_CFUpdate Proto: StepRequestType_CFUpdate
type StepRequestType_CFUpdate struct {
	Spec   *CFStackInput  `json:"spec,omitempty"`
	Output *CFStackOutput `json:"output,omitempty"`
}

// StackEventType_Configured Proto: StackEventType_Configured
type StackEventType_Configured struct {
	ApplicationName string `json:"applicationName,omitempty"`
	EnvironmentId   string `json:"environmentId,omitempty"`
	EnvironmentName string `json:"environmentName,omitempty"`
}

// TriggerSource_InlineSource Proto: TriggerSource_InlineSource
type TriggerSource_InlineSource struct {
	Version     string                   `json:"version"`
	Application *application.Application `json:"application"`
}

// TriggerSource_GithubSource Proto: TriggerSource_GithubSource
type TriggerSource_GithubSource struct {
	Owner  string `json:"owner"`
	Repo   string `json:"repo"`
	Branch string `json:"branch,omitempty"`
	Tag    string `json:"tag,omitempty"`
	Commit string `json:"commit,omitempty"`
}

// TriggerDeploymentResponse
type TriggerDeploymentResponse struct {
	DeploymentId  string `json:"deploymentId,omitempty"`
	EnvironmentId string `json:"environmentId,omitempty"`
	ClusterId     string `json:"clusterId,omitempty"`
}

// ListEnvironmentEventsResponse
type ListEnvironmentEventsResponse struct {
	Events []*EnvironmentEvent `json:"events,omitempty"`
	Page   *list.PageResponse  `json:"page,omitempty"`
}

func (s ListEnvironmentEventsResponse) GetPageToken() *string {
	if s.Page == nil {
		return nil
	}
	return s.Page.NextToken
}

func (s ListEnvironmentEventsResponse) GetItems() []*EnvironmentEvent {
	return s.Events
}
