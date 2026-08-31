package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kimjooyoon/gooo-foundation-issuer/internal/issuer"
)

func main() {
	requestPath := flag.String("request", "", "exact authorization request JSON")
	rotationPath := flag.String("rotation", "", "immutable rotation input JSON")
	outDir := flag.String("out", "", "evidence output directory")
	flag.Parse()
	if *requestPath == "" || *rotationPath == "" || *outDir == "" {
		fatal("request, rotation, and out are required")
	}
	if err := os.MkdirAll(*outDir, 0o700); err != nil {
		fatal("create output directory: %v", err)
	}
	var request issuer.Request
	var rotation issuer.RotationInput
	if err := issuer.ReadJSON(*requestPath, &request); err != nil {
		fatal("read request: %v", err)
	}
	if err := issuer.ReadJSON(*rotationPath, &rotation); err != nil {
		fatal("read rotation: %v", err)
	}
	private, err := issuer.ParsePrivateKey([]byte(os.Getenv("FOUNDATION_ISSUER_PRIVATE_KEY")))
	if err != nil {
		fatal("read Actions secret: %v", err)
	}
	receipt, err := issuer.Issue(request, rotation, os.Getenv("ACTIONS_ID_TOKEN"), private, time.Now())
	if err != nil {
		fatal("issue receipt: %v", err)
	}
	if err := issuer.WriteJSON(filepath.Join(*outDir, "signed-receipt.json"), receipt); err != nil {
		fatal("write receipt: %v", err)
	}
	report := issuer.VerifyReceipt(receipt, rotation, time.Now())
	if err := issuer.WriteJSON(filepath.Join(*outDir, "issuer-verification-report.json"), report); err != nil {
		fatal("write report: %v", err)
	}
	if err := issuer.WriteJSON(filepath.Join(*outDir, "oidc-actor-claim-observation.json"), receipt.Payload.OIDC); err != nil {
		fatal("write claims: %v", err)
	}
	if err := issuer.WriteJSON(filepath.Join(*outDir, "replay-revocation-evidence.json"), issuer.BuildReplayEvidence(receipt)); err != nil {
		fatal("write replay evidence: %v", err)
	}
	dossier := fmt.Sprintf("# Foundation authorization issuer dossier\n\n- receipt schema: `%s`\n- target: `%s` PR #%d\n- issuer state: `%s`\n- overall decision: `%s`\n- external human independence: `%s`\n- Guardian integration: `%s`\n- private key: Actions secret only; not persisted\n- authority claim: `%s`\n", receipt.Schema, receipt.Payload.Repository, receipt.Payload.PullRequest, receipt.Payload.IssuanceState, report.Decision, receipt.Payload.HumanIndependenceState, receipt.Payload.IntegrationState, receipt.Payload.ExternalAuthorityClaim)
	if err := os.WriteFile(filepath.Join(*outDir, "dossier.md"), []byte(dossier), 0o600); err != nil {
		fatal("write dossier: %v", err)
	}
	if report.Decision == issuer.Refuted {
		fatal("receipt verification is REFUTED: %s", report.Reason)
	}
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
