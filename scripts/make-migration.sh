#!/usr/bin/env bash

set -e

timestamp() {
  date +"%Y%m%d_%H%M%S"
}

echo "-- +migrate Up

-- +migrate Down" > ./migrations/"$(timestamp)"_"$name".sql
