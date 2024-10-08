#!/bin/bash

# Bash script for building Go executables

SOURCE_FILE="./src/main.go"
BUILD_FLAGS="-s -w"
UPX_COMPRESSION_FLAGS="-6"
ENABLE_UPX=false

# Create build folder
mkdir -p build

# Linux
# AMD64
GOOS=linux GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o "build/AEUSTNetworkAutoLogin-Linux-x86_64" $SOURCE_FILE
if [ "$ENABLE_UPX" = true ] ; then
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-Linux-x86_64"
    mv "build/AEUSTNetworkAutoLogin-Linux-x86_64" "build/AEUSTNetworkAutoLogin-Linux-x86_64-compressed"
fi

# 386
GOOS=linux GOARCH=386 go build -ldflags "$BUILD_FLAGS" -o "build/AEUSTNetworkAutoLogin-Linux-x86" $SOURCE_FILE
if [ "$ENABLE_UPX" = true ] ; then
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-Linux-x86"
    mv "build/AEUSTNetworkAutoLogin-Linux-x86" "build/AEUSTNetworkAutoLogin-Linux-x86-compressed"
fi

# macOS
# AMD64
GOOS=darwin GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o "build/AEUSTNetworkAutoLogin-macOS-x86_64" $SOURCE_FILE
if [ "$ENABLE_UPX" = true ] ; then
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-macOS-x86_64"
    mv "build/AEUSTNetworkAutoLogin-macOS-x86_64" "build/AEUSTNetworkAutoLogin-macOS-x86_64-compressed"
fi

# ARM64
GOOS=darwin GOARCH=arm64 go build -ldflags "$BUILD_FLAGS" -o "build/AEUSTNetworkAutoLogin-macOS-arm64" $SOURCE_FILE
if [ "$ENABLE_UPX" = true ] ; then
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-macOS-arm64"
    mv "build/AEUSTNetworkAutoLogin-macOS-arm64" "build/AEUSTNetworkAutoLogin-macOS-arm64-compressed"
fi

# Windows
# AMD64
GOOS=windows GOARCH=amd64 go build -ldflags "$BUILD_FLAGS" -o "build/AEUSTNetworkAutoLogin-windows-x86_64.exe" $SOURCE_FILE
if [ "$ENABLE_UPX" = true ] ; then
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-windows-x86_64.exe"
    mv "build/AEUSTNetworkAutoLogin-windows-x86_64.exe" "build/AEUSTNetworkAutoLogin-windows-x86_64-compressed.exe"
fi

# 386
GOOS=windows GOARCH=386 go build -ldflags "$BUILD_FLAGS" -o "build/AEUSTNetworkAutoLogin-windows-x86.exe" $SOURCE_FILE
if [ "$ENABLE_UPX" = true ] ; then
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-windows-x86.exe"
    mv "build/AEUSTNetworkAutoLogin-windows-x86.exe" "build/AEUSTNetworkAutoLogin-windows-x86-compressed.exe"
fi

echo "All builds completed!"
