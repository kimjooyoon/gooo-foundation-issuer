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

regular_files=$(find . -path './.git' -prune -o -type f -print | wc -l | tr -d ' ')
descendant_directories=$(find . -path './.git' -prune -o -mindepth 1 -type d -print | wc -l | tr -d ' ')
go_files=$(find . -path './.git' -prune -o -type f -name '*.go' -print | wc -l | tr -d ' ')
gooo_files=$(find . -path './.git' -prune -o -type f -name '*.gooo' -print | wc -l | tr -d ' ')
go_lines=$(find . -path './.git' -prune -o -type f -name '*.go' -print0 | xargs -0 cat | wc -l | tr -d ' ')
gooo_lines=$(find . -path './.git' -prune -o -type f -name '*.gooo' -print0 | xargs -0 cat | wc -l | tr -d ' ')
output_files=$(find "$root" -type f | wc -l | tr -d ' ')
output_bytes=$(find "$root" -type f -exec stat -c '%s' {} + | awk '{s+=$1} END {print s+0}')
test_total=$(jq -s '[.[] | select(.Action == "run")] | length' "$root/test.json")
test_failed=$(jq -s '[.[] | select(.Action == "fail")] | length' "$root/test.json")
build_peak_rss_kib=$(cat "$root/build-rss")
test_peak_rss_kib=$(cat "$root/test-rss")
peak_rss_kib=$((build_peak_rss_kib > test_peak_rss_kib ? build_peak_rss_kib : test_peak_rss_kib))

cat > "$root/ci-report.json" <<EOF
{
  "schema": "gooo/foundation-issuer/ci-report/v1",
  "scope": {"issuer_protocol": "CLOSED", "live_issuer_execution": "UNKNOWN", "meta_ontology_go_pr_619_integration": "UNKNOWN", "external_human_independence": "UNKNOWN"},
  "inventory": {"go_files": $go_files, "go_lines": $go_lines, "gooo_files": $gooo_files, "gooo_lines": $gooo_lines, "regular_files": $regular_files, "descendant_directories": $descendant_directories, "root_readme_excluded": false},
  "timing": {"build_wall_ms": $((build_end-build_start)), "test_wall_ms": $((test_end-test_start)), "conformance_wall_ms": 0, "peak_rss_kib": $peak_rss_kib},
  "tests": {"total": $test_total, "executed": $test_total, "reused": 0, "failed": $test_failed, "unknown": 0},
  "outputs": {"files": $output_files, "bytes": $output_bytes},
  "authority": {"repository_writes": 0, "local_test_executions": 0, "direct_main_commits": 1, "post_bootstrap_direct_main": 0, "bootstrap_correction_count": 0}
}
EOF
jq -e '.scope.issuer_protocol == "CLOSED" and .tests.failed == 0 and .authority.local_test_executions == 0' "$root/ci-report.json"
