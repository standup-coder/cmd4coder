# Changelog

All notable changes to cmd4coder project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-01-07

### Added
- **12个Kubernetes生态工具分类**
  - K8s集群管理 (container_k8s_cluster)
  - K8s容器运行时 (container_k8s_runtime)
  - K8s监控日志 (container_k8s_monitor)
  - K8s网络插件 (container_k8s_network)
  - K8s存储管理 (container_k8s_storage)
  - K8s持续集成 (container_k8s_cicd)
  - K8s配置管理 (container_k8s_config)
  - K8s备份恢复 (container_k8s_backup)
  - K8s安全工具 (container_k8s_security)
  - K8s辅助工具 (container_k8s_utilities)
  - K8s云平台工具 (container_k8s_cloud)
  - K8s开发调试 (container_k8s_dev)

- **128条Kubernetes生态工具命令**
  - 集群管理工具：kubeadm init/join/upgrade, kubelet, etcdctl等 (12条)
  - 容器运行时：crictl, ctr, containerd等 (9条)
  - 监控日志：prometheus, grafana, loki, promtail, fluentd, fluent-bit等 (11条)
  - 网络插件：calicoctl, cilium等 (7条)
  - 存储管理：helm repo/install/upgrade/uninstall等 (8条)
  - CI/CD工具：argocd, flux, tekton等 (11条)
  - 配置管理：ansible-playbook, terraform等 (7条)
  - 备份恢复：velero, restic等 (8条)
  - 安全工具：trivy, kube-bench, falco等 (7条)
  - 辅助工具：k9s, kubectx, kubens, stern, popeye等 (5条)
  - 云平台工具：eksctl, az aks, gcloud container clusters等 (9条)
  - 开发调试：skaffold, tilt, telepresence等 (8条)

- **12个新的YAML数据文件**
  - container/k8s-cluster.yaml
  - container/k8s-runtime.yaml
  - container/k8s-monitor.yaml
  - container/k8s-network.yaml
  - container/k8s-storage.yaml
  - container/k8s-cicd.yaml
  - container/k8s-config.yaml
  - container/k8s-backup.yaml
  - container/k8s-security.yaml
  - container/k8s-utilities.yaml
  - container/k8s-cloud.yaml
  - container/k8s-dev.yaml

### Changed
- 更新metadata.yaml版本号从1.0.0到1.1.0
- 更新README.md，反映新增的Kubernetes生态工具
- 命令总数从220条增加到350+条
- 分类总数从20+个增加到32个

### Improved
- 完善了Kubernetes运维工具链的覆盖范围
- 为每条命令提供了详细的使用说明、选项、示例和风险评估
- 所有命令包含安装方法和版本检查命令

### Documentation
- 更新核心文档(README.md)
- 新增Kubernetes工具专题说明
- 更新版本号和命令统计信息

## [1.0.0] - 2025-12-14

### Added
- 初始版本发布
- 220个精选命令
- 20+个命令分类
- CLI和TUI双模式交互
- 命令搜索和查询功能
- Markdown和JSON导出功能
- 配置管理和历史记录
- 跨平台支持（Linux、macOS、Windows）

### Features
- 4级优先级搜索算法
- LRU缓存优化
- 多维度命令索引
- 风险等级标注
- 详细的使用示例
