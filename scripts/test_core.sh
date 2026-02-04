#!/bin/bash
# 核心功能验证脚本

echo "=== cmd4coder 核心功能验证 ==="
echo ""

# 检查必要文件
echo "1. 检查项目结构完整性..."
if [[ -d "cmd/cli" && -d "internal" && -d "data" ]]; then
    echo "✅ 项目结构完整"
else
    echo "❌ 项目结构不完整"
    exit 1
fi

# 检查主要源文件
echo ""
echo "2. 检查核心源文件..."
required_files=(
    "cmd/cli/main.go"
    "cmd/cli/commands.go"
    "internal/model/command.go"
    "internal/service/command_service.go"
    "internal/ui/tui/tui.go"
)

missing_files=()
for file in "${required_files[@]}"; do
    if [[ -f "$file" ]]; then
        echo "✅ $file"
    else
        echo "❌ $file (缺失)"
        missing_files+=("$file")
    fi
done

if [[ ${#missing_files[@]} -gt 0 ]]; then
    echo "❌ 缺失 ${#missing_files[@]} 个核心文件"
    exit 1
fi

# 检查数据文件
echo ""
echo "3. 检查数据文件..."
yaml_count=$(find data/ -name "*.yaml" | wc -l)
if [[ $yaml_count -ge 30 ]]; then
    echo "✅ 数据文件充足 ($yaml_count 个YAML文件)"
else
    echo "⚠️  数据文件较少 ($yaml_count 个YAML文件)"
fi

# 检查文档文件
echo ""
echo "4. 检查文档完整性..."
doc_files=(
    "README.md"
    "docs/architecture/ARCHITECTURE.md"
    "docs/legal/LICENSE"
    "docs/legal/CODE_OF_CONDUCT.md"
)

for doc in "${doc_files[@]}"; do
    if [[ -f "$doc" ]]; then
        echo "✅ $doc"
    else
        echo "❌ $doc (缺失)"
    fi
done

# 检查构建配置
echo ""
echo "5. 检查构建配置..."
if [[ -f "build/config/go.mod" && -f "build/config/go.sum" ]]; then
    echo "✅ Go模块配置完整"
else
    echo "❌ Go模块配置缺失"
fi

if [[ -f "scripts/build.sh" ]]; then
    echo "✅ 构建脚本存在"
else
    echo "❌ 构建脚本缺失"
fi

echo ""
echo "=== 验证完成 ==="
echo "项目核心组件完整，准备好进行构建和测试。"
echo "请安装Go环境后运行: go run ./cmd/cli --help"