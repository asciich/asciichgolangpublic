#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$( dirname "$( readlink -f "${BASH_SOURCE[0]}" )" )"
SCRIPT_NAME="$( basename "$( readlink -f "${BASH_SOURCE[0]}" )" )"

export DEBIAN_FRONTEND=noninteractive

function show_space_info() {
    echo "show_space_info started."
    dpkg-query -W -f='${Package} ${Version}\n' | sort
    df -h
    echo "show_space_info finished."
}

# Prints the names of installed packages matching the given dpkg glob patterns.
function list_installed_packages() {
    dpkg-query -W -f='${Status} ${Package}\n' "$@" 2>/dev/null \
        | awk '$3 == "installed" { print $4 }' \
        || true
}

# Removes only the packages that are really installed; skips the rest.
function remove_packages_if_installed() {
    local packages
    mapfile -t packages < <(list_installed_packages "$@")

    if [ "${#packages[@]}" -eq 0 ] ; then
        echo "No installed packages matching '$*'. Skip remove."
        return 0
    fi

    echo "Removing: ${packages[*]}"
    sudo apt-get remove -y --purge "${packages[@]}"
}

function remove_unused_packages() {
    echo "remove_unused_packages started."

    remove_packages_if_installed google-chrome-stable microsoft-edge-stable firefox mono-devel 'ruby*'

    # Remove all PHP packages together, including php-pear / php-common,
    # so apt has no leftover dependencies it can't satisfy.
    remove_packages_if_installed 'php*'

    sudo apt-get autoremove -y --purge
    sudo apt-get clean

    echo "remove_unused_packages finished."
}

function remove_unused_directories() {
    echo "remove_unused_directories started."

    local dirs=(
        "/usr/share/dotnet"
    )

    for dir in "${dirs[@]}" ; do
        if [ -e "${dir}" ] ; then
            sudo rm -rf "${dir}"
            echo "${dir} deleted."
        else
            echo "${dir} already absent. Skip delete."
        fi
    done

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
