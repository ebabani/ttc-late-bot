#!/bin/sh

set -ex
GOOS=linux GOARCH=amd64 go build
cf push -b binary_buildpack
