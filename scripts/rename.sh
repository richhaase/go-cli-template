#!/usr/bin/env bash
# Compatible with bash 3.2 (stock macOS) — avoid bash 4+ features.
set -euo pipefail

# Colors (disabled if not a TTY)
if [[ -t 1 ]]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED='' GREEN='' YELLOW='' BLUE='' NC=''
fi

usage() {
    cat <<EOF
Usage: $(basename "$0") [OPTIONS]

Bootstrap a new Go CLI project from the template.

Options:
    -o, --owner OWNER       GitHub owner/organization (e.g., "myuser")
    -r, --repo REPO         Repository name (e.g., "my-awesome-cli")
    -b, --binary NAME       Binary name (defaults to repo name)
    -d, --description DESC  Short description for the CLI
    -c, --copyright NAME    Copyright holder for LICENSE (defaults to owner)
    -y, --yes               Skip confirmation prompt
    -h, --help              Show this help message

Examples:
    $(basename "$0") -o myuser -r my-cli
    $(basename "$0") --owner myorg --repo toolname --binary tool
    $(basename "$0") -o myuser -r my-cli -d "A tool for doing things"

EOF
    exit 0
}

error() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

info() {
    echo -e "${BLUE}→${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

warn() {
    echo -e "${YELLOW}!${NC} $1"
}

# In-place sed that works on both BSD (macOS) and GNU sed,
# without leaving backup files behind.
sed_inplace() {
    if [[ "$(uname)" == "Darwin" ]]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

# Parse arguments
OWNER=""
REPO=""
BINARY=""
DESCRIPTION=""
COPYRIGHT=""
SKIP_CONFIRM=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -o|--owner)
            OWNER="$2"
            shift 2
            ;;
        -r|--repo)
            REPO="$2"
            shift 2
            ;;
        -b|--binary)
            BINARY="$2"
            shift 2
            ;;
        -d|--description)
            DESCRIPTION="$2"
            shift 2
            ;;
        -c|--copyright)
            COPYRIGHT="$2"
            shift 2
            ;;
        -y|--yes)
            SKIP_CONFIRM=true
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            error "Unknown option: $1\nRun '$(basename "$0") --help' for usage."
            ;;
    esac
done

# Interactive prompts for missing values
if [[ -z "$OWNER" ]]; then
    read -rp "GitHub owner/organization: " OWNER
fi

if [[ -z "$REPO" ]]; then
    read -rp "Repository name: " REPO
fi

if [[ -z "$BINARY" ]]; then
    BINARY="$REPO"
    read -rp "Binary name [$BINARY]: " input
    BINARY="${input:-$BINARY}"
fi

if [[ -z "$DESCRIPTION" ]]; then
    read -rp "Short description (optional): " DESCRIPTION
fi

if [[ -z "$COPYRIGHT" && "$SKIP_CONFIRM" != true ]]; then
    read -rp "Copyright holder for LICENSE [$OWNER]: " COPYRIGHT
fi
COPYRIGHT="${COPYRIGHT:-$OWNER}"

# Validate inputs
[[ -z "$OWNER" ]] && error "Owner is required"
[[ -z "$REPO" ]] && error "Repository name is required"
[[ -z "$BINARY" ]] && error "Binary name is required"

# Validate binary name (must be valid Go identifier-ish)
if [[ ! "$BINARY" =~ ^[a-zA-Z][a-zA-Z0-9_-]*$ ]]; then
    error "Binary name must start with a letter and contain only letters, numbers, hyphens, and underscores"
fi

# Env-var prefix: uppercased binary name with hyphens as underscores
# (e.g. my-cli → MY_CLI, so MYCLI_DEBUG becomes MY_CLI_DEBUG)
BINARY_UPPER=$(printf '%s' "$BINARY" | tr '[:lower:]-' '[:upper:]_')

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Show summary and confirm
echo ""
echo -e "${BLUE}Project Configuration:${NC}"
echo "  Module:      github.com/$OWNER/$REPO"
echo "  Binary:      $BINARY"
echo "  Env prefix:  ${BINARY_UPPER}_"
echo "  Description: ${DESCRIPTION:-"(none)"}"
echo "  Copyright:   $COPYRIGHT"
echo "  Location:    $PROJECT_ROOT"
echo ""

