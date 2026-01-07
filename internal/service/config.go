package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cmd4coder/cmd4coder/internal/model"
)

// ConfigService 配置管理服务
type ConfigService struct {
	config       *model.Config
	userData     *model.UserData
	configPath   string
	userDataPath string
	mu           sync.RWMutex
}

// NewConfigService 创建配置服务
func NewConfigService() (*ConfigService, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户目录失败: %w", err)
	}

	configDir := filepath.Join(homeDir, ".cmd4coder")
	configPath := filepath.Join(configDir, "config.json")
	userDataPath := filepath.Join(configDir, "userdata.json")

	// 加载配置
	config, err := model.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	// 加载用户数据
	userData, err := model.LoadUserData(userDataPath)
	if err != nil {
		return nil, fmt.Errorf("加载用户数据失败: %w", err)
	}

	return &ConfigService{
		config:       config,
		userData:     userData,
		configPath:   configPath,
		userDataPath: userDataPath,
	}, nil
}

// GetConfig 获取配置
func (s *ConfigService) GetConfig() *model.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// UpdateConfig 更新配置
func (s *ConfigService) UpdateConfig(updater func(*model.Config)) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	updater(s.config)

	if err := s.config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	if err := s.config.Save(s.configPath); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	return nil
}

// SaveConfig 保存配置
func (s *ConfigService) SaveConfig() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if err := s.config.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	return s.config.Save(s.configPath)
}

// GetUserData 获取用户数据
func (s *ConfigService) GetUserData() *model.UserData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userData
}

// SaveUserData 保存用户数据
func (s *ConfigService) SaveUserData() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.userData.Save(s.userDataPath)
}

// AddFavorite 添加收藏
func (s *ConfigService) AddFavorite(commandName, category, note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userData.AddFavorite(commandName, category, note)
	return s.userData.Save(s.userDataPath)
}

// RemoveFavorite 删除收藏
func (s *ConfigService) RemoveFavorite(commandName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userData.RemoveFavorite(commandName)
	return s.userData.Save(s.userDataPath)
}

// IsFavorite 检查是否已收藏
func (s *ConfigService) IsFavorite(commandName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.userData.IsFavorite(commandName)
}

// GetFavorites 获取所有收藏
func (s *ConfigService) GetFavorites() []model.Favorite {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.userData.Favorites
}

// AddHistory 添加历史记录
func (s *ConfigService) AddHistory(commandName, category string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userData.AddHistory(commandName, category)
	return s.userData.Save(s.userDataPath)
}

// GetRecentHistory 获取最近的历史记录
func (s *ConfigService) GetRecentHistory(limit int) []model.HistoryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.userData.GetRecentHistory(limit)
}

// ClearHistory 清空历史记录
func (s *ConfigService) ClearHistory() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userData.ClearHistory()
	return s.userData.Save(s.userDataPath)
}

// SetLanguage 设置语言
func (s *ConfigService) SetLanguage(lang string) error {
	return s.UpdateConfig(func(c *model.Config) {
		c.Language = lang
	})
}

// SetTheme 设置主题
func (s *ConfigService) SetTheme(theme string) error {
	return s.UpdateConfig(func(c *model.Config) {
		c.Theme = theme
	})
}

// SetPageSize 设置每页显示数量
func (s *ConfigService) SetPageSize(size int) error {
	return s.UpdateConfig(func(c *model.Config) {
		c.PageSize = size
	})
}

// EnableTUI 启用/禁用TUI
func (s *ConfigService) EnableTUI(enabled bool) error {
	return s.UpdateConfig(func(c *model.Config) {
		c.TUI.Enabled = enabled
	})
}

// SetDefaultExportFormat 设置默认导出格式
func (s *ConfigService) SetDefaultExportFormat(format string) error {
	return s.UpdateConfig(func(c *model.Config) {
		c.Export.DefaultFormat = format
	})
}
