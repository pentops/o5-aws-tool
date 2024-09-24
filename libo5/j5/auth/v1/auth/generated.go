package auth

// Code generated by jsonapi. DO NOT EDIT.
// Source: github.com/pentops/o5-aws-tool/libo5/j5/auth/v1/auth

import (
	time "time"
)

// Claim Proto: Claim
type Claim struct {
	RealmId    string   `json:"realmId,omitempty"`
	TenantType string   `json:"tenantType,omitempty"`
	TenantId   string   `json:"tenantId,omitempty"`
	Scopes     []string `json:"scopes,omitempty"`
}

// Actor Proto: Actor
type Actor struct {
	SubjectId            string                `json:"subjectId,omitempty"`
	SubjectType          string                `json:"subjectType,omitempty"`
	AuthenticationMethod *AuthenticationMethod `json:"authenticationMethod,omitempty"`
	Claim                *Claim                `json:"claim"`
	ActorTags            map[string]string     `json:"actorTags,omitempty"`
}

// MethodAuthType Proto Message: MethodAuthType
type MethodAuthType struct {
	J5TypeKey string                    `json:"!type,omitempty"`
	None      *MethodAuthType_None      `json:"none,omitempty"`
	JwtBearer *MethodAuthType_JWTBearer `json:"jwtBearer,omitempty"`
	Custom    *MethodAuthType_Custom    `json:"custom,omitempty"`
	Cookie    *MethodAuthType_Cookie    `json:"cookie,omitempty"`
}

func (s MethodAuthType) OneofKey() string {
	if s.None != nil {
		return "none"
	}
	if s.JwtBearer != nil {
		return "jwtBearer"
	}
	if s.Custom != nil {
		return "custom"
	}
	if s.Cookie != nil {
		return "cookie"
	}
	return ""
}

func (s MethodAuthType) Type() interface{} {
	if s.None != nil {
		return s.None
	}
	if s.JwtBearer != nil {
		return s.JwtBearer
	}
	if s.Custom != nil {
		return s.Custom
	}
	if s.Cookie != nil {
		return s.Cookie
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

// MethodAuthType_None Proto: MethodAuthType_None
type MethodAuthType_None struct {
}

// Fingerprint Proto: Fingerprint
type Fingerprint struct {
	IpAddress *string `json:"ipAddress,omitempty"`
	UserAgent *string `json:"userAgent,omitempty"`
}

// AuthenticationMethod_Session Proto: AuthenticationMethod_Session
type AuthenticationMethod_Session struct {
	SessionManager  string     `json:"sessionManager,omitempty"`
	SessionId       string     `json:"sessionId,omitempty"`
	VerifiedAt      *time.Time `json:"verifiedAt,omitempty"`
	AuthenticatedAt *time.Time `json:"authenticatedAt,omitempty"`
}

// MethodAuthType_Cookie Proto: MethodAuthType_Cookie
type MethodAuthType_Cookie struct {
}

// Action Proto: Action
type Action struct {
	Method      string       `json:"method"`
	Actor       *Actor       `json:"actor"`
	Fingerprint *Fingerprint `json:"fingerprint,omitempty"`
}

// AuthenticationMethod Proto Message: AuthenticationMethod
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

// MethodAuthType_JWTBearer Proto: MethodAuthType_JWTBearer
type MethodAuthType_JWTBearer struct {
}

// MethodAuthType_Custom Proto: MethodAuthType_Custom
type MethodAuthType_Custom struct {
	PassThroughHeaders []string `json:"passThroughHeaders,omitempty"`
}
