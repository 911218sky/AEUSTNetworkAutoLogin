# PowerShell Script for Building Go Executables

$SOURCE_FILE = "./src/main.go"
$BUILD_FLAGS = "-s -w"
$UPX_COMPRESSION_FLAGS = "-6"
$ENABLE_UPX = $true
$BUILD_FOLDER = "build"

# Create build folder
New-Item -ItemType Directory -Force -Path $BUILD_FOLDER

# Linux
# AMD64
$env:GOOS = "linux"
$env:GOARCH = "amd64"
if ($BUILD_FLAGS -eq "") {
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64" $SOURCE_FILE
} else {
    go build -ldflags $BUILD_FLAGS -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64" $SOURCE_FILE
}
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64"
    Move-Item "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86_64-compressed" -Force
}

# 386
$env:GOARCH = "386"
if ($BUILD_FLAGS -eq "") {
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86" $SOURCE_FILE
} else {
    go build -ldflags $BUILD_FLAGS -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86" $SOURCE_FILE
}
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86"
    Move-Item "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-Linux-x86-compressed" -Force
}

# macOS
# AMD64
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
if ($BUILD_FLAGS -eq "") {
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64" $SOURCE_FILE
} else {
    go build -ldflags $BUILD_FLAGS -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64" $SOURCE_FILE
}
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64"
    Move-Item "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-x86_64-compressed" -Force
}

# ARM64
$env:GOARCH = "arm64"
if ($BUILD_FLAGS -eq "") {
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64" $SOURCE_FILE
} else {
    go build -ldflags $BUILD_FLAGS -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64" $SOURCE_FILE
}
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64"
    Move-Item "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-macOS-arm64-compressed" -Force
}

# Windows
# AMD64
$env:GOOS = "windows"
$env:GOARCH = "amd64"
if ($BUILD_FLAGS -eq "") {
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe" $SOURCE_FILE
} else {
    go build -ldflags $BUILD_FLAGS -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe" $SOURCE_FILE
}
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe"
    Move-Item "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64.exe" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86_64-compressed.exe" -Force
}

# 386
$env:GOARCH = "386"
if ($BUILD_FLAGS -eq "") {
    go build -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe" $SOURCE_FILE
} else {
    go build -ldflags $BUILD_FLAGS -o "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe" $SOURCE_FILE
}
if ($ENABLE_UPX -eq $true) {
    upx $UPX_COMPRESSION_FLAGS "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe"
    Move-Item "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86.exe" "$BUILD_FOLDER/AEUSTNetworkAutoLogin-windows-x86-compressed.exe" -Force
}

Write-Host "All builds completed!"
