package xfs

import (
	"fmt"
	"time"
)

// QuotaType 配额类型
type QuotaType uint8

const (
	UserQuota    QuotaType = iota + 1 // 用户配额
	GroupQuota                        // 组配额
	ProjectQuota                      // 项目配额
)

func (q QuotaType) String() string {
	switch q {
	case UserQuota:
		return "user"
	case GroupQuota:
		return "group"
	case ProjectQuota:
		return "project"
	default:
		return "unknown"
	}
}

// QuotaInfo 配额信息结构
type QuotaInfo struct {
	ID          uint32    `json:"id"`           // 用户ID/组ID/项目ID
	Type        QuotaType `json:"type"`         // 配额类型
	Path        string    `json:"path"`         // 路径
	Device      string    `json:"device"`       // 设备
	BlockUsed   uint64    `json:"block_used"`   // 已使用块数 (KB)
	BlockSoft   uint64    `json:"block_soft"`   // 块软限制 (KB)
	BlockHard   uint64    `json:"block_hard"`   // 块硬限制 (KB)
	InodeUsed   uint64    `json:"inode_used"`   // 已使用inode数
	InodeSoft   uint64    `json:"inode_soft"`   // inode软限制
	InodeHard   uint64    `json:"inode_hard"`   // inode硬限制
	LastUpdated time.Time `json:"last_updated"` // 最后更新时间
}

// IsBlockExceeded 检查块使用是否超限
func (q *QuotaInfo) IsBlockExceeded() bool {
	return q.BlockHard > 0 && q.BlockUsed >= q.BlockHard
}

// IsInodeExceeded 检查inode使用是否超限
func (q *QuotaInfo) IsInodeExceeded() bool {
	return q.InodeHard > 0 && q.InodeUsed >= q.InodeHard
}

// BlockUsagePercent 获取块使用百分比
func (q *QuotaInfo) BlockUsagePercent() float64 {
	if q.BlockHard == 0 {
		return 0.0
	}
	return float64(q.BlockUsed) / float64(q.BlockHard) * 100.0
}

// InodeUsagePercent 获取inode使用百分比
func (q *QuotaInfo) InodeUsagePercent() float64 {
	if q.InodeHard == 0 {
		return 0.0
	}
	return float64(q.InodeUsed) / float64(q.InodeHard) * 100.0
}

// QuotaLimits 配额限制结构
type QuotaLimits struct {
	BlockSoft uint64 `json:"block_soft"` // 块软限制 (KB)
	BlockHard uint64 `json:"block_hard"` // 块硬限制 (KB)
	InodeSoft uint64 `json:"inode_soft"` // inode软限制
	InodeHard uint64 `json:"inode_hard"` // inode硬限制
}

// ProjectInfo 项目信息结构
type ProjectInfo struct {
	ID   uint32 `json:"id"`   // 项目ID
	Name string `json:"name"` // 项目名称
	Path string `json:"path"` // 项目路径
}

// QuotaReport 配额报告结构
type QuotaReport struct {
	Filesystem    string      `json:"filesystem"`     // 文件系统路径
	TotalQuotas   int         `json:"total_quotas"`   // 总配额数
	OverQuotas    int         `json:"over_quotas"`    // 超限配额数
	WarningQuotas int         `json:"warning_quotas"` // 警告配额数
	GeneratedAt   time.Time   `json:"generated_at"`   // 生成时间
	Quotas        []QuotaInfo `json:"quotas"`         // 配额详情
}

// QuotaError 配额操作错误
type QuotaError struct {
	Op   string // 操作类型
	Path string // 文件路径
	Err  error  // 底层错误
}

func (e *QuotaError) Error() string {
	return fmt.Sprintf("quota %s %s: %v", e.Op, e.Path, e.Err)
}

func (e *QuotaError) Unwrap() error {
	return e.Err
}

// FormatSize 格式化字节数为可读格式
func FormatSize(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}

	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
