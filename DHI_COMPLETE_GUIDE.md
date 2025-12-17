# Docker Hardened Images (DHI) - Complete Build Guide

## Overview

The Moby project has migrated to Docker Hardened Images (DHI) for enhanced security. This guide covers all build methods and how the automatic Debian version mapping works.

## Quick Start

### Using Make (Recommended)
```bash
# Default (bookworm/debian12)
make binary

# Using Debian Trixie
BASE_DEBIAN_DISTRO=trixie make binary

# Using Debian Bullseye
BASE_DEBIAN_DISTRO=bullseye make binary
```

### Using Docker Bake
```bash
# Default (bookworm/debian12)
docker buildx bake binary

# Using different Debian version
BASE_DEBIAN_DISTRO=trixie docker buildx bake binary
```

### Using Docker Build with Wrapper
```bash
# Via environment variable
BASE_DEBIAN_DISTRO=trixie ./hack/dhi-build-wrapper.sh build .

# Via build argument
./hack/dhi-build-wrapper.sh build --build-arg BASE_DEBIAN_DISTRO=trixie .
```

### Using Docker Build Directly (Manual)
```bash
# You must specify BOTH arguments
docker build \
  --build-arg BASE_DEBIAN_DISTRO=trixie \
  --build-arg DEBIAN_VERSION=debian13 \
  .
```

## Supported Debian Versions

| Codename | DHI Version | Support Status |
|----------|-------------|----------------|
| bullseye | debian11    | Supported      |
| bookworm | debian12    | Default        |
| trixie   | debian13    | Supported      |

## How the Mapping Works

DHI uses inconsistent tagging across different image types:
- **debian-base images** use codenames: `bookworm`, `trixie`, `bullseye`
- **golang images** use version numbers: `debian12`, `debian13`, `debian11`

To provide a seamless single-argument experience, we have two mapping mechanisms:

### 1. docker-bake.hcl (For Bake Builds)

The `docker-bake.hcl` file contains a `debian_version()` function:

```hcl
function "debian_version" {
  params = [codename]
  result = codename == "bullseye" ? "debian11" : 
           codename == "bookworm" ? "debian12" : 
           codename == "trixie" ? "debian13" : 
           "debian12"
}
```

When you use `make` or `docker buildx bake`:
1. You pass `BASE_DEBIAN_DISTRO` as an environment variable
2. The bake file automatically calculates `DEBIAN_VERSION`
3. Both arguments are passed to the Dockerfile

**Works with:**
- `make` commands
- `docker buildx bake` commands
- GitHub Actions (uses docker/bake-action)

### 2. dhi-build-wrapper.sh (For Direct Docker Build)

The wrapper script `hack/dhi-build-wrapper.sh` provides the same mapping for `docker build`:

```bash
# Wrapper automatically adds DEBIAN_VERSION
./hack/dhi-build-wrapper.sh build --build-arg BASE_DEBIAN_DISTRO=trixie .

# Equivalent to:
docker build \
  --build-arg BASE_DEBIAN_DISTRO=trixie \
  --build-arg DEBIAN_VERSION=debian13 \
  .
```

**Works with:**
- Direct `docker build` commands
- Scripts that call `docker build`
- Local development builds

## Detailed Usage

### Make Targets

All make targets that build Docker images support the `BASE_DEBIAN_DISTRO` variable:

```bash
# Build static binary
BASE_DEBIAN_DISTRO=trixie make binary

# Build dynamic binary
BASE_DEBIAN_DISTRO=trixie make dynbinary

# Start development container
BASE_DEBIAN_DISTRO=bullseye make dev

# Run tests
BASE_DEBIAN_DISTRO=bookworm make test-integration
```

### Docker Bake Targets

```bash
# Build specific target
BASE_DEBIAN_DISTRO=trixie docker buildx bake binary

# Build all targets
BASE_DEBIAN_DISTRO=trixie docker buildx bake all

# Development container
BASE_DEBIAN_DISTRO=bookworm docker buildx bake dev
```

### Using the Wrapper Script

The wrapper can be used as a drop-in replacement for `docker build`:

```bash
# Basic usage
./hack/dhi-build-wrapper.sh build .

# With build args
./hack/dhi-build-wrapper.sh build \
  --build-arg BASE_DEBIAN_DISTRO=trixie \
  --tag myimage:latest \
  .

# With environment variable
BASE_DEBIAN_DISTRO=bullseye ./hack/dhi-build-wrapper.sh build .

# Suppress wrapper messages
QUIET=1 ./hack/dhi-build-wrapper.sh build .
```

**Note:** The wrapper is only needed for `docker build`. For `docker buildx bake`, use the command directly as the bake file handles mapping automatically.

### GitHub Actions

No changes needed! Workflows already use `docker/bake-action`:

```yaml
- name: Build
  uses: docker/bake-action@v6
  with:
    targets: binary
  env:
    BASE_DEBIAN_DISTRO: trixie  # Optional: override default
```

## Verification

### Verify Bake Mapping

