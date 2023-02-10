#!/usr/bin/env bash

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