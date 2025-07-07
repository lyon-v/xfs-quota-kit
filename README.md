# XFS Quota Kit

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-green.svg)](https://golang.org/)
[![Platform](https://img.shields.io/badge/Platform-Linux-blue.svg)](https://www.kernel.org/)

一个先进的 XFS 文件系统配额管理工具包，提供全面的用户、组和项目配额管理功能。

## 功能特性

### 🎯 核心功能
- **多类型配额管理**: 支持用户、组和项目配额
- **批量操作**: 支持批量设置和管理配额
- **实时监控**: 配额使用情况实时监控和告警
- **报告生成**: 全面的配额使用报告
- **REST API**: 提供完整的 REST API 接口

### 🛠 技术特性
- **原生支持**: 直接使用系统调用，无需外部依赖
- **配置文件**: 支持 YAML 配置文件
- **日志系统**: 结构化日志记录
- **命令行友好**: 丰富的 CLI 命令和选项
- **Docker 支持**: 提供 Docker 镜像

### 📊 监控与报告
- 配额使用率监控
- 超限告警
- 自动报告生成
- 多种输出格式（表格、JSON）

## 安装

### 从源码编译

```bash
git clone https://github.com/yourusername/xfs-quota-kit.git
cd xfs-quota-kit
make build
sudo make install
```

### 使用 Docker

```bash
docker run --rm -it --privileged \
  -v /mnt:/mnt \
  xfs-quota-kit:latest --help
```

## 快速开始

### 1. 检查文件系统

```bash
# 检查是否为 XFS 文件系统
xfs-quota-kit report filesystem /mnt/xfs
```

### 2. 设置用户配额

```bash
# 为用户 ID 1001 设置配额
xfs-quota-kit quota set /mnt/xfs \
  --type user \
  --id 1001 \
  --block-hard 2GB \
  --inode-hard 100000
```

### 3. 查看配额信息

```bash
# 获取用户配额信息
xfs-quota-kit quota get /mnt/xfs --type user --id 1001

# 列出所有用户配额
xfs-quota-kit quota list /mnt/xfs --type user
```

### 4. 创建项目配额

```bash
# 创建项目配额
xfs-quota-kit project create myproject /mnt/xfs/projects/myproject

# 为项目设置配额
xfs-quota-kit quota set /mnt/xfs \
  --type project \
  --id 1000 \
  --block-hard 10GB
```

### 5. 生成报告

```bash
# 生成配额使用报告
xfs-quota-kit report generate /mnt/xfs

# 以 JSON 格式输出
xfs-quota-kit report generate /mnt/xfs --format json
```

### 6. 启动监控

```bash
# 启动实时监控
xfs-quota-kit monitor start /mnt/xfs --interval 5m --threshold 80
```

## 配置

### 配置文件

创建配置文件 `/etc/xfs-quota-kit/config.yaml`:

```yaml
# 基本配置
xfs:
  default_path: "/mnt/xfs"
  projects_file: "/etc/projects"
  projid_file: "/etc/projid"

# 默认限制
default_limits:
  user_block_soft: "1GB"
  user_block_hard: "2GB"
  user_inode_soft: 100000
  user_inode_hard: 200000

# 监控配置
monitor:
  enabled: true
  interval: "5m"
  alert_threshold: 80
```

### 环境变量

```bash
export XFS_QUOTA_DEFAULT_PATH="/mnt/xfs"
export XFS_QUOTA_LOGGING_LEVEL="info"
export XFS_QUOTA_SERVER_PORT="8080"
```

## API 服务器

启动 REST API 服务器：

```bash
xfs-quota-kit server --port 8080
```

### API 端点

- `GET /api/v1/quotas` - 获取配额列表
- `POST /api/v1/quotas` - 创建配额
- `GET /api/v1/quotas/{id}` - 获取特定配额
- `PUT /api/v1/quotas/{id}` - 更新配额
- `DELETE /api/v1/quotas/{id}` - 删除配额
- `GET /api/v1/reports` - 生成报告

## 命令参考

### 配额管理

```bash
# 获取配额
xfs-quota-kit quota get [path] --type [user|group|project] --id [ID]

# 设置配额
xfs-quota-kit quota set [path] --type [user|group|project] --id [ID] \
  --block-soft [SIZE] --block-hard [SIZE] \
  --inode-soft [COUNT] --inode-hard [COUNT]

# 移除配额
xfs-quota-kit quota remove [path] --type [user|group|project] --id [ID]

# 列出配额
xfs-quota-kit quota list [path] --type [user|group|project] --format [table|json]
```

### 项目管理

```bash
# 创建项目
xfs-quota-kit project create [name] [path]

# 删除项目
xfs-quota-kit project remove [name]

# 列出项目
xfs-quota-kit project list
```

### 报告和监控

```bash
# 生成报告
xfs-quota-kit report generate [path] --format [table|json]

# 文件系统信息
xfs-quota-kit report filesystem [path]

# 开始监控
xfs-quota-kit monitor start [path] --interval [DURATION] --threshold [PERCENT]

# 监控状态
xfs-quota-kit monitor status
```

### 服务器

```bash
# 启动 API 服务器
xfs-quota-kit server --host [HOST] --port [PORT]
```

## 开发

### 构建

```bash
# 安装依赖
make deps

# 代码检查
make lint

# 运行测试
make test

# 构建
make build

# 构建所有平台
make build-all
```

### 开发模式

```bash
# 开发模式运行
make dev

# 生成文档
make docs
```

## 系统要求

- **操作系统**: Linux (内核 3.0+)
- **文件系统**: XFS
- **权限**: 需要 root 权限操作配额
- **Go 版本**: 1.21 或更高（仅构建时需要）

## XFS 配额设置

确保您的 XFS 文件系统启用了配额支持：

```bash
# 挂载时启用用户和组配额
mount -o uquota,gquota /dev/sdb1 /mnt/xfs

# 挂载时启用项目配额
mount -o pquota /dev/sdb1 /mnt/xfs

# 在 /etc/fstab 中配置
/dev/sdb1 /mnt/xfs xfs defaults,pquota 0 0
```

## 贡献

我们欢迎各种形式的贡献！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解详细信息。

## 许可证

本项目采用 Apache 2.0 许可证。详情请查看 [LICENSE](LICENSE) 文件。

## 致谢

本项目参考了以下开源项目：
- [silenceper/xfsquota](https://github.com/silenceper/xfsquota) - 基础的 XFS 配额工具
- [anexia-it/fsquota](https://github.com/anexia-it/fsquota) - 原生 Go 配额库
- [cirocosta/xfsvol](https://github.com/cirocosta/xfsvol) - XFS 项目配额实现

## 支持

如果您遇到问题或有建议，请：
1. 查看 [文档](docs/)
2. 搜索现有的 [Issues](https://github.com/yourusername/xfs-quota-kit/issues)
3. 创建新的 Issue
