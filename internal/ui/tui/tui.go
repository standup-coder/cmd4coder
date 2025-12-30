package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cmd4coder/cmd4coder/internal/service"
)

// Run 运行TUI应用
func Run(cmdService *service.CommandService, cfgService *service.ConfigService) error {
	m := NewModel(cmdService, cfgService)
	
	p := tea.NewProgram(m, tea.WithAltScreen())
	
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("运行TUI失败: %w", err)
	}
	
	return nil
}
