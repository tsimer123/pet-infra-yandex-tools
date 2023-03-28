#!/bin/bash

set -eu

script_dir="$(dirname "${BASH_SOURCE[0]}")"
project_dir="${script_dir}/.."

function fill_env_file() {
    echo "" >${project_dir}/.env
}

rm -rf .env
touch .env

fill_env_file
