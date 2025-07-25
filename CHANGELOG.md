# Changelog

本文档记录了XFS Quota Kit项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
- 初始项目架构和核心功能
- 用户、组和项目配额管理命令
- REST API 服务器
- 配额使用监控和报告生成
- 配置文件支持（YAML格式）
- Docker 支持和容器化部署
- 命令行自动补全功能
- 批量配额操作支持
- 多种输出格式（表格、JSON）

### 文档
- 完整的README文档
- API文档和使用指南  
- 快速入门指南
- 贡献指南
- Docker Compose示例
- 生产环境配置示例

### 技术实现
- 基于Cobra的CLI框架
- 原生Go实现的XFS配额管理
- 使用Viper进行配置管理
- 结构化日志记录
- 完整的测试覆盖
- Makefile构建系统

### 部署支持
- 安装脚本
- systemd服务配置
- Docker镜像构建
- 多平台编译支持

## [1.0.0] - TBD

### 计划功能
- 完整的XFS配额系统调用实现
- Web界面
- 高级监控和告警
- 数据库后端支持
- 集群部署支持
- 认证和授权系统

---

## 版本说明

### 版本号规则
- `MAJOR.MINOR.PATCH` (例如: 1.0.0)
- **MAJOR**: 不兼容的API更改
- **MINOR**: 向后兼容的功能添加  
- **PATCH**: 向后兼容的问题修复

### 发布类型
- `[未发布]` - 开发中的功能
- `[版本号] - 日期` - 正式发布版本

### 变更类型
- `新增` - 新功能
- `变更` - 现有功能的变更
- `废弃` - 即将移除的功能
- `移除` - 已移除的功能
- `修复` - 问题修复
- `安全` - 安全相关的修复 