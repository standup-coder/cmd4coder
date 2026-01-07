package data

import (
	"strings"
	"sync"

	"github.com/cmd4coder/cmd4coder/internal/model"
)

// Index 命令索引
type Index struct {
	// 命令名称 -> 命令对象
	nameIndex map[string]*model.Command

	// 分类 -> 命令列表
	categoryIndex map[string][]*model.Command

	// 关键词 -> 命令列表（倒排索引）
	keywordIndex map[string][]*model.Command

	// 平台 -> 命令列表
	platformIndex map[string][]*model.Command

	mu sync.RWMutex
}

// NewIndex 创建索引
func NewIndex() *Index {
	return &Index{
		nameIndex:     make(map[string]*model.Command),
		categoryIndex: make(map[string][]*model.Command),
		keywordIndex:  make(map[string][]*model.Command),
		platformIndex: make(map[string][]*model.Command),
	}
}

// BuildIndex 构建索引
func (idx *Index) BuildIndex(commands []*model.Command) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// 清空现有索引
	idx.nameIndex = make(map[string]*model.Command)
	idx.categoryIndex = make(map[string][]*model.Command)
	idx.keywordIndex = make(map[string][]*model.Command)
	idx.platformIndex = make(map[string][]*model.Command)

	for _, cmd := range commands {
		// 检查命令名称是否重复
		if _, exists := idx.nameIndex[cmd.Name]; exists {
			return model.ErrDuplicateCommand{Name: cmd.Name}
		}

		// 构建名称索引
		idx.nameIndex[cmd.Name] = cmd

		// 构建分类索引
		idx.categoryIndex[cmd.Category] = append(idx.categoryIndex[cmd.Category], cmd)

		// 构建平台索引
		for _, platform := range cmd.Platforms {
			idx.platformIndex[platform] = append(idx.platformIndex[platform], cmd)
		}

		// 构建关键词索引（从命令名称、描述中提取）
		idx.buildKeywordIndex(cmd)
	}

	return nil
}

// buildKeywordIndex 为单个命令构建关键词索引
func (idx *Index) buildKeywordIndex(cmd *model.Command) {
	keywords := make(map[string]bool)

	// 从命令名称提取
	for _, word := range tokenize(cmd.Name) {
		keywords[word] = true
	}

	// 从描述提取
	for _, word := range tokenize(cmd.Description) {
		keywords[word] = true
	}

	// 从分类提取
	for _, word := range tokenize(cmd.Category) {
		keywords[word] = true
	}

	// 添加到倒排索引
	for keyword := range keywords {
		idx.keywordIndex[keyword] = append(idx.keywordIndex[keyword], cmd)
	}
}

// tokenize 分词（简单实现）
func tokenize(text string) []string {
	text = strings.ToLower(text)
	// 按空格、斜杠、下划线等分割
	text = strings.ReplaceAll(text, "/", " ")
	text = strings.ReplaceAll(text, "_", " ")
	text = strings.ReplaceAll(text, "-", " ")

	words := strings.Fields(text)
	var result []string
	for _, word := range words {
		if len(word) > 1 { // 过滤单字符
			result = append(result, word)
		}
	}
	return result
}

// GetByName 根据名称获取命令
func (idx *Index) GetByName(name string) (*model.Command, error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	cmd, ok := idx.nameIndex[name]
	if !ok {
		return nil, model.ErrCommandNotFound{Name: name}
	}
	return cmd, nil
}

// GetByCategory 根据分类获取命令列表
func (idx *Index) GetByCategory(category string) []*model.Command {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return idx.categoryIndex[category]
}

// GetByPlatform 根据平台获取命令列表
func (idx *Index) GetByPlatform(platform string) []*model.Command {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	return idx.platformIndex[platform]
}

// Search 搜索命令
func (idx *Index) Search(query string) []*model.Command {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	query = strings.ToLower(query)
	resultMap := make(map[string]*searchResult)

	// 1. 精确匹配命令名称
	for name, cmd := range idx.nameIndex {
		if strings.ToLower(name) == query {
			resultMap[cmd.Name] = &searchResult{
				command:  cmd,
				priority: 100, // 最高优先级
			}
		}
	}

	// 2. 前缀匹配命令名称
	for name, cmd := range idx.nameIndex {
		if strings.HasPrefix(strings.ToLower(name), query) {
			if _, exists := resultMap[cmd.Name]; !exists {
				resultMap[cmd.Name] = &searchResult{
					command:  cmd,
					priority: 80,
				}
			}
		}
	}

	// 3. 命令名称包含查询词
	for name, cmd := range idx.nameIndex {
		if strings.Contains(strings.ToLower(name), query) {
			if _, exists := resultMap[cmd.Name]; !exists {
				resultMap[cmd.Name] = &searchResult{
					command:  cmd,
					priority: 60,
				}
			}
		}
	}

	// 4. 关键词匹配
	queryWords := tokenize(query)
	for _, word := range queryWords {
		if commands, ok := idx.keywordIndex[word]; ok {
			for _, cmd := range commands {
				if _, exists := resultMap[cmd.Name]; !exists {
					resultMap[cmd.Name] = &searchResult{
						command:  cmd,
						priority: 40,
					}
				}
			}
		}
	}

	// 转换为数组并排序
	var results []*searchResult
	for _, r := range resultMap {
		results = append(results, r)
	}

	// 按优先级排序
	sortSearchResults(results)

	// 提取命令
	var commands []*model.Command
	for _, r := range results {
		commands = append(commands, r.command)
	}

	return commands
}

// GetAllCategories 获取所有分类
func (idx *Index) GetAllCategories() []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	categories := make([]string, 0, len(idx.categoryIndex))
	for category := range idx.categoryIndex {
		categories = append(categories, category)
	}
	return categories
}

// GetAllCommands 获取所有命令
func (idx *Index) GetAllCommands() []*model.Command {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	commands := make([]*model.Command, 0, len(idx.nameIndex))
	for _, cmd := range idx.nameIndex {
		commands = append(commands, cmd)
	}
	return commands
}

// searchResult 搜索结果
type searchResult struct {
	command  *model.Command
	priority int
}

// sortSearchResults 对搜索结果排序
func sortSearchResults(results []*searchResult) {
	// 简单冒泡排序（按优先级降序）
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].priority < results[j].priority {
				results[i], results[j] = results[j], results[i]
			} else if results[i].priority == results[j].priority {
				// 同优先级按名称排序
				if results[i].command.Name > results[j].command.Name {
					results[i], results[j] = results[j], results[i]
				}
			}
		}
	}
}
