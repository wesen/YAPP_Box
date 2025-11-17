#!/usr/bin/env bash

# Generate SCAD for every example YAML (excluding intentional error fixtures).
# Usage:
#   scripts/generate_all_examples.sh [output_dir]
#
# Optional environment variables:
#   YAPPCTL_BIN  - path to an existing yappctl binary. If unset, the script builds one.
#   SKIP_REGEX   - extended regex of example paths to skip (default skips *with-errors* fixtures).

set -euo pipefail

ROOT="$(git rev-parse --show-toplevel)"
EXAMPLES_DIR="${ROOT}/examples"
OUT_DIR="${1:-/tmp/yapp-task19-examples}"
YAPPCTL_BIN="${YAPPCTL_BIN:-${ROOT}/bin/yappctl-task19}"
SKIP_REGEX="${SKIP_REGEX:-with-errors|90-unresolved}"

mkdir -p "${OUT_DIR}"
mkdir -p "$(dirname "${YAPPCTL_BIN}")"

if [[ ! -x "${YAPPCTL_BIN}" ]]; then
  echo "==> building yappctl at ${YAPPCTL_BIN}"
  (cd "${ROOT}" && go build -o "${YAPPCTL_BIN}" ./cmd/yappctl)
fi

mapfile -t EXAMPLES < <(cd "${ROOT}" && find "${EXAMPLES_DIR#"${ROOT}/"}" -name '*.yaml' ! -name '*.resolved.yaml' | sort)

FAILURES=()
SKIPPED=()

for example in "${EXAMPLES[@]}"; do
  rel="${example#${EXAMPLES_DIR}/}"
  if [[ -n "${SKIP_REGEX}" && "${rel}" =~ ${SKIP_REGEX} ]]; then
    echo "Skipping ${rel} (matches SKIP_REGEX)"
    SKIPPED+=("${rel}")
    continue
  fi

  out_file="${OUT_DIR}/${rel//\//-}.scad"
  mkdir -p "$(dirname "${out_file}")"

  echo "Generating ${rel} -> ${out_file}"
  if ! "${YAPPCTL_BIN}" generate -i "${ROOT}/${example}" -o "${out_file}" > "${OUT_DIR}/generate.log" 2>&1; then
    echo "FAILED: ${rel}"
    cat "${OUT_DIR}/generate.log"
    FAILURES+=("${rel}")
  fi
done

rm -f "${OUT_DIR}/generate.log"

echo
echo "Skipped (${#SKIPPED[@]}): ${SKIPPED[*]:-none}"

if [[ ${#FAILURES[@]} -gt 0 ]]; then
  echo "Example generation failures (${#FAILURES[@]}):"
  printf ' - %s\n' "${FAILURES[@]}"
  exit 1
fi

echo "All example YAMLs generated successfully."

