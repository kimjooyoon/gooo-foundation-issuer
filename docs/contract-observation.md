# #619 contract observation

The observed PR head is `90f5fdf2198a7da5cde405f9f675cc62975b9226` for PR
`619`, based on `dev@ac3a56b933d9a9b934fe26709485dc2f36edd916`.

The current Guardian workflow is `pull_request_target` only. It grants
`contents: read`, `pull-requests: read`, and `actions: read`; it does not grant
`id-token: write`, does not define a `workflow_dispatch` input contract, and
does not consume an external issuer repository.

The current Foundation dispatch policy is hard-coded to PR #609, branch
`agent/dev-main-sync-20260831-rerun`, head
`8b47db349315c02933296423b0ae7fa80ffeb1dc`, and the v6 receipt schema. Its
live validator compares the event pull request against that exact PR #609
tuple. Applying it to PR #619 therefore produces the observed
`CI-FOUNDATION-AUTHORIZATION-001: Guardian dispatch live candidate tuple is not exact`.

This issuer records a signed exact #619 tuple, but it does not claim that the
current Guardian can consume it. The integration decision remains `UNKNOWN`
until the target repository exposes an exact external-receipt consumer path.
