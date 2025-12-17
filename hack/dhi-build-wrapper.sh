#!/usr/bin/env bash
#
# Docker Hardened Images (DHI) Build Wrapper
#
# This wrapper script automatically maps Debian distribution codenames to DHI version numbers
# for 'docker build' commands.
#
# Usage:
#   ./hack/dhi-build-wrapper.sh build --build-arg BASE_DEBIAN_DISTRO=trixie ...
#   BASE_DEBIAN_DISTRO=trixie ./hack/dhi-build-wrapper.sh build ...
#
# The script will automatically add the matching DEBIAN_VERSION argument.
#
# For 'docker buildx bake', use the native command directly as docker-bake.hcl
# handles the mapping automatically:
#   BASE_DEBIAN_DISTRO=trixie docker buildx bake binary

set -e

# Debian codename to DHI version mapping
declare -A DEBIAN_VERSIONS=(
    ["bullseye"]="debian11"
    ["bookworm"]="debian12"
    ["trixie"]="debian13"
)

DEFAULT_DISTRO="bookworm"
DEFAULT_VERSION="debian12"

# Function to get DEBIAN_VERSION from codename
get_debian_version() {
    local codename="$1"
    echo "${DEBIAN_VERSIONS[$codename]:-$DEFAULT_VERSION}"
}

# Function to extract BASE_DEBIAN_DISTRO from command line arguments
extract_distro_from_args() {
    local distro=""
    
    # Check environment variable first
    if [ -n "$BASE_DEBIAN_DISTRO" ]; then
        distro="$BASE_DEBIAN_DISTRO"
        echo "$distro"
        return
    fi
    
    # Parse command line arguments
    for arg in "$@"; do
        if [[ "$arg" =~ ^BASE_DEBIAN_DISTRO=(.+)$ ]]; then
            distro="${BASH_REMATCH[1]}"
            break
        fi
    done
    
    echo "${distro:-$DEFAULT_DISTRO}"
}

# Function to check if DEBIAN_VERSION is already specified
has_debian_version() {
    for arg in "$@"; do
        if [[ "$arg" =~ DEBIAN_VERSION ]]; then
            return 0
        fi
    done
    return 1
}

# Main logic
main() {
    if [ $# -eq 0 ]; then
        echo "Error: No command specified" >&2
        echo "Usage: $0 <docker-command> [args...]" >&2
        echo "" >&2
        echo "Examples:" >&2
        echo "  $0 build --build-arg BASE_DEBIAN_DISTRO=trixie ." >&2
        echo "  BASE_DEBIAN_DISTRO=trixie $0 build ." >&2
        exit 1
    fi
    
    # Extract the distro from args or env
    local distro=$(extract_distro_from_args "$@")
    local version=$(get_debian_version "$distro")
    
    # Check if this is a bake command
    if [[ "$1" == "buildx" && "$2" == "bake" ]] || [[ "$1" == "bake" ]]; then
        echo "ℹ️  Note: 'docker buildx bake' has built-in mapping in docker-bake.hcl" >&2
        echo "ℹ️  Use: BASE_DEBIAN_DISTRO=$distro docker buildx bake [target]" >&2
        echo "" >&2
        echo "Executing your command as-is..." >&2
        exec docker "$@"
    fi
    
    # For 'docker build', inject DEBIAN_VERSION if not already specified
    if ! has_debian_version "$@"; then
        # Show what we're doing (can be silenced with QUIET=1)
        if [ "${QUIET:-0}" != "1" ]; then
            echo "🔧 DHI Build Wrapper: $distro → $version" >&2
        fi
        
        # Check if BASE_DEBIAN_DISTRO was passed as --build-arg
        local has_distro_arg=false
        for arg in "$@"; do
            if [[ "$arg" =~ BASE_DEBIAN_DISTRO ]]; then
                has_distro_arg=true
                break
            fi
        done
        
        # Find where to insert the arguments
        # We want to insert right after 'build' and before other args
        local new_args=()
        local inserted=false
        
        for arg in "$@"; do
            new_args+=("$arg")
            
            # After we see 'build', insert our args
            if [ "$arg" = "build" ] && [ "$inserted" = false ]; then
                # If BASE_DEBIAN_DISTRO wasn't passed as arg but is in env, add it
                if [ "$has_distro_arg" = false ] && [ -n "$BASE_DEBIAN_DISTRO" ]; then
                    new_args+=("--build-arg" "BASE_DEBIAN_DISTRO=$distro")
                fi
                new_args+=("--build-arg" "DEBIAN_VERSION=$version")
                inserted=true
            fi
        done
        
        exec docker "${new_args[@]}"
    else
        # DEBIAN_VERSION already specified, just pass through
        exec docker "$@"
    fi
}

main "$@"
