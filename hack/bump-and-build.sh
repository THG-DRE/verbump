#!/bin/bash

VERBUMP_VERSION=$(cat version)

if ! test -f verbump; then
    if [ "$(uname)" == "Darwin" ]; then
        OS=darwin
        ARCH=arm64
    elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
        OS=linux
        ARCH=amd64
    fi

    wget https://github.com/THG-DRE/verbump/releases/download/v$VERBUMP_VERSION/verbump-$OS-$ARCH-$VERBUMP_VERSION.tar.gz
    tar -xvf verbump-$OS-$ARCH-$VERBUMP_VERSION.tar.gz
    rm verbump-$OS-$ARCH-$VERBUMP_VERSION.tar.gz
fi

VERSION=$(cat version)
echo "Previous Version: $VERSION"

./verbump bump --repository="." --version-file="version" --include="."

VERSION=$(cat version)
echo "New Version: $VERSION"

make build