# 贡献指南

首先,感谢您考虑为cmd4coder做出贡献!正是像您这样的人使cmd4coder成为如此出色的工具。

## 📋 目录

- [行为准则](#行为准则)
- [我能做什么贡献](#我能做什么贡献)
- [开发环境搭建](#开发环境搭建)
- [提交流程](#提交流程)
- [代码规范](#代码规范)
- [命令数据贡献](#命令数据贡献)
- [问题报告](#问题报告)
- [功能请求](#功能请求)

## 行为准则

本项目及其所有参与者都受[行为准则](CODE_OF_CONDUCT.md)约束。参与即表示您同意遵守其条款。

## 我能做什么贡献

### 🐛 报告Bug

如果您发现了Bug,请创建一个Issue并提供详细信息:
- Bug的清晰描述
- 复现步骤
- 预期行为和实际行为
- 您的环境信息(操作系统、Go版本等)
- 相关日志或截图

### 💡 提出功能建议

我们欢迎新功能的想法!请创建一个Issue并说明:
- 功能的详细描述
- 为什么需要这个功能
- 如何实现(可选)
- 可能的替代方案

### 📝 改进文档

文档永远不嫌完善!您可以:
- 修复错别字或语法错误
- 改进现有说明
- 添加更多示例
- 翻译文档

### ➕ 添加命令数据

这是最常见的贡献方式!请参阅[命令数据贡献](#命令数据贡献)章节。

### 🔧 修复代码问题

浏览[Issue列表](https://github.com/cmd4coder/cmd4coder/issues),查找标记为`good first issue`或`help wanted`的问题。

## 开发环境搭建

### 环境要求

- Go 1.21 或更高版本
- Git
- 代码编辑器(推荐VS Code、GoLand或Vim)

### 克隆仓库

```bash
# Fork仓库后克隆您的Fork
git clone https://github.com/YOUR_USERNAME/cmd4coder.git
cd cmd4coder

# 添加上游仓库
git remote add upstream https://github.com/cmd4coder/cmd4coder.git
```

### 安装依赖

```bash
go mod download
go mod tidy
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行测试并查看覆盖率
go test -cover ./...

# 运行竞态检测
go test -race ./...
```

### 本地构建

```bash
# Linux/macOS
./build.sh

# Windows
.\build.ps1

# 或使用go build
go build -o bin/cmd4coder ./cmd/cli
```

### 运行程序

```bash
# CLI模式 (开发阶段推荐使用 go run)
go run ./cmd/cli list -d ./data

# TUI模式
go run ./cmd/cli -d ./data

# 构建后使用
./bin/cmd4coder list -d ./data     # Linux/macOS
.\build\cmd4coder-v1.0.0-windows-amd64.exe list -d ./data  # Windows
```

## 提交流程

### 1. Fork仓库

点击GitHub页面右上角的"Fork"按钮。

### 2. 创建特性分支

```bash
git checkout -b feature/your-feature-name
# 或
git checkout -b fix/your-bug-fix
```

分支命名建议:
- `feature/` - 新功能
- `fix/` - Bug修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

### 3. 进行修改

按照[代码规范](#代码规范)进行开发。

### 4. 提交更改

```bash
git add .
git commit -m "类型: 简短描述

详细描述(可选)

关联Issue(可选): #123"
```

提交信息类型:
- `feat`: 新功能
- `fix`: Bug修复
- `docs`: 文档更新
- `style`: 代码格式(不影响功能)
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

示例:
```
feat: 添加TUI交互式界面

实现了基于bubbletea的TUI界面,支持:
- 三栏布局
- 分类浏览
- 命令搜索
- 收藏功能

Closes #45
```

### 5. 推送到Fork仓库

```bash
git push origin feature/your-feature-name
```

### 6. 创建Pull Request

1. 访问您的Fork仓库页面
2. 点击"Pull Request"按钮
3. 填写PR描述:
   - 简要说明更改内容
   - 关联相关Issue
   - 说明测试情况
   - 添加截图(如适用)

### 7. 代码审查

维护者会审查您的PR并可能提出修改建议。请及时回应并进行必要的更新。

## 代码规范

### Go代码规范

我们遵循Go官方代码规范和最佳实践:

1. **代码格式化**: 使用`gofmt`或`goimports`
   ```bash
   go fmt ./...
   ```

2. **代码检查**: 使用`go vet`
   ```bash
   go vet ./...
   ```

3. **Lint检查**: 使用`golangci-lint`(推荐)
   ```bash
   golangci-lint run
   ```

### 命名约定

- **包名**: 小写,单个单词,如`model`、`service`
- **文件名**: 小写,下划线分隔,如`command_service.go`
- **变量/函数**: 驼峰命名,如`commandList`、`GetCommand`
- **常量**: 大写驼峰或全大写,如`MaxCacheSize`、`VERSION`
- **接口**: 通常以`-er`结尾,如`Commander`、`Loader`

### 注释规范

1. **包注释**: 每个包都应有包级注释
   ```go
   // Package model 定义核心数据模型
   package model
   ```

2. **函数注释**: 公开函数必须有注释
   ```go
   // GetCommand 根据名称获取命令信息
   // 如果找不到命令,返回nil和错误
   func GetCommand(name string) (*Command, error) {
       // ...
   }
   ```

3. **复杂逻辑注释**: 复杂算法或逻辑应添加解释

### 错误处理

1. 不要忽略错误
   ```go
   // ❌ 错误示例
   data, _ := loadData()
   
   // ✅ 正确示例
   data, err := loadData()
   if err != nil {
       return fmt.Errorf("failed to load data: %w", err)
   }
   ```

2. 使用`%w`包装错误以保留错误链
3. 错误信息应该小写开头,不以句号结尾

### 测试要求

1. **单元测试**: 新功能必须包含测试
   ```go
   func TestGetCommand(t *testing.T) {
       // 测试代码
   }
   ```

2. **测试命名**: `Test<FunctionName>`格式
3. **覆盖率**: 新代码应保持≥80%覆盖率
4. **表格驱动测试**: 推荐使用表格驱动测试
   ```go
   func TestAdd(t *testing.T) {
       tests := []struct{
           name     string
           a, b     int
           expected int
       }{
           {"positive", 1, 2, 3},
           {"negative", -1, -2, -3},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               result := Add(tt.a, tt.b)
               if result != tt.expected {
                   t.Errorf("got %d, want %d", result, tt.expected)
               }
           })
       }
   }
   ```

## 命令数据贡献

添加或修改命令数据是最常见的贡献方式。

### YAML格式规范

每个命令应包含以下字段:

```yaml
- name: command-name                # 必填: 命令名称
  category: "类别/子类别"           # 必填: 所属分类
  install_required: true/false      # 必填: 是否需要单独安装
  description: "简短描述"           # 必填: 功能简述
  platforms:                        # 必填: 支持的平台
    - linux
    - darwin
    - windows
  usage:                            # 必填: 使用方式
    - "command [选项] [参数]"
  options:                          # 推荐: 常用选项
    - flag: "-l"
      description: "长格式显示"
  examples:                         # 必填: 使用示例(至少2个)
    - command: "command -l"
      description: "示例说明"
      output: "预期输出(可选)"
  risks:                            # 推荐: 风险说明
    - level: low/medium/high/critical
      description: "风险描述"
  notes:                            # 可选: 注意事项
    - "注意事项1"
  install_method: "安装方法"        # 如果install_required为true则必填
  version_check: "版本检查命令"     # 推荐
  related_commands:                 # 可选: 相关命令
    - related-command
  references:                       # 可选: 参考链接
    - "https://example.com/docs"
```

### 必填字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| name | 命令名称 | `ls`、`docker run` |
| category | 分类路径 | `操作系统/通用Linux命令` |
| install_required | 是否需安装 | `true` / `false` |
| description | 简短描述 | `列出目录内容` |
| platforms | 支持平台 | `linux`, `darwin`, `windows` |
| usage | 使用方式 | `ls [选项] [目录]` |
| examples | 使用示例 | 至少提供2个实际示例 |

### 数据质量要求

1. **准确性**: 所有信息必须准确无误
2. **完整性**: 必填字段必须填写
3. **实用性**: 提供实际有用的示例
4. **风险标注**: 危险命令必须标注风险
5. **格式规范**: 严格遵循YAML格式

### 验证方法

提交前请运行数据验证工具:

```bash
go run ./cmd/validator -d ./data
```

确保:
- ✅ YAML格式正确
- ✅ 所有必填字段存在
- ✅ 风险级别正确
- ✅ 无重复命令

### 提交前检查清单

- [ ] 命令信息准确无误
- [ ] 所有必填字段已填写
- [ ] 提供了至少2个实用示例
- [ ] 危险命令已标注风险
- [ ] YAML格式正确(无语法错误)
- [ ] 运行验证工具通过
- [ ] 文件编码为UTF-8(无BOM)

## 问题报告

使用以下模板报告Bug:

```markdown
### 问题描述
简要描述遇到的问题

### 复现步骤
1. 执行命令 `...`
2. 看到错误 `...`
3. ...

### 预期行为
描述您期望发生什么

### 实际行为
描述实际发生了什么

### 环境信息
- OS: [e.g. Ubuntu 22.04]
- Go版本: [e.g. 1.21.0]
- cmd4coder版本: [e.g. v1.0.0]

### 日志和截图
如果适用,添加相关日志或截图

### 附加信息
其他任何相关信息
```

## 功能请求

使用以下模板提出新功能:

```markdown
### 功能描述
清晰描述您想要的功能

### 使用场景
为什么需要这个功能?它解决什么问题?

### 预期效果
描述功能应该如何工作

### 可能的实现方案
如果您有实现想法,请分享

### 替代方案
是否考虑过其他解决方案?

### 附加信息
其他任何相关信息或参考
```

## 📞 联系方式

如有疑问,可以通过以下方式联系:

- **GitHub Issues**: [创建Issue](https://github.com/cmd4coder/cmd4coder/issues)
- **Discussions**: [参与讨论](https://github.com/cmd4coder/cmd4coder/discussions)

## 🙏 致谢

感谢所有为cmd4coder做出贡献的人!

## 📄 许可证

通过贡献代码,您同意您的贡献将按照项目的[MIT许可证](LICENSE)进行许可。
