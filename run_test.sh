#!/usr/bin/env bash
set -euo pipefail

# run_tests.sh
# Runs integration tests, generates a coverage badge via shields.io,
# and pushes the result to the current repo.


COVERAGE_FILE="coverage.out"
COVERAGE_SVG="coverage.svg"

# Run integration test
	echo "==> Running integration tests"
	go test -tags=integration -v -coverprofile="${COVERAGE_FILE}" -covermode=atomic ./...

	echo "==> Coverage summary"
	go tool cover -func="${COVERAGE_FILE}"

# Extract total coverage percentage
	COVERAGE=$(go tool cover -func="${COVERAGE_FILE}" | tail -1 | awk '{print $3}' | tr -d '%')
	echo "==> Total coverage: ${COVERAGE}%"

# Pick badge color based on coverage percentage
	if (( $(echo "${COVERAGE} >= 80" | bc -l) )); then
		COLOR="brightgreen"
	elif (( $(echo "${COVERAGE} >= 60" | bc -l) )); then
		COLOR="yellow"
	elif (( $(echo "${COVERAGE} >= 40" | bc -l) )); then
		COLOR="orange"
	else
		COLOR="red"
	fi

# Fetch coverage badge
	echo "==> Fetching coverage badge from shields.io"
	curl -s "https://img.shields.io/badge/coverage-${COVERAGE}%25-${COLOR}" > "${COVERAGE_SVG}"

# Commit it to repo
	echo "==> Committing and pushing results"
	git add "${COVERAGE_SVG}"
	if git diff --cached --quiet; then
		echo "==> No changes to commit (coverage unchanged)"
	else
		git commit -m "chore: update coverage badge (${COVERAGE}%)"
		git push
	fi

	echo "==> Done. Badge saved to ${COVERAGE_SVG}"