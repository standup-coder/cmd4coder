# cmd4coder 架构文档

## 系统架构概述

cmd4coder采用清晰的分层架构设计,将应用分为四个主要层次:模型层(Model)、数据层(Data)、服务层(Service)和界面层(UI)。

```
┌─────────────────────────────────────────────────────────┐
│                    用户界面层 (UI)                        │
│  ┌─────────────────┐    ┌─────────────────────────┐     │
│  │   CLI Interface │    │    TUI Interface        │     │
│  │  (Cobra based)  │    │  (Bubbletea based)      │     │
│  └─────────────────┘    └─────────────────────────┘     │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                  业务逻辑层 (Service)                     │
│  ┌──────────────────────────────────────────────┐       │
│  │         Command Service                      │       │
│  │  - 命令查询  - 搜索  - 分类管理               │       │
│  └──────────────────────────────────────────────┘       │
│  ┌──────────────────────────────────────────────┐       │
│  │         Config Service (未来)                 │       │
│  │  - 配置管理  - 用户数据                       │       │
│  └──────────────────────────────────────────────┘       │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                    数据访问层 (Data)                      │
│  ┌────────────┐   ┌────────────┐   ┌─────────────┐     │
│  │   Loader   │   │   Index    │   │   Cache     │     │
│  │ YAML加载器 │   │ 索引管理器  │   │  LRU缓存    │     │
│  └────────────┘   └────────────┘   └─────────────┘     │
└──────────────────────┬──────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                    数据模型层 (Model)                     │
│  ┌────────┐  ┌──────────┐  ┌────────┐  ┌──────────┐    │
│  │Command │  │ Category │  │ Config │  │  Errors  │    │
│  └────────┘  └──────────┘  └────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────┐
│                    数据存储 (Data Files)                  │
│  ┌────────────────────────────────────────────┐         │
│  │  YAML文件 (data/)                          │         │
│  │  - os/           - lang/        - network/       │         │
│  │  - container/    - database/   - vcs/           │         │
│  │  - diagnostic/   - build/                       │         │
│  │                                                  │         │
│  │  容器编排目录 (container/):                      │         │
│  │  - docker.yaml                                   │         │
│  │  - kubernetes.yaml                               │         │
│  │  - k8s-cluster.yaml    (集群管理)               │         │
│  │  - k8s-runtime.yaml    (容器运行时)             │         │
│  │  - k8s-monitor.yaml    (监控日志)               │         │
│  │  - k8s-network.yaml    (网络插件)               │         │
│  │  - k8s-storage.yaml    (存储管理)               │         │
│  │  - k8s-cicd.yaml       (CI/CD)                  │         │
│  │  - k8s-config.yaml     (配置管理)               │         │
│  │  - k8s-backup.yaml     (备份恢复)               │         │
│  │  - k8s-security.yaml   (安全工具)               │         │
│  │  - k8s-utilities.yaml  (辅助工具)               │         │
│  │  - k8s-cloud.yaml      (云平台工具)             │         │
│  │  - k8s-dev.yaml        (开发调试)               │         │
│  └────────────────────────────────────────────┘         │
└─────────────────────────────────────────────────────────┘
```

## 模块职责说明

### 1. 模型层 (internal/model)

**职责**: 定义核心数据结构和业务规则

**主要文件**:
- `command.go`: 命令数据结构、风险级别枚举、验证方法
- `category.go`: 分类和元数据结构
- `config.go`: 配置和用户数据结构
- `errors.go`: 自定义错误类型

**关键设计**:
- 使用结构体标签支持YAML序列化/反序列化
- 实现`Validate()`方法进行数据完整性检查
- 定义风险级别枚举确保数据一致性

### 2. 数据访问层 (internal/data)

**职责**: 数据加载、索引构建、缓存管理

**主要文件**:
- `loader.go`: YAML文件并行加载
- `index.go`: 多级索引系统(名称/分类/关键词/平台)
- `cache.go`: LRU缓存实现

**关键设计**:

#### 2.1 并行加载
```go
// 使用goroutine并行加载多个YAML文件
func LoadAllCommands(dataDir string) ([]*model.Command, error) {
    // 并发加载提高性能
}
```

#### 2.2 多级索引
```
NameIndex (map[string]*Command)     // O(1)精确查找
├─ "ls"       -> Command
├─ "docker"   -> Command
└─ "kubectl"  -> Command

CategoryIndex (map[string][]*Command)
├─ "操作系统/通用Linux命令" -> [Command, ...]
└─ "容器编排/Docker命令"   -> [Command, ...]

KeywordIndex (map[string][]*Command)  // 倒排索引
├─ "文件"     -> [Command, ...]
├─ "网络"     -> [Command, ...]
└─ "容器"     -> [Command, ...]

PlatformIndex (map[string][]*Command)
├─ "linux"    -> [Command, ...]
├─ "darwin"   -> [Command, ...]
└─ "windows"  -> [Command, ...]
```

