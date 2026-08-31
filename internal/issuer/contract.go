package issuer

import "fmt"

const (
	ExpectedRequestSchema              = "gooo/foundation-authorization/request/v1"
	ExpectedRepository                 = "kimjooyoon/meta-ontology-go"
	ExpectedPullRequest                = 619
	ExpectedBaseRef                    = "dev"
	ExpectedBaseSHA                    = "ac3a56b933d9a9b934fe26709485dc2f36edd916"
	ExpectedHeadRef                    = "agent/meta-policy-compilation-semantic-authority-20260901"
	ExpectedHeadSHA                    = "90f5fdf2198a7da5cde405f9f675cc62975b9226"
	ExpectedCandidateDigest            = "sha256:0df9610f6300503ac39f18a9389695de33ac2b0edb3523a25fc4df77c0fa166e"
	ExpectedProtectedScopeDigest       = "sha256:3e09937e0e48ec2fcf4c6c605ae5937756648b4532a1735e5d4cd6b9c04e44fb"
	ExpectedRotationRepository         = "kimjooyoon/gooo-foundation-rotation"
	ExpectedRotationVersion            = "v0.1.2"
	ExpectedRotationReleaseID    int64 = 380045288
	ExpectedRotationTagObject          = "cc311eb5b14f54b5e467f9402c6f7a32e9b5b3dc"
	ExpectedRotationTarget             = "5534a6814fdd9c8c6ddd957d5c26fd5c15258e02"
)

var ExpectedChangedPaths = []string{
	".github/agent-scope-table.md",
	".github/ci-governance.json",
	".github/workflows/meta-policy-compilation.yml",
	".github/workflows/transformation-effect.yml",
	"cmd/meta-policy-compilation-consumer/main.go",
	"cmd/meta-policy-compilation-witness/main.go",
	"docs/language/meta-policy-compilation.md",
	"examples/language-syntax-roundtrip/corpus.json",
	"examples/meta-policy-compilation/README.md",
	"examples/meta-policy-compilation/cases.json",
	"examples/meta-policy-compilation/policy.gooo",
	"internal/meta/languagereadiness/languagesyntax/conformance/evaluate_test.go",
	"internal/meta/languagereadiness/languagesyntax/model.go",
	"internal/meta/languagereadiness/languagesyntax/registry.go",
	"internal/meta/policycompilation/compile.go",
	"internal/meta/policycompilation/constants.go",
	"internal/meta/policycompilation/digest.go",
	"internal/meta/policycompilation/evaluate.go",
	"internal/meta/policycompilation/judge.go",
	"internal/meta/policycompilation/model.go",
	"internal/meta/policycompilation/policycompilation_test.go",
	"internal/meta/policycompilation/receipt.go",
	"internal/meta/policycompilation/verify.go",
	"internal/verify/scope_meta_policy_compilation_semantic_authority_20260901.go",
}

var ExpectedProtectedScope = []string{
	".github/agent-scope-table.md",
	".github/ci-governance.json",
	".github/workflows/meta-policy-compilation.yml",
	".github/workflows/transformation-effect.yml",
	"internal/verify/scope_meta_policy_compilation_semantic_authority_20260901.go",
}

func ValidateRequest(request Request) error {
	if request.Schema != ExpectedRequestSchema || request.Repository != ExpectedRepository || request.PullRequest != ExpectedPullRequest || request.BaseRef != ExpectedBaseRef || request.BaseSHA != ExpectedBaseSHA || request.HeadRef != ExpectedHeadRef || request.HeadSHA != ExpectedHeadSHA {
		return fmt.Errorf("exact PR #619 tuple mismatch")
	}
	if request.CandidateDigest != ExpectedCandidateDigest || request.ProtectedScopeDigest != ExpectedProtectedScopeDigest {
		return fmt.Errorf("candidate or protected-scope digest mismatch")
	}
	if !equalStrings(request.ChangedPaths, ExpectedChangedPaths) || !equalStrings(request.ProtectedScope, ExpectedProtectedScope) {
		return fmt.Errorf("changed path or protected scope list mismatch")
	}
	if request.ActorIdentity == "" || request.IssuerIdentity == "" || request.Nonce == "" {
		return fmt.Errorf("authority identity or nonce is missing")
	}
	return nil
}

func ValidateRotation(rotation RotationInput) error {
	if rotation.Repository != ExpectedRotationRepository || rotation.Version != ExpectedRotationVersion || rotation.ReleaseID != ExpectedRotationReleaseID || !rotation.Immutable || rotation.TagObjectSHA != ExpectedRotationTagObject || rotation.TargetCommit != ExpectedRotationTarget {
		return fmt.Errorf("v0.1.2 immutable rotation input mismatch")
	}
	expected := []RotationAsset{
		{ID: 538516119, Name: "gooo-foundation-rotation-evidence-v0.1.2.tar.gz", Size: 4767, SHA256: "sha256:d259722c20cb3e525575d4bfb1488424ab4d23eed964fe68334345711635fe6d"},
		{ID: 538516120, Name: "gooo-foundation-rotation-linux-amd64", Size: 4604245, SHA256: "sha256:82bd491f106abbe729242bcb0bb8b5a47d8dd5d9ceaa35be73722e4fc90d6b36"},
	}
	if len(rotation.Assets) != len(expected) {
		return fmt.Errorf("rotation asset count is not exact")
	}
	for index := range expected {
		if rotation.Assets[index] != expected[index] {
			return fmt.Errorf("rotation asset %d mismatch", index+1)
		}
	}
	return nil
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
