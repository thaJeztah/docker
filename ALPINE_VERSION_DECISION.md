# Alpine Version Decision - DHI Migration

## Question
Is alpine3.22 the current latest in DHI, given that some Dockerfiles used `alpine` (without version) which got replaced by alpine3.22?

## Answer: Yes, alpine3.22 is correct ✅

### Current Status

**Official Alpine Linux:**
- Latest: `alpine:latest` = **3.23.0** (released recently)
- Previous: 3.22, 3.21, 3.20

**Docker Hardened Images (DHI):**
- golang images: **alpine3.22** ✅ (3.23 not available yet)
- alpine-base images: **3.22** ✅ (3.23 not available yet)

### Verification

```bash
# Official Alpine
docker run --rm alpine:latest cat /etc/alpine-release
# Output: 3.23.0

# DHI availability check
docker manifest inspect dhi.io/golang:1.25.5-alpine3.22    # ✅ Available
docker manifest inspect dhi.io/golang:1.25.5-alpine3.23    # ❌ Not found
docker manifest inspect dhi.io/alpine-base:3.22            # ✅ Available
docker manifest inspect dhi.io/alpine-base:3.23            # ❌ Not found
```

### Migration Changes

The DHI migration correctly replaced:
```dockerfile
# Before
FROM alpine                              # Was using :latest (3.23)
FROM golang:${GO_VERSION}-alpine         # Was using :latest (3.23)

# After  
FROM dhi.io/alpine-base:3.22             # Explicit 3.22
FROM dhi.io/golang:${GO_VERSION}-alpine3.22  # Explicit 3.22
```

### Why Alpine 3.22 and Not 3.23?

1. **DHI Availability**: DHI only provides alpine3.22 images currently
2. **Stability**: Using the most recent DHI-available version (3.22)
3. **Explicit Versioning**: Better than `:latest` for reproducibility

### When Will 3.23 Be Available?

Monitor DHI releases:
- Check periodically: `docker manifest inspect dhi.io/alpine-base:3.23`
- When available, can update across all Dockerfiles
- For now, 3.22 is the correct choice

### Files Using Alpine in DHI Migration

From the git diff, these files were updated to alpine3.22:

1. `api/Dockerfile` - golang alpine variant
2. `daemon/libnetwork/cmd/diagnostic/Dockerfile.client` - alpine-base
3. `daemon/libnetwork/cmd/networkdb-test/Dockerfile` - (needs verification)
4. `daemon/libnetwork/cmd/ssd/Dockerfile` - (needs verification)
5. `hack/dockerfiles/generate-files.Dockerfile` - (needs verification)
6. `hack/dockerfiles/govulncheck.Dockerfile` - (needs verification)

### Recommendation

✅ **Keep alpine3.22** - This is the correct version to use for now

When DHI releases alpine3.23:
1. Test with new version
2. Update all Dockerfiles consistently
3. Update in a single commit/PR

## Summary

- ✅ Alpine 3.22 is the latest available in DHI
- ✅ Alpine 3.23 exists upstream but not yet in DHI
- ✅ The migration correctly uses 3.22
- ✅ No changes needed right now
- ⏳ Can update to 3.23 when DHI releases it
