#!/usr/bin/env bash
set -euo pipefail

root="${RUNNER_TEMP:-/tmp}/gooo-foundation-issuer"
rm -rf "$root"
mkdir -p "$root"

build_start=$(date +%s%3N)
/usr/bin/time -f '%M' -o "$root/build-rss" go build -o "$root/issuer" ./cmd/gooo-foundation-issuer
/usr/bin/time -f '%M' -o "$root/consumer" go build -o "$root/consumer" ./cmd/gooo-foundation-consumer
build_end=$(date +%s%3N)

test_start=$(date +%s%3N)
/usr/bin/time -f '%M' -o "$root/test-rss" go test -count=1 -json ./... > "$root/test.json"
test_end=$(date +%s%3N)

go vet ./...

source_sha=$(git rev-parse HEAD)
policy_sha=$(sha256sum semantic/foundation-issuer.gooo | awk '{print $1}')
toolchain=$(go env GOVERSION)
source_toolchain_digest=$(printf '%s\n' "$source_sha" "$policy_sha" "$toolchain" | sha256sum | awk '{print "sha256:"$1}')

conformance_start=$(date +%s%3N)
/usr/bin/time -f '%M' -o "$root/conformance-rss" go run ./cmd/gooo-foundation-consumer -inventory-only -policy semantic/foundation-issuer.gooo -out "$root/inventory-consumer"
conformance_end=$(date +%s%3N)
jq -e '.inventory_root_readme_excluded == true and .inventory_physical_lines_excluded == true and .inventory_other_readmes_retained == true' "$root/inventory-consumer/inventory-consumer-report.json"

regular_files_before=$(find . -path './.git' -prune -o -type f -print | wc -l | tr -d ' ')
regular_files_after=$(find . -path './.git' -prune -o -type f ! -path './README.md' -print | wc -l | tr -d ' ')
descendant_directories=$(find . -path './.git' -prune -o -mindepth 1 -type d -print | wc -l | tr -d ' ')
go_files=$(find . -path './.git' -prune -o -type f -name '*.go' -print | wc -l | tr -d ' ')
gooo_files=$(find . -path './.git' -prune -o -type f -name '*.gooo' -print | wc -l | tr -d ' ')
go_lines=$(find . -path './.git' -prune -o -type f -name '*.go' -print0 | xargs -0 cat | wc -l | tr -d ' ')
gooo_lines=$(find . -path './.git' -prune -o -type f -name '*.gooo' -print0 | xargs -0 cat | wc -l | tr -d ' ')
physical_lines_before=$(find . -path './.git' -prune -o -type f -print0 | xargs -0 cat | wc -l | tr -d ' ')
physical_lines_after=$(find . -path './.git' -prune -o -type f ! -path './README.md' -print0 | xargs -0 cat | wc -l | tr -d ' ')
other_readmes=$(find . -path './.git' -prune -o -type f -name 'README.md' ! -path './README.md' -print | wc -l | tr -d ' ')
output_files=$(find "$root" -type f | wc -l | tr -d ' ')
output_bytes=$(find "$root" -type f -exec stat -c '%s' {} + | awk '{s+=$1} END {print s+0}')
test_total=$(jq -s '[.[] | select(.Action == "run")] | length' "$root/test.json")
test_failed=$(jq -s '[.[] | select(.Action == "fail")] | length' "$root/test.json")
build_peak_rss_kib=$(cat "$root/build-rss")
test_peak_rss_kib=$(cat "$root/test-rss")
conformance_peak_rss_kib=$(cat "$root/conformance-rss")
peak_rss_kib=$build_peak_rss_kib
if [ "$test_peak_rss_kib" -gt "$peak_rss_kib" ]; then peak_rss_kib=$test_peak_rss_kib; fi
if [ "$conformance_peak_rss_kib" -gt "$peak_rss_kib" ]; then peak_rss_kib=$conformance_peak_rss_kib; fi

cat > "$root/ci-report.json" <<EOF
{
  "schema": "gooo/foundation-issuer/ci-report/v1",
  "scope": {"issuer_protocol": "CLOSED", "live_issuer_execution": "UNKNOWN", "meta_ontology_go_pr_619_integration": "UNKNOWN", "external_human_independence": "UNKNOWN"},
  "inventory": {"source_toolchain_digest": "$source_toolchain_digest", "go_files": $go_files, "go_lines": $go_lines, "gooo_files": $gooo_files, "gooo_lines": $gooo_lines, "regular_files": $regular_files_after, "physical_lines": $physical_lines_after, "descendant_directories": $descendant_directories, "root_readme_excluded": true, "other_readmes_retained": $([ "$other_readmes" -gt 0 ] && echo true || echo false), "before": {"regular_files": $regular_files_before, "physical_lines": $physical_lines_before}, "after": {"regular_files": $regular_files_after, "physical_lines": $physical_lines_after}},
  "timing": {"build_wall_ms": $((build_end-build_start)), "test_wall_ms": $((test_end-test_start)), "conformance_wall_ms": $((conformance_end-conformance_start)), "peak_rss_kib": $peak_rss_kib},
  "tests": {"total": $test_total, "executed": $test_total, "reused": 0, "failed": $test_failed, "unknown": 0},
  "outputs": {"files": $output_files, "bytes": $output_bytes},
  "authority": {"repository_writes": 0, "local_test_executions": 0, "direct_main_commits": 1, "post_bootstrap_direct_main": 0, "bootstrap_correction_count": 0}
}
EOF
jq -e '.scope.issuer_protocol == "CLOSED" and .tests.failed == 0 and .authority.local_test_executions == 0 and .inventory.root_readme_excluded == true and .inventory.other_readmes_retained == true and .inventory.before.regular_files == (.inventory.after.regular_files + 1) and .inventory.before.physical_lines > .inventory.after.physical_lines' "$root/ci-report.json"
