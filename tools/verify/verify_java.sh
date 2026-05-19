#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
JAVA_DIR="$PROJECT_ROOT/generated/java"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "========================================="
echo "  Java SDK 编译验证"
echo "========================================="

if ! command -v mvn &>/dev/null; then
    echo -e "${YELLOW}⚠️  Maven 未安装，跳过 Java 验证${NC}"
    exit 0
fi

if ! command -v java &>/dev/null; then
    echo -e "${YELLOW}⚠️  Java 未安装，跳过 Java 验证${NC}"
    exit 0
fi

echo "Maven: $(mvn --version 2>&1 | head -1)"
echo "Java: $(java -version 2>&1 | head -1)"

if [ ! -f "$JAVA_DIR/pom.xml" ]; then
    echo -e "${RED}❌ Java SDK 的 pom.xml 不存在: $JAVA_DIR/pom.xml${NC}"
    exit 1
fi

cd "$JAVA_DIR"

ERRORS=0

echo ""
echo "🔨 编译 Java SDK..."
if mvn compile -q -B 2>&1; then
    echo -e "${GREEN}✅ Java 编译通过${NC}"
else
    echo -e "${RED}❌ Java 编译失败${NC}"
    ERRORS=$((ERRORS + 1))
fi

echo ""
echo "🧪 运行 Java 测试..."
TEST_DIR="$JAVA_DIR/src/test"
if [ -d "$TEST_DIR" ]; then
    if mvn test -q -B 2>&1; then
        echo -e "${GREEN}✅ Java 测试通过${NC}"
    else
        echo -e "${YELLOW}⚠️  Java 测试失败（可能缺少 JUnit 依赖）${NC}"
        ERRORS=$((ERRORS + 1))
    fi
else
    echo -e "${YELLOW}⚠️  未找到 Java 测试目录，跳过${NC}"
fi

echo ""
echo "========================================="
if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ Java SDK 验证全部通过${NC}"
else
    echo -e "${RED}❌ Java SDK 验证失败 ($ERRORS 个错误)${NC}"
    exit 1
fi
