#!/usr/bin/env bash
set -euo pipefail

VERSION=${1:-"1.0.0"}
DIST_DIR="dist"

echo "=== Démarrage du build de gopack v${VERSION} ==="

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"

TARGETS=(
    "linux/amd64"
    "windows/amd64"
    "linux/arm64"
)

LDFLAGS="-s -w -X main.version=${VERSION}"

for TARGET in "${TARGETS[@]}"; do
    GOOS=${TARGET%/*}
    GOARCH=${TARGET#*/}
    
    OUTPUT_NAME="gopack-${GOOS}-${GOARCH}"
    if [ "${GOOS}" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi

    echo "Building pour ${GOOS}/${GOARCH} -> ${DIST_DIR}/${OUTPUT_NAME}..."
    GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags "${LDFLAGS}" -o "${DIST_DIR}/${OUTPUT_NAME}" main.go
done

echo "Génération des sommes de contrôle SHA256..."
cd "${DIST_DIR}"
sha256sum gopack-* > SHA256SUMS
cd ..

echo "=== Build terminé avec succès ==="
ls -lh "${DIST_DIR}"