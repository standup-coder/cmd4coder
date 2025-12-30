package main

import (
	"fmt"
	"strings"

	"github.com/cmd4coder/cmd4coder/internal/model"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [category]",
	Short: "列出命令",
	Long:  `列出指定分类下的所有命令，如果不指定分类则列出所有命令`,
	Example: `  cmd4coder list
  cmd4coder list "操作系统/Ubuntu系统命令"
  cmd4coder list "编程语言/Java工具链"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var commands []*model.Command
		var title string

		if len(args) == 0 {
			// 列出所有命令
			commands = cmdService.GetAllCommands()
			title = "所有命令"
		} else {
			// 列出指定分类的命令
			category := args[0]
			commands = cmdService.ListCommandsByCategory(category)
			title = fmt.Sprintf("分类: %s", category)
		}

		if len(commands) == 0 {
			fmt.Println("未找到命令")
			return nil
		}

		// 输出命令列表
		fmt.Printf("\n%s (共 %d 个命令)\n", title, len(commands))
		fmt.Println(strings.Repeat("=", 80))

		for _, cmd := range commands {
			riskIndicator := getRiskIndicator(cmd.GetHighestRisk())
			installIndicator := ""
			if cmd.InstallRequired {
				installIndicator = "[需安装]"
			}

			fmt.Printf("%-20s %s %s %s\n",
				cmd.Name,
				riskIndicator,
				installIndicator,
				cmd.Description)
		}

		fmt.Println()
		fmt.Println("使用 'cmd4coder show <命令名>' 查看详细信息")

		return nil
	},
}

var showCmd = &cobra.Command{
	Use:   "show <command>",
	Short: "显示命令详细信息",
	Long:  `显示指定命令的完整信息，包括用法、选项、示例、注意事项和风险说明`,
	Example: `  cmd4coder show ls
  cmd4coder show docker
  cmd4coder show git`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdName := args[0]
		command, err := cmdService.GetCommand(cmdName)
		if err != nil {
			return fmt.Errorf("命令 '%s' 未找到", cmdName)
		}

		printCommandDetail(command)
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "搜索命令",
	Long:  `根据关键词搜索命令，支持模糊匹配和多关键词`,
	Example: `  cmd4coder search file
  cmd4coder search network
  cmd4coder search "java 诊断"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		commands := cmdService.SearchCommands(query)

		if len(commands) == 0 {
			fmt.Printf("未找到与 '%s' 相关的命令\n", query)
			return nil
		}

		fmt.Printf("\n搜索结果: '%s' (共 %d 个命令)\n", query, len(commands))
		fmt.Println(strings.Repeat("=", 80))

		for _, command := range commands {
			riskIndicator := getRiskIndicator(command.GetHighestRisk())
			fmt.Printf("%-20s %s %s\n",
				command.Name,
				riskIndicator,
				command.Description)
		}

		fmt.Println()
		fmt.Println("使用 'cmd4coder show <命令名>' 查看详细信息")

		return nil
	},
}

var categoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "列出所有分类",
	Long:  `显示所有可用的命令分类`,
	RunE: func(cmd *cobra.Command, args []string) error {
		categories := cmdService.GetAllCategories()

		fmt.Printf("\n所有分类 (共 %d 个)\n", len(categories))
		fmt.Println(strings.Repeat("=", 80))

		for _, category := range categories {
			commands := cmdService.ListCommandsByCategory(category)
			fmt.Printf("%-40s (%d 个命令)\n", category, len(commands))
		}

		fmt.Println()
		fmt.Println("使用 'cmd4coder list <分类名>' 查看分类下的命令")

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cmd4coder version %s\n", Version)
		fmt.Printf("Build time: %s\n", BuildTime)
		fmt.Printf("Commit: %s\n", CommitHash)

		if cmdService != nil {
			metadata := cmdService.GetMetadata()
			if metadata != nil {
				fmt.Printf("Data version: %s\n", metadata.Version)
				fmt.Printf("Data updated: %s\n", metadata.UpdatedAt)
			}
			fmt.Printf("Total commands: %d\n", cmdService.GetCommandCount())
			fmt.Printf("Total categories: %d\n", cmdService.GetCategoryCount())
		}
	},
}

// Helper functions

func getRiskIndicator(risk model.RiskLevel) string {
	switch risk {
	case model.RiskLevelLow:
		return "🟢"
	case model.RiskLevelMedium:
		return "🟡"
	case model.RiskLevelHigh:
		return "🟠"
	case model.RiskLevelCritical:
		return "🔴"
	default:
		return "  "
	}
}

func printCommandDetail(cmd *model.Command) {
	fmt.Printf("\n命令: %s\n", cmd.Name)
	fmt.Println(strings.Repeat("=", 80))

	fmt.Printf("\n📝 描述:\n  %s\n", cmd.Description)
	fmt.Printf("\n📂 分类: %s\n", cmd.Category)
	fmt.Printf("💻 平台: %s\n", strings.Join(cmd.Platforms, ", "))

	if cmd.InstallRequired {
		fmt.Printf("\n📦 安装方式:\n  %s\n", cmd.InstallMethod)
	}

	// 使用方式
	fmt.Printf("\n💡 使用方式:\n")
	for _, usage := range cmd.Usage {
		fmt.Printf("  %s\n", usage)
	}

	// 常用选项
	if len(cmd.Options) > 0 {
		fmt.Printf("\n⚙️  常用选项:\n")
		for _, opt := range cmd.Options {
			fmt.Printf("  %-20s %s\n", opt.Flag, opt.Description)
		}
	}

	// 示例
	if len(cmd.Examples) > 0 {
		fmt.Printf("\n📋 使用示例:\n")
		for i, example := range cmd.Examples {
			fmt.Printf("\n  示例 %d: %s\n", i+1, example.Description)
			fmt.Printf("  $ %s\n", example.Command)
			if example.Output != "" {
				fmt.Printf("  输出: %s\n", example.Output)
			}
		}
	}

	// 注意事项
	if len(cmd.Notes) > 0 {
		fmt.Printf("\n⚠️  注意事项:\n")
		for _, note := range cmd.Notes {
			fmt.Printf("  • %s\n", note)
		}
	}

	// 风险说明
	if len(cmd.Risks) > 0 {
		fmt.Printf("\n⚡ 风险说明:\n")
		for _, risk := range cmd.Risks {
			indicator := getRiskIndicator(risk.Level)
			fmt.Printf("  %s [%s] %s\n", indicator, risk.Level, risk.Description)
		}
	}

	// 相关命令
	if len(cmd.RelatedCommands) > 0 {
		fmt.Printf("\n🔗 相关命令: %s\n", strings.Join(cmd.RelatedCommands, ", "))
	}

	// 参考链接
	if len(cmd.References) > 0 {
		fmt.Printf("\n📚 参考链接:\n")
		for _, ref := range cmd.References {
			fmt.Printf("  %s\n", ref)
		}
	}

	fmt.Println()
}
