package export

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cmd4coder/cmd4coder/internal/model"
)

func TestExportToMarkdown(t *testing.T) {
	// 创建测试命令
	commands := []*model.Command{
		{
			Name:        "test-cmd",
			Category:    "test",
			Description: "Test command",
			Platforms:   []string{"linux"},
			Usage:       []string{"test-cmd [options]"},
			Options: []model.Option{
				{Flag: "-h", Description: "Show help"},
			},
			Examples: []model.Example{
				{Command: "test-cmd -h", Description: "Show help"},
			},
			Risks: []model.Risk{
				{Level: model.RiskLevelLow, Description: "No risk"},
			},
		},
	}

	// 创建临时目录
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")

	// 测试导出
	err := ExportToMarkdown(commands, testFile)
	if err != nil {
		t.Fatalf("ExportToMarkdown() error = %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	// 读取文件内容
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// 验证内容包含关键信息
	contentStr := string(content)
	expectedStrings := []string{
		"# 命令行工具大全",
		"test-cmd",
		"Test command",
		"linux",
	}

	for _, expected := range expectedStrings {
		if !contains(contentStr, expected) {
			t.Errorf("Output does not contain expected string: %s", expected)
		}
	}
}

func TestExportToJSON(t *testing.T) {
	commands := []*model.Command{
		{
			Name:        "test-cmd",
			Category:    "test",
			Description: "Test command",
			Platforms:   []string{"linux"},
		},
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	err := ExportToJSON(commands, testFile)
	if err != nil {
		t.Fatalf("ExportToJSON() error = %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	// 读取并验证JSON格式
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "test-cmd") {
		t.Errorf("JSON output does not contain command name")
	}
	if !contains(contentStr, "version") {
		t.Errorf("JSON output does not contain version field")
	}
}

func TestExportToJSONCompact(t *testing.T) {
	commands := []*model.Command{
		{
			Name:        "test-cmd",
			Category:    "test",
			Description: "Test command",
		},
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-compact.json")

	err := ExportToJSONCompact(commands, testFile)
	if err != nil {
		t.Fatalf("ExportToJSONCompact() error = %v", err)
	}

	// 验证文件存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created")
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// 紧凑格式应该更小（没有缩进）
	if len(content) == 0 {
		t.Errorf("Compact JSON output is empty")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr))))
}
