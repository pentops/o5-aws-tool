package auth

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/gen/j5/auth/v1/auth

import (
	time "time"
)

// Action Proto: Action
type Action struct {
	Method      string       `json:"method"`
	Actor       *Actor       `json:"actor"`
	Fingerprint *Fingerprint `json:"fingerprint,omitempty"`
}

// Actor Proto: Actor
type Actor struct {
	SubjectId            string                `json:"subjectId,omitempty"`
	SubjectType          string                `json:"subjectType,omitempty"`
	AuthenticationMethod *AuthenticationMethod `json:"authenticationMethod,omitempty"`
	Claim                *Claim                `json:"claim"`
	ActorTags            map[string]string     `json:"actorTags,omitempty"`
}

// AuthenticationMethod Proto Oneof: j5.auth.v1.AuthenticationMethod
type AuthenticationMethod struct {
	J5TypeKey string                         `json:"!type,omitempty"`
	Jwt       *AuthenticationMethod_JWT      `json:"jwt,omitempty"`
	Session   *AuthenticationMethod_Session  `json:"session,omitempty"`
	External  *AuthenticationMethod_External `json:"external,omitempty"`
}

func (s AuthenticationMethod) OneofKey() string {
	if s.Jwt != nil {
		return "jwt"
	}
	if s.Session != nil {
		return "session"
	}
	if s.External != nil {
		return "external"
	}
	return ""
}

func (s AuthenticationMethod) Type() interface{} {
	if s.Jwt != nil {
		return s.Jwt
	}
	if s.Session != nil {
		return s.Session
	}
	if s.External != nil {
		return s.External
	}
	return nil
}

// AuthenticationMethod_External Proto: AuthenticationMethod_External
type AuthenticationMethod_External struct {
	SystemName string            `json:"systemName,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// AuthenticationMethod_JWT Proto: AuthenticationMethod_JWT
type AuthenticationMethod_JWT struct {
	JwtId    string     `json:"jwtId,omitempty"`
	Issuer   string     `json:"issuer,omitempty"`
	IssuedAt *time.Time `json:"issuedAt,omitempty"`
}

// AuthenticationMethod_Session Proto: AuthenticationMethod_Session
type AuthenticationMethod_Session struct {
	SessionManager  string     `json:"sessionManager,omitempty"`
	SessionId       string     `json:"sessionId,omitempty"`
	VerifiedAt      *time.Time `json:"verifiedAt,omitempty"`
	AuthenticatedAt *time.Time `json:"authenticatedAt,omitempty"`
}

// Claim Proto: Claim
type Claim struct {
	RealmId    string   `json:"realmId,omitempty"`
	TenantType string   `json:"tenantType,omitempty"`
	TenantId   string   `json:"tenantId,omitempty"`
	Scopes     []string `json:"scopes,omitempty"`
}

// Fingerprint Proto: Fingerprint
type Fingerprint struct {
	IpAddress *string `json:"ipAddress,omitempty"`
	UserAgent *string `json:"userAgent,omitempty"`
}
