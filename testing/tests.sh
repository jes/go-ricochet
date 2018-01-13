#!/bin/bash

set -e
pwd
go test ${1} -coverprofile=main.cover.out -v .
go test ${1} -coverprofile=utils.cover.out -v ./utils
go test ${1} -coverprofile=channels.cover.out -v ./channels
go test ${1} -coverprofile=connection.cover.out -v ./connection
go test ${1} -coverprofile=policies.cover.out -v ./policies
go test ${1} -coverprofile=policies.cover.out -v ./identity
echo "mode: set" > coverage.out && cat *.cover.out | grep -v mode: | sort -r | \
awk '{if($1 != last) {print $0;last=$1}}' >> coverage.out
rm -rf *.cover.out
