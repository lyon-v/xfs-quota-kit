package xfs

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const (
	// quotactl 命令常量
	QCMD_CMD  = 0xff
	QCMD_TYPE = 0x0f00

	// XFS配额命令
	Q_XQUOTAON  = 0x800
	Q_XQUOTAOFF = 0x801
	Q_XGETQUOTA = 0x800
	Q_XSETQLIM  = 0x803
	Q_XGETQSTAT = 0x805

	// 配额类型
	USRQUOTA = 0
	GRPQUOTA = 1
	PRJQUOTA = 2

	// XFS magic number
	XFS_SUPER_MAGIC = 0x58465342
)

// DqBlk 配额块结构体（简化版本）
type DqBlk struct {
	BHardLimit uint64 // 块硬限制
	BSoftLimit uint64 // 块软限制
	CurSpace   uint64 // 当前使用空间
	IHardLimit uint64 // inode硬限制
	ISoftLimit uint64 // inode软限制
	CurInodes  uint64 // 当前使用inode数
	BTime      int64  // 块宽限时间
	ITime      int64  // inode宽限时间
}

// QuotaManager 配额管理器接口
type QuotaManager interface {
	// 基础操作
	GetQuota(quotaType QuotaType, id uint32, path string) (*QuotaInfo, error)
	SetQuota(quotaType QuotaType, id uint32, path string, limits QuotaLimits) error
	RemoveQuota(quotaType QuotaType, id uint32, path string) error

	// 批量操作
	GetAllQuotas(quotaType QuotaType, path string) ([]QuotaInfo, error)
	SetBatchQuotas(quotaType QuotaType, path string, quotas map[uint32]QuotaLimits) error

	// 项目配额特殊操作
	CreateProject(name string, path string) (*ProjectInfo, error)
	RemoveProject(name string) error
	GetProjects() ([]ProjectInfo, error)

	// 报告和监控
	GenerateReport(path string) (*QuotaReport, error)
	CheckQuotaStatus(path string) error

	// 文件系统操作
	IsXFSFilesystem(path string) (bool, error)
	GetFilesystemInfo(path string) (map[string]interface{}, error)
}

// quotaManager 配额管理器实现
type quotaManager struct{}

// NewQuotaManager 创建新的配额管理器
func NewQuotaManager() QuotaManager {
	return &quotaManager{}
}

// GetQuota 获取配额信息
func (q *quotaManager) GetQuota(quotaType QuotaType, id uint32, path string) (*QuotaInfo, error) {
	device, err := q.getDeviceFromPath(path)
	if err != nil {
		return nil, &QuotaError{Op: "get", Path: path, Err: err}
	}

	var dqblk DqBlk
	_ = makeQuotaCmd(Q_XGETQUOTA, quotaType) // 模拟系统调用

	// 模拟系统调用（实际实现需要使用真实的quotactl系统调用）
	// 这里为了演示，返回模拟数据
	quota := &QuotaInfo{
		ID:          id,
		Type:        quotaType,
		Path:        path,
		Device:      device,
		BlockUsed:   dqblk.CurSpace / 1024, // 转换为KB
		BlockSoft:   dqblk.BSoftLimit / 1024,
		BlockHard:   dqblk.BHardLimit / 1024,
		InodeUsed:   dqblk.CurInodes,
		InodeSoft:   dqblk.ISoftLimit,
		InodeHard:   dqblk.IHardLimit,
		LastUpdated: time.Now(),
	}

	return quota, nil
}

// SetQuota 设置配额限制
func (q *quotaManager) SetQuota(quotaType QuotaType, id uint32, path string, limits QuotaLimits) error {
	device, err := q.getDeviceFromPath(path)
	if err != nil {
		return &QuotaError{Op: "set", Path: path, Err: err}
	}

	var dqblk DqBlk
	dqblk.BSoftLimit = limits.BlockSoft * 1024 // 转换为字节
	dqblk.BHardLimit = limits.BlockHard * 1024
	dqblk.ISoftLimit = limits.InodeSoft
	dqblk.IHardLimit = limits.InodeHard

	_ = makeQuotaCmd(Q_XSETQLIM, quotaType) // 模拟系统调用

	// 模拟系统调用
	_ = device
	_ = dqblk

	return nil
}

// RemoveQuota 删除配额限制
func (q *quotaManager) RemoveQuota(quotaType QuotaType, id uint32, path string) error {
	// 通过设置所有限制为0来删除配额
	limits := QuotaLimits{
		BlockSoft: 0,
		BlockHard: 0,
		InodeSoft: 0,
		InodeHard: 0,
	}

	return q.SetQuota(quotaType, id, path, limits)
}

// GetAllQuotas 获取所有配额
func (q *quotaManager) GetAllQuotas(quotaType QuotaType, path string) ([]QuotaInfo, error) {
	// 这里应该遍历所有用户/组/项目ID，实际实现需要读取系统配置
	var quotas []QuotaInfo

	// 模拟返回一些配额数据
	for i := uint32(1000); i <= 1005; i++ {
		quota, err := q.GetQuota(quotaType, i, path)
		if err != nil {
			continue // 跳过不存在的配额
		}
		quotas = append(quotas, *quota)
	}

	return quotas, nil
}

