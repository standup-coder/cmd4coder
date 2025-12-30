# Build script for cmd4coder (PowerShell)
# Builds executables for multiple platforms

$VERSION = "1.0.0"
$BUILD_DIR = "build"
$APP_NAME = "cmd4coder"
$COMMIT_HASH = git rev-parse --short HEAD 2>$null
if (-not $COMMIT_HASH) { $COMMIT_HASH = "unknown" }
$BUILD_TIME = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

Write-Host "Building cmd4coder v$VERSION" -ForegroundColor Cyan
Write-Host "Commit: $COMMIT_HASH" -ForegroundColor Cyan
Write-Host "Build Time: $BUILD_TIME" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan

# Run tests before building
Write-Host "Running tests..." -ForegroundColor Yellow
$testResult = go test ./... -cover
if ($LASTEXITCODE -ne 0) {
    Write-Host "Tests failed! Aborting build." -ForegroundColor Red
    exit 1
}
Write-Host "All tests passed!" -ForegroundColor Green
Write-Host "" -ForegroundColor Cyan

# Clean build directory
if (Test-Path $BUILD_DIR) {
    Write-Host "Cleaning build directory..." -ForegroundColor Yellow
    Remove-Item -Recurse -Force $BUILD_DIR
}
New-Item -ItemType Directory -Force -Path $BUILD_DIR | Out-Null

# Build for different platforms
$platforms = @(
    @{OS="linux"; ARCH="amd64"},
    @{OS="linux"; ARCH="arm64"},
    @{OS="darwin"; ARCH="amd64"},
    @{OS="darwin"; ARCH="arm64"},
    @{OS="windows"; ARCH="amd64"}
)

$BuildTime = $BUILD_TIME

foreach ($platform in $platforms) {
    $GOOS = $platform.OS
    $GOARCH = $platform.ARCH
    
    $output_name = "${APP_NAME}-v${VERSION}-${GOOS}-${GOARCH}"
    if ($GOOS -eq "windows") {
        $output_name = "${output_name}.exe"
    }
    
    Write-Host "Building for $GOOS/$GOARCH..." -ForegroundColor Green
    
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    $env:CGO_ENABLED = "0"
    
    $ldflags = "-s -w -X 'main.Version=$VERSION' -X 'main.BuildTime=$BUILD_TIME' -X 'main.CommitHash=$COMMIT_HASH'"
    
    & go build -ldflags $ldflags -trimpath -o "$BUILD_DIR\$output_name" .\cmd\cli
    
    if ($LASTEXITCODE -eq 0) {
        $fileSize = (Get-Item "$BUILD_DIR\$output_name").Length
        $fileSizeMB = [math]::Round($fileSize / 1MB, 2)
        Write-Host "✓ Successfully built: $output_name ($fileSizeMB MB)" -ForegroundColor Green
        
        # Calculate SHA256 checksum
        $hash = (Get-FileHash "$BUILD_DIR\$output_name" -Algorithm SHA256).Hash
        Add-Content -Path "$BUILD_DIR\checksums.txt" -Value "${hash}  $output_name"
        
        # Create archive
        Push-Location $BUILD_DIR
        if ($GOOS -eq "windows") {
            $archiveName = $output_name -replace '\.exe$',''
            Compress-Archive -Path $output_name -DestinationPath "${archiveName}.zip" -Force
            Write-Host "  Created archive: ${archiveName}.zip" -ForegroundColor Gray
        } else {
            # For Linux/Mac, create tar.gz if tar is available
            Write-Host "  Binary saved (create .tar.gz manually if needed)" -ForegroundColor Gray
        }
        Pop-Location
    } else {
        Write-Host "✗ Failed to build for $GOOS/$GOARCH" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
Write-Host "================================" -ForegroundColor Cyan
Write-Host "Build completed successfully!" -ForegroundColor Green
Write-Host "Version: v$VERSION" -ForegroundColor Cyan
Write-Host "Commit: $COMMIT_HASH" -ForegroundColor Cyan
Write-Host "Artifacts are in the $BUILD_DIR\ directory" -ForegroundColor Cyan
Write-Host ""
Write-Host "Build artifacts:" -ForegroundColor Cyan
Get-ChildItem -Path $BUILD_DIR | Format-Table Name, Length, LastWriteTime

Write-Host ""
Write-Host "Checksums saved to $BUILD_DIR\checksums.txt" -ForegroundColor Cyan

# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""
$env:CGO_ENABLED = ""
