#!/usr/bin/env bash

set -e
set -o pipefail

SCRIPT_DIR="$( dirname "$( readlink -f "${BASH_SOURCE[0]}" )" )"
SCRIPT_NAME="$( basename "$( readlink -f "${BASH_SOURCE[0]}" )" )"

function main() {
    cd "${SCRIPT_DIR}"
    echo "Add your functionality here"
}

main

echo "" >&2
echo "Bash script '${SCRIPT_NAME}' finished successfully" >&2
