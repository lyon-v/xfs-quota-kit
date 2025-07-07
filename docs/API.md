# REST API 文档

XFS Quota Kit 提供完整的 REST API 接口，用于远程管理配额。

## 基础信息

- **Base URL**: `http://localhost:8080/api/v1`
- **Content-Type**: `application/json`
- **认证**: 暂不需要（未来版本将支持）

## 启动 API 服务器

```bash
# 默认端口 8080
xfs-quota-kit server

# 指定端口和主机
xfs-quota-kit server --host 0.0.0.0 --port 9090

# 使用配置文件
xfs-quota-kit server --config /etc/xfs-quota-kit/config.yaml
```

## API 端点

### 配额管理

#### 获取配额列表

```http
GET /api/v1/quotas?type={user|group|project}&path={filesystem_path}
```

**参数:**
- `type`: 配额类型（user, group, project）
- `path`: 文件系统路径

**响应:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1001,
      "type": "user",
      "path": "/mnt/xfs",
      "device": "/dev/sdb1",
      "block_used": 1048576,
      "block_soft": 1073741824,
      "block_hard": 2147483648,
      "inode_used": 1250,
      "inode_soft": 100000,
      "inode_hard": 200000,
      "last_updated": "2024-01-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

**示例:**
```bash
curl "http://localhost:8080/api/v1/quotas?type=user&path=/mnt/xfs"
```

#### 获取特定配额

```http
GET /api/v1/quotas/{type}/{id}?path={filesystem_path}
```

**参数:**
- `type`: 配额类型
- `id`: 用户/组/项目 ID
- `path`: 文件系统路径

**响应:**
```json
{
  "status": "success",
  "data": {
    "id": 1001,
    "type": "user",
    "path": "/mnt/xfs",
    "device": "/dev/sdb1",
    "block_used": 1048576,
    "block_soft": 1073741824,
    "block_hard": 2147483648,
    "inode_used": 1250,
    "inode_soft": 100000,
    "inode_hard": 200000,
    "block_usage_percent": 48.8,
    "inode_usage_percent": 1.25,
    "is_block_exceeded": false,
    "is_inode_exceeded": false,
    "last_updated": "2024-01-15T10:30:00Z"
  }
}
```

**示例:**
```bash
curl "http://localhost:8080/api/v1/quotas/user/1001?path=/mnt/xfs"
```

#### 创建/更新配额

```http
POST /api/v1/quotas
PUT /api/v1/quotas/{type}/{id}
```

**请求体:**
```json
{
  "type": "user",
  "id": 1001,
  "path": "/mnt/xfs",
  "limits": {
    "block_soft": 1073741824,
    "block_hard": 2147483648,
    "inode_soft": 100000,
    "inode_hard": 200000
  }
}
```

**响应:**
```json
{
  "status": "success",
  "message": "Quota set successfully",
  "data": {
    "id": 1001,
    "type": "user",
    "path": "/mnt/xfs"
  }
}
```

**示例:**
```bash
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

# 更新配额
curl -X PUT http://localhost:8080/api/v1/quotas/user/1001 \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/mnt/xfs",
    "limits": {
      "block_hard": 4294967296,
      "inode_hard": 400000
    }
  }'
```

#### 删除配额

```http
DELETE /api/v1/quotas/{type}/{id}?path={filesystem_path}
```

**响应:**
```json
{
  "status": "success",
  "message": "Quota removed successfully"
}
```

**示例:**
```bash
curl -X DELETE "http://localhost:8080/api/v1/quotas/user/1001?path=/mnt/xfs"
```

### 项目管理

#### 获取项目列表

```http
GET /api/v1/projects
```

