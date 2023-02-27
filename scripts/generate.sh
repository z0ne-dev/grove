#!/usr/bin/env bash

#
# generate.sh Copyright (c) 2023 z0ne.
# All Rights Reserved.
# Licensed under the EUPL 1.2 License.
# See LICENSE the project root for license information.
#
# SPDX-License-Identifier: EUPL-1.2
#

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

# make sure `node` is installed
if ! command -v node &>/dev/null; then
  fatal "node is not installed"
fi

info "installing node dependencies"
pushd frontend > /dev/null
pnpm i

info "building frontend"
pnpm run build
popd > /dev/null