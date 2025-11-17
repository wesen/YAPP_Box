#!/usr/bin/env bash

# Run the standard test/build sequence required by the Path 1 (A→B) refactor.
# This matches the rhythm called out in the playbooks (module tests -> build -> full test suite).

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
cd "${ROOT}"

echo "==> go test ./pkg/yappgen/modules/..."
go test ./pkg/yappgen/modules/...

echo "==> go build ./..."
go build ./...

echo "==> go test ./..."
go test ./...

echo "Full suite finished successfully."

