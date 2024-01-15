#!/bin/sh
VERSION="latest-$(date +%Y%m%dT%H%M%S)"
sed "s/VERSION/${VERSION}/" cmd/version.template >cmd/version.go 
go build -o patch/akamai-review && cp patch/akamai-review ~/.akamai-cli/src/akamai-review/akamai-review
sed "s/VERSION/latest/" cmd/version.template >cmd/version.go 