#!/usr/bin/env bash

set -euo pipefail

shopt -s globstar

if ! [[ "$0" =~ scripts/genproto.sh ]]; then
    echo "must be run from repository root"
    exit 255
fi

source ./scripts/lib.sh

API_ROOT="/api"

function dirs() {
    dirs=()
    while IFS= read -r dir; do
        dirs+=("$dir")
    done < <(find . -type f -name "*.proto" -exec dirname {} \; | xargs -n1 basename | sort -u)
    echo "${dirs[@]}"
}

function pb_files() {
    pb_files=$(find . -type f -name '*.proto')
    echo "${pb_files[@]}"
}

function get_for_modules() {
    local go_out="internal/common/genproto"
    if []
}

echo "directories containing protos to be built: $(dirs)"
echo "found pb_files: $(pb_files)"
gen_for_modules