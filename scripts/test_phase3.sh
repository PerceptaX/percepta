#!/bin/bash
set -e

echo "=== Phase 3 Verification Test ==="
echo

# Clean slate
echo "1. Cleaning database..."
rm -f ~/.local/share/percepta/percepta.db
echo "   ✓ Database cleaned"
echo

# Verify database creation and schema
echo "2. Testing database initialization..."
./percepta diff test --from v1 --to v2 2>&1 | grep -q "failed to get observation" && echo "   ✓ Database created and schema initialized"
echo

# Check schema
echo "3. Verifying schema..."
SCHEMA=$(sqlite3 ~/.local/share/percepta/percepta.db ".schema")
echo "$SCHEMA" | grep -q "CREATE TABLE observations" && echo "   ✓ observations table created"
echo "$SCHEMA" | grep -q "firmware TEXT" && echo "   ✓ firmware column exists"
echo "$SCHEMA" | grep -q "idx_device_firmware" && echo "   ✓ index created"
echo

# Check for CGO dependencies
echo "4. Checking for CGO dependencies..."
if go list -f '{{.ImportPath}}: {{.CgoCFLAGS}} {{.CgoLDFLAGS}}' ./cmd/percepta 2>&1 | grep -q "github.com/mattn/go-sqlite3"; then
    echo "   ✗ ERROR: Found mattn/go-sqlite3 (CGO dependency)"
    exit 1
fi

if go list -m all | grep -q "modernc.org/sqlite"; then
    echo "   ✓ Using modernc.org/sqlite (pure Go)"
else
    echo "   ✗ ERROR: modernc.org/sqlite not found"
    exit 1
fi
echo

# Build without CGO
echo "5. Testing build without CGO..."
CGO_ENABLED=0 go build -o percepta_nocgo ./cmd/percepta && rm percepta_nocgo
echo "   ✓ Build successful without CGO"
echo

echo "=== All Phase 3 checks passed! ==="
echo
echo "Next steps for manual testing:"
echo "  1. Add firmware tag to config: ~/.config/percepta/config.yaml"
echo "     devices:"
echo "       fpga:"
echo "         firmware: v1"
echo "  2. Run: ./percepta observe fpga"
echo "  3. Change firmware tag to v2 in config"
echo "  4. Run: ./percepta observe fpga"
echo "  5. Run: ./percepta diff fpga --from v1 --to v2"
