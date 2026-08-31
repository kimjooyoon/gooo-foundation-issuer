package issuer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func ParseOIDC(token string) (OIDCClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return OIDCClaims{}, fmt.Errorf("OIDC token is not a JWT")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return OIDCClaims{}, fmt.Errorf("OIDC payload decode failed: %w", err)
	}
	var raw struct {
		Issuer          string          `json:"iss"`
		Subject         string          `json:"sub"`
		Audience        json.RawMessage `json:"aud"`
		IssuedAt        int64           `json:"iat"`
		ExpiresAt       int64           `json:"exp"`
		Repository      string          `json:"repository"`
		RepositoryOwner string          `json:"repository_owner"`
		Workflow        string          `json:"workflow"`
		Ref             string          `json:"ref"`
		SHA             string          `json:"sha"`
		Actor           string          `json:"actor"`
	}
	if err := json.Unmarshal(payload, &raw); err != nil {
		return OIDCClaims{}, fmt.Errorf("OIDC claims decode failed: %w", err)
	}
	var audience string
	if err := json.Unmarshal(raw.Audience, &audience); err != nil {
		var audiences []string
		if listErr := json.Unmarshal(raw.Audience, &audiences); listErr != nil || len(audiences) == 0 {
			return OIDCClaims{}, fmt.Errorf("OIDC audience is missing")
		}
		audience = strings.Join(audiences, ",")
	}
	if raw.Issuer == "" || raw.Subject == "" || audience == "" || raw.IssuedAt == 0 || raw.ExpiresAt == 0 {
		return OIDCClaims{}, fmt.Errorf("required OIDC claim is missing")
	}
	return OIDCClaims{
		Issuer: raw.Issuer, Subject: raw.Subject, Audience: audience,
		IssuedAt: raw.IssuedAt, ExpiresAt: raw.ExpiresAt,
		Repository: raw.Repository, RepositoryOwner: raw.RepositoryOwner,
		Workflow: raw.Workflow, Ref: raw.Ref, SHA: raw.SHA, Actor: raw.Actor,
		RawObserved: true, SignatureState: Unknown,
	}, nil
}
