package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cmd4coder/cmd4coder/internal/model"
	"github.com/cmd4coder/cmd4coder/internal/service"
)

// Model TUI模型
type Model struct {
	// 服务
	commandService *service.CommandService
	configService  *service.ConfigService

	// 数据
	categories  []string
	commands    []*model.Command
	selectedCmd *model.Command

	// UI组件
	searchInput  textinput.Model
	categoryList list.Model
	commandList  list.Model

	// 状态
	activePanel int // 0: search, 1: category, 2: command, 3: detail
	width       int
	height      int
	ready       bool

	// 键盘绑定
	keys keyMap
}

// keyMap 键盘映射
type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	Tab      key.Binding
	Search   key.Binding
	Favorite key.Binding
	Export   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

// 默认键盘绑定
var defaultKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "向上"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "向下"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "向左"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "向右"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "选择"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "切换面板"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "搜索"),
	),
	Favorite: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "收藏"),
	),
	Export: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "导出"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "帮助"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "退出"),
	),
}

// NewModel 创建新的TUI模型
func NewModel(cmdService *service.CommandService, cfgService *service.ConfigService) *Model {
	// 搜索输入框
	ti := textinput.New()
	ti.Placeholder = "搜索命令..."
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	return &Model{
		commandService: cmdService,
		configService:  cfgService,
		searchInput:    ti,
		activePanel:    0,
		keys:           defaultKeys,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update 更新模型
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.setupLists()
			m.ready = true
		}
		return m, nil

	case tea.KeyMsg:
		// 全局快捷键
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Tab):
			m.activePanel = (m.activePanel + 1) % 3
			m.updateFocus()
			return m, nil

		case key.Matches(msg, m.keys.Search):
			m.activePanel = 0
			m.searchInput.Focus()
			return m, nil
		}

		// 面板特定的键盘处理
		return m.handlePanelInput(msg)
	}

	// 更新搜索输入框
	m.searchInput, cmd = m.searchInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// handlePanelInput 处理面板输入
func (m Model) handlePanelInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.activePanel {
	case 0: // 搜索面板
		switch {
		case key.Matches(msg, m.keys.Enter):
			// 执行搜索
			m.performSearch()
			m.activePanel = 1
			return m, nil
		case key.Matches(msg, m.keys.Down):
			m.activePanel = 1
			m.updateFocus()
			return m, nil
		}

	case 1: // 分类列表
		switch {
		case key.Matches(msg, m.keys.Up):
			if m.categoryList.Index() == 0 {
				m.activePanel = 0
				m.searchInput.Focus()
			}
		case key.Matches(msg, m.keys.Enter), key.Matches(msg, m.keys.Right):
			m.loadCategoryCommands()
			m.activePanel = 2
			m.updateFocus()
			return m, nil
		}
		m.categoryList, cmd = m.categoryList.Update(msg)

	case 2: // 命令列表
		switch {
		case key.Matches(msg, m.keys.Left):
			m.activePanel = 1
			m.updateFocus()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			m.loadCommandDetail()
			return m, nil
		case key.Matches(msg, m.keys.Favorite):
			m.toggleFavorite()
			return m, nil
		}
		m.commandList, cmd = m.commandList.Update(msg)
	}

	return m, cmd
}

// View 渲染视图
func (m Model) View() string {
	if !m.ready {
		return "初始化中..."
	}

	// 样式
	docStyle := lipgloss.NewStyle().Padding(1, 2)

	// 标题
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Render("CMD4Coder - 命令速查工具")

	// 搜索栏
	searchBar := m.renderSearchBar()

	// 三栏布局
	panels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.renderCategoryPanel(),
		m.renderCommandPanel(),
		m.renderDetailPanel(),
	)

	// 状态栏
	statusBar := m.renderStatusBar()

	// 帮助信息
	helpBar := m.renderHelpBar()

	// 组合所有部分
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		searchBar,
		panels,
		statusBar,
		helpBar,
	)

	return docStyle.Render(content)
}

// renderSearchBar 渲染搜索栏
func (m Model) renderSearchBar() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(0, 1).
		Width(m.width - 6)

	if m.activePanel == 0 {
		style = style.BorderForeground(lipgloss.Color("170"))
	}

	return style.Render(m.searchInput.View())
}

// renderCategoryPanel 渲染分类面板
func (m Model) renderCategoryPanel() string {
	panelWidth := (m.width - 6) / 3
	panelHeight := m.height - 12

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(panelWidth).
		Height(panelHeight)

	if m.activePanel == 1 {
		style = style.BorderForeground(lipgloss.Color("170"))
	}

	title := lipgloss.NewStyle().Bold(true).Render("📁 分类")

	if len(m.categories) == 0 {
		return style.Render(title + "\n\n无数据")
	}

	return style.Render(title + "\n" + m.categoryList.View())
}

// renderCommandPanel 渲染命令面板
func (m Model) renderCommandPanel() string {
	panelWidth := (m.width - 6) / 3
	panelHeight := m.height - 12

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(panelWidth).
		Height(panelHeight)

	if m.activePanel == 2 {
		style = style.BorderForeground(lipgloss.Color("170"))
	}

	title := lipgloss.NewStyle().Bold(true).Render("📝 命令")

	if len(m.commands) == 0 {
		return style.Render(title + "\n\n请选择分类")
	}

	return style.Render(title + "\n" + m.commandList.View())
}

