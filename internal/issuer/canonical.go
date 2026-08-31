package issuer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func CanonicalJSON(value any) ([]byte, error) {
	return json.Marshal(value)
}

func Digest(value []byte) string {
	sum := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func PayloadBytes(payload Payload) ([]byte, error) {
	return CanonicalJSON(payload)
}
