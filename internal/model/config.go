package model

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config 应用程序配置
type Config struct {
	// 应用设置
	Language     string `json:"language"`      // 语言: zh/en
	Theme        string `json:"theme"`         // 主题: default/dark/light
	Editor       string `json:"editor"`        // 默认编辑器
	PageSize     int    `json:"page_size"`     // 每页显示数量
	
	// TUI设置
	TUI TUIConfig `json:"tui"`
	
	// 搜索设置
	Search SearchConfig `json:"search"`
	
	// 导出设置
	Export ExportConfig `json:"export"`
	
	// 用户数据路径
	UserDataPath string `json:"user_data_path"`
	
	// 版本
	Version string `json:"version"`
}

// TUIConfig TUI界面配置
type TUIConfig struct {
	Enabled         bool   `json:"enabled"`           // 是否启用TUI
	ColorScheme     string `json:"color_scheme"`      // 颜色方案
	ShowHelp        bool   `json:"show_help"`         // 显示帮助信息
	ShowLineNumbers bool   `json:"show_line_numbers"` // 显示行号
	WrapText        bool   `json:"wrap_text"`         // 文本换行
}

// SearchConfig 搜索配置
type SearchConfig struct {
	CaseSensitive bool `json:"case_sensitive"` // 大小写敏感
	MaxResults    int  `json:"max_results"`    // 最大结果数
	UseCache      bool `json:"use_cache"`      // 使用缓存
	CacheSize     int  `json:"cache_size"`     // 缓存大小
}

// ExportConfig 导出配置
type ExportConfig struct {
	DefaultFormat string `json:"default_format"` // 默认格式: markdown/json
	OutputDir     string `json:"output_dir"`     // 输出目录
	IncludeDate   bool   `json:"include_date"`   // 包含日期
}

// UserData 用户数据
type UserData struct {
	// 收藏的命令
	Favorites []Favorite `json:"favorites"`
	
	// 历史记录
	History []HistoryEntry `json:"history"`
	
	// 最后更新时间
	LastUpdated time.Time `json:"last_updated"`
}

// Favorite 收藏的命令
type Favorite struct {
	CommandName string    `json:"command_name"` // 命令名称
	Category    string    `json:"category"`     // 分类
	Note        string    `json:"note"`         // 备注
	AddedAt     time.Time `json:"added_at"`     // 添加时间
}

// HistoryEntry 历史记录条目
type HistoryEntry struct {
	CommandName string    `json:"command_name"` // 命令名称
	Category    string    `json:"category"`     // 分类
	AccessedAt  time.Time `json:"accessed_at"`  // 访问时间
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	userDataPath := filepath.Join(homeDir, ".cmd4coder")
	
	return &Config{
		Language:     "zh",
		Theme:        "default",
		Editor:       "vim",
		PageSize:     20,
		UserDataPath: userDataPath,
		Version:      "1.0.0",
		
		TUI: TUIConfig{
			Enabled:         true,
			ColorScheme:     "default",
			ShowHelp:        true,
			ShowLineNumbers: false,
			WrapText:        true,
		},
		
		Search: SearchConfig{
			CaseSensitive: false,
			MaxResults:    50,
			UseCache:      true,
			CacheSize:     100,
		},
		
		Export: ExportConfig{
			DefaultFormat: "markdown",
			OutputDir:     ".",
			IncludeDate:   true,
		},
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	// 如果文件不存在，返回默认配置
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	
	return &config, nil
}

// Save 保存配置到文件
func (c *Config) Save(path string) error {
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	
	return nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.PageSize <= 0 {
		return fmt.Errorf("page_size 必须大于0")
	}
	
	if c.Search.MaxResults <= 0 {
		return fmt.Errorf("search.max_results 必须大于0")
	}
	
	if c.Search.CacheSize <= 0 {
		return fmt.Errorf("search.cache_size 必须大于0")
	}
	
	if c.Language != "zh" && c.Language != "en" {
		return fmt.Errorf("language 必须是 zh 或 en")
	}
	
	validFormats := map[string]bool{"markdown": true, "json": true}
	if !validFormats[c.Export.DefaultFormat] {
		return fmt.Errorf("export.default_format 必须是 markdown 或 json")
	}
	
	return nil
}

// NewUserData 创建新的用户数据
func NewUserData() *UserData {
	return &UserData{
		Favorites:   make([]Favorite, 0),
		History:     make([]HistoryEntry, 0),
		LastUpdated: time.Now(),
	}
}

// LoadUserData 从文件加载用户数据
func LoadUserData(path string) (*UserData, error) {
	// 如果文件不存在，返回新的用户数据
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return NewUserData(), nil
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取用户数据文件失败: %w", err)
	}
	
	var userData UserData
	if err := json.Unmarshal(data, &userData); err != nil {
		return nil, fmt.Errorf("解析用户数据文件失败: %w", err)
	}
	
	return &userData, nil
}

// Save 保存用户数据到文件
func (u *UserData) Save(path string) error {
	// 更新时间戳
	u.LastUpdated = time.Now()
	
	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建用户数据目录失败: %w", err)
	}
	
	data, err := json.MarshalIndent(u, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化用户数据失败: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入用户数据文件失败: %w", err)
	}
	
	return nil
}

// AddFavorite 添加收藏
func (u *UserData) AddFavorite(commandName, category, note string) {
	// 检查是否已存在
	for _, fav := range u.Favorites {
		if fav.CommandName == commandName {
			return
		}
	}
	
	u.Favorites = append(u.Favorites, Favorite{
		CommandName: commandName,
		Category:    category,
		Note:        note,
		AddedAt:     time.Now(),
	})
}

// RemoveFavorite 删除收藏
func (u *UserData) RemoveFavorite(commandName string) {
	for i, fav := range u.Favorites {
		if fav.CommandName == commandName {
			u.Favorites = append(u.Favorites[:i], u.Favorites[i+1:]...)
			return
		}
	}
}

// IsFavorite 检查是否已收藏
func (u *UserData) IsFavorite(commandName string) bool {
	for _, fav := range u.Favorites {
		if fav.CommandName == commandName {
			return true
		}
	}
	return false
}

// AddHistory 添加历史记录
func (u *UserData) AddHistory(commandName, category string) {
	// 移除旧的相同记录
	for i := 0; i < len(u.History); i++ {
		if u.History[i].CommandName == commandName {
			u.History = append(u.History[:i], u.History[i+1:]...)
			i--
		}
	}
	
	// 添加到开头
	u.History = append([]HistoryEntry{{
		CommandName: commandName,
		Category:    category,
		AccessedAt:  time.Now(),
	}}, u.History...)
	
	// 限制历史记录数量
	maxHistory := 100
	if len(u.History) > maxHistory {
		u.History = u.History[:maxHistory]
	}
}

// GetRecentHistory 获取最近的历史记录
func (u *UserData) GetRecentHistory(limit int) []HistoryEntry {
	if limit > len(u.History) {
		limit = len(u.History)
	}
	return u.History[:limit]
}

// ClearHistory 清空历史记录
func (u *UserData) ClearHistory() {
	u.History = make([]HistoryEntry, 0)
}
