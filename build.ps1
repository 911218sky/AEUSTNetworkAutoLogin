# PowerShell Script for Building Go Executables

$SOURCE_FILE = "./src/main.go"
$BUILD_FLAGS = "-s -w"
$UPX_COMPRESSION_FLAGS = "-6"
$ENABLE_UPX = $false

# Create build folder
New-Item -ItemType Directory -Force -Path "build"

# Linux
# AMD64
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags $BUILD_FLAGS -o "build/AEUSTNetworkAutoLogin-Linux-x86_64" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-Linux-x86_64"
    Move-Item "build/AEUSTNetworkAutoLogin-Linux-x86_64" "build/AEUSTNetworkAutoLogin-Linux-x86_64-compressed" -Force
}

# 386
$env:GOARCH = "386"
go build -ldflags $BUILD_FLAGS -o "build/AEUSTNetworkAutoLogin-Linux-x86" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-Linux-x86"
    Move-Item "build/AEUSTNetworkAutoLogin-Linux-x86" "build/AEUSTNetworkAutoLogin-Linux-x86-compressed" -Force
}

# macOS
# AMD64
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags $BUILD_FLAGS -o "build/AEUSTNetworkAutoLogin-macOS-x86_64" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-macOS-x86_64"
    Move-Item "build/AEUSTNetworkAutoLogin-macOS-x86_64" "build/AEUSTNetworkAutoLogin-macOS-x86_64-compressed" -Force
}

# ARM64
$env:GOARCH = "arm64"
go build -ldflags $BUILD_FLAGS -o "build/AEUSTNetworkAutoLogin-macOS-arm64" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-macOS-arm64"
    Move-Item "build/AEUSTNetworkAutoLogin-macOS-arm64" "build/AEUSTNetworkAutoLogin-macOS-arm64-compressed" -Force
}

# Windows
# AMD64
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags $BUILD_FLAGS -o "build/AEUSTNetworkAutoLogin-windows-x86_64.exe" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-windows-x86_64.exe"
    Move-Item "build/AEUSTNetworkAutoLogin-windows-x86_64.exe" "build/AEUSTNetworkAutoLogin-windows-x86_64-compressed.exe" -Force
}

# 386
$env:GOARCH = "386"
go build -ldflags $BUILD_FLAGS -o "build/AEUSTNetworkAutoLogin-windows-x86.exe" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "build/AEUSTNetworkAutoLogin-windows-x86.exe"
    Move-Item "build/AEUSTNetworkAutoLogin-windows-x86.exe" "build/AEUSTNetworkAutoLogin-windows-x86-compressed.exe" -Force
}

Write-Host "All builds completed!"
