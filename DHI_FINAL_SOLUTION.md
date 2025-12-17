# DHI Migration - Complete Solution Summary

## Problem

After migrating to Docker Hardened Images (DHI), users needed to specify TWO build arguments instead of one:
- `BASE_DEBIAN_DISTRO` (codename: bookworm, trixie, bullseye)
- `DEBIAN_VERSION` (DHI version: debian12, debian13, debian11)

This was necessary because DHI images use inconsistent tagging:
- `dhi.io/debian-base` uses codenames
- `dhi.io/golang` uses version numbers

## Solution

We implemented **two complementary mapping mechanisms** to restore the single-argument interface:

### 1. docker-bake.hcl (For Bake Builds)
**Automatic mapping** using a function in the bake file

### 2. hack/dhi-build-wrapper.sh (For Direct Docker Builds)  
**Wrapper script** that injects the correct DEBIAN_VERSION

## What Was Implemented

### Changes to docker-bake.hcl

```hcl
# New variable
variable "BASE_DEBIAN_DISTRO" {
  default = "bookworm"
}

# Mapping function
function "debian_version" {
  params = [codename]
  result = codename == "bullseye" ? "debian11" : 
           codename == "bookworm" ? "debian12" : 
           codename == "trixie" ? "debian13" : 
           "debian12"
}

# Updated _common target to pass both args
target "_common" {
  args = {
    BASE_DEBIAN_DISTRO = BASE_DEBIAN_DISTRO
    DEBIAN_VERSION = debian_version(BASE_DEBIAN_DISTRO)
    ...
  }
}
```

### New File: hack/dhi-build-wrapper.sh

A bash script that:
1. Reads `BASE_DEBIAN_DISTRO` from env var or `--build-arg`
2. Maps it to the corresponding `DEBIAN_VERSION`
3. Injects both arguments into the `docker build` command

### Documentation

- **DHI_COMPLETE_GUIDE.md** - Comprehensive user guide
- **test_dhi_wrapper.sh** - Automated tests for the wrapper

## User Experience

### Before (Two Arguments Required)
```bash
# ❌ Would fail
docker build --build-arg BASE_DEBIAN_DISTRO=trixie .

# ✅ Had to do this
docker build \
  --build-arg BASE_DEBIAN_DISTRO=trixie \
  --build-arg DEBIAN_VERSION=debian13 \
  .
```

### After (Single Argument)

**Option 1: Make (recommended)**
```bash
BASE_DEBIAN_DISTRO=trixie make binary
```

**Option 2: Docker Bake**
```bash
BASE_DEBIAN_DISTRO=trixie docker buildx bake binary
```

**Option 3: Docker Build with Wrapper**
```bash
BASE_DEBIAN_DISTRO=trixie ./hack/dhi-build-wrapper.sh build .
# or
./hack/dhi-build-wrapper.sh build --build-arg BASE_DEBIAN_DISTRO=trixie .
```

**Option 4: Direct Docker Build (manual)**
```bash
# Still works but requires both args
docker build \
  --build-arg BASE_DEBIAN_DISTRO=trixie \
  --build-arg DEBIAN_VERSION=debian13 \
  .
```

## Coverage

| Build Method | Mapping Solution | Status |
|--------------|------------------|--------|
| `make` | docker-bake.hcl | ✅ Automatic |
| `docker buildx bake` | docker-bake.hcl | ✅ Automatic |
| `docker build` | dhi-build-wrapper.sh | ✅ Automatic (with wrapper) |
| GitHub Actions | docker-bake.hcl | ✅ Automatic |
| Direct `docker build` | Manual | ⚠️ Manual (both args required) |

## Zero Breaking Changes

- ✅ Makefile - No changes required
- ✅ GitHub Actions - No changes required
- ✅ Existing scripts - Continue to work
- ✅ CI/CD pipelines - No changes required

## Verification

### Test Bake Mapping
```bash
# Default
docker buildx bake --print binary | jq '.target.binary.args | {BASE_DEBIAN_DISTRO, DEBIAN_VERSION}'
# Output: {"BASE_DEBIAN_DISTRO": "bookworm", "DEBIAN_VERSION": "debian12"}

# Trixie
BASE_DEBIAN_DISTRO=trixie docker buildx bake --print binary | jq '.target.binary.args | {BASE_DEBIAN_DISTRO, DEBIAN_VERSION}'
# Output: {"BASE_DEBIAN_DISTRO": "trixie", "DEBIAN_VERSION": "debian13"}
```

