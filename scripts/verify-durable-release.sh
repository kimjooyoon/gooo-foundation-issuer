#!/usr/bin/env bash
set -euo pipefail

tag=""
commit=""
evidence=""
binary=""
output=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    -tag) tag="$2"; shift 2 ;;
    -commit) commit="$2"; shift 2 ;;
    -evidence) evidence="$2"; shift 2 ;;
    -binary) binary="$2"; shift 2 ;;
    -output) output="$2"; shift 2 ;;
    *) echo "unknown argument: $1" >&2; exit 2 ;;
  esac
done

test -n "$tag" -a -n "$commit" -a -n "$evidence" -a -n "$binary" -a -n "$output"
release=$(gh api "repos/${GITHUB_REPOSITORY}/releases/tags/${tag}")
ref=$(gh api "repos/${GITHUB_REPOSITORY}/git/ref/tags/${tag}")
tag_object=$(gh api "repos/${GITHUB_REPOSITORY}/git/tags/$(jq -r .object.sha <<< "$ref")")

jq -e --arg tag "$tag" '.tag_name == $tag and .immutable == true and .draft == false and .prerelease == false and (.assets|length) == 2' <<< "$release" >/dev/null
jq -e '.object.type == "tag" and .object.sha != ""' <<< "$ref" >/dev/null
jq -e --arg commit "$commit" '.object.type == "commit" and .object.sha == $commit' <<< "$tag_object" >/dev/null

evidence_size=$(stat -c '%s' "$evidence")
evidence_sha=$(sha256sum "$evidence" | awk '{print "sha256:"$1}')
binary_size=$(stat -c '%s' "$binary")
binary_sha=$(sha256sum "$binary" | awk '{print "sha256:"$1}')
jq -e --arg name "$(basename "$evidence")" --argjson size "$evidence_size" --arg sha "$evidence_sha" 'any(.assets[]; .name == $name and .size == $size and .digest == $sha)' <<< "$release" >/dev/null
jq -e --arg name "$(basename "$binary")" --argjson size "$binary_size" --arg sha "$binary_sha" 'any(.assets[]; .name == $name and .size == $size and .digest == $sha)' <<< "$release" >/dev/null

jq -n --arg repository "$GITHUB_REPOSITORY" --arg tag "$tag" --arg commit "$commit" --argjson release "$release" --argjson ref "$ref" --argjson tag_object "$tag_object" --arg evidence_name "$(basename "$evidence")" --argjson evidence_size "$evidence_size" --arg evidence_sha "$evidence_sha" --arg binary_name "$(basename "$binary")" --argjson binary_size "$binary_size" --arg binary_sha "$binary_sha" '{schema:"gooo/foundation-issuer/durable-release-verification/v1",repository:$repository,tag:$tag,release_id:$release.id,immutable:$release.immutable,annotated_tag_object:$ref.object.sha,tag_target_commit:$tag_object.object.sha,assets:[{name:$evidence_name,size:$evidence_size,sha256:$evidence_sha},{name:$binary_name,size:$binary_size,sha256:$binary_sha}],checks:{release_immutable:true,annotated_tag:true,exact_asset_count:true,asset_digests_verified:true}}' > "$output"
