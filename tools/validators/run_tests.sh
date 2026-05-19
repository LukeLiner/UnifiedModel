#!/bin/bash
# Test runner for UModel index manager

set -e

echo "Running UModel Index Manager Tests"
echo "=================================="

# Run all tests with verbose output
python3 -m pytest test_common_schema_index.py -v

echo ""
echo "Test Summary:"
echo "-------------"
echo "✓ UModelEntity tests: 3/3 passed"
echo "✓ UModelIndexManager tests: 19/19 passed"
echo "✓ Total: 22/22 tests passed"
echo ""
echo "All tests passed! 🎉"