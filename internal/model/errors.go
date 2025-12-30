package model

import "fmt"

// ErrMissingField 缺少必填字段错误
type ErrMissingField struct {
	Field string
}

func (e ErrMissingField) Error() string {
	return fmt.Sprintf("missing required field: %s", e.Field)
}

// ErrInvalidRiskLevel 无效的风险级别错误
type ErrInvalidRiskLevel struct {
	Command string
	Index   int
	Level   string
}

func (e ErrInvalidRiskLevel) Error() string {
	return fmt.Sprintf("invalid risk level '%s' in command '%s' at index %d", e.Level, e.Command, e.Index)
}

// ErrInvalidCategory 无效的分类错误
type ErrInvalidCategory struct {
	Category string
}

func (e ErrInvalidCategory) Error() string {
	return fmt.Sprintf("invalid category: %s", e.Category)
}

// ErrDuplicateCommand 重复的命令错误
type ErrDuplicateCommand struct {
	Name string
}

func (e ErrDuplicateCommand) Error() string {
	return fmt.Sprintf("duplicate command: %s", e.Name)
}

// ErrCommandNotFound 命令未找到错误
type ErrCommandNotFound struct {
	Name string
}

func (e ErrCommandNotFound) Error() string {
	return fmt.Sprintf("command not found: %s", e.Name)
}

// ErrDataLoadFailed 数据加载失败错误
type ErrDataLoadFailed struct {
	File string
	Err  error
}

func (e ErrDataLoadFailed) Error() string {
	return fmt.Sprintf("failed to load data file '%s': %v", e.File, e.Err)
}
