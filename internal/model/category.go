package model

// Category 分类信息
type Category struct {
	ID          string   `yaml:"id" json:"id"`                                 // 分类ID
	Name        string   `yaml:"name" json:"name"`                             // 分类名称
	Description string   `yaml:"description" json:"description"`               // 分类描述
	Parent      string   `yaml:"parent,omitempty" json:"parent,omitempty"`     // 父分类ID（空表示顶级分类）
	Order       int      `yaml:"order" json:"order"`                           // 排序顺序
	Icon        string   `yaml:"icon,omitempty" json:"icon,omitempty"`         // 图标（用于TUI）
	Children    []string `yaml:"children,omitempty" json:"children,omitempty"` // 子分类ID列表
}

// IsTopLevel 判断是否为顶级分类
func (c *Category) IsTopLevel() bool {
	return c.Parent == ""
}

// CategoryTree 分类树结构
type CategoryTree struct {
	Category *Category
	Children []*CategoryTree
}

// Metadata 元数据
type Metadata struct {
	Version     string              `yaml:"version" json:"version"`         // 数据版本
	UpdatedAt   string              `yaml:"updated_at" json:"updated_at"`   // 更新时间
	Categories  map[string]Category `yaml:"categories" json:"categories"`   // 分类索引
	DataFiles   []string            `yaml:"data_files" json:"data_files"`   // 数据文件列表
	Description string              `yaml:"description" json:"description"` // 元数据描述
}

// Validate 验证元数据
func (m *Metadata) Validate() error {
	if m.Version == "" {
		return ErrMissingField{Field: "version"}
	}
	if len(m.Categories) == 0 {
		return ErrMissingField{Field: "categories"}
	}
	if len(m.DataFiles) == 0 {
		return ErrMissingField{Field: "data_files"}
	}
	return nil
}
