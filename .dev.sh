#!/bin/zsh

script_dir="$(dirname "$0")"
# Change directory to the location of the script and push the current directory onto a stack
pushd $script_dir
# Build the code
go build -o slv-dev ./cli/main
# Go back to the original directory
popd
# Run slv from current directory
$script_dir/slv-dev "$@"