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
    sudo apt list --installed | grep "^php8" | cut -d '/' -f1 | sudo apt-get remove -y
    sudo apt-get autoremove -y
    sudo apt-get clean

    echo "remove_unused_packages finished."
}

function remove_unused_directories() {
    echo "remove_unused_directories started."

    DOT_NET_DIR="/usr/share/dotnet/"
    if [ -e "${DOT_NET_DIR}" ] ; then
        rm -rf "${DOT_NET_DIR}"
        echo "${DOT_NET_DIR} deleted."
    else
        echo "${DOT_NET_DIR} already absent. Skip delete."
    fi

    echo "remove_unused_directories finished."
}

function main() {
    cd "${SCRIPT_DIR}"

    show_space_info

    remove_unused_packages
    remove_unused_directories

    show_space_info
}

main

echo ""
echo "Bash script '${SCRIPT_NAME}' finished successfully"
