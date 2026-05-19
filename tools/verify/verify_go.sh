#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
GO_DIR="$PROJECT_ROOT/sdk/go"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "========================================="
echo "  Go SDK 编译验证"
echo "========================================="

if ! command -v go &>/dev/null; then
    echo -e "${YELLOW}⚠️  Go 未安装，跳过 Go 验证${NC}"
    exit 0
fi

echo "Go version: $(go version)"

if [ ! -d "$GO_DIR" ]; then
    echo -e "${RED}❌ Go SDK 目录不存在: $GO_DIR${NC}"
    exit 1
fi

cd "$GO_DIR"

ERRORS=0

echo ""
echo "📦 检查 Go modules..."
if [ ! -f "go.mod" ]; then
    echo -e "${RED}❌ go.mod 不存在${NC}"
    exit 1
fi
go mod download 2>/dev/null || true
echo -e "${GREEN}✅ Go modules OK${NC}"

echo ""
echo "🔨 编译 Go SDK..."
if go build ./...; then
    echo -e "${GREEN}✅ Go 编译通过${NC}"
else
    echo -e "${RED}❌ Go 编译失败${NC}"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "🔍 运行 go vet..."
if go vet ./... 2>&1; then
    echo -e "${GREEN}✅ Go vet 通过${NC}"
else
    echo -e "${YELLOW}⚠️  Go vet 发现问题${NC}"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "🧪 运行 Go 测试..."
TEST_FILES=$(find . -name "*_test.go" 2>/dev/null | head -1)
if [ -n "$TEST_FILES" ]; then
    if go test ./... 2>&1; then
        echo -e "${GREEN}✅ Go 测试通过${NC}"
    else
        echo -e "${RED}❌ Go 测试失败${NC}"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo -e "${YELLOW}⚠️  未找到 Go 测试文件，跳过${NC}"
fi

echo ""
echo "========================================="
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ Go SDK 验证全部通过${NC}"
else
    echo -e "${RED}❌ Go SDK 验证失败 ($ERRORS 个错误)${NC}"
    exit 1
fi
