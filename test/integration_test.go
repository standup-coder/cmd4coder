package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cmd4coder/cmd4coder/internal/model"
	"github.com/cmd4coder/cmd4coder/internal/service"
)

// TestCommandServiceIntegration 命令服务集成测试
func TestCommandServiceIntegration(t *testing.T) {
	// 获取数据目录
	dataDir := filepath.Join("..", "data")
	
	// 检查数据目录是否存在
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skip("数据目录不存在，跳过集成测试")
	}
	
	// 创建命令服务
	cmdService, err := service.NewCommandService(dataDir)
	if err != nil {
		t.Fatalf("创建命令服务失败: %v", err)
	}
	
	// 测试获取所有命令
	t.Run("GetAllCommands", func(t *testing.T) {
		commands := cmdService.GetAllCommands()
		if len(commands) == 0 {
			t.Error("应该至少有一个命令")
		}
		
		t.Logf("总命令数: %d", len(commands))
	})
	
	// 测试获取所有分类
	t.Run("GetAllCategories", func(t *testing.T) {
		categories := cmdService.GetAllCategories()
		if len(categories) == 0 {
			t.Error("应该至少有一个分类")
		}
		
		t.Logf("总分类数: %d", len(categories))
	})
	
	// 测试搜索功能
	t.Run("SearchCommands", func(t *testing.T) {
		results := cmdService.SearchCommands("ls")
		if len(results) == 0 {
			t.Error("搜索'ls'应该有结果")
		}
		
		t.Logf("搜索'ls'找到 %d 个命令", len(results))
	})
	
	// 测试根据分类获取命令
	t.Run("ListCommandsByCategory", func(t *testing.T) {
		categories := cmdService.GetAllCategories()
		if len(categories) == 0 {
			t.Skip("没有分类可测试")
		}
		
		category := categories[0]
		commands := cmdService.ListCommandsByCategory(category)
		
		t.Logf("分类 '%s' 有 %d 个命令", category, len(commands))
		
		// 验证每个命令的分类
		for _, cmd := range commands {
			if cmd.Category != category {
				t.Errorf("命令 %s 的分类应该是 %s，实际是 %s", cmd.Name, category, cmd.Category)
			}
		}
	})
	
	// 测试获取特定命令
	t.Run("GetCommand", func(t *testing.T) {
		allCommands := cmdService.GetAllCommands()
		if len(allCommands) == 0 {
			t.Skip("没有命令可测试")
		}
		
		cmdName := allCommands[0].Name
		cmd, err := cmdService.GetCommand(cmdName)
		if err != nil {
			t.Fatalf("获取命令 %s 失败: %v", cmdName, err)
		}
		
		if cmd.Name != cmdName {
			t.Errorf("期望命令名 %s，实际是 %s", cmdName, cmd.Name)
		}
	})
	
	// 测试缓存功能
	t.Run("SearchCache", func(t *testing.T) {
		query := "test_cache_query_12345"
		
		// 第一次搜索
		results1 := cmdService.SearchCommands(query)
		
		// 第二次搜索（应该从缓存获取）
		results2 := cmdService.SearchCommands(query)
		
		if len(results1) != len(results2) {
			t.Error("缓存结果应该与第一次搜索相同")
		}
	})
	
	// 测试高风险命令
	t.Run("GetHighRiskCommands", func(t *testing.T) {
		highRiskCommands := cmdService.GetHighRiskCommands()
		
		t.Logf("高风险命令数: %d", len(highRiskCommands))
		
		// 验证每个命令都是高风险
		for _, cmd := range highRiskCommands {
			risk := cmd.GetHighestRisk()
			if risk != model.RiskLevelHigh && risk != model.RiskLevelCritical {
				t.Errorf("命令 %s 不是高风险命令", cmd.Name)
			}
		}
	})
}

