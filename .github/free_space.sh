#!/usr/bin/env bash

set -e
set -o pipefail

SCRIPT_DIR="$( dirname "$( readlink -f "${BASH_SOURCE[0]}" )" )"
SCRIPT_NAME="$( basename "$( readlink -f "${BASH_SOURCE[0]}" )" )"

function show_space_info() {
    echo "show_space_info started."
    apt list --installed
    df -h
    echo "show_space_info finished."
}

function remove_unused_packages() {
    echo "remove_unused_packages started."

    sudo apt-get remove -y google-chrome-stable firefox powershell mono-devel ruby-full
    sudo apt-get autoremove -y
    sudo apt-get clean

    echo "remove_unused_packages finished."
}

function main() {
    cd "${SCRIPT_DIR}"

    show_space_info

    remove_unused_packages

    show_space_info
}

main

echo ""
echo "Bash script '${SCRIPT_NAME}' finished successfully"