if [[ "$SKIP_CONFIRM" != true ]]; then
    read -rp "Proceed with renaming? [Y/n] " confirm
    # bash 3.2-safe lowercase comparison (no ${var,,})
    confirm=$(printf '%s' "$confirm" | tr '[:upper:]' '[:lower:]')
    if [[ "$confirm" == "n" || "$confirm" == "no" ]]; then
        echo "Aborted."
        exit 0
    fi
fi

echo ""

# Perform replacements
cd "$PROJECT_ROOT"

# Files to process (exclude .git, binaries, LICENSE — handled separately)
FILES=$(find . -type f \
    -not -path "./.git/*" \
    -not -path "./bin/*" \
    -not -path "./dist/*" \
    -not -path "./scripts/rename.sh" \
    -not -name "LICENSE" \
    -not -name "*.sum" \
    -not -name "*.exe")

info "Replacing placeholders in files..."

for file in $FILES; do
    if [[ -f "$file" ]]; then
        # Replace OWNER/REPO with actual values
        sed_inplace "s|github\.com/OWNER/REPO|github.com/$OWNER/$REPO|g" "$file"
        sed_inplace "s|OWNER/REPO|$OWNER/$REPO|g" "$file"
        sed_inplace "s|OWNER|$OWNER|g" "$file"

        # Replace env-var prefix (MYCLI_*) before the lowercase binary name
        sed_inplace "s|MYCLI|$BINARY_UPPER|g" "$file"

        # Replace mycli with binary name
        sed_inplace "s|mycli|$BINARY|g" "$file"

        # Replace description if provided
        if [[ -n "$DESCRIPTION" ]]; then
            sed_inplace "s|A brief description of your CLI|$DESCRIPTION|g" "$file"
            sed_inplace "s|Brief description of what this CLI does\.|$DESCRIPTION|g" "$file"
        fi
    fi
done

success "Replaced placeholders in source files"

# Strip the template-usage section from the generated project's README
if grep -q '^## Using This Template' README.md; then
    info "Removing template-usage section from README.md..."
    sed_inplace '/^## Using This Template$/,/^---$/d' README.md
    success "Cleaned README.md"
fi

# Update LICENSE copyright line (year + holder)
if [[ -f LICENSE ]]; then
    info "Updating LICENSE copyright..."
    YEAR=$(date +%Y)
    sed_inplace "s|Copyright (c) [0-9][0-9][0-9][0-9] OWNER|Copyright (c) $YEAR $COPYRIGHT|" LICENSE
    success "Updated LICENSE"
fi

# Rename cmd/mycli directory
if [[ -d "cmd/mycli" && "$BINARY" != "mycli" ]]; then
    info "Renaming cmd/mycli → cmd/$BINARY..."
    mv "cmd/mycli" "cmd/$BINARY"
    success "Renamed command directory"
fi

# Update go.mod
info "Updating go.mod..."
sed_inplace "s|^module .*|module github.com/$OWNER/$REPO|" go.mod
success "Updated go.mod"

# Run go mod tidy
info "Running go mod tidy..."
go mod tidy
success "Dependencies updated"

# Verify build
info "Verifying build..."
if go build -o "bin/$BINARY" "./cmd/$BINARY"; then
    success "Build successful"
    rm -rf bin
else
    warn "Build failed - you may need to fix some issues manually"
fi

# Remove this script
info "Cleaning up..."
rm -f "$SCRIPT_DIR/rename.sh"
rmdir "$SCRIPT_DIR" 2>/dev/null || true
success "Removed bootstrap script"

echo ""
echo -e "${GREEN}Done!${NC} Your project is ready."
echo ""
echo "Next steps:"
echo "  1. git add -A && git commit -m 'Initial commit'"
echo "  2. gh repo create $OWNER/$REPO --public --source=. --push"
echo "  3. just deps-update   # keep dependencies current"
echo "  4. Start coding!"
echo ""
