package issuer

import "time"

type State string

const (
	Closed  State = "CLOSED"
	Unknown State = "UNKNOWN"
	Refuted State = "REFUTED"
)

const (
	ReceiptSchema = "gooo/foundation-authorization/issuer-receipt/v1"
	PayloadType   = "application/vnd.gooo.foundation-authorization+json;version=1"
)

type Request struct {
	Schema                 string   `json:"schema"`
	Repository             string   `json:"repository"`
	PullRequest            int      `json:"pull_request"`
	BaseRef                string   `json:"base_ref"`
	BaseSHA                string   `json:"base_sha"`
	HeadRef                string   `json:"head_ref"`
	HeadSHA                string   `json:"head_sha"`
	CandidateDigest        string   `json:"candidate_digest"`
	ChangedPaths           []string `json:"changed_paths"`
	ProtectedScope         []string `json:"protected_scope"`
	ProtectedScopeDigest   string   `json:"protected_scope_digest"`
	ActorIdentity          string   `json:"actor_identity"`
	IssuerIdentity         string   `json:"issuer_identity"`
	Nonce                  string   `json:"nonce"`
	ExternalHumanAuthority State    `json:"external_human_independence"`
}

type RotationAsset struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Size   int64  `json:"size_bytes"`
	SHA256 string `json:"sha256"`
}

type RotationInput struct {
	Repository   string          `json:"repository"`
	Version      string          `json:"version"`
	ReleaseID    int64           `json:"release_id"`
	Immutable    bool            `json:"immutable"`
	TagObjectSHA string          `json:"tag_object_sha"`
	TargetCommit string          `json:"target_commit"`
	Assets       []RotationAsset `json:"assets"`
}

type OIDCClaims struct {
	Issuer          string `json:"iss"`
	Subject         string `json:"sub"`
	Audience        string `json:"aud"`
	IssuedAt        int64  `json:"iat"`
	ExpiresAt       int64  `json:"exp"`
	Repository      string `json:"repository"`
	RepositoryOwner string `json:"repository_owner"`
	Workflow        string `json:"workflow"`
	Ref             string `json:"ref"`
	SHA             string `json:"sha"`
	Actor           string `json:"actor"`
	RawObserved     bool   `json:"raw_observed"`
	SignatureState  State  `json:"signature_state"`
}

type Payload struct {
	Schema                 string        `json:"schema"`
	Repository             string        `json:"repository"`
	PullRequest            int           `json:"pull_request"`
	BaseRef                string        `json:"base_ref"`
	BaseSHA                string        `json:"base_sha"`
	HeadRef                string        `json:"head_ref"`
	HeadSHA                string        `json:"head_sha"`
	CandidateDigest        string        `json:"candidate_digest"`
	ProtectedScope         []string      `json:"protected_scope"`
	ProtectedScopeDigest   string        `json:"protected_scope_digest"`
	ActorIdentity          string        `json:"actor_identity"`
	IssuerIdentity         string        `json:"issuer_identity"`
	Nonce                  string        `json:"nonce"`
	IssuedAt               string        `json:"issued_at"`
	ExpiresAt              string        `json:"expires_at"`
	Generation             int           `json:"generation"`
	PriorReceiptDigest     string        `json:"prior_receipt_digest"`
	Decision               State         `json:"decision"`
	IssuanceState          State         `json:"issuance_state"`
	HumanIndependenceState State         `json:"human_independence_state"`
	IntegrationState       State         `json:"integration_state"`
	ExternalAuthorityClaim string        `json:"external_authority_claim"`
	Rotation               RotationInput `json:"rotation"`
	OIDC                   OIDCClaims    `json:"oidc"`
}

type Receipt struct {
	Schema        string  `json:"schema"`
	PayloadType   string  `json:"payload_type"`
	Payload       Payload `json:"payload"`
	PayloadDigest string  `json:"payload_digest"`
	PublicKey     string  `json:"public_key"`
	Signature     string  `json:"signature"`
}

type VerificationReport struct {
	Schema                 string `json:"schema"`
	Decision               State  `json:"decision"`
	SignatureValid         bool   `json:"signature_valid"`
	TupleExact             bool   `json:"tuple_exact"`
	RotationExact          bool   `json:"rotation_exact"`
	OIDCClaimsObserved     bool   `json:"oidc_claims_observed"`
	OIDCSignatureState     State  `json:"oidc_signature_state"`
	HumanIndependenceState State  `json:"human_independence_state"`
	IntegrationState       State  `json:"integration_state"`
	Reason                 string `json:"reason"`
	Stage                  string `json:"stage"`
	Step                   string `json:"step"`
	UnknownClass           string `json:"unknown_class"`
	NextOperation          string `json:"next_operation"`
	BlockedBy              string `json:"blocked_by"`
}

type ReplayEvidence struct {
	Schema        string `json:"schema"`
	ReceiptDigest string `json:"receipt_digest"`
	FirstUse      State  `json:"first_use"`
	SecondUse     State  `json:"second_use"`
	RevokedKeyUse State  `json:"revoked_key_use"`
	Decision      State  `json:"decision"`
	Reason        string `json:"reason"`
}

func Combine(states ...State) State {
	for _, state := range states {
		if state == Refuted {
			return Refuted
		}
	}
	for _, state := range states {
		if state == Unknown {
			return Unknown
		}
	}
	return Closed
}

func (p Payload) Expired(now time.Time) bool {
	expires, err := time.Parse(time.RFC3339, p.ExpiresAt)
	return err == nil && !now.Before(expires)
}