// renderDetailPanel 渲染详情面板
func (m Model) renderDetailPanel() string {
	panelWidth := (m.width - 6) / 3
	panelHeight := m.height - 12

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Width(panelWidth).
		Height(panelHeight)

	title := lipgloss.NewStyle().Bold(true).Render("📖 详情")

	if m.selectedCmd == nil {
		return style.Render(title + "\n\n请选择命令")
	}

	detail := m.formatCommandDetail()
	return style.Render(title + "\n" + detail)
}

// renderStatusBar 渲染状态栏
func (m Model) renderStatusBar() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render

	totalCmds := m.commandService.Count()
	status := fmt.Sprintf("总命令数: %d | 当前分类: %d 个命令", totalCmds, len(m.commands))

	return style(status)
}

// renderHelpBar 渲染帮助栏
func (m Model) renderHelpBar() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		Render

	help := "tab:切换 /:搜索 f:收藏 e:导出 q:退出"
	return style(help)
}

// formatCommandDetail 格式化命令详情
func (m Model) formatCommandDetail() string {
	cmd := m.selectedCmd

	detail := fmt.Sprintf("名称: %s\n\n", cmd.Name)
	detail += fmt.Sprintf("描述: %s\n\n", cmd.Description)

	if len(cmd.Usage) > 0 {
		detail += "用法:\n"
		for _, u := range cmd.Usage {
			detail += fmt.Sprintf("  %s\n", u)
		}
		detail += "\n"
	}

	if len(cmd.Examples) > 0 {
		detail += "示例:\n"
		for i, ex := range cmd.Examples {
			if i >= 3 {
				break // 只显示前3个
			}
			detail += fmt.Sprintf("  %s\n  %s\n\n", ex.Command, ex.Description)
		}
	}

	return detail
}

// setupLists 设置列表
func (m *Model) setupLists() {
	// 加载分类
	cats := m.commandService.GetCategories()
	m.categories = cats

	// 设置分类列表
	items := make([]list.Item, len(cats))
	for i, cat := range cats {
		items[i] = listItem{title: cat, desc: ""}
	}

	m.categoryList = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.categoryList.Title = ""
	m.categoryList.SetShowStatusBar(false)
	m.categoryList.SetFilteringEnabled(false)
	m.categoryList.SetShowHelp(false)

	// 设置命令列表
	m.commandList = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.commandList.Title = ""
	m.commandList.SetShowStatusBar(false)
	m.commandList.SetFilteringEnabled(false)
	m.commandList.SetShowHelp(false)
}

// performSearch 执行搜索
func (m *Model) performSearch() {
	query := m.searchInput.Value()
	if query == "" {
		return
	}

	results := m.commandService.Search(query)
	m.commands = results

	// 更新命令列表
	items := make([]list.Item, len(results))
	for i, cmd := range results {
		items[i] = listItem{
			title: cmd.Name,
			desc:  cmd.Description,
		}
	}
	m.commandList.SetItems(items)
}

// loadCategoryCommands 加载分类下的命令
func (m *Model) loadCategoryCommands() {
	if len(m.categories) == 0 {
		return
	}

	selectedIdx := m.categoryList.Index()
	if selectedIdx < 0 || selectedIdx >= len(m.categories) {
		return
	}

	category := m.categories[selectedIdx]
	cmds := m.commandService.GetByCategory(category)
	m.commands = cmds

	// 更新命令列表
	items := make([]list.Item, len(cmds))
	for i, cmd := range cmds {
		items[i] = listItem{
			title: cmd.Name,
			desc:  cmd.Description,
		}
	}
	m.commandList.SetItems(items)
}

// loadCommandDetail 加载命令详情
func (m *Model) loadCommandDetail() {
	if len(m.commands) == 0 {
		return
	}

	selectedIdx := m.commandList.Index()
	if selectedIdx < 0 || selectedIdx >= len(m.commands) {
		return
	}

	m.selectedCmd = m.commands[selectedIdx]

	// 添加到历史记录
	if m.configService != nil {
		m.configService.AddHistory(m.selectedCmd.Name, m.selectedCmd.Category)
	}
}

// toggleFavorite 切换收藏状态
func (m *Model) toggleFavorite() {
	if m.selectedCmd == nil || m.configService == nil {
		return
	}

	if m.configService.IsFavorite(m.selectedCmd.Name) {
		m.configService.RemoveFavorite(m.selectedCmd.Name)
	} else {
		m.configService.AddFavorite(m.selectedCmd.Name, m.selectedCmd.Category, "")
	}
}

// updateFocus 更新焦点
func (m *Model) updateFocus() {
	if m.activePanel == 0 {
		m.searchInput.Focus()
	} else {
		m.searchInput.Blur()
	}
}

// listItem 列表项
type listItem struct {
	title string
	desc  string
}

func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string { return i.desc }
func (i listItem) FilterValue() string { return i.title }
