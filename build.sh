#!/bin/bash

# Bash script for building Go executables

SOURCE_FILE="./src/main.go"
BUILD_FLAGS="-s -w"
UPX_COMPRESSION_FLAGS="-6"
ENABLE_UPX=false
BUILD_FOLDER="build"

# Create build folder
mkdir -p $BUILD_FOLDER

# Linux
# AMD64
GOOS=linux GOARCH=amd64
if [ -z "$BUILD_FLAGS" ]; then
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64" $SOURCE_FILE
else
    go build -ldflags "$BUILD_FLAGS" -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64" $SOURCE_FILE
fi
if [ "$ENABLE_UPX" = true ]; then
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64"
    mv "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64-compressed"
fi

# 386
GOOS=linux GOARCH=386
if [ -z "$BUILD_FLAGS" ]; then
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86" $SOURCE_FILE
else
    go build -ldflags "$BUILD_FLAGS" -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86" $SOURCE_FILE
fi
if [ "$ENABLE_UPX" = true ]; then
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86"
    mv "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86-compressed"
fi

# macOS
# AMD64
GOOS=darwin GOARCH=amd64
if [ -z "$BUILD_FLAGS" ]; then
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64" $SOURCE_FILE
else
    go build -ldflags "$BUILD_FLAGS" -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64" $SOURCE_FILE
fi
if [ "$ENABLE_UPX" = true ]; then
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64"
    mv "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64-compressed"
fi

# ARM64
GOOS=darwin GOARCH=arm64
if [ -z "$BUILD_FLAGS" ]; then
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64" $SOURCE_FILE
else
    go build -ldflags "$BUILD_FLAGS" -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64" $SOURCE_FILE
fi
if [ "$ENABLE_UPX" = true ]; then
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64"
    mv "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64-compressed"
fi

# Windows
# AMD64
GOOS=windows GOARCH=amd64
if [ -z "$BUILD_FLAGS" ]; then
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe" $SOURCE_FILE
else
    go build -ldflags "$BUILD_FLAGS" -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe" $SOURCE_FILE
fi
if [ "$ENABLE_UPX" = true ]; then
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe"
    mv "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64-compressed.exe"
fi

# 386
GOOS=windows GOARCH=386
if [ -z "$BUILD_FLAGS" ]; then
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe" $SOURCE_FILE
else
    go build -ldflags "$BUILD_FLAGS" -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe" $SOURCE_FILE
fi
if [ "$ENABLE_UPX" = true ]; then
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe"
    mv "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86-compressed.exe"
fi

echo "All builds completed!"
