# 快速入门指南

本指南将帮助您快速开始使用 XFS Quota Kit。

## 前提条件

### 系统要求
- Linux 操作系统（内核 3.0+）
- XFS 文件系统
- Root 权限（用于配额操作）

### XFS 配额支持检查

首先，确保您的 XFS 文件系统启用了配额支持：

```bash
# 检查当前挂载选项
mount | grep xfs

# 如果没有配额选项，需要重新挂载
umount /mnt/xfs
mount -o pquota /dev/sdb1 /mnt/xfs

# 或者在 /etc/fstab 中永久配置
echo "/dev/sdb1 /mnt/xfs xfs defaults,pquota 0 0" >> /etc/fstab
```

## 安装

### 方法 1: 从源码构建

```bash
# 克隆仓库
git clone https://github.com/yourusername/xfs-quota-kit.git
cd xfs-quota-kit

# 构建
make build

# 安装
sudo make install
```

### 方法 2: 使用安装脚本

```bash
# 下载并运行安装脚本
curl -fsSL https://raw.githubusercontent.com/yourusername/xfs-quota-kit/main/scripts/install.sh | sudo bash
```

### 方法 3: Docker

```bash
# 拉取镜像
docker pull xfs-quota-kit:latest

# 运行（需要特权模式访问设备）
docker run --rm -it --privileged \
  -v /mnt:/mnt \
  xfs-quota-kit:latest --help
```

## 基本使用

### 1. 验证安装

```bash
# 检查版本
xfs-quota-kit version

# 显示帮助
xfs-quota-kit --help
```

### 2. 检查文件系统

```bash
# 检查路径是否为 XFS 文件系统
xfs-quota-kit report filesystem /mnt/xfs
```

### 3. 用户配额管理

#### 设置用户配额

```bash
# 为用户 1001 设置配额限制
xfs-quota-kit quota set /mnt/xfs \
  --type user \
  --id 1001 \
  --block-soft 1GB \
  --block-hard 2GB \
  --inode-soft 100000 \
  --inode-hard 200000
```

#### 查看用户配额

```bash
# 查看特定用户配额
xfs-quota-kit quota get /mnt/xfs --type user --id 1001

# 列出所有用户配额
xfs-quota-kit quota list /mnt/xfs --type user
```

#### 移除用户配额

```bash
# 移除用户配额
xfs-quota-kit quota remove /mnt/xfs --type user --id 1001
```

### 4. 组配额管理

```bash
# 设置组配额
xfs-quota-kit quota set /mnt/xfs \
  --type group \
  --id 100 \
  --block-hard 10GB \
  --inode-hard 1000000

# 查看组配额
xfs-quota-kit quota get /mnt/xfs --type group --id 100

# 列出所有组配额
xfs-quota-kit quota list /mnt/xfs --type group
```

### 5. 项目配额管理

#### 创建项目

```bash
# 创建新项目
xfs-quota-kit project create myproject /mnt/xfs/projects/myproject

# 为项目设置配额
xfs-quota-kit quota set /mnt/xfs \
  --type project \
  --id 1000 \
  --block-hard 5GB \
  --inode-hard 500000
```

#### 管理项目

```bash
# 列出所有项目
xfs-quota-kit project list

# 删除项目
xfs-quota-kit project remove myproject
```

### 6. 生成报告

```bash
# 生成配额使用报告
xfs-quota-kit report generate /mnt/xfs

# 以 JSON 格式输出
xfs-quota-kit report generate /mnt/xfs --format json

# 保存到文件
xfs-quota-kit report generate /mnt/xfs --format json > quota-report.json
```

### 7. 实时监控

```bash
# 启动监控（在后台运行）
xfs-quota-kit monitor start /mnt/xfs \
  --interval 5m \
  --threshold 80 &

# 查看监控状态
xfs-quota-kit monitor status
```

## 配置文件

### 创建配置文件

```bash
# 创建配置目录
sudo mkdir -p /etc/xfs-quota-kit

# 创建基本配置文件
sudo tee /etc/xfs-quota-kit/config.yaml << EOF
# XFS Quota Kit 配置
xfs:
  default_path: "/mnt/xfs"
  projects_file: "/etc/projects"
  projid_file: "/etc/projid"
  
  default_limits:
    user_block_soft: "1GB"
    user_block_hard: "2GB"
    user_inode_soft: 100000
    user_inode_hard: 200000

monitor:
  enabled: true
  interval: "5m"
  alert_threshold: 80
  report_path: "/var/log/xfs-quota-kit/reports"

logging:
  level: "info"
  format: "json"
  output: "file"
  file: "/var/log/xfs-quota-kit/app.log"
EOF
```

### 使用配置文件

```bash
# 使用指定配置文件
xfs-quota-kit --config /etc/xfs-quota-kit/config.yaml quota list /mnt/xfs --type user
```

## 环境变量

可以使用环境变量覆盖配置：

```bash
# 设置默认路径
export XFS_QUOTA_DEFAULT_PATH="/mnt/xfs"

# 设置日志级别
export XFS_QUOTA_LOGGING_LEVEL="debug"

# 设置服务器端口
export XFS_QUOTA_SERVER_PORT="9090"
```

## 启动 API 服务器

```bash
# 启动 REST API 服务器
xfs-quota-kit server --port 8080

# 后台运行
nohup xfs-quota-kit server --port 8080 > /var/log/xfs-quota-kit/server.log 2>&1 &
```

### API 使用示例

```bash
# 获取配额列表
curl http://localhost:8080/api/v1/quotas

# 创建配额
curl -X POST http://localhost:8080/api/v1/quotas \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user",
    "id": 1001,
    "path": "/mnt/xfs",
    "limits": {
      "block_hard": 2147483648,
      "inode_hard": 200000
    }
  }'
```

## 常见问题

### 1. 权限错误

```bash
# 确保以 root 用户运行
sudo xfs-quota-kit quota list /mnt/xfs --type user
```

### 2. 文件系统不支持配额

```bash
# 检查挂载选项
mount | grep xfs

# 重新挂载启用配额
sudo umount /mnt/xfs
sudo mount -o pquota /dev/sdb1 /mnt/xfs
```

### 3. 项目配额设置失败

```bash
# 确保项目文件存在
sudo touch /etc/projects /etc/projid

# 检查项目目录权限
sudo chown root:root /etc/projects /etc/projid
sudo chmod 644 /etc/projects /etc/projid
```

## 下一步

- 查看 [配置指南](CONFIG.md) 了解详细配置选项
- 阅读 [API 文档](API.md) 了解 REST API 使用
- 参考 [故障排除](TROUBLESHOOTING.md) 解决常见问题 