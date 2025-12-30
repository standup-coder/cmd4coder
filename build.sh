#!/bin/bash
# Build script for cmd4coder
# Builds executables for multiple platforms

VERSION="1.0.0"
BUILD_DIR="build"
APP_NAME="cmd4coder"
COMMIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

echo "Building cmd4coder v${VERSION}"
echo "Commit: ${COMMIT_HASH}"
echo "Build Time: ${BUILD_TIME}"
echo "================================"
echo ""

# Run tests before building
echo "Running tests..."
go test ./... -cover
if [ $? -ne 0 ]; then
    echo "Tests failed! Aborting build."
    exit 1
fi
echo "All tests passed!"
echo ""

# Clean build directory
if [ -d "$BUILD_DIR" ]; then
    echo "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
fi
mkdir -p "$BUILD_DIR"

# Build for different platforms
platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r GOOS GOARCH <<< "$platform"
    
    output_name="${APP_NAME}-v${VERSION}-${GOOS}-${GOARCH}"
    if [ "$GOOS" = "windows" ]; then
        output_name="${output_name}.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags "-s -w -X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.CommitHash=${COMMIT_HASH}'" \
        -trimpath \
        -o "${BUILD_DIR}/${output_name}" \
        ./cmd/cli
    
    if [ $? -eq 0 ]; then
        file_size=$(ls -lh "${BUILD_DIR}/${output_name}" | awk '{print $5}')
        echo "✓ Successfully built: ${output_name} (${file_size})"
        
        # Calculate SHA256 checksum
        if command -v sha256sum >/dev/null 2>&1; then
            (cd "$BUILD_DIR" && sha256sum "$output_name" >> checksums.txt)
        elif command -v shasum >/dev/null 2>&1; then
            (cd "$BUILD_DIR" && shasum -a 256 "$output_name" >> checksums.txt)
        fi
        
        # Create archive
        cd "$BUILD_DIR"
        if [ "$GOOS" = "windows" ]; then
            archive_name="${output_name%.exe}"
            if command -v zip >/dev/null 2>&1; then
                zip "${archive_name}.zip" "$output_name" >/dev/null
                echo "  Created archive: ${archive_name}.zip"
            fi
        else
            if command -v tar >/dev/null 2>&1; then
                tar -czf "${output_name}.tar.gz" "$output_name"
                echo "  Created archive: ${output_name}.tar.gz"
            fi
        fi
        cd ..
    else
        echo "✗ Failed to build for $GOOS/$GOARCH"
        exit 1
    fi
done

echo ""
echo "================================"
echo "Build completed successfully!"
echo "Version: v${VERSION}"
echo "Commit: ${COMMIT_HASH}"
echo "Artifacts are in the ${BUILD_DIR}/ directory"
echo ""
echo "Build artifacts:"
ls -lh "${BUILD_DIR}/"
echo ""
if [ -f "${BUILD_DIR}/checksums.txt" ]; then
    echo "Checksums saved to ${BUILD_DIR}/checksums.txt"
fi
