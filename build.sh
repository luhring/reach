#!/bin/bash

# This script takes an argument for which OS to build for: darwin, linux, or windows.
# If no argument is provided, the script builds for all three.

# To build for a specific version, set the `REACH_VERSION` variable to something like "2.0.1" before running the script.

set -e

export REACH_VERSION=${REACH_VERSION:-"0.0.0"}

set -u

export CGO_ENABLED=0
export GOARCH=amd64

set -x

mkdir -p ./build

pushd ./build
  for CURRENT_OS in "darwin" "linux" "windows"
  do
    if [ -z $1 ] || [ $1 == $CURRENT_OS ]
    then
      export GOOS=$CURRENT_OS

      if [ $CURRENT_OS == "windows" ]
      then
        REACH_EXECUTABLE="reach.exe"
      else
        REACH_EXECUTABLE="reach"
      fi

      REACH_DIR_FOR_OS=$(printf "reach_%s_%s_amd64" $REACH_VERSION $CURRENT_OS)
      mkdir -p ./$REACH_DIR_FOR_OS

      go build -a -v -tags netgo -o "./$REACH_DIR_FOR_OS/$REACH_EXECUTABLE" ..
      cp -nv ../LICENSE ../README.md "./$REACH_DIR_FOR_OS/"

      if [ $CURRENT_OS == "windows" ]
      then
        zip $REACH_DIR_FOR_OS.zip ./$REACH_DIR_FOR_OS/*
        openssl dgst -sha256 ./$REACH_DIR_FOR_OS.zip >> ./checksums.txt
      else
        tar -cvzf $REACH_DIR_FOR_OS.tar.gz ./$REACH_DIR_FOR_OS/*
        openssl dgst -sha256 ./$REACH_DIR_FOR_OS.tar.gz >> ./checksums.txt
      fi
    fi

  done

  cat ./checksums.txt
popd

set +eux