#### 2.3 LRU缓存
```go
// 使用双向链表+哈希表实现O(1)访问
type LRUCache struct {
    capacity int
    cache    map[string]*list.Element
    list     *list.List
}
```

### 3. 服务层 (internal/service)

**职责**: 业务逻辑处理、命令查询、搜索算法

**主要文件**:
- `command_service.go`: 命令查询服务

**关键功能**:

#### 3.1 搜索算法
4级优先级匹配:
1. 精确匹配 (优先级100): 命令名称完全匹配
2. 前缀匹配 (优先级80): 命令名称前缀匹配
3. 包含匹配 (优先级60): 命令名称或描述包含关键词
4. 关键词匹配 (优先级40): 通过关键词索引匹配

```go
// 搜索结果排序
type SearchResult struct {
    Command  *model.Command
    Score    int  // 匹配分数
}
```

#### 3.2 缓存策略
```go
// 热门查询结果缓存
func (s *CommandService) Search(query string) ([]*model.Command, error) {
    // 1. 检查缓存
    if cached := s.cache.Get(query); cached != nil {
        return cached, nil
    }
    
    // 2. 执行搜索
    results := s.index.Search(query)
    
    // 3. 更新缓存
    s.cache.Put(query, results)
    
    return results, nil
}
```

### 4. 界面层

#### 4.1 CLI界面 (cmd/cli)

**技术栈**: Cobra框架

**命令结构**:
```
cmd4coder
├── list [category]       # 列出命令
├── show <name>           # 显示详情
├── search <query>        # 搜索命令
├── categories            # 列出分类
├── export <name>         # 导出命令
├── export-all            # 导出所有
└── version               # 版本信息
```

**设计原则**:
- 命令简洁明了
- 支持长短参数
- 提供详细帮助信息
- 彩色输出增强可读性

#### 4.2 TUI界面 (internal/ui/tui - 规划中)

**技术栈**: Bubbletea + Lipgloss + Bubbles

**状态管理**:
```go
type Model struct {
    // 视图状态
    viewMode    ViewMode    // 浏览/搜索/详情
    focusPanel  Panel       // 左/中/右
    
    // 数据状态
    categories  []Category
    commands    []*Command
    selected    *Command
    
    // 用户数据
    favorites   []string
    history     []string
    
    // 缓存
    service     *service.CommandService
}
```

**消息处理**:
```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        return m.handleKeyPress(msg)
    case searchResultMsg:
        return m.handleSearchResult(msg)
    // ...
    }
}
```

### 5. 导出功能 (pkg/export)

**职责**: 命令数据导出为Markdown/JSON

**主要文件**:
- `markdown.go`: Markdown格式导出
- `json.go`: JSON格式导出

**导出流程**:
```
1. 获取命令列表 -> 2. 格式化数据 -> 3. 写入文件
```

## 数据流图

### 查询命令流程

```
User Input
    │
    ▼
CLI/TUI Interface
    │
    ▼
Command Service
    │
    ├─> Check Cache ─> Cache Hit ─> Return
    │                      │
    │                   Cache Miss
    │                      │
    ▼                      ▼
Index System          Load from Index
    │                      │
    ├─> Name Index         │
    ├─> Category Index     │
    ├─> Keyword Index      │
    └─> Platform Index     │
           │               │
           ▼               ▼
       Search & Sort   Update Cache
           │               │
           ▼               ▼
        Results ──────> Return
```

### 搜索命令流程

```
Search Query
    │
    ▼
Parse & Tokenize
    │
    ├─> Exact Match (Score: 100)
    ├─> Prefix Match (Score: 80)
    ├─> Contains Match (Score: 60)
    └─> Keyword Match (Score: 40)
    │
    ▼
Merge & Sort by Score
    │
    ▼
Filter & Limit
    │
    ▼
Return Results
```

### 数据加载流程

```
Application Start
    │
    ▼
Load Metadata
    │
    ▼
Scan Data Directory
    │
    ▼
Parallel Load YAML Files
    │
    ├─> Worker 1: os/*.yaml
    ├─> Worker 2: lang/*.yaml
    ├─> Worker 3: network/*.yaml
    └─> Worker N: ...
    │
    ▼
Validate Commands
    │
    ▼
Build Indexes
    │
    ├─> Name Index
    ├─> Category Index
    ├─> Keyword Index
    └─> Platform Index
    │
    ▼
Initialize Cache
    │
    ▼
Ready to Serve
```

## 关键设计决策

### 1. 为什么选择YAML而非JSON/数据库?

**优势**:
- 可读性强,便于人工编辑
- 支持注释,方便说明
- 版本控制友好
- 无需数据库环境
- 加载速度快(本地文件)

**劣势**:
- 不适合大规模数据(>10000条)
- 查询性能不如数据库
- 需要额外的索引层

**结论**: 对于命令工具这种相对静态、规模可控(数百条)的数据,YAML是最佳选择。

### 2. 为什么实现自定义LRU缓存?

