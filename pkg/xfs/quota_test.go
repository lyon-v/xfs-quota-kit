package xfs

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuotaType_String(t *testing.T) {
	tests := []struct {
		qType    QuotaType
		expected string
	}{
		{UserQuota, "user"},
		{GroupQuota, "group"},
		{ProjectQuota, "project"},
		{QuotaType(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.qType.String())
		})
	}
}

func TestQuotaInfo_IsBlockExceeded(t *testing.T) {
	tests := []struct {
		name     string
		quota    QuotaInfo
		expected bool
	}{
		{
			name: "not exceeded",
			quota: QuotaInfo{
				BlockUsed: 500,
				BlockHard: 1000,
			},
			expected: false,
		},
		{
			name: "exactly at limit",
			quota: QuotaInfo{
				BlockUsed: 1000,
				BlockHard: 1000,
			},
			expected: true,
		},
		{
			name: "exceeded",
			quota: QuotaInfo{
				BlockUsed: 1500,
				BlockHard: 1000,
			},
			expected: true,
		},
		{
			name: "no hard limit",
			quota: QuotaInfo{
				BlockUsed: 1500,
				BlockHard: 0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.quota.IsBlockExceeded())
		})
	}
}

func TestQuotaInfo_IsInodeExceeded(t *testing.T) {
	tests := []struct {
		name     string
		quota    QuotaInfo
		expected bool
	}{
		{
			name: "not exceeded",
			quota: QuotaInfo{
				InodeUsed: 500,
				InodeHard: 1000,
			},
			expected: false,
		},
		{
			name: "exactly at limit",
			quota: QuotaInfo{
				InodeUsed: 1000,
				InodeHard: 1000,
			},
			expected: true,
		},
		{
			name: "exceeded",
			quota: QuotaInfo{
				InodeUsed: 1500,
				InodeHard: 1000,
			},
			expected: true,
		},
		{
			name: "no hard limit",
			quota: QuotaInfo{
				InodeUsed: 1500,
				InodeHard: 0,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.quota.IsInodeExceeded())
		})
	}
}

func TestQuotaInfo_BlockUsagePercent(t *testing.T) {
	tests := []struct {
		name     string
		quota    QuotaInfo
		expected float64
	}{
		{
			name: "50% usage",
			quota: QuotaInfo{
				BlockUsed: 500,
				BlockHard: 1000,
			},
			expected: 50.0,
		},
		{
			name: "100% usage",
			quota: QuotaInfo{
				BlockUsed: 1000,
				BlockHard: 1000,
			},
			expected: 100.0,
		},
		{
			name: "no hard limit",
			quota: QuotaInfo{
				BlockUsed: 1000,
				BlockHard: 0,
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.quota.BlockUsagePercent())
		})
	}
}

func TestQuotaInfo_InodeUsagePercent(t *testing.T) {
	tests := []struct {
		name     string
		quota    QuotaInfo
		expected float64
	}{
		{
			name: "25% usage",
			quota: QuotaInfo{
				InodeUsed: 250,
				InodeHard: 1000,
			},
			expected: 25.0,
		},
		{
			name: "100% usage",
			quota: QuotaInfo{
				InodeUsed: 1000,
				InodeHard: 1000,
			},
			expected: 100.0,
		},
		{
			name: "no hard limit",
			quota: QuotaInfo{
				InodeUsed: 1000,
				InodeHard: 0,
			},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.quota.InodeUsagePercent())
		})
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		bytes    uint64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
		{1099511627776, "1.0 TB"},
		{1125899906842624, "1.0 PB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, FormatSize(tt.bytes))
		})
	}
}

func TestQuotaError_Error(t *testing.T) {
	err := &QuotaError{
		Op:   "get",
		Path: "/mnt/xfs",
		Err:  assert.AnError,
	}

	expected := "quota get /mnt/xfs: assert.AnError general error for testing"
	assert.Equal(t, expected, err.Error())
}

func TestQuotaError_Unwrap(t *testing.T) {
	originalErr := assert.AnError
	err := &QuotaError{
		Op:   "set",
		Path: "/mnt/xfs",
		Err:  originalErr,
	}

	assert.Equal(t, originalErr, err.Unwrap())
}

// Mock implementation for testing
type mockQuotaManager struct {
	quotas   map[string]*QuotaInfo
	projects map[string]*ProjectInfo
}

func newMockQuotaManager() *mockQuotaManager {
	return &mockQuotaManager{
		quotas:   make(map[string]*QuotaInfo),
		projects: make(map[string]*ProjectInfo),
	}
}

func (m *mockQuotaManager) GetQuota(quotaType QuotaType, id uint32, path string) (*QuotaInfo, error) {
	key := makeQuotaKey(quotaType, id, path)
	if quota, exists := m.quotas[key]; exists {
		return quota, nil
	}
	return &QuotaInfo{
		ID:          id,
		Type:        quotaType,
		Path:        path,
		LastUpdated: time.Now(),
	}, nil
}