// TestConfigServiceIntegration 配置服务集成测试
func TestConfigServiceIntegration(t *testing.T) {
	// 创建配置服务
	cfgService, err := service.NewConfigService()
	if err != nil {
		t.Fatalf("创建配置服务失败: %v", err)
	}
	
	// 测试获取配置
	t.Run("GetConfig", func(t *testing.T) {
		config := cfgService.GetConfig()
		if config == nil {
			t.Fatal("配置不应该为nil")
		}
		
		if config.Language == "" {
			t.Error("语言不应该为空")
		}
	})
	
	// 测试用户数据
	t.Run("UserData", func(t *testing.T) {
		userData := cfgService.GetUserData()
		if userData == nil {
			t.Fatal("用户数据不应该为nil")
		}
		
		// 测试添加收藏
		err := cfgService.AddFavorite("test_cmd", "test_category", "test note")
		if err != nil {
			t.Fatalf("添加收藏失败: %v", err)
		}
		
		// 验证收藏
		if !cfgService.IsFavorite("test_cmd") {
			t.Error("test_cmd 应该在收藏中")
		}
		
		// 测试删除收藏
		err = cfgService.RemoveFavorite("test_cmd")
		if err != nil {
			t.Fatalf("删除收藏失败: %v", err)
		}
		
		// 验证已删除
		if cfgService.IsFavorite("test_cmd") {
			t.Error("test_cmd 不应该在收藏中")
		}
	})
	
	// 测试历史记录
	t.Run("History", func(t *testing.T) {
		// 添加历史记录
		err := cfgService.AddHistory("cmd1", "cat1")
		if err != nil {
			t.Fatalf("添加历史记录失败: %v", err)
		}
		
		err = cfgService.AddHistory("cmd2", "cat2")
		if err != nil {
			t.Fatalf("添加历史记录失败: %v", err)
		}
		
		// 获取最近历史
		recent := cfgService.GetRecentHistory(10)
		if len(recent) < 2 {
			t.Errorf("应该至少有2条历史记录，实际有 %d 条", len(recent))
		}
		
		// 最新的应该在前面
		if recent[0].CommandName != "cmd2" {
			t.Errorf("最新的历史记录应该是 cmd2，实际是 %s", recent[0].CommandName)
		}
		
		// 清空历史
		err = cfgService.ClearHistory()
		if err != nil {
			t.Fatalf("清空历史失败: %v", err)
		}
		
		recent = cfgService.GetRecentHistory(10)
		if len(recent) != 0 {
			t.Errorf("清空后应该没有历史记录，实际有 %d 条", len(recent))
		}
	})
}

// TestDataLoadingPerformance 数据加载性能测试
func TestDataLoadingPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}
	
	dataDir := filepath.Join("..", "data")
	
	// 检查数据目录是否存在
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skip("数据目录不存在，跳过性能测试")
	}
	
	t.Run("StartupTime", func(t *testing.T) {
		// 测试启动时间
		_, err := service.NewCommandService(dataDir)
		if err != nil {
			t.Fatalf("创建服务失败: %v", err)
		}
		
		// 实际性能会由 testing 框架报告
		t.Log("启动完成")
	})
}

// TestSearchPerformance 搜索性能测试
func TestSearchPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过性能测试")
	}
	
	dataDir := filepath.Join("..", "data")
	
	// 检查数据目录是否存在
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skip("数据目录不存在，跳过性能测试")
	}
	
	cmdService, err := service.NewCommandService(dataDir)
	if err != nil {
		t.Fatalf("创建服务失败: %v", err)
	}
	
	queries := []string{"ls", "docker", "git", "java", "maven"}
	
	for _, query := range queries {
		t.Run(query, func(t *testing.T) {
			results := cmdService.SearchCommands(query)
			t.Logf("搜索 '%s' 找到 %d 个结果", query, len(results))
		})
	}
}

// BenchmarkSearch 搜索基准测试
func BenchmarkSearch(b *testing.B) {
	dataDir := filepath.Join("..", "data")
	
	cmdService, err := service.NewCommandService(dataDir)
	if err != nil {
		b.Fatalf("创建服务失败: %v", err)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		cmdService.SearchCommands("ls")
	}
}

// BenchmarkGetByCategory 分类查询基准测试
func BenchmarkGetByCategory(b *testing.B) {
	dataDir := filepath.Join("..", "data")
	
	cmdService, err := service.NewCommandService(dataDir)
	if err != nil {
		b.Fatalf("创建服务失败: %v", err)
	}
	
	categories := cmdService.GetAllCategories()
	if len(categories) == 0 {
		b.Skip("没有分类")
	}
	
	category := categories[0]
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		cmdService.ListCommandsByCategory(category)
	}
}
