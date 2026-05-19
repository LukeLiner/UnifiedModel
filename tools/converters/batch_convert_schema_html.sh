#!/bin/bash

# 批量转换 schema YAML 文件为 HTML（递归目录结构）
# 用法: 
#   ./batch_convert_schema_html.sh <源目录> <目标HTML目录>              # 生成混合版本
#   ./batch_convert_schema_html.sh <源目录> <目标HTML目录> <语言模式>    # 生成指定语言版本
#   支持的语言模式: mixed, cn, en

set -e  # 遇到错误时退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查参数
if [ $# -lt 2 ] || [ $# -gt 3 ]; then
    print_error "用法: $0 <源目录> <目标HTML目录> [语言模式]"
    print_error "语言模式: mixed(默认), cn, en"
    exit 1
fi

SRC_DIR="$1"
DST_DIR="$2"
LANGUAGE_MODE="${3:-mixed}"  # 默认为mixed模式

# 验证语言参数
case "$LANGUAGE_MODE" in
    mixed|cn|en)
        ;;
    *)
        print_error "无效的语言参数: $LANGUAGE_MODE"
        print_error "支持的参数: mixed, cn, en"
        exit 1
        ;;
esac

print_info "源目录: $SRC_DIR"
print_info "目标HTML目录: $DST_DIR"
print_info "语言模式: $LANGUAGE_MODE"

# CD 到 sh文件所在目录的上一级
cd "$(dirname "$0")/.."

echo "当前工作目录: $(pwd)"

# 检查源目录
if [ ! -d "$SRC_DIR" ]; then
    print_error "源目录不存在: $SRC_DIR"
    exit 1
fi

# 创建目标目录（如果不存在）
if [ ! -d "$DST_DIR" ]; then
    print_info "创建目标目录: $DST_DIR"
    mkdir -p "$DST_DIR"
fi

# 检查转换脚本
if [ ! -f "scripts/converters/schema_to_table_html_dynamic.py" ]; then
    print_error "转换脚本不存在: scripts/converters/schema_to_table_html_dynamic.py"
    exit 1
fi

# 统计变量
file_count=0
success_count=0
error_count=0

# 遍历所有 .yaml 文件
while IFS= read -r -d '' yaml_file; do
    # 计算相对路径
    rel_path="${yaml_file#$SRC_DIR/}"
    # 替换扩展名为 .html
    html_rel_path="${rel_path%.yaml}.html"
    # 目标 HTML 路径
    html_file="$DST_DIR/$html_rel_path"
    # 目标 HTML 所在目录
    html_dir="$(dirname "$html_file")"
    # 创建目标目录
    mkdir -p "$html_dir"

    print_info "转换: $rel_path -> $html_rel_path ($LANGUAGE_MODE 模式)"

    # 注意：schema_to_table_html_dynamic.py 目前还不支持语言参数
    # 未来需要添加 -l "$LANGUAGE_MODE" 参数支持
    if python scripts/converters/schema_to_table_html_dynamic.py "$yaml_file" -o "$html_file"; then
        print_success "✓ 转换成功: $html_rel_path"
        ((++success_count))
    else
        print_error "✗ 转换失败: $rel_path"
        ((++error_count))
    fi
    ((++file_count))
    echo

done < <(find "$SRC_DIR" -type f -name "*.yaml" -print0)

echo "生成 index.html..."
# 注意：generate_common_schema_index.py 目前还不支持语言参数
# 未来需要添加语言参数支持
python scripts/converters/generate_common_schema_index.py --doc-dir "$DST_DIR"
echo "生成 index.html 完成"

# 输出统计结果
echo "=================================="
print_info "转换完成统计:"
echo "  总文件数: $file_count"
echo "  成功转换: $success_count"
echo "  转换失败: $error_count"

if [ $error_count -eq 0 ]; then
    print_success "🎉 所有文件转换成功！($LANGUAGE_MODE 模式)"
    print_info "HTML文件已生成到: $DST_DIR/"
    exit 0
else
    print_warning "⚠️  有 $error_count 个文件转换失败"
    exit 1
fi 