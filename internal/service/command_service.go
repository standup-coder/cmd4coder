package service

import (
	"github.com/cmd4coder/cmd4coder/internal/data"
	"github.com/cmd4coder/cmd4coder/internal/model"
)

// CommandService 命令查询服务
type CommandService struct {
	loader *data.Loader
	index  *data.Index
	cache  *data.SearchCache
}

// NewCommandService 创建命令服务
func NewCommandService(dataDir string) (*CommandService, error) {
	loader := data.NewLoader(dataDir)
	index := data.NewIndex()
	cache := data.NewSearchCache(100) // 缓存最近100次搜索

	// 加载所有命令并构建索引
	commands, err := loader.LoadAllCommands()
	if err != nil {
		return nil, err
	}

	if err := index.BuildIndex(commands); err != nil {
		return nil, err
	}

	return &CommandService{
		loader: loader,
		index:  index,
		cache:  cache,
	}, nil
}

// GetCommand 根据名称获取命令
func (s *CommandService) GetCommand(name string) (*model.Command, error) {
	return s.index.GetByName(name)
}

// ListCommandsByCategory 根据分类列出命令
func (s *CommandService) ListCommandsByCategory(category string) []*model.Command {
	return s.index.GetByCategory(category)
}

// ListCommandsByPlatform 根据平台列出命令
func (s *CommandService) ListCommandsByPlatform(platform string) []*model.Command {
	return s.index.GetByPlatform(platform)
}

// SearchCommands 搜索命令
func (s *CommandService) SearchCommands(query string) []*model.Command {
	// 先检查缓存
	if commands, ok := s.cache.GetSearchResult(query); ok {
		return commands
	}

	// 执行搜索
	commands := s.index.Search(query)

	// 缓存结果
	s.cache.SetSearchResult(query, commands)

	return commands
}

// GetAllCategories 获取所有分类
func (s *CommandService) GetAllCategories() []string {
	return s.index.GetAllCategories()
}

// GetAllCommands 获取所有命令
func (s *CommandService) GetAllCommands() []*model.Command {
	return s.index.GetAllCommands()
}

// GetMetadata 获取元数据
func (s *CommandService) GetMetadata() *model.Metadata {
	return s.loader.GetMetadata()
}

// GetCommandCount 获取命令总数
func (s *CommandService) GetCommandCount() int {
	return len(s.index.GetAllCommands())
}

// GetCategoryCount 获取分类总数
func (s *CommandService) GetCategoryCount() int {
	return len(s.index.GetAllCategories())
}

// FilterCommandsByRisk 根据风险级别过滤命令
func (s *CommandService) FilterCommandsByRisk(riskLevel model.RiskLevel) []*model.Command {
	allCommands := s.index.GetAllCommands()
	var filtered []*model.Command

	for _, cmd := range allCommands {
		if cmd.GetHighestRisk() == riskLevel {
			filtered = append(filtered, cmd)
		}
	}

	return filtered
}

// GetHighRiskCommands 获取高风险命令
func (s *CommandService) GetHighRiskCommands() []*model.Command {
	allCommands := s.index.GetAllCommands()
	var highRisk []*model.Command

	for _, cmd := range allCommands {
		risk := cmd.GetHighestRisk()
		if risk == model.RiskLevelHigh || risk == model.RiskLevelCritical {
			highRisk = append(highRisk, cmd)
		}
	}

	return highRisk
}

// Reload 重新加载数据
func (s *CommandService) Reload() error {
	commands, err := s.loader.LoadAllCommands()
	if err != nil {
		return err
	}

	if err := s.index.BuildIndex(commands); err != nil {
		return err
	}

	// 清空缓存
	s.cache.Clear()

	return nil
}

// Count 获取命令总数（别名）
func (s *CommandService) Count() int {
	return s.GetCommandCount()
}

// GetCategories 获取所有分类（别名）
func (s *CommandService) GetCategories() []string {
	return s.GetAllCategories()
}

// Search 搜索命令（别名）
func (s *CommandService) Search(query string) []*model.Command {
	return s.SearchCommands(query)
}

// GetByCategory 根据分类获取命令（别名）
func (s *CommandService) GetByCategory(category string) []*model.Command {
	return s.ListCommandsByCategory(category)
}
