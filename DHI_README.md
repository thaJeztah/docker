# DHI Build Solution - README

## Quick Summary

This solution restores the single-argument build interface after the DHI migration by providing automatic mapping from Debian codenames to DHI version numbers.

## What You Need to Know

### For Daily Development (Using Make)
```bash
BASE_DEBIAN_DISTRO=trixie make binary
```
✅ Works automatically - no changes needed!

### For Docker Bake Users
```bash
BASE_DEBIAN_DISTRO=trixie docker buildx bake binary
```
✅ Works automatically via docker-bake.hcl

### For Direct Docker Build Users
```bash
BASE_DEBIAN_DISTRO=trixie ./hack/dhi-build-wrapper.sh build .
```
✅ Use the wrapper script for automatic mapping

### For GitHub Actions
```yaml
env:
  BASE_DEBIAN_DISTRO: trixie
```
✅ Works automatically via docker/bake-action

## Files to Review

1. **docker-bake.hcl** - Added `debian_version()` function for automatic mapping
2. **hack/dhi-build-wrapper.sh** - Wrapper script for direct docker build
3. **DHI_COMPLETE_GUIDE.md** - Full documentation
4. **test_dhi_wrapper.sh** - Automated tests

## Quick Test

```bash
# Test bake mapping
BASE_DEBIAN_DISTRO=trixie docker buildx bake --print binary | jq '.target.binary.args'

# Test wrapper
bash test_dhi_wrapper.sh
```

## Supported Mappings

- `bookworm` → `debian12` (default)
- `trixie` → `debian13`
- `bullseye` → `debian11`

## Zero Breaking Changes

✅ All existing workflows continue to work  
✅ No Makefile changes required  
✅ No CI/CD changes required  
✅ Optional wrapper for docker build users  

## Read More

- **DHI_COMPLETE_GUIDE.md** - Comprehensive usage guide
- **DHI_FINAL_SOLUTION.md** - Technical solution details