**考虑因素**:
- 减少外部依赖
- 完全控制缓存策略
- 轻量级实现满足需求
- 学习和展示缓存实现

**实现细节**:
- 使用`container/list`双向链表
- O(1)访问和更新
- 线程安全(使用sync.Mutex)

### 3. 为什么采用4级搜索优先级?

**设计理念**:
- 精确匹配优先,符合用户期望
- 前缀匹配次之,支持自动补全场景
- 包含匹配提供模糊搜索
- 关键词匹配兜底,扩大召回

**权重设计**:
- 差距明显(100/80/60/40),避免混淆
- 可调整,预留扩展空间

### 4. 为什么分离CLI和TUI?

**优势**:
- 关注点分离
- 独立开发和测试
- 根据场景选择合适界面
- CLI更适合脚本化
- TUI更适合交互式浏览

**实现方式**:
- 共享Service层和Data层
- 不同入口点
- 启动时自动判断模式

## 扩展点说明

### 1. 新增数据源

```go
// 定义数据源接口
type DataSource interface {
    Load() ([]*model.Command, error)
    Reload() error
}

// 实现新的数据源
type RemoteDataSource struct {
    url string
}

func (r *RemoteDataSource) Load() ([]*model.Command, error) {
    // 从远程API加载
}
```

### 2. 自定义索引

```go
// 在Index Manager中添加新索引
func (idx *IndexManager) BuildTagIndex() {
    idx.tagIndex = make(map[string][]*model.Command)
    for _, cmd := range idx.commands {
        for _, tag := range cmd.Tags {
            idx.tagIndex[tag] = append(idx.tagIndex[tag], cmd)
        }
    }
}
```

### 3. 插件系统(未来)

```go
// 插件接口
type Plugin interface {
    Name() string
    Init(*service.CommandService) error
    Execute(args []string) error
}

// 插件管理器
type PluginManager struct {
    plugins map[string]Plugin
}
```

### 4. 导出格式扩展

```go
// 导出接口
type Exporter interface {
    Export(commands []*model.Command, output string) error
}

// 新增HTML导出
type HTMLExporter struct{}

func (h *HTMLExporter) Export(commands []*model.Command, output string) error {
    // 生成HTML
}
```

## 性能优化

### 已实现的优化

1. **并行加载**: 使用goroutine并行加载YAML文件
2. **索引预构建**: 启动时构建索引,查询时O(1)访问
3. **LRU缓存**: 缓存热门查询结果
4. **分词优化**: 简化的分词算法,降低复杂度

### 可能的优化方向

1. **延迟加载**: 按需加载分类数据
2. **增量索引**: 支持增量更新而非全量重建
3. **并行搜索**: 并行搜索多个索引
4. **结果缓存预热**: 启动时预加载热门查询
5. **二进制缓存**: 预编译数据为二进制格式

## 测试策略

### 单元测试

- 模型层: 数据验证逻辑测试
- 数据层: 加载器、索引、缓存独立测试
- 服务层: 查询、搜索逻辑测试
- 导出层: 格式正确性测试

### 集成测试

- 完整数据流测试
- CLI命令端到端测试
- 性能基准测试

### 测试覆盖率目标

- 总体: ≥80%
- 核心模块: ≥85%
- UI层: ≥50%

## 部署和分发

### 构建

```bash
# 单平台构建
go build -o bin/cmd4coder ./cmd/cli

# 多平台构建
GOOS=linux GOARCH=amd64 go build -o bin/cmd4coder-linux-amd64 ./cmd/cli
GOOS=darwin GOARCH=amd64 go build -o bin/cmd4coder-darwin-amd64 ./cmd/cli
GOOS=windows GOARCH=amd64 go build -o bin/cmd4coder-windows-amd64.exe ./cmd/cli
```

### 发布

1. 创建Git标签
2. GitHub Actions自动构建多平台二进制
3. 创建GitHub Release
4. 上传构建产物

## 总结

cmd4coder采用清晰的分层架构,各层职责明确,模块间解耦良好。通过合理的索引设计和缓存策略,在保持简洁的同时提供了良好的性能。架构设计预留了足够的扩展点,便于未来功能增强。

# 多平台构建
GOOS=linux GOARCH=amd64 go build -o bin/cmd4coder-linux-amd64 ./cmd/cli
GOOS=darwin GOARCH=amd64 go build -o bin/cmd4coder-darwin-amd64 ./cmd/cli
GOOS=windows GOARCH=amd64 go build -o bin/cmd4coder-windows-amd64.exe ./cmd/cli
```

### 发布

1. 创建Git标签
2. GitHub Actions自动构建多平台二进制
3. 创建GitHub Release
4. 上传构建产物

## 总结

cmd4coder采用清晰的分层架构,各层职责明确,模块间解耦良好。通过合理的索引设计和缓存策略,在保持简洁的同时提供了良好的性能。架构设计预留了足够的扩展点,便于未来功能增强。
