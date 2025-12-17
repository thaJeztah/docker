#!/usr/bin/env bash
#
# Test script for DHI build wrapper
#
set -e

echo "Testing DHI Build Wrapper"
echo "=========================="
echo

# Create a test Dockerfile that just echoes the build args
cat > /tmp/Dockerfile.dhi-test << 'EOF'
FROM alpine:latest
ARG BASE_DEBIAN_DISTRO
ARG DEBIAN_VERSION
RUN echo "BASE_DEBIAN_DISTRO=${BASE_DEBIAN_DISTRO}"
RUN echo "DEBIAN_VERSION=${DEBIAN_VERSION}"
EOF

echo "Test 1: Default (bookworm)"
echo "---"
./hack/dhi-build-wrapper.sh build \
    -f /tmp/Dockerfile.dhi-test \
    --progress=plain \
    --no-cache \
    /tmp 2>&1 | grep -E "DHI Build Wrapper|BASE_DEBIAN_DISTRO|DEBIAN_VERSION" || true
echo

echo "Test 2: Trixie via --build-arg"
echo "---"
./hack/dhi-build-wrapper.sh build \
    --build-arg BASE_DEBIAN_DISTRO=trixie \
    -f /tmp/Dockerfile.dhi-test \
    --progress=plain \
    --no-cache \
    /tmp 2>&1 | grep -E "DHI Build Wrapper|BASE_DEBIAN_DISTRO|DEBIAN_VERSION" || true
echo

echo "Test 3: Bullseye via environment variable"
echo "---"
BASE_DEBIAN_DISTRO=bullseye ./hack/dhi-build-wrapper.sh build \
    -f /tmp/Dockerfile.dhi-test \
    --progress=plain \
    --no-cache \
    /tmp 2>&1 | grep -E "DHI Build Wrapper|BASE_DEBIAN_DISTRO|DEBIAN_VERSION" || true
echo

echo "Test 4: Both args specified (should pass through)"
echo "---"
QUIET=1 ./hack/dhi-build-wrapper.sh build \
    --build-arg BASE_DEBIAN_DISTRO=trixie \
    --build-arg DEBIAN_VERSION=debian13 \
    -f /tmp/Dockerfile.dhi-test \
    --progress=plain \
    --no-cache \
    /tmp 2>&1 | grep -E "DHI Build Wrapper|BASE_DEBIAN_DISTRO|DEBIAN_VERSION" || true
echo

# Cleanup
rm /tmp/Dockerfile.dhi-test

echo "=========================="
echo "All tests completed!"
