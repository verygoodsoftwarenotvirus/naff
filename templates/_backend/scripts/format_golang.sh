#!/usr/bin/env bash
set -euo pipefail

# Format all Go files using gofmt
# Usage: format_golang.sh <project_root> <gofmt_command>

PROJECT_ROOT="${1:-$(pwd)}"

if sed --version 2>/dev/null | grep -q 'GNU'; then
  SED_INPLACE=(sed -i)
else
  SED_INPLACE=(sed -i '')
fi

while IFS= read -r -d '' file; do
  # Ensure a blank line after the sentinel `_ struct{}` field
  # shellcheck disable=SC2016
  "${SED_INPLACE[@]}" '/_ struct{} `json:"-"`/{
N
/\n[[:space:]]*$/!s/\n/\n\n/
}' "${file}"
  # GO_FORMAT contains a command with arguments, so we use eval
  # shellcheck disable=SC2086
  eval "gofmt -s -w \"${file}\""
done < <(find "${PROJECT_ROOT}" -type f -not -path '*/vendor/*' -name "*.go" -print0)
