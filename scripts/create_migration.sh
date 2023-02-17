#!/usr/bin/env bash

#
# create_migration.sh Copyright (c) 2023 z0ne.
# All Rights Reserved.
# Licensed under the EUPL 1.2 License.
# See LICENSE the project root for license information.
#
# SPDX-License-Identifier: EUPL-1.2
#

set -e

name=$1
if [ -z "$name" ]; then
  read -r -p "Migration name: " name
fi

slug=$(echo "$name" | tr '[:upper:]' '[:lower:]' | tr -c '[:alnum:]' '-' | tr -s '-' | sed -e 's/^-//' -e 's/-$//')
filename=$(date +%Y%m%d%H%M%S)_$slug.up.sql

echo "-- $name" > "migrations/$filename"

if [ -z "$URL" ]; then
  echo "Migration created: migrations/$filename"
else
  echo "Migration created: file://$(realpath "migrations/$filename")"
fi