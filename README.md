# gooo-foundation-issuer

Owner-operated external issuer for exact Foundation authorization tuples.

The issuer consumes the immutable `gooo-foundation-rotation v0.1.2` input,
binds PR #619's observed repository/base/head/path tuple, obtains the current
GitHub Actions OIDC claims, and signs a receipt with an Ed25519 key held only
as a GitHub Actions secret. The public key is the only key material committed
to this repository.

The `.gooo` graph owns the issuance state machine, exact tuple rules,
single-use/replay/revocation precedence, and the human/consumer authority
boundary. Go 1.27 is the parser, signer, verifier, and serialization backend.

This is intentionally an owner-operated external-repository authority. It is
not independent human approval. `external_human_independence` and the current
`meta-ontology-go` Guardian consumer integration remain `UNKNOWN`.
