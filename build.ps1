# PowerShell Script for Building Go Executables

$SOURCE_FILE = "./src/main.go"
$BUILD_FLAGS = "-s -w"
$UPX_COMPRESSION_FLAGS = "-6"
$ENABLE_UPX = $false

# Linux
# AMD64
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags $BUILD_FLAGS -o "AEUSTNetworkAutoLogin-Linux-x86_64" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "AEUSTNetworkAutoLogin-Linux-x86_64"
}

# 386
$env:GOARCH = "386"
go build -ldflags $BUILD_FLAGS -o "AEUSTNetworkAutoLogin-Linux-x86" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "AEUSTNetworkAutoLogin-Linux-x86"
}

# macOS
# AMD64
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags $BUILD_FLAGS -o "AEUSTNetworkAutoLogin-macOS-x86_64" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "AEUSTNetworkAutoLogin-macOS-x86_64"
}

# ARM64
$env:GOARCH = "arm64"
go build -ldflags $BUILD_FLAGS -o "AEUSTNetworkAutoLogin-macOS-arm64" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "AEUSTNetworkAutoLogin-macOS-arm64"
}

# Windows
# AMD64
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags $BUILD_FLAGS -o "AEUSTNetworkAutoLogin-windows-x86_64.exe" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "AEUSTNetworkAutoLogin-windows-x86_64.exe"
}

# 386
$env:GOARCH = "386"
go build -ldflags $BUILD_FLAGS -o "AEUSTNetworkAutoLogin-windows-x86.exe" $SOURCE_FILE
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "AEUSTNetworkAutoLogin-windows-x86.exe"
}

Write-Host "ALL DONE!"
