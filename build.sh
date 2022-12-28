#!/bin/bash

CMD=akamai-review
MAIN=${PWD}
BUILD=${PWD}/build
VERSION=${1:-latest}
echo $VERSION

mkdir -p build
rm build/*

cd ${MAIN}
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
