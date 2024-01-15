#!/bin/bash

CMD=akamai-review
MAIN=${PWD}
BUILD=${PWD}/build
VERSION=$(jq '.commands[0].version' -r <cli.json)
echo $VERSION
sed "s/VERSION/${VERSION}/" cmd/version.template >cmd/version.go 

mkdir -p build
rm build/*

cd ${MAIN}
GOOS=darwin GOARCH=arm64 go build -o ${BUILD}/${CMD}-${VERSION}-macarm64 .
shasum -a 256 ${BUILD}/${CMD}-${VERSION}-macarm64 | awk '{print $1'} > ${BUILD}/${CMD}-${VERSION}-macarm64.sig

GOOS=darwin GOARCH=amd64 go build -o ${BUILD}/${CMD}-${VERSION}-macamd64 .
shasum -a 256 ${BUILD}/${CMD}-${VERSION}-macamd64 | awk '{print $1'} > ${BUILD}/${CMD}-${VERSION}-macamd64.sig

GOOS=linux GOARCH=amd64 go build -o ${BUILD}/${CMD}-${VERSION}-linuxamd64 .
shasum -a 256 ${BUILD}/${CMD}-${VERSION}-linuxamd64 | awk '{print $1}' > ${BUILD}/${CMD}-${VERSION}-linuxamd64.sig

GOOS=linux GOARCH=386 go build -o ${BUILD}/${CMD}-${VERSION}-linux386
shasum -a 256 ${BUILD}/${CMD}-${VERSION}-linux386 | awk '{print $1}' > ${BUILD}/${CMD}-${VERSION}-linux386.sig

GOOS=windows GOARCH=386 go build -o ${BUILD}/${CMD}-${VERSION}-windows386.exe .
shasum -a 256 ${BUILD}/${CMD}-${VERSION}-windows386.exe | awk '{print $1}' > ${BUILD}/${CMD}-${VERSION}-windows386.exe.sig

GOOS=windows GOARCH=amd64 go build -o ${BUILD}/${CMD}-${VERSION}-windowsamd64.exe .
shasum -a 256 ${BUILD}/${CMD}-${VERSION}-windowsamd64.exe | awk '{print $1}' > ${BUILD}/${CMD}-${VERSION}-windowsamd64.exe.sig
