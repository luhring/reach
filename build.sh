#!/bin/bash

set -ex

export REACH_VERSION=${REACH_VERSION:-"0.0.0"}

set -u

export CGO_ENABLED=0
export GOARCH=amd64
export REACH_DIR_DARWIN=$(printf "reach_%s_darwin_amd64" $REACH_VERSION)
export REACH_DIR_LINUX=$(printf "reach_%s_linux_amd64" $REACH_VERSION)
export REACH_DIR_WINDOWS=$(printf "reach_%s_windows_amd64" $REACH_VERSION)

mkdir -p ./build

GOOS=darwin go build -a -tags netgo -o "./build/$REACH_DIR_DARWIN/reach"
GOOS=linux go build -a -tags netgo -o "./build/$REACH_DIR_LINUX/reach"
GOOS=windows go build -a -tags netgo -o "./build/$REACH_DIR_WINDOWS/reach.exe"

cp -nv ./LICENSE ./README.md "./build/$REACH_DIR_DARWIN"
cp -nv ./LICENSE ./README.md "./build/$REACH_DIR_LINUX"
cp -nv ./LICENSE ./README.md "./build/$REACH_DIR_WINDOWS"

pushd ./build
  tar -cvzf $REACH_DIR_DARWIN.tar.gz ./$REACH_DIR_DARWIN/*
  tar -cvzf $REACH_DIR_LINUX.tar.gz ./$REACH_DIR_LINUX/*
  tar -cvzf $REACH_DIR_WINDOWS.tar.gz ./$REACH_DIR_WINDOWS/*

  openssl dgst -sha256 ./$REACH_DIR_DARWIN.tar.gz >> ./checksums.txt
  openssl dgst -sha256 ./$REACH_DIR_LINUX.tar.gz >> ./checksums.txt
  openssl dgst -sha256 ./$REACH_DIR_WINDOWS.tar.gz >> ./checksums.txt

  cat ./checksums.txt
popd

set +eux