// SetBatchQuotas 批量设置配额
func (q *quotaManager) SetBatchQuotas(quotaType QuotaType, path string, quotas map[uint32]QuotaLimits) error {
	var lastErr error
	for id, limits := range quotas {
		if err := q.SetQuota(quotaType, id, path, limits); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// CreateProject 创建项目配额
func (q *quotaManager) CreateProject(name string, path string) (*ProjectInfo, error) {
	// 检查项目是否已存在
	projects, err := q.GetProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return nil, fmt.Errorf("project %s already exists", name)
		}
	}

	// 分配新的项目ID（实际实现需要更智能的ID分配）
	projectID := uint32(1000 + len(projects))

	// 创建项目目录
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	// 写入项目配置文件（实际实现需要处理 /etc/projects 和 /etc/projid）
	project := &ProjectInfo{
		ID:   projectID,
		Name: name,
		Path: path,
	}

	return project, nil
}

// RemoveProject 删除项目配额
func (q *quotaManager) RemoveProject(name string) error {
	// 实际实现需要从 /etc/projects 和 /etc/projid 文件中删除项目配置
	return nil
}

// GetProjects 获取所有项目
func (q *quotaManager) GetProjects() ([]ProjectInfo, error) {
	// 实际实现需要读取 /etc/projects 和 /etc/projid 文件
	var projects []ProjectInfo

	// 模拟返回一些项目数据
	projects = append(projects, ProjectInfo{
		ID:   1000,
		Name: "example-project",
		Path: "/mnt/xfs/projects/example",
	})

	return projects, nil
}

// GenerateReport 生成配额报告
func (q *quotaManager) GenerateReport(path string) (*QuotaReport, error) {
	report := &QuotaReport{
		Filesystem:  path,
		GeneratedAt: time.Now(),
	}

	// 收集所有类型的配额信息
	allQuotas := []QuotaInfo{}

	for _, qType := range []QuotaType{UserQuota, GroupQuota, ProjectQuota} {
		quotas, err := q.GetAllQuotas(qType, path)
		if err != nil {
			continue
		}
		allQuotas = append(allQuotas, quotas...)
	}

	report.Quotas = allQuotas
	report.TotalQuotas = len(allQuotas)

	// 统计超限和警告的配额
	for _, quota := range allQuotas {
		if quota.IsBlockExceeded() || quota.IsInodeExceeded() {
			report.OverQuotas++
		} else if quota.BlockUsagePercent() > 80 || quota.InodeUsagePercent() > 80 {
			report.WarningQuotas++
		}
	}

	return report, nil
}

// CheckQuotaStatus 检查配额状态
func (q *quotaManager) CheckQuotaStatus(path string) error {
	isXFS, err := q.IsXFSFilesystem(path)
	if err != nil {
		return err
	}
	if !isXFS {
		return fmt.Errorf("path %s is not on an XFS filesystem", path)
	}
	return nil
}

// IsXFSFilesystem 检查是否为XFS文件系统
func (q *quotaManager) IsXFSFilesystem(path string) (bool, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return false, err
	}

	// 检查文件系统类型是否为XFS
	return stat.Type == XFS_SUPER_MAGIC, nil
}

// GetFilesystemInfo 获取文件系统信息
func (q *quotaManager) GetFilesystemInfo(path string) (map[string]interface{}, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return nil, err
	}

	info := map[string]interface{}{
		"type":         fmt.Sprintf("0x%X", stat.Type),
		"block_size":   stat.Bsize,
		"total_size":   FormatSize(uint64(stat.Blocks) * uint64(stat.Bsize)),
		"free_size":    FormatSize(uint64(stat.Bavail) * uint64(stat.Bsize)),
		"used_size":    FormatSize(uint64(stat.Blocks-stat.Bavail) * uint64(stat.Bsize)),
		"total_inodes": stat.Files,
		"free_inodes":  stat.Ffree,
	}

	return info, nil
}

// 辅助函数

// getDeviceFromPath 从路径获取设备名
func (q *quotaManager) getDeviceFromPath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// 读取 /proc/mounts 查找挂载点
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 简化实现：假设设备名为 /dev/sdb1
	// 实际实现需要解析 /proc/mounts 文件
	_ = absPath
	return "/dev/sdb1", nil
}

// makeQuotaCmd 构造配额命令
func makeQuotaCmd(cmd int, quotaType QuotaType) uint32 {
	var qtype int
	switch quotaType {
	case UserQuota:
		qtype = USRQUOTA
	case GroupQuota:
		qtype = GRPQUOTA
	case ProjectQuota:
		qtype = PRJQUOTA
	}

	return uint32((cmd << 8) | qtype)
}

// quotactl 系统调用的模拟实现
func quotactl(cmd uint32, device string, id uint32, data uintptr) error {
	// 这里是模拟实现，实际需要调用真实的 quotactl 系统调用
	// 在真实实现中，需要使用 syscall.Syscall6 或类似方法
	_ = cmd
	_ = device
	_ = id
	_ = data

	return nil
}
