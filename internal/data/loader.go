package data

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cmd4coder/cmd4coder/internal/model"
	"gopkg.in/yaml.v3"
)

// Loader 数据加载器
type Loader struct {
	dataDir  string
	metadata *model.Metadata
	mu       sync.RWMutex
}

// NewLoader 创建数据加载器
func NewLoader(dataDir string) *Loader {
	return &Loader{
		dataDir: dataDir,
	}
}

// LoadMetadata 加载元数据
func (l *Loader) LoadMetadata() (*model.Metadata, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	metadataPath := filepath.Join(l.dataDir, "metadata.yaml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, model.ErrDataLoadFailed{File: metadataPath, Err: err}
	}

	var metadata model.Metadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, model.ErrDataLoadFailed{File: metadataPath, Err: err}
	}

	if err := metadata.Validate(); err != nil {
		return nil, err
	}

	l.metadata = &metadata
	return &metadata, nil
}

// LoadCommandList 加载单个命令列表文件
func (l *Loader) LoadCommandList(filePath string) (*model.CommandList, error) {
	fullPath := filepath.Join(l.dataDir, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, model.ErrDataLoadFailed{File: fullPath, Err: err}
	}

	var cmdList model.CommandList
	if err := yaml.Unmarshal(data, &cmdList); err != nil {
		return nil, model.ErrDataLoadFailed{File: fullPath, Err: err}
	}

	if err := cmdList.Validate(); err != nil {
		return nil, fmt.Errorf("validation errors in %s: %v", fullPath, err)
	}

	return &cmdList, nil
}

// LoadAllCommands 加载所有命令
func (l *Loader) LoadAllCommands() ([]*model.Command, error) {
	// 先加载元数据
	metadata, err := l.LoadMetadata()
	if err != nil {
		return nil, err
	}

	var allCommands []*model.Command
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, len(metadata.DataFiles))

	// 并行加载所有数据文件
	for _, dataFile := range metadata.DataFiles {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()

			cmdList, err := l.LoadCommandList(file)
			if err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			allCommands = append(allCommands, cmdList.Commands...)
			mu.Unlock()
		}(dataFile)
	}

	wg.Wait()
	close(errCh)

	// 检查是否有错误
	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return allCommands, nil
}

// GetMetadata 获取元数据
func (l *Loader) GetMetadata() *model.Metadata {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.metadata
}