func (m *mockQuotaManager) SetQuota(quotaType QuotaType, id uint32, path string, limits QuotaLimits) error {
	key := makeQuotaKey(quotaType, id, path)
	m.quotas[key] = &QuotaInfo{
		ID:          id,
		Type:        quotaType,
		Path:        path,
		BlockSoft:   limits.BlockSoft,
		BlockHard:   limits.BlockHard,
		InodeSoft:   limits.InodeSoft,
		InodeHard:   limits.InodeHard,
		LastUpdated: time.Now(),
	}
	return nil
}

func (m *mockQuotaManager) RemoveQuota(quotaType QuotaType, id uint32, path string) error {
	key := makeQuotaKey(quotaType, id, path)
	delete(m.quotas, key)
	return nil
}

func (m *mockQuotaManager) GetAllQuotas(quotaType QuotaType, path string) ([]QuotaInfo, error) {
	var quotas []QuotaInfo
	for _, quota := range m.quotas {
		if quota.Type == quotaType && quota.Path == path {
			quotas = append(quotas, *quota)
		}
	}
	return quotas, nil
}

func (m *mockQuotaManager) SetBatchQuotas(quotaType QuotaType, path string, quotas map[uint32]QuotaLimits) error {
	for id, limits := range quotas {
		if err := m.SetQuota(quotaType, id, path, limits); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockQuotaManager) CreateProject(name string, path string) (*ProjectInfo, error) {
	project := &ProjectInfo{
		ID:   uint32(len(m.projects) + 1000),
		Name: name,
		Path: path,
	}
	m.projects[name] = project
	return project, nil
}

func (m *mockQuotaManager) RemoveProject(name string) error {
	delete(m.projects, name)
	return nil
}

func (m *mockQuotaManager) GetProjects() ([]ProjectInfo, error) {
	var projects []ProjectInfo
	for _, project := range m.projects {
		projects = append(projects, *project)
	}
	return projects, nil
}

func (m *mockQuotaManager) GenerateReport(path string) (*QuotaReport, error) {
	return &QuotaReport{
		Filesystem:  path,
		GeneratedAt: time.Now(),
	}, nil
}

func (m *mockQuotaManager) CheckQuotaStatus(path string) error {
	return nil
}

func (m *mockQuotaManager) IsXFSFilesystem(path string) (bool, error) {
	return true, nil
}

func (m *mockQuotaManager) GetFilesystemInfo(path string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"type":         "0x58465342",
		"total_size":   "100 GB",
		"free_size":    "50 GB",
		"used_size":    "50 GB",
		"total_inodes": 1000000,
		"free_inodes":  500000,
	}, nil
}

func makeQuotaKey(quotaType QuotaType, id uint32, path string) string {
	return fmt.Sprintf("%s:%d:%s", quotaType.String(), id, path)
}

func TestMockQuotaManager(t *testing.T) {
	manager := newMockQuotaManager()

	// Test SetQuota and GetQuota
	limits := QuotaLimits{
		BlockSoft: 1024,
		BlockHard: 2048,
		InodeSoft: 1000,
		InodeHard: 2000,
	}

	err := manager.SetQuota(UserQuota, 1001, "/mnt/xfs", limits)
	require.NoError(t, err)

	quota, err := manager.GetQuota(UserQuota, 1001, "/mnt/xfs")
	require.NoError(t, err)
	assert.Equal(t, uint32(1001), quota.ID)
	assert.Equal(t, UserQuota, quota.Type)
	assert.Equal(t, "/mnt/xfs", quota.Path)
	assert.Equal(t, uint64(1024), quota.BlockSoft)
	assert.Equal(t, uint64(2048), quota.BlockHard)

	// Test GetAllQuotas
	quotas, err := manager.GetAllQuotas(UserQuota, "/mnt/xfs")
	require.NoError(t, err)
	assert.Len(t, quotas, 1)

	// Test RemoveQuota
	err = manager.RemoveQuota(UserQuota, 1001, "/mnt/xfs")
	require.NoError(t, err)

	quotas, err = manager.GetAllQuotas(UserQuota, "/mnt/xfs")
	require.NoError(t, err)
	assert.Len(t, quotas, 0)

	// Test CreateProject
	project, err := manager.CreateProject("test-project", "/mnt/xfs/projects/test")
	require.NoError(t, err)
	assert.Equal(t, "test-project", project.Name)
	assert.Equal(t, "/mnt/xfs/projects/test", project.Path)

	// Test GetProjects
	projects, err := manager.GetProjects()
	require.NoError(t, err)
	assert.Len(t, projects, 1)

	// Test RemoveProject
	err = manager.RemoveProject("test-project")
	require.NoError(t, err)

	projects, err = manager.GetProjects()
	require.NoError(t, err)
	assert.Len(t, projects, 0)
}
