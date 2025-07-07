# 贡献指南

感谢您对 XFS Quota Kit 项目的关注！我们欢迎各种形式的贡献。

## 开发环境设置

### 前提条件

- Go 1.21 或更高版本
- Linux 操作系统（用于XFS配额功能）
- Git
- Make

### 克隆仓库

```bash
git clone https://github.com/yourusername/xfs-quota-kit.git
cd xfs-quota-kit
```

### 安装依赖

```bash
make deps
```

### 构建和测试

```bash
# 构建
make build

# 运行测试
make test

# 代码检查
make lint

# 格式化代码
make format
```

## 贡献流程

### 1. 创建 Issue

在开始编写代码之前，请先创建一个 Issue 来讨论您想要贡献的功能或修复的问题。

### 2. Fork 和分支

```bash
# Fork 仓库到您的账户
# 然后克隆您的 fork

git clone https://github.com/yourusername/xfs-quota-kit.git
cd xfs-quota-kit

# 添加上游仓库
git remote add upstream https://github.com/originalowner/xfs-quota-kit.git

# 创建功能分支
git checkout -b feature/your-feature-name
```

### 3. 开发

- 遵循现有的代码风格
- 添加适当的测试
- 更新文档（如果需要）
- 确保所有测试通过

### 4. 提交

```bash
# 提交您的更改
git add .
git commit -m "Add: your feature description"

# 推送到您的分支
git push origin feature/your-feature-name
```

### 5. 创建 Pull Request

在 GitHub 上创建 Pull Request，包含：

- 清晰的标题和描述
- 相关 Issue 的链接
- 测试结果截图（如果适用）

## 代码规范

### Go 代码风格

- 使用 `gofmt` 格式化代码
- 遵循 Go 官方风格指南
- 使用有意义的变量和函数名
- 添加适当的注释

### 提交消息格式

使用以下格式：

```
<type>: <subject>

<body>

<footer>
```

类型包括：
- `feat`: 新功能
- `fix`: 修复问题
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 代码重构
- `test`: 添加测试
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat: add project quota management commands

Add support for creating, listing, and removing XFS project quotas
through the CLI interface.

Closes #123
```

## 测试

### 单元测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./pkg/xfs/

# 运行测试并生成覆盖率报告
make test-coverage
```

### 集成测试

由于 XFS 配额功能需要特殊的文件系统环境，集成测试需要：

1. XFS 文件系统
2. Root 权限
3. 启用配额的挂载点

```bash
# 设置测试环境
sudo make test-setup

# 运行集成测试
sudo make test-integration
```

## 文档

### 代码文档

- 所有公开的函数和类型都应该有文档注释
- 使用 `godoc` 格式
- 包含使用示例

### 用户文档

如果您的贡献涉及用户可见的功能，请更新：

- `README.md`
- `docs/` 目录下的相关文档
- 命令行帮助文本

## 发布流程

### 版本号

我们使用语义化版本控制 (SemVer)：

- `MAJOR.MINOR.PATCH`
- MAJOR：不兼容的 API 更改
- MINOR：向后兼容的功能添加
- PATCH：向后兼容的问题修复

### 发布步骤

1. 更新版本号
2. 更新 CHANGELOG.md
3. 创建 Git 标签
4. 构建发布二进制文件
5. 创建 GitHub Release

## 社区行为准则

### 我们的承诺

为了营造一个开放和友好的环境，我们承诺：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最好的事情
- 对其他社区成员表示同理心

### 不可接受的行为

- 使用性化的语言或图像
- 挑衅、侮辱或贬损的评论
- 公开或私下的骚扰
- 发布他人的私人信息
- 其他不道德或不专业的行为

## 获得帮助

如果您需要帮助或有疑问：

1. 查看现有的 [Issues](https://github.com/yourusername/xfs-quota-kit/issues)
2. 查看 [文档](docs/)
3. 创建新的 Issue 寻求帮助
4. 联系维护者

## 许可证

通过贡献代码，您同意您的贡献将在与项目相同的 Apache 2.0 许可证下获得许可。

感谢您的贡献！🎉 