#!/bin/bash -e
pushd tools/check-structfield-settings && go run . "$@"
