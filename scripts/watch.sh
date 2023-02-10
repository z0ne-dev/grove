#!/usr/bin/env bash

set -e

# echo message and exit faulty
function fatal() {
  echo -e "\e[31m! $*\e[0m"
  exit 1
}

# echo info message
function info() {
  echo -e "\e[36m> $*\e[0m"
}

# make sure `go` is installed
if ! command -v go &>/dev/null; then
  fatal "go is not installed"
fi

# make sure `node` is installed
if ! command -v node &>/dev/null; then
  fatal "node is not installed"
fi

# install `gow` if not installed
if ! command -v gow &>/dev/null; then
  info "installing gow"
  go install github.com/mitranim/gow@latest > /dev/null
fi

info "installing node dependencies"
pushd frontend > /dev/null
pnpm i

info "starting node watchers"
pnpm run watch > /dev/null &
_node_watch_pid=$!
popd > /dev/null

info "starting gow"
gow -c run -tags dev ./cmd/grove

info "exiting node"
kill $_node_watch_pid