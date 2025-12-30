package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cmd4coder/cmd4coder/internal/service"
	"github.com/cmd4coder/cmd4coder/internal/ui/tui"
	"github.com/spf13/cobra"
)

var (
	// Version 版本信息
	Version = "1.0.0"
	// BuildTime 构建时间
	BuildTime = "unknown"
	// CommitHash Git提交哈希
	CommitHash = "unknown"

	// 全局命令服务
	cmdService *service.CommandService
	// 配置服务
	cfgService *service.ConfigService

	// 数据目录
	dataDir string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "cmd4coder",
	Short: "命令行工具大全 - 面向运维和开发者的命令行参考工具",
	Long: `cmd4coder 是一个简单优雅的命令行工具大全。

它提供了完整的命令清单，包括：
  - Linux 命令（Ubuntu/CentOS/通用）
  - 编程语言工具链（Java/Go/Python/Node.js等）
  - 诊断工具（Arthas/tsar等）
  - 网络工具（dig/curl/tcpdump等）
  - 容器编排（Docker/Kubernetes）
  - 数据库工具（MySQL/Redis/PostgreSQL）
  - 版本控制（Git/SVN）
  - 构建工具（Maven/Gradle/Make）

支持两种使用模式：
  1. CLI模式：通过命令行参数快速查询
  2. TUI模式：交互式文本界面浏览

更多信息请访问: https://github.com/cmd4coder/cmd4coder`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 初始化命令服务
		if dataDir == "" {
			// 默认数据目录
			execPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %w", err)
			}
			dataDir = filepath.Join(filepath.Dir(execPath), "data")
		}

		var err error
		cmdService, err = service.NewCommandService(dataDir)
		if err != nil {
			return fmt.Errorf("failed to initialize command service: %w", err)
		}

		// 初始化配置服务
		cfgService, _ = service.NewConfigService()

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// 无子命令时启动TUI模式
		if err := tui.Run(cmdService, cfgService); err != nil {
			fmt.Fprintf(os.Stderr, "TUI错误: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dataDir, "data-dir", "d", "", "数据目录路径")
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(categoriesCmd)
	rootCmd.AddCommand(versionCmd)
}