```bash
# Check default mapping
docker buildx bake --print binary | jq '.target.binary.args | {BASE_DEBIAN_DISTRO, DEBIAN_VERSION}'

# Output:
# {
#   "BASE_DEBIAN_DISTRO": "bookworm",
#   "DEBIAN_VERSION": "debian12"
# }

# Check with different distro
BASE_DEBIAN_DISTRO=trixie docker buildx bake --print binary | \
  jq '.target.binary.args | {BASE_DEBIAN_DISTRO, DEBIAN_VERSION}'

# Output:
# {
#   "BASE_DEBIAN_DISTRO": "trixie",
#   "DEBIAN_VERSION": "debian13"
# }
```

### Verify Wrapper Script

Run the automated test:

```bash
bash test_dhi_wrapper.sh
```

Or manually verify:

```bash
# Watch the wrapper inject DEBIAN_VERSION
./hack/dhi-build-wrapper.sh build \
  --build-arg BASE_DEBIAN_DISTRO=trixie \
  --progress=plain \
  .
# Look for: 🔧 DHI Build Wrapper: trixie → debian13
```

## Troubleshooting

### Error: Mismatched Debian Versions

**Problem:** Build fails with image pull errors like "manifest unknown"

**Solution:** Ensure you're using one of these methods:
- ✅ `make` or `docker buildx bake` (automatic mapping)
- ✅ `./hack/dhi-build-wrapper.sh build` (automatic mapping)
- ✅ Manual `docker build` with BOTH arguments specified correctly

### Error: Unknown Flag

**Problem:** `unknown flag: --build-arg`

**Solution:** Make sure `--build-arg` comes after the `build` command:
```bash
# ❌ Wrong
docker --build-arg BASE_DEBIAN_DISTRO=trixie build .

# ✅ Correct
docker build --build-arg BASE_DEBIAN_DISTRO=trixie .
```

### Wrapper Not Working

**Problem:** Wrapper doesn't seem to add DEBIAN_VERSION

**Symptoms:** Build uses wrong version or fails

**Debug:**
```bash
# Enable wrapper messages (default)
./hack/dhi-build-wrapper.sh build ...
# Should show: 🔧 DHI Build Wrapper: <distro> → <version>

# Check if DEBIAN_VERSION is already specified
grep "DEBIAN_VERSION" your-command
# If already specified, wrapper passes through as-is
```

### Bake Mapping Not Applied

**Problem:** Using `--set` doesn't trigger mapping

**Solution:** Pass `BASE_DEBIAN_DISTRO` as an environment variable, not via `--set`:
```bash
# ❌ Wrong - won't trigger mapping
docker buildx bake binary --set *.args.BASE_DEBIAN_DISTRO=trixie

# ✅ Correct - triggers mapping
BASE_DEBIAN_DISTRO=trixie docker buildx bake binary
```

## Benefits of DHI

- **Enhanced Security**: Minimal, hardened base images with reduced attack surface
- **Non-root by Default**: Images run as non-root users where applicable  
- **TLS Certificates**: Included by default
- **Reproducible Builds**: Specific versions ensure consistency
- **Regular Updates**: Maintained with security patches

## Migration Notes

### For Existing Scripts

If you have scripts that call `docker build`:

**Option 1: Use the wrapper (recommended)**
```bash
# Before
docker build --build-arg BASE_DEBIAN_DISTRO=trixie .

# After
./hack/dhi-build-wrapper.sh build --build-arg BASE_DEBIAN_DISTRO=trixie .
```

**Option 2: Add DEBIAN_VERSION manually**
```bash
# Map yourself
case "$BASE_DEBIAN_DISTRO" in
  bullseye) DEBIAN_VERSION="debian11" ;;
  bookworm) DEBIAN_VERSION="debian12" ;;
  trixie)   DEBIAN_VERSION="debian13" ;;
esac

docker build \
  --build-arg BASE_DEBIAN_DISTRO="$BASE_DEBIAN_DISTRO" \
  --build-arg DEBIAN_VERSION="$DEBIAN_VERSION" \
  .
```

### For CI/CD Pipelines

If using `docker/bake-action` → **No changes needed!**

If using `docker build` directly:
1. Switch to `docker buildx bake` (recommended)
2. Or wrap calls with `./hack/dhi-build-wrapper.sh`
3. Or add DEBIAN_VERSION mapping to your pipeline

## Summary

| Build Method | Mapping | User Action |
|--------------|---------|-------------|
| `make` | Automatic (via bake) | Pass `BASE_DEBIAN_DISTRO` env var |
| `docker buildx bake` | Automatic (via bake file) | Pass `BASE_DEBIAN_DISTRO` env var |
| `docker build` with wrapper | Automatic (via script) | Use `./hack/dhi-build-wrapper.sh` |
| `docker build` direct | Manual | Pass BOTH arguments |
| GitHub Actions (bake) | Automatic (via bake) | Pass `BASE_DEBIAN_DISTRO` env var |

**Recommendation:** Use `make` or `docker buildx bake` for the best experience.
