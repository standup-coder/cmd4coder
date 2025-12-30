package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cmd4coder/cmd4coder/internal/data"
	"github.com/cmd4coder/cmd4coder/internal/model"
)

// ValidationReport 验证报告
type ValidationReport struct {
	TotalFiles      int
	TotalCommands   int
	SuccessFiles    int
	FailedFiles     int
	Errors          []ValidationError
	Warnings        []ValidationWarning
	CommandsByCategory map[string]int
}

// ValidationError 验证错误
type ValidationError struct {
	File    string
	Command string
	Field   string
	Message string
}

// ValidationWarning 验证警告
type ValidationWarning struct {
	File    string
	Command string
	Message string
}

func main() {
	dataDir := flag.String("d", "./data", "数据目录路径")
	verbose := flag.Bool("v", false, "详细输出")
	flag.Parse()

	fmt.Println("CMD4Coder 数据验证工具")
	fmt.Println("====================")
	fmt.Printf("数据目录: %s\n\n", *dataDir)

	// 加载元数据
	loader := data.NewLoader(*dataDir)
	
	metadata, err := loader.LoadMetadata()
	if err != nil {
		fmt.Printf("❌ 无法加载元数据: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ 元数据加载成功\n")
	fmt.Printf("  - 定义分类数: %d\n", len(metadata.Categories))
	fmt.Printf("  - 数据文件数: %d\n\n", len(metadata.DataFiles))

	// 验证所有数据文件
	report := &ValidationReport{
		CommandsByCategory: make(map[string]int),
	}

	fmt.Println("开始验证数据文件...")
	fmt.Println("----------------------------------------")

	for _, dataFile := range metadata.DataFiles {
		report.TotalFiles++
		filePath := filepath.Join(*dataDir, dataFile)
		
		cmdList, err := loader.LoadCommandList(filePath)
		if err != nil {
			report.FailedFiles++
			report.Errors = append(report.Errors, ValidationError{
				File:    dataFile,
				Message: err.Error(),
			})
			fmt.Printf("❌ %s - 加载失败: %v\n", dataFile, err)
			continue
		}

		// 验证命令列表
		if err := cmdList.Validate(); err != nil {
			report.FailedFiles++
			report.Errors = append(report.Errors, ValidationError{
				File:    dataFile,
				Message: err.Error(),
			})
			fmt.Printf("❌ %s - 验证失败: %v\n", dataFile, err)
			if *verbose {
				fmt.Printf("   - %v\n", err)
			}
			continue
		}

		// 统计命令
		report.SuccessFiles++
		report.TotalCommands += len(cmdList.Commands)
		report.CommandsByCategory[cmdList.Category] += len(cmdList.Commands)

		// 检查警告项
		for _, cmd := range cmdList.Commands {
			// 检查是否缺少install_method
			if cmd.InstallMethod == "" {
				report.Warnings = append(report.Warnings, ValidationWarning{
					File:    dataFile,
					Command: cmd.Name,
					Message: "缺少 install_method 字段",
				})
			}

			// 检查是否缺少version_check
			if cmd.Versions == nil && len(cmd.Usage) > 0 {
				// 建议添加版本检查方法
			}

			// 检查examples数量
			if len(cmd.Examples) < 2 {
				report.Warnings = append(report.Warnings, ValidationWarning{
					File:    dataFile,
					Command: cmd.Name,
					Message: fmt.Sprintf("示例数量较少 (%d)，建议至少2个", len(cmd.Examples)),
				})
			}

			// 检查高风险命令是否有足够的风险说明
			if cmd.GetHighestRisk() >= model.RiskLevelHigh {
				if len(cmd.Risks) < 2 {
					report.Warnings = append(report.Warnings, ValidationWarning{
						File:    dataFile,
						Command: cmd.Name,
						Message: "高风险命令建议提供详细的风险说明",
					})
				}
			}
		}

		fmt.Printf("✓ %s - 验证通过 (%d 个命令)\n", dataFile, len(cmdList.Commands))
	}

	// 输出报告
	fmt.Println("\n========================================")
	fmt.Println("验证报告")
	fmt.Println("========================================")
	fmt.Printf("总文件数: %d\n", report.TotalFiles)
	fmt.Printf("成功: %d\n", report.SuccessFiles)
	fmt.Printf("失败: %d\n", report.FailedFiles)
	fmt.Printf("总命令数: %d\n", report.TotalCommands)
	
	if len(report.Errors) > 0 {
		fmt.Printf("\n❌ 错误 (%d):\n", len(report.Errors))
		for i, e := range report.Errors {
			if i < 20 { // 只显示前20个错误
				fmt.Printf("  [%s] %s: %s\n", e.File, e.Command, e.Message)
			}
		}
		if len(report.Errors) > 20 {
			fmt.Printf("  ... 还有 %d 个错误\n", len(report.Errors)-20)
		}
	}

	if len(report.Warnings) > 0 && *verbose {
		fmt.Printf("\n⚠️  警告 (%d):\n", len(report.Warnings))
		for i, w := range report.Warnings {
			if i < 10 {
				fmt.Printf("  [%s] %s: %s\n", w.File, w.Command, w.Message)
			}
		}
		if len(report.Warnings) > 10 {
			fmt.Printf("  ... 还有 %d 个警告\n", len(report.Warnings)-10)
		}
	}

	// 分类统计
	fmt.Println("\n分类统计:")
	fmt.Println("----------------------------------------")
	for category, count := range report.CommandsByCategory {
		fmt.Printf("  %-40s %3d 个命令\n", category, count)
	}

	// 质量评分
	fmt.Println("\n数据质量评分:")
	fmt.Println("----------------------------------------")
	
	completeness := float64(report.TotalCommands) / 350.0 * 100
	fmt.Printf("完整度: %.1f%% (%d/350)\n", completeness, report.TotalCommands)
	
	accuracy := float64(report.SuccessFiles) / float64(report.TotalFiles) * 100
	fmt.Printf("准确率: %.1f%% (%d/%d 文件通过验证)\n", accuracy, report.SuccessFiles, report.TotalFiles)

	warningRate := float64(len(report.Warnings)) / float64(report.TotalCommands) * 100
	fmt.Printf("警告率: %.1f%% (%d 个警告)\n", warningRate, len(report.Warnings))

	// 总体评分
	overallScore := (accuracy*0.6 + (100-warningRate)*0.2 + completeness*0.2)
	fmt.Printf("\n总体评分: %.1f/100\n", overallScore)

	if overallScore >= 90 {
		fmt.Println("评级: ⭐⭐⭐⭐⭐ 优秀")
	} else if overallScore >= 80 {
		fmt.Println("评级: ⭐⭐⭐⭐ 良好")
	} else if overallScore >= 70 {
		fmt.Println("评级: ⭐⭐⭐ 中等")
	} else {
		fmt.Println("评级: ⭐⭐ 需要改进")
	}

	// 根据结果设置退出码
	if report.FailedFiles > 0 {
		os.Exit(1)
	}
}