**响应:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1000,
      "name": "myproject",
      "path": "/mnt/xfs/projects/myproject"
    }
  ],
  "count": 1
}
```

#### 创建项目

```http
POST /api/v1/projects
```

**请求体:**
```json
{
  "name": "myproject",
  "path": "/mnt/xfs/projects/myproject"
}
```

**响应:**
```json
{
  "status": "success",
  "message": "Project created successfully",
  "data": {
    "id": 1000,
    "name": "myproject",
    "path": "/mnt/xfs/projects/myproject"
  }
}
```

#### 删除项目

```http
DELETE /api/v1/projects/{name}
```

**响应:**
```json
{
  "status": "success",
  "message": "Project removed successfully"
}
```

### 报告

#### 生成配额报告

```http
GET /api/v1/reports?path={filesystem_path}&format={json|table}
```

**参数:**
- `path`: 文件系统路径
- `format`: 输出格式（可选，默认 json）

**响应:**
```json
{
  "status": "success",
  "data": {
    "filesystem": "/mnt/xfs",
    "total_quotas": 5,
    "over_quotas": 1,
    "warning_quotas": 2,
    "generated_at": "2024-01-15T10:30:00Z",
    "quotas": [
      {
        "id": 1001,
        "type": "user",
        "block_used": 1048576,
        "block_hard": 2147483648,
        "inode_used": 1250,
        "inode_hard": 200000,
        "block_usage_percent": 48.8,
        "inode_usage_percent": 1.25
      }
    ]
  }
}
```

#### 获取文件系统信息

```http
GET /api/v1/filesystem?path={filesystem_path}
```

**响应:**
```json
{
  "status": "success",
  "data": {
    "path": "/mnt/xfs",
    "type": "0x58465342",
    "is_xfs": true,
    "block_size": 4096,
    "total_size": "100 GB",
    "used_size": "45 GB",
    "free_size": "55 GB",
    "total_inodes": 26214400,
    "free_inodes": 26201150
  }
}
```

### 监控

#### 获取监控状态

```http
GET /api/v1/monitor/status
```

**响应:**
```json
{
  "status": "success",
  "data": {
    "enabled": true,
    "running": true,
    "interval": "5m",
    "threshold": 80,
    "last_check": "2024-01-15T10:25:00Z",
    "alerts": [
      {
        "type": "warning",
        "quota_type": "user",
        "quota_id": 1002,
        "path": "/mnt/xfs",
        "usage_percent": 85.2,
        "timestamp": "2024-01-15T10:20:00Z"
      }
    ]
  }
}
```

#### 启动监控

```http
POST /api/v1/monitor/start
```

**请求体:**
```json
{
  "path": "/mnt/xfs",
  "interval": "5m",
  "threshold": 80
}
```

#### 停止监控

```http
POST /api/v1/monitor/stop
```

## 批量操作

### 批量设置配额

```http
POST /api/v1/quotas/batch
```

**请求体:**
```json
{
  "type": "user",
  "path": "/mnt/xfs",
  "quotas": {
    "1001": {
      "block_hard": 2147483648,
      "inode_hard": 200000
    },
    "1002": {
      "block_hard": 1073741824,
      "inode_hard": 100000
    }
  }
}
```

**响应:**
```json
{
  "status": "success",
  "message": "Batch quotas set successfully",
  "data": {
    "successful": 2,
    "failed": 0,
    "errors": []
  }
}
```

## 错误处理

### 错误响应格式

```json
{
  "status": "error",
  "error": {
    "code": "QUOTA_NOT_FOUND",
    "message": "Quota not found for user 1001",
    "details": {
      "type": "user",
      "id": 1001,
      "path": "/mnt/xfs"
    }
  }
}
```

### 错误代码

| 代码 | HTTP状态 | 描述 |
|------|----------|------|
| `INVALID_REQUEST` | 400 | 请求参数无效 |
| `QUOTA_NOT_FOUND` | 404 | 配额不存在 |
| `FILESYSTEM_NOT_FOUND` | 404 | 文件系统不存在 |
| `PERMISSION_DENIED` | 403 | 权限不足 |
| `INTERNAL_ERROR` | 500 | 内部服务器错误 |
| `NOT_XFS_FILESYSTEM` | 422 | 不是XFS文件系统 |
| `QUOTA_OPERATION_FAILED` | 422 | 配额操作失败 |

## WebSocket 支持

### 实时监控

```javascript
// 连接 WebSocket
const ws = new WebSocket('ws://localhost:8080/api/v1/monitor/ws');

// 监听配额告警
ws.onmessage = function(event) {
  const alert = JSON.parse(event.data);
  console.log('Quota alert:', alert);
};

// 订阅特定路径的监控
ws.send(JSON.stringify({
  action: 'subscribe',
  path: '/mnt/xfs'
}));
```

### 消息格式

```json
{
  "type": "quota_alert",
  "data": {
    "quota_type": "user",
    "quota_id": 1001,
    "path": "/mnt/xfs",
    "usage_percent": 92.5,
    "threshold": 80,
    "timestamp": "2024-01-15T10:30:00Z"
  }
}
```

## SDK 示例

### JavaScript/Node.js

```javascript
const axios = require('axios');

class XFSQuotaKit {
  constructor(baseURL = 'http://localhost:8080/api/v1') {
    this.client = axios.create({ baseURL });
  }

  async getQuotas(type, path) {
    const response = await this.client.get('/quotas', {
      params: { type, path }
    });
    return response.data;
  }

  async setQuota(type, id, path, limits) {
    const response = await this.client.post('/quotas', {
      type, id, path, limits
    });
    return response.data;
  }

  async generateReport(path) {
    const response = await this.client.get('/reports', {
      params: { path }
    });
    return response.data;
  }
}

// 使用示例
const quota = new XFSQuotaKit();
const userQuotas = await quota.getQuotas('user', '/mnt/xfs');
```

### Python

```python
import requests

class XFSQuotaKit:
    def __init__(self, base_url='http://localhost:8080/api/v1'):
        self.base_url = base_url

    def get_quotas(self, quota_type, path):
        response = requests.get(f'{self.base_url}/quotas', 
                              params={'type': quota_type, 'path': path})
        return response.json()

    def set_quota(self, quota_type, quota_id, path, limits):
        data = {
            'type': quota_type,
            'id': quota_id,
            'path': path,
            'limits': limits
        }
        response = requests.post(f'{self.base_url}/quotas', json=data)
        return response.json()

# 使用示例
quota = XFSQuotaKit()
user_quotas = quota.get_quotas('user', '/mnt/xfs')
```

## 性能考虑

- API 响应通常在 100ms 内
- 批量操作可能需要更长时间
- 建议对频繁请求使用缓存
- WebSocket 连接适用于实时监控需求 