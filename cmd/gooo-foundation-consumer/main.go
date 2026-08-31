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
	receiptPath := flag.String("receipt", "", "signed receipt JSON")
	rotationPath := flag.String("rotation", "", "immutable rotation input JSON")
	publicKeyPath := flag.String("public-key", "", "raw base64 public key")
	outDir := flag.String("out", "", "consumer output directory")
	flag.Parse()
	if *receiptPath == "" || *rotationPath == "" || *publicKeyPath == "" || *outDir == "" {
		fatal("receipt, rotation, public-key, and out are required")
	}
	if err := os.MkdirAll(*outDir, 0o700); err != nil {
		fatal("create output directory: %v", err)
	}
	var receipt issuer.Receipt
	var rotation issuer.RotationInput
	if err := issuer.ReadJSON(*receiptPath, &receipt); err != nil {
		fatal("read receipt: %v", err)
	}
	if err := issuer.ReadJSON(*rotationPath, &rotation); err != nil {
		fatal("read rotation: %v", err)
	}
	keyBytes, err := os.ReadFile(*publicKeyPath)
	if err != nil {
		fatal("read public key: %v", err)
	}
	if _, err := issuer.ParsePublicKey(string(keyBytes)); err != nil {
		fatal("read public key: %v", err)
	}
	receipt.PublicKey = string(keyBytes)
	report := issuer.VerifyReceipt(receipt, rotation, time.Now())
	if err := issuer.WriteJSON(filepath.Join(*outDir, "independent-consumer-report.json"), report); err != nil {
		fatal("write consumer report: %v", err)
	}
	if report.Decision == issuer.Refuted {
		fatal("consumer decision is REFUTED: %s", report.Reason)
	}
	fmt.Printf("consumer decision=%s signature_valid=%t tuple_exact=%t rotation_exact=%t\n", report.Decision, report.SignatureValid, report.TupleExact, report.RotationExact)
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
