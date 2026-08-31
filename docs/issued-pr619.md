# PR #619 issuance record

The issuer workflow run `33441369615` observed and signed the exact current
PR #619 tuple after consuming the immutable `gooo-foundation-rotation v0.1.2`
release. The signature and independent consumer verification are valid, but
the overall decision is `UNKNOWN` because independent human approval and an
exact target Guardian consumer are not observed.

```text
target_repository: kimjooyoon/meta-ontology-go
pull_request: 619
base: dev@ac3a56b933d9a9b934fe26709485dc2f36edd916
head: agent/meta-policy-compilation-semantic-authority-20260901@90f5fdf2198a7da5cde405f9f675cc62975b9226
candidate_digest: sha256:0df9610f6300503ac39f18a9389695de33ac2b0edb3523a25fc4df77c0fa166e
protected_scope_digest: sha256:3e09937e0e48ec2fcf4c6c605ae5937756648b4532a1735e5d4cd6b9c04e44fb
receipt_digest: sha256:2b83603b1551f090d8d92799f8f156ead5358684951c02afccd760219c187f5d
public_key: 5p0KdsDkbj1wPJ+fshp0sTb9lf8GvbFeagbdG8kCU7w=
issuer_state: CLOSED
decision: UNKNOWN
human_independence: UNKNOWN
guardian_integration: UNKNOWN
```

The observed OIDC claims were:

```text
iss: https://token.actions.githubusercontent.com
sub: repo:kimjooyoon@115961382/gooo-foundation-issuer@1352865272:ref:refs/heads/main
aud: gooo-foundation-issuer
iat: 1788211608
exp: 1788211908
repository: kimjooyoon/gooo-foundation-issuer
workflow: Issue Foundation authorization
ref: refs/heads/main
sha: 3b618b53d4b95ffa29279d171645a584dd4ccabf
actor: kimjooyoon
oidc_signature_state: UNKNOWN
```

The current target Guardian cannot consume this receipt exactly. Its
`ci-guardian.yml` is `pull_request_target` only with read permissions and no
OIDC or dispatch input path. Its v6 dispatch validator is hard-coded to PR
#609 and rejects the #619 event with
`CI-FOUNDATION-AUTHORIZATION-001: Guardian dispatch live candidate tuple is not exact`.

Issuer evidence:

```text
run_id: 33441369615
job_id: 99650135450
artifact_id: 9776360555
artifact_name: foundation-issuer-dossier-33441369615
artifact_size_bytes: 6804
artifact_sha256: bb49f3efbc34fc1a0de2750f4bb8b9818eafccbf89878307b1402be0f33d6c70
```

The artifact contains the signed receipt, issuer report, independent consumer
report, OIDC/actor observation, replay/revocation evidence, dossier, semantic
graph, and durable-release verification. Replay evidence is
`first_use=CLOSED`, `second_use=REFUTED`, `revoked_key_use=REFUTED`.

The immutable release is:

```text
release_url: https://github.com/kimjooyoon/gooo-foundation-issuer/releases/tag/v0.1.0
release_id: 380062555
immutable: true
annotated_tag_object: fcafea8d59e44d30d7fe7e9e967698747f3827a1
tag_target_commit: 3b618b53d4b95ffa29279d171645a584dd4ccabf
asset: 538547113 gooo-foundation-issuer-evidence-v0.1.0.tar.gz 3245 sha256:911b907eebb4ecee12ea96c41e4f02328c186c8b01d045f43822e3c16105c411
asset: 538547119 gooo-foundation-issuer-linux-amd64 5834241 sha256:33b8dc45c8c64cff8394efe8e9610a2567331b536d94050828fa81210b23ad64
```

The v0.1.0 release verification checks are all true: immutable release,
annotated tag, exact asset count, and asset digest verification. No target
repository commit, PR update, force merge, or target workflow dispatch was
performed because the exact consumer interface is unavailable.
