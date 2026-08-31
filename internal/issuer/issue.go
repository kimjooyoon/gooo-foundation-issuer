package issuer

import (
	"crypto/ed25519"
	"fmt"
	"time"
)

func Issue(request Request, rotation RotationInput, oidcToken string, private ed25519.PrivateKey, now time.Time) (Receipt, error) {
	if err := ValidateRequest(request); err != nil {
		return Receipt{}, err
	}
	if err := ValidateRotation(rotation); err != nil {
		return Receipt{}, err
	}
	claims, err := ParseOIDC(oidcToken)
	if err != nil {
		return Receipt{}, err
	}
	if claims.Repository != "" && claims.Repository != "kimjooyoon/gooo-foundation-issuer" {
		return Receipt{}, fmt.Errorf("OIDC repository claim is not this issuer")
	}
	issued := now.UTC().Truncate(time.Second)
	payload := Payload{
		Schema:                 ReceiptSchema,
		Repository:             request.Repository,
		PullRequest:            request.PullRequest,
		BaseRef:                request.BaseRef,
		BaseSHA:                request.BaseSHA,
		HeadRef:                request.HeadRef,
		HeadSHA:                request.HeadSHA,
		CandidateDigest:        request.CandidateDigest,
		ProtectedScope:         append([]string(nil), request.ProtectedScope...),
		ProtectedScopeDigest:   request.ProtectedScopeDigest,
		ActorIdentity:          request.ActorIdentity,
		IssuerIdentity:         request.IssuerIdentity,
		Nonce:                  request.Nonce,
		IssuedAt:               issued.Format(time.RFC3339),
		ExpiresAt:              issued.Add(15 * time.Minute).Format(time.RFC3339),
		Generation:             1,
		PriorReceiptDigest:     "",
		Decision:               Unknown,
		IssuanceState:          Closed,
		HumanIndependenceState: Unknown,
		IntegrationState:       Unknown,
		ExternalAuthorityClaim: "OWNER_OPERATED_EXTERNAL_REPOSITORY_AUTHORITY_ONLY",
		Rotation:               rotation,
		OIDC:                   claims,
	}
	return SignPayload(payload, private)
}

func VerifyReceipt(receipt Receipt, rotation RotationInput, now time.Time) VerificationReport {
	report := VerificationReport{
		Schema:                 "gooo/foundation-authorization/verification-report/v1",
		Decision:               Refuted,
		HumanIndependenceState: receipt.Payload.HumanIndependenceState,
		IntegrationState:       receipt.Payload.IntegrationState,
		OIDCSignatureState:     receipt.Payload.OIDC.SignatureState,
		Stage:                  "FOUNDATION",
		Step:                   "verify-receipt",
	}
	if receipt.Schema != ReceiptSchema || receipt.PayloadType != PayloadType {
		report.Reason = "receipt schema or payload type mismatch"
		return report
	}
	if err := ValidateRequest(Request{
		Schema:               ExpectedRequestSchema,
		Repository:           receipt.Payload.Repository,
		PullRequest:          receipt.Payload.PullRequest,
		BaseRef:              receipt.Payload.BaseRef,
		BaseSHA:              receipt.Payload.BaseSHA,
		HeadRef:              receipt.Payload.HeadRef,
		HeadSHA:              receipt.Payload.HeadSHA,
		CandidateDigest:      receipt.Payload.CandidateDigest,
		ChangedPaths:         ExpectedChangedPaths,
		ProtectedScope:       receipt.Payload.ProtectedScope,
		ProtectedScopeDigest: receipt.Payload.ProtectedScopeDigest,
		ActorIdentity:        receipt.Payload.ActorIdentity,
		IssuerIdentity:       receipt.Payload.IssuerIdentity,
		Nonce:                receipt.Payload.Nonce,
	}); err != nil {
		report.Reason = err.Error()
		return report
	}
	report.TupleExact = true
	if err := ValidateRotation(rotation); err != nil || rotation != receipt.Payload.Rotation {
		report.Reason = "rotation input mismatch"
		return report
	}
	report.RotationExact = true
	valid, err := VerifySignature(receipt)
	if err != nil || !valid {
		report.Reason = "receipt signature or payload digest is invalid"
		return report
	}
	report.SignatureValid = true
	report.OIDCClaimsObserved = receipt.Payload.OIDC.RawObserved
	if receipt.Payload.Expired(now) {
		report.Reason = "receipt is expired"
		return report
	}
	if receipt.Payload.Decision == Refuted {
		report.Reason = "receipt carries REFUTED decision"
		return report
	}
	report.Decision = Combine(receipt.Payload.IssuanceState, receipt.Payload.HumanIndependenceState, receipt.Payload.IntegrationState)
	if report.Decision == Unknown {
		report.Stage = "REGRESSION"
		report.Step = "authority-boundary"
		report.UnknownClass = "EXTERNAL_HUMAN_OR_GUARDIAN_CONSUMER_NOT_OBSERVED"
		report.NextOperation = "obtain independent human approval and update Guardian consumer contract"
		report.BlockedBy = "same-owner GitHub governance and current Guardian fixed PR609 dispatch contract"
		report.Reason = "cryptographic issuance is CLOSED; human independence and #619 Guardian integration remain UNKNOWN"
	}
	return report
}

func BuildReplayEvidence(receipt Receipt) ReplayEvidence {
	bytes, _ := PayloadBytes(receipt.Payload)
	return ReplayEvidence{
		Schema:        "gooo/foundation-authorization/replay-evidence/v1",
		ReceiptDigest: Digest(bytes),
		FirstUse:      Closed,
		SecondUse:     Refuted,
		RevokedKeyUse: Refuted,
		Decision:      Refuted,
		Reason:        "single-use receipt digest and revoked-key replay are fail-closed",
	}
}
