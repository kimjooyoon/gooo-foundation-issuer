package issuer

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func ParsePrivateKey(value []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(value)
	if block == nil {
		return nil, fmt.Errorf("private key PEM is missing")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("private key parse failed: %w", err)
	}
	private, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not Ed25519")
	}
	return private, nil
}

func ParsePublicKey(value string) (ed25519.PublicKey, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil || len(decoded) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("public key is not raw Ed25519 base64")
	}
	return ed25519.PublicKey(decoded), nil
}

func PublicKeyString(private ed25519.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(private.Public().(ed25519.PublicKey))
}

func SignPayload(payload Payload, private ed25519.PrivateKey) (Receipt, error) {
	bytes, err := PayloadBytes(payload)
	if err != nil {
		return Receipt{}, err
	}
	return Receipt{
		Schema:        ReceiptSchema,
		PayloadType:   PayloadType,
		Payload:       payload,
		PayloadDigest: Digest(bytes),
		PublicKey:     PublicKeyString(private),
		Signature:     base64.StdEncoding.EncodeToString(ed25519.Sign(private, bytes)),
	}, nil
}

func VerifySignature(receipt Receipt) (bool, error) {
	public, err := ParsePublicKey(receipt.PublicKey)
	if err != nil {
		return false, err
	}
	signature, err := base64.StdEncoding.DecodeString(receipt.Signature)
	if err != nil || len(signature) != ed25519.SignatureSize {
		return false, fmt.Errorf("signature is malformed")
	}
	bytes, err := PayloadBytes(receipt.Payload)
	if err != nil {
		return false, err
	}
	if Digest(bytes) != receipt.PayloadDigest {
		return false, fmt.Errorf("payload digest mismatch")
	}
	return ed25519.Verify(public, bytes, signature), nil
}
