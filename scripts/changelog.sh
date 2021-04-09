#!/usr/bin/env bash

set -o pipefail

tfile=$(mktemp /tmp/bl-CHANGELOG-XXXXXX)
github-changelog-generator -org binarylane -repo bl-cli >"$tfile"

GO111MODULE=on go mod tidy

echo "$tfile"
