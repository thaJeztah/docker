# Docker Hardened Images (DHI) Build Guide

## Overview

The Moby project Dockerfiles have been migrated to use Docker Hardened Images (DHI) for enhanced security. The build system automatically handles the mapping between Debian release codenames and DHI version numbers.

## Building with Different Debian Versions

### Using Make (Recommended)

The Makefile uses `docker buildx bake` which automatically handles the Debian version mapping:

```bash
# Default (bookworm/debian12)
make binary

# Using Debian Trixie
BASE_DEBIAN_DISTRO=trixie make binary

# Using Debian Bullseye
BASE_DEBIAN_DISTRO=bullseye make binary
```

### Using Docker Bake Directly

```bash
# Default (bookworm/debian12)
docker buildx bake binary

# Using Debian Trixie
BASE_DEBIAN_DISTRO=trixie docker buildx bake binary

# Using Debian Bullseye
BASE_DEBIAN_DISTRO=bullseye docker buildx bake binary
```

### Using Docker Build Directly

#### Option 1: Using the Wrapper Script (Recommended)

The wrapper script automatically maps the Debian codename to the correct DHI version:

```bash
# Via environment variable
BASE_DEBIAN_DISTRO=trixie ./hack/dhi-build-wrapper.sh build .

# Via build argument  
./hack/dhi-build-wrapper.sh build --build-arg BASE_DEBIAN_DISTRO=trixie .

# Default (bookworm)
./hack/dhi-build-wrapper.sh build .
```

The wrapper works as a drop-in replacement for `docker build` and automatically injects the correct `DEBIAN_VERSION` argument.

#### Option 2: Specifying Both Arguments Manually

When using `docker build` directly (without wrapper or bake), you must specify both arguments manually:

```bash
# Bookworm
docker build --build-arg BASE_DEBIAN_DISTRO=bookworm --build-arg DEBIAN_VERSION=debian12 .

# Trixie
docker build --build-arg BASE_DEBIAN_DISTRO=trixie --build-arg DEBIAN_VERSION=debian13 .

# Bullseye
docker build --build-arg BASE_DEBIAN_DISTRO=bullseye --build-arg DEBIAN_VERSION=debian11 .
```

## Supported Debian Versions

| Codename | DHI Version | Support Status |
|----------|-------------|----------------|
| bullseye | debian11    | Supported      |
| bookworm | debian12    | Default        |
| trixie   | debian13    | Supported      |

## How It Works

### The Mapping Problem

DHI uses inconsistent tagging across different image types:
- `dhi.io/debian-base` uses codenames: `bookworm`, `trixie`, `bullseye`
- `dhi.io/golang` uses version numbers: `debian12`, `debian13`, `debian11`

### The Solution

The `docker-bake.hcl` file contains a `debian_version()` function that automatically maps codenames to version numbers:

```hcl
function "debian_version" {
  params = [codename]
  result = codename == "bullseye" ? "debian11" : 
           codename == "bookworm" ? "debian12" : 
           codename == "trixie" ? "debian13" : 
           "debian12"  # default
}
```

When you pass `BASE_DEBIAN_DISTRO` as an environment variable, the bake file automatically:
1. Sets `BASE_DEBIAN_DISTRO` to your chosen codename
2. Calculates `DEBIAN_VERSION` using the mapping function
3. Passes both as build arguments to the Dockerfile

### User Experience

**Before DHI migration:**
```bash
# Single argument
make binary  # uses bookworm
BASE_DEBIAN_DISTRO=trixie make binary  # uses trixie
```

**After DHI migration (with bake file solution):**
```bash
# Still single argument!
make binary  # uses bookworm (debian12)
BASE_DEBIAN_DISTRO=trixie make binary  # uses trixie (debian13)
```

The user experience remains unchanged - you only need to specify `BASE_DEBIAN_DISTRO`.

## GitHub Actions

GitHub Actions workflows already use `docker/bake-action`, so they automatically benefit from the mapping. To change the Debian version in CI:

```yaml
- name: Build
  uses: docker/bake-action@v6
  with:
    targets: binary
  env:
    BASE_DEBIAN_DISTRO: trixie
```

## Verification

To verify the mapping is working correctly:

```bash
# Check what will be built
docker buildx bake --print binary

# Check with different distro
BASE_DEBIAN_DISTRO=trixie docker buildx bake --print binary

# Look for these lines in the output:
#   "BASE_DEBIAN_DISTRO": "trixie",
#   "DEBIAN_VERSION": "debian13",
```

## Troubleshooting

### "Invalid build argument" error

If you see errors about invalid build arguments, make sure you're using the bake file:
- ✅ `make binary` (uses bake)
- ✅ `docker buildx bake binary` (uses bake)
- ❌ `docker build .` (direct build, needs both args)

### Wrong Debian version being used

Verify your environment variable is being passed:
```bash
# This should show your distro in the output
BASE_DEBIAN_DISTRO=trixie docker buildx bake --print binary | grep BASE_DEBIAN_DISTRO
```

### Mapping not working

The mapping only works when passing `BASE_DEBIAN_DISTRO` as an environment variable to the bake command. Using `--set` won't trigger re-evaluation:
- ✅ `BASE_DEBIAN_DISTRO=trixie docker buildx bake binary`
- ❌ `docker buildx bake binary --set *.args.BASE_DEBIAN_DISTRO=trixie`

## Benefits of DHI

- **Enhanced Security**: Minimal, hardened base images with reduced attack surface
- **Non-root by Default**: Images run as non-root users where applicable
- **TLS Certificates**: Included by default
- **Reproducible Builds**: Specific versions ensure consistency
- **Maintained Compatibility**: All existing build processes continue to work
