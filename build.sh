#!/bin/bash

# This script takes an argument for which OS to build for: darwin, linux, or windows.
# If no argument is provided, the script builds for all three.

# To build for a specific version, set the `REACH_VERSION` variable to something like "2.0.1" before running the script.

set -e

export REACH_VERSION=${REACH_VERSION:-"0.0.0"}
export SPECIFIED_OS=""

if [[ -z "$1" ]]
then
  export SPECIFIED_OS="$1"
fi

set -u

export CGO_ENABLED=0
export GOARCH=amd64

set -x

function build_for_os {
  local GOOS="$1"
  local REACH_EXECUTABLE

  if [[ "$GOOS" == "windows" ]]
  then
    REACH_EXECUTABLE="reach.exe"
  else
    REACH_EXECUTABLE="reach"
  fi

  local REACH_DIR_FOR_OS
  REACH_DIR_FOR_OS=$(printf "reach_%s_%s_amd64" "$REACH_VERSION" "$GOOS")

  mkdir -p "./$REACH_DIR_FOR_OS"

  GOOS=$GOOS go build -a -v -tags netgo -o "./$REACH_DIR_FOR_OS/$REACH_EXECUTABLE" ..
  cp -nv ../COPYING ../COPYING.LESSER ../README.md "./$REACH_DIR_FOR_OS/"

  if [[ "$GOOS" == "windows" ]]
  then
    zip "$REACH_DIR_FOR_OS.zip" "./$REACH_DIR_FOR_OS"/*
    openssl dgst -sha256 "./$REACH_DIR_FOR_OS.zip" >> ./checksums.txt
  else
    tar -cvzf "$REACH_DIR_FOR_OS.tar.gz" "./$REACH_DIR_FOR_OS"/*
    openssl dgst -sha256 "./$REACH_DIR_FOR_OS.tar.gz" >> ./checksums.txt
  fi
}

rm -rf ./build
mkdir -p ./build

pushd ./build
  if [[ -n "$SPECIFIED_OS" ]]
  then
    build_for_os "$SPECIFIED_OS"
  else
    for CURRENT_OS in "darwin" "linux" "windows"
    do
        build_for_os "$CURRENT_OS"
    done
  fi

  set +x

  cat ./checksums.txt
popd

set +eu