### Test Wrapper Script
```bash
bash test_dhi_wrapper.sh
# Runs 4 automated tests covering all scenarios
```

### Test Real Build
```bash
# Using make
BASE_DEBIAN_DISTRO=trixie make binary

# Using bake
BASE_DEBIAN_DISTRO=trixie docker buildx bake binary

# Using wrapper
BASE_DEBIAN_DISTRO=trixie ./hack/dhi-build-wrapper.sh build .
```

## Files Modified

1. **docker-bake.hcl**
   - Added `BASE_DEBIAN_DISTRO` variable
   - Added `debian_version()` function
   - Updated `_common` target to pass both arguments

2. **Dockerfile**
   - Added `DEBIAN_VERSION` argument with default
   - Updated comments to document the mapping
   - Changed golang image to use `DEBIAN_VERSION`

## Files Added

1. **hack/dhi-build-wrapper.sh**
   - Bash script for docker build mapping
   - Handles env vars and --build-arg
   - Shows informative messages

2. **DHI_COMPLETE_GUIDE.md**
   - Complete documentation
   - Usage examples for all methods
   - Troubleshooting guide

3. **test_dhi_wrapper.sh**
   - Automated test suite
   - Validates all wrapper scenarios

## Supported Scenarios

✅ Environment variable: `BASE_DEBIAN_DISTRO=trixie make binary`  
✅ Build argument: `./hack/dhi-build-wrapper.sh build --build-arg BASE_DEBIAN_DISTRO=trixie .`  
✅ Default (no args): `make binary` uses bookworm/debian12  
✅ Multiple targets: All bake targets inherit the mapping  
✅ CI/CD: GitHub Actions using docker/bake-action work automatically  
✅ Pass-through: Wrapper doesn't interfere if DEBIAN_VERSION already specified  

## Design Decisions

### Why Two Solutions?

1. **docker-bake.hcl** (for bake builds)
   - Native HCL functionality
   - No external scripts needed
   - Works with make and bake
   - Used by GitHub Actions

2. **dhi-build-wrapper.sh** (for docker build)
   - Covers direct `docker build` usage
   - Works with existing scripts
   - Optional (users can still specify both args manually)
   - Clean, simple bash script

### Why Not Dockerfile-Only?

Docker BuildKit doesn't support computed ARG defaults based on other ARGs. We can't do:
```dockerfile
ARG BASE_DEBIAN_DISTRO="bookworm"
ARG DEBIAN_VERSION="${BASE_DEBIAN_DISTRO == 'trixie' ? 'debian13' : 'debian12'}"  # ❌ Not possible
```

Therefore, the mapping must happen **before** the Dockerfile receives the arguments.

### Why Not a Single Wrapper for Everything?

- `make` already uses `docker buildx bake` internally
- GitHub Actions already uses `docker/bake-action`
- Wrapping these would be redundant and add complexity
- The bake file provides clean, native mapping for these use cases

## Recommendations

1. **For daily development:** Use `make` commands
2. **For custom builds:** Use `docker buildx bake`
3. **For scripts using docker build:** Use `./hack/dhi-build-wrapper.sh`
4. **For CI/CD:** Use `docker buildx bake` or `docker/bake-action`

## Next Steps

1. Review the implementation in `docker-bake.hcl` and `hack/dhi-build-wrapper.sh`
2. Run verification tests:
   ```bash
   # Test bake mapping
   BASE_DEBIAN_DISTRO=trixie docker buildx bake --print binary | jq
   
   # Test wrapper
   bash test_dhi_wrapper.sh
   ```
3. Test a real build:
   ```bash
   BASE_DEBIAN_DISTRO=trixie make binary
   ```
4. Read the complete guide: `DHI_COMPLETE_GUIDE.md`
5. Commit the changes when satisfied

## Benefits

✅ Single build argument restored  
✅ Zero breaking changes  
✅ Clean, maintainable solution  
✅ Works across all build methods  
✅ Well documented  
✅ Automated tests included  
✅ No external dependencies  
