package model

import "time"

// RiskLevel 风险级别
type RiskLevel string

const (
	RiskLevelLow      RiskLevel = "low"      // 低风险，一般操作
	RiskLevelMedium   RiskLevel = "medium"   // 中风险，需谨慎操作
	RiskLevelHigh     RiskLevel = "high"     // 高风险，可能影响系统稳定性
	RiskLevelCritical RiskLevel = "critical" // 严重风险，可能导致数据丢失或系统崩溃
)

// IsValid 检查风险级别是否有效
func (r RiskLevel) IsValid() bool {
	switch r {
	case RiskLevelLow, RiskLevelMedium, RiskLevelHigh, RiskLevelCritical:
		return true
	default:
		return false
	}
}

// Risk 风险描述
type Risk struct {
	Level       RiskLevel `yaml:"level" json:"level"`             // 风险级别
	Description string    `yaml:"description" json:"description"` // 风险描述
}

// Option 命令选项
type Option struct {
	Flag        string `yaml:"flag" json:"flag"`               // 选项标志，如 -a, --all
	Description string `yaml:"description" json:"description"` // 选项描述
}

// Example 使用示例
type Example struct {
	Command     string `yaml:"command" json:"command"`                   // 示例命令
	Description string `yaml:"description" json:"description"`           // 示例说明
	Output      string `yaml:"output,omitempty" json:"output,omitempty"` // 预期输出（可选）
}

// VersionInfo 版本信息
type VersionInfo struct {
	MinVersion string `yaml:"min_version,omitempty" json:"min_version,omitempty"` // 最低版本
	MaxVersion string `yaml:"max_version,omitempty" json:"max_version,omitempty"` // 最高版本
	Notes      string `yaml:"notes,omitempty" json:"notes,omitempty"`             // 版本说明
}

// Command 命令结构
type Command struct {
	Name            string       `yaml:"name" json:"name"`                                             // 命令名称
	Category        string       `yaml:"category" json:"category"`                                     // 所属分类
	InstallRequired bool         `yaml:"install_required" json:"install_required"`                     // 是否需要单独安装
	InstallMethod   string       `yaml:"install_method,omitempty" json:"install_method,omitempty"`     // 安装方式说明
	Description     string       `yaml:"description" json:"description"`                               // 命令功能简述
	Usage           []string     `yaml:"usage" json:"usage"`                                           // 常用使用方式
	Options         []Option     `yaml:"options" json:"options"`                                       // 常用选项说明
	Examples        []Example    `yaml:"examples" json:"examples"`                                     // 使用示例
	Notes           []string     `yaml:"notes,omitempty" json:"notes,omitempty"`                       // 注意事项
	Risks           []Risk       `yaml:"risks,omitempty" json:"risks,omitempty"`                       // 风险说明
	RelatedCommands []string     `yaml:"related_commands,omitempty" json:"related_commands,omitempty"` // 相关命令
	Platforms       []string     `yaml:"platforms" json:"platforms"`                                   // 支持的平台
	Versions        *VersionInfo `yaml:"versions,omitempty" json:"versions,omitempty"`                 // 版本兼容性说明
	References      []string     `yaml:"references,omitempty" json:"references,omitempty"`             // 参考链接
}

// Validate 验证命令数据完整性
func (c *Command) Validate() error {
	if c.Name == "" {
		return ErrMissingField{Field: "name"}
	}
	if c.Category == "" {
		return ErrMissingField{Field: "category"}
	}
	if c.Description == "" {
		return ErrMissingField{Field: "description"}
	}
	if len(c.Usage) == 0 {
		return ErrMissingField{Field: "usage"}
	}
	if len(c.Examples) == 0 {
		return ErrMissingField{Field: "examples"}
	}
	if len(c.Platforms) == 0 {
		return ErrMissingField{Field: "platforms"}
	}

	// 验证风险级别
	for i, risk := range c.Risks {
		if !risk.Level.IsValid() {
			return ErrInvalidRiskLevel{
				Command: c.Name,
				Index:   i,
				Level:   string(risk.Level),
			}
		}
	}

	return nil
}

// GetRiskLevel 获取命令的最高风险级别
func (c *Command) GetRiskLevel() RiskLevel {
	if len(c.Risks) == 0 {
		return RiskLevelLow
	}

	highest := RiskLevelLow
	for _, risk := range c.Risks {
		if riskLevelValue(risk.Level) > riskLevelValue(highest) {
			highest = risk.Level
		}
	}
	return highest
}

// GetHighestRisk 获取命令的最高风险级别（别名）
func (c *Command) GetHighestRisk() RiskLevel {
	return c.GetRiskLevel()
}

// riskLevelValue 获取风险级别的数值表示
func riskLevelValue(level RiskLevel) int {
	switch level {
	case RiskLevelLow:
		return 1
	case RiskLevelMedium:
		return 2
	case RiskLevelHigh:
		return 3
	case RiskLevelCritical:
		return 4
	default:
		return 0
	}
}

// SupportsPlatform 检查命令是否支持指定平台
func (c *Command) SupportsPlatform(platform string) bool {
	for _, p := range c.Platforms {
		if p == platform {
			return true
		}
	}
	return false
}

// HasPlatform 检查命令是否支持指定平台（别名）
func (c *Command) HasPlatform(platform string) bool {
	return c.SupportsPlatform(platform)
}

// CommandList 命令列表的包装类型
type CommandList struct {
	Category    string     `yaml:"category" json:"category"`                         // 分类名称
	Description string     `yaml:"description" json:"description"`                   // 分类描述
	Commands    []*Command `yaml:"commands" json:"commands"`                         // 命令列表
	UpdatedAt   time.Time  `yaml:"updated_at,omitempty" json:"updated_at,omitempty"` // 更新时间
}

// Validate 验证命令列表
func (cl *CommandList) Validate() error {
	if cl.Category == "" {
		return ErrMissingField{Field: "category"}
	}
	if cl.Description == "" {
		return ErrMissingField{Field: "description"}
	}

	// 验证每个命令
	for _, cmd := range cl.Commands {
		if err := cmd.Validate(); err != nil {
			return err
		}
	}

	return nil
}
