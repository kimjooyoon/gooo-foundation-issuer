package issuer

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

func testRequest() Request {
	return Request{
		Schema: ExpectedRequestSchema, Repository: ExpectedRepository, PullRequest: ExpectedPullRequest,
		BaseRef: ExpectedBaseRef, BaseSHA: ExpectedBaseSHA, HeadRef: ExpectedHeadRef, HeadSHA: ExpectedHeadSHA,
		CandidateDigest: ExpectedCandidateDigest, ChangedPaths: append([]string(nil), ExpectedChangedPaths...),
		ProtectedScope: append([]string(nil), ExpectedProtectedScope...), ProtectedScopeDigest: ExpectedProtectedScopeDigest,
		ActorIdentity: "kimjooyoon", IssuerIdentity: "github.com/kimjooyoon/gooo-foundation-issuer", Nonce: "test-nonce",
	}
}

func testRotation() RotationInput {
	return RotationInput{
		Repository: ExpectedRotationRepository, Version: ExpectedRotationVersion, ReleaseID: ExpectedRotationReleaseID,
		Immutable: true, TagObjectSHA: ExpectedRotationTagObject, TargetCommit: ExpectedRotationTarget,
		Assets: []RotationAsset{
			{ID: 538516119, Name: "gooo-foundation-rotation-evidence-v0.1.2.tar.gz", Size: 4767, SHA256: "sha256:d259722c20cb3e525575d4bfb1488424ab4d23eed964fe68334345711635fe6d"},
			{ID: 538516120, Name: "gooo-foundation-rotation-linux-amd64", Size: 4604245, SHA256: "sha256:82bd491f106abbe729242bcb0bb8b5a47d8dd5d9ceaa35be73722e4fc90d6b36"},
		},
	}
}

func testOIDCToken(t *testing.T) string {
	t.Helper()
	header, err := json.Marshal(map[string]string{"alg": "none", "typ": "JWT"})
	if err != nil {
		t.Fatal(err)
	}
	claims, err := json.Marshal(map[string]any{
		"iss": "https://token.actions.githubusercontent.com", "sub": "repo:kimjooyoon/gooo-foundation-issuer:ref:refs/heads/main",
		"aud": "gooo-foundation-issuer", "iat": 1700000000, "exp": 4102444800,
		"repository": "kimjooyoon/gooo-foundation-issuer", "repository_owner": "kimjooyoon",
		"workflow": "Issue Foundation authorization", "ref": "refs/heads/main", "sha": ExpectedRotationTarget, "actor": "kimjooyoon",
	})
	if err != nil {
		t.Fatal(err)
	}
	encode := base64.RawURLEncoding.EncodeToString
	return encode(header) + "." + encode(claims) + ".fixture"
}

func TestIssueAndIndependentConsumer(t *testing.T) {
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := Issue(testRequest(), testRotation(), testOIDCToken(t), private, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if receipt.PublicKey != base64.StdEncoding.EncodeToString(public) {
		t.Fatalf("public key binding mismatch")
	}
	report := VerifyReceipt(receipt, testRotation(), time.Date(2026, 9, 1, 0, 1, 0, 0, time.UTC))
	if !report.SignatureValid || !report.TupleExact || !report.RotationExact || report.Decision != Unknown {
		t.Fatalf("unexpected verification report: %+v", report)
	}
}

func TestStatePrecedenceAndRefutation(t *testing.T) {
	if Combine(Closed, Unknown, Closed) != Unknown {
		t.Fatal("UNKNOWN must outrank CLOSED")
	}
	if Combine(Unknown, Refuted, Closed) != Refuted {
		t.Fatal("REFUTED must outrank UNKNOWN")
	}
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := Issue(testRequest(), testRotation(), testOIDCToken(t), private, time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	_ = public
	receipt.Payload.HeadSHA = "0000000000000000000000000000000000000000"
	report := VerifyReceipt(receipt, testRotation(), time.Date(2026, 9, 1, 0, 1, 0, 0, time.UTC))
	if report.Decision != Refuted {
		t.Fatalf("tampered receipt must be REFUTED: %+v", report)
	}
}
