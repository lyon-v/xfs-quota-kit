#!/bin/bash

# XFS Quota Kit 基本使用示例

set -e

# 配置变量
XFS_PATH="/mnt/xfs"
BINARY="./build/xfs-quota-kit"

echo "XFS Quota Kit 基本使用示例"
echo "============================"

# 检查二进制文件是否存在
if [ ! -f "$BINARY" ]; then
    echo "错误：未找到二进制文件 $BINARY"
    echo "请先运行 'make build' 构建项目"
    exit 1
fi

# 检查是否以 root 权限运行
if [ "$EUID" -ne 0 ]; then
    echo "警告：配额操作通常需要 root 权限"
    echo "某些命令可能会失败"
fi

echo
echo "1. 检查版本信息"
echo "=================="
$BINARY version

echo
echo "2. 显示帮助信息"
echo "=================="
$BINARY --help

echo
echo "3. 检查文件系统信息"
echo "===================="
echo "检查路径：$XFS_PATH"
if [ -d "$XFS_PATH" ]; then
    $BINARY report filesystem "$XFS_PATH" || echo "无法获取文件系统信息（可能需要 root 权限）"
else
    echo "路径 $XFS_PATH 不存在，跳过文件系统检查"
fi

echo
echo "4. 用户配额操作示例"
echo "===================="

# 设置用户配额（需要 root 权限）
echo "设置用户 1001 的配额..."
$BINARY quota set "$XFS_PATH" --type user --id 1001 \
    --block-hard 2GB --inode-hard 100000 || echo "设置配额失败（可能需要 root 权限或 XFS 文件系统）"

# 查看用户配额
echo "查看用户 1001 的配额..."
$BINARY quota get "$XFS_PATH" --type user --id 1001 || echo "获取配额失败"

# 列出所有用户配额
echo "列出所有用户配额..."
$BINARY quota list "$XFS_PATH" --type user || echo "列出配额失败"

echo
echo "5. 项目配额操作示例"
echo "===================="

# 创建项目
echo "创建项目 'example-project'..."
$BINARY project create example-project "$XFS_PATH/projects/example" || echo "创建项目失败"

# 列出项目
echo "列出所有项目..."
$BINARY project list || echo "列出项目失败"

echo
echo "6. 报告生成示例"
echo "=================="

# 生成配额报告
echo "生成配额使用报告..."
$BINARY report generate "$XFS_PATH" || echo "生成报告失败"

# 以 JSON 格式生成报告
echo "生成 JSON 格式报告..."
$BINARY report generate "$XFS_PATH" --format json || echo "生成 JSON 报告失败"

echo
echo "7. 自动补全示例"
echo "=================="
echo "生成 bash 自动补全脚本："
$BINARY completion bash > /tmp/xfs-quota-kit-completion.bash
echo "自动补全脚本已保存到 /tmp/xfs-quota-kit-completion.bash"
echo "可以使用以下命令加载："
echo "source /tmp/xfs-quota-kit-completion.bash"

echo
echo "示例完成！"
echo "=========="
echo "注意："
echo "- 大多数配额操作需要 root 权限"
echo "- 需要在 XFS 文件系统上运行"
echo "- 文件系统需要启用配额支持"
echo ""
echo "更多信息请查看："
echo "- README.md"
echo "- docs/GETTING_STARTED.md"
echo "- docs/API.md" 