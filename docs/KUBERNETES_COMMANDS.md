# Kubernetes生态工具命令清单

本文档汇总了cmd4coder项目中新增的Kubernetes生态系统工具命令。

## 版本信息
- **发布版本**: 1.1.0
- **发布日期**: 2026-01-07
- **新增分类**: 12个
- **新增命令**: 128条

---

## 目录

- [1. 集群管理工具 (12条)](#1-集群管理工具)
- [2. 容器运行时工具 (9条)](#2-容器运行时工具)
- [3. 监控与日志工具 (11条)](#3-监控与日志工具)
- [4. 网络插件工具 (7条)](#4-网络插件工具)
- [5. 存储管理工具 (8条)](#5-存储管理工具)
- [6. CI/CD工具 (11条)](#6-cicd工具)
- [7. 配置管理工具 (7条)](#7-配置管理工具)
- [8. 备份恢复工具 (8条)](#8-备份恢复工具)
- [9. 安全工具 (7条)](#9-安全工具)
- [10. 辅助工具 (5条)](#10-辅助工具)
- [11. 云平台工具 (9条)](#11-云平台工具)
- [12. 开发调试工具 (8条)](#12-开发调试工具)

---

## 1. 集群管理工具

**分类**: Kubernetes集群管理和控制平面工具  
**数据文件**: `container/k8s-cluster.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| kubeadm init | 初始化Kubernetes控制平面节点 | critical | linux |
| kubeadm join | 将工作节点加入集群 | high | linux |
| kubeadm upgrade | 升级Kubernetes集群版本 | critical | linux |
| kubeadm reset | 重置节点状态 | critical | linux |
| kubeadm token list | 列出集群加入令牌 | low | linux |
| kubeadm config view | 查看集群配置 | low | linux |
| systemctl status kubelet | 检查kubelet服务状态 | low | linux |
| journalctl -u kubelet | 查看kubelet日志 | low | linux |
| etcdctl snapshot save | 创建etcd数据备份 | low | linux |
| etcdctl snapshot restore | 从备份恢复etcd集群 | critical | linux |
| etcdctl member list | 列出etcd集群成员 | low | linux |
| etcdctl get | 获取etcd键值 | low | linux |

---

## 2. 容器运行时工具

**分类**: Kubernetes容器运行时管理  
**数据文件**: `container/k8s-runtime.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| crictl images | 列出容器镜像 | low | linux |
| crictl pods | 列出所有Pod | low | linux |
| crictl ps | 列出运行中的容器 | low | linux |
| crictl logs | 获取容器日志 | low | linux |
| crictl exec | 在容器中执行命令 | high | linux |
| crictl stats | 显示容器资源使用统计 | low | linux |
| ctr images list | 列出containerd镜像 | low | linux |
| ctr containers list | 列出containerd容器 | low | linux |
| systemctl status containerd | 检查containerd服务状态 | low | linux |

---

## 3. 监控与日志工具

**分类**: Kubernetes监控和日志收集  
**数据文件**: `container/k8s-monitor.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| promtool check config | 验证Prometheus配置 | low | linux/darwin/windows |
| systemctl status prometheus | 检查Prometheus服务状态 | low | linux |
| systemctl status grafana-server | 检查Grafana服务状态 | low | linux |
| grafana-cli plugins list-remote | 列出可用Grafana插件 | low | linux/darwin/windows |
| grafana-cli plugins install | 安装Grafana插件 | medium | linux/darwin/windows |
| systemctl status loki | 检查Loki服务状态 | low | linux |
| systemctl status promtail | 检查Promtail服务状态 | low | linux |
| fluentd -c | 启动Fluentd | medium | linux/darwin/windows |
| systemctl status fluentd | 检查Fluentd服务状态 | low | linux |
| fluent-bit -c | 启动Fluent Bit | medium | linux/darwin/windows |
| systemctl status kube-state-metrics | 检查kube-state-metrics状态 | low | linux |

---

## 4. 网络插件工具

**分类**: Kubernetes网络插件管理  
**数据文件**: `container/k8s-network.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| calicoctl get nodes | 列出Calico节点 | low | linux/darwin/windows |
| calicoctl get networkpolicies | 列出网络策略 | low | linux/darwin/windows |
| calicoctl apply | 应用网络策略 | high | linux/darwin/windows |
| calicoctl delete networkpolicy | 删除网络策略 | critical | linux/darwin/windows |
| cilium status | 检查Cilium代理状态 | low | linux |
| cilium connectivity test | 运行Pod间连接测试 | low | linux |
| cilium hubble status | 检查Hubble可观测性状态 | low | linux |

---

## 5. 存储管理工具

**分类**: Kubernetes存储和Helm包管理  
**数据文件**: `container/k8s-storage.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| helm repo add | 添加Helm仓库 | low | linux/darwin/windows |
| helm repo update | 更新Helm仓库索引 | low | linux/darwin/windows |
| helm install | 安装Helm chart | high | linux/darwin/windows |
| helm upgrade | 升级Helm release | high | linux/darwin/windows |
| helm uninstall | 卸载Helm release | critical | linux/darwin/windows |
| helm list | 列出Helm releases | low | linux/darwin/windows |
| helm status | 显示release状态 | low | linux/darwin/windows |
| helm template | 本地渲染chart模板 | low | linux/darwin/windows |

---

## 6. CI/CD工具

**分类**: Kubernetes持续集成与部署  
**数据文件**: `container/k8s-cicd.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| argocd login | 登录ArgoCD服务器 | low | linux/darwin/windows |
| argocd app create | 创建ArgoCD应用 | high | linux/darwin/windows |
| argocd app list | 列出ArgoCD应用 | low | linux/darwin/windows |
| argocd app sync | 同步应用状态 | high | linux/darwin/windows |
| argocd app delete | 删除ArgoCD应用 | critical | linux/darwin/windows |
| flux bootstrap git | 引导Flux到Git仓库 | high | linux/darwin/windows |
| flux get kustomizations | 列出Kustomization资源 | low | linux/darwin/windows |
| flux reconcile source git | 触发Git源同步 | medium | linux/darwin/windows |
| tkn pipeline list | 列出Tekton pipelines | low | linux/darwin/windows |
| tkn pipeline start | 启动pipeline执行 | high | linux/darwin/windows |
| tkn pipelinerun logs | 查看pipeline运行日志 | low | linux/darwin/windows |

---

## 7. 配置管理工具

**分类**: 配置管理与IaC  
**数据文件**: `container/k8s-config.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| ansible-playbook | 执行Ansible playbook | high | linux/darwin/windows |
| ansible all -m ping | 测试主机连通性 | low | linux/darwin/windows |
| terraform init | 初始化Terraform工作目录 | low | linux/darwin/windows |
| terraform plan | 生成执行计划 | low | linux/darwin/windows |
| terraform apply | 应用配置变更 | critical | linux/darwin/windows |
| terraform destroy | 销毁基础设施 | critical | linux/darwin/windows |
| terraform output | 显示输出值 | low | linux/darwin/windows |

---

## 8. 备份恢复工具

**分类**: 灾难恢复与备份  
**数据文件**: `container/k8s-backup.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| velero backup create | 创建Kubernetes资源备份 | low | linux/darwin/windows |
| velero backup list | 列出所有备份 | low | linux/darwin/windows |
| velero restore create | 从备份恢复 | critical | linux/darwin/windows |
| velero schedule create | 创建定时备份 | low | linux/darwin/windows |
| restic init | 初始化restic仓库 | low | linux/darwin/windows |
| restic backup | 创建restic备份 | low | linux/darwin/windows |
| restic restore | 从restic备份恢复 | high | linux/darwin/windows |
| restic snapshots | 列出备份快照 | low | linux/darwin/windows |

---

## 9. 安全工具

**分类**: 安全扫描与合规检查  
**数据文件**: `container/k8s-security.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| trivy image | 扫描容器镜像漏洞 | low | linux/darwin/windows |
| trivy fs | 扫描文件系统 | low | linux/darwin/windows |
| trivy k8s | 扫描Kubernetes集群 | low | linux/darwin/windows |
| trivy config | 扫描配置文件 | low | linux/darwin/windows |
| kube-bench run | 运行CIS基准检查 | low | linux |
| systemctl status falco | 检查Falco服务状态 | low | linux |
| falco | 启动Falco运行时监控 | low | linux |

---

## 10. 辅助工具

**分类**: Kubernetes集群管理辅助工具  
**数据文件**: `container/k8s-utilities.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| k9s | 交互式Kubernetes终端UI | medium | linux/darwin/windows |
| kubectx | 切换Kubernetes上下文 | high | linux/darwin/windows |
| kubens | 切换Kubernetes命名空间 | medium | linux/darwin/windows |
| stern | 多Pod日志聚合 | low | linux/darwin/windows |
| popeye | 集群资源健康检查 | low | linux/darwin/windows |

---

## 11. 云平台工具

**分类**: 云平台专属Kubernetes工具  
**数据文件**: `container/k8s-cloud.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| eksctl create cluster | 创建AWS EKS集群 | critical | linux/darwin/windows |
| eksctl get cluster | 列出EKS集群 | low | linux/darwin/windows |
| eksctl delete cluster | 删除EKS集群 | critical | linux/darwin/windows |
| az aks create | 创建Azure AKS集群 | critical | linux/darwin/windows |
| az aks list | 列出AKS集群 | low | linux/darwin/windows |
| az aks delete | 删除AKS集群 | critical | linux/darwin/windows |
| gcloud container clusters create | 创建GKE集群 | critical | linux/darwin/windows |
| gcloud container clusters list | 列出GKE集群 | low | linux/darwin/windows |
| gcloud container clusters delete | 删除GKE集群 | critical | linux/darwin/windows |

---

## 12. 开发调试工具

**分类**: Kubernetes开发与调试  
**数据文件**: `container/k8s-dev.yaml`

| 命令名称 | 用途 | 风险等级 | 平台支持 |
|---------|------|---------|---------|
| skaffold dev | 启动持续开发模式 | medium | linux/darwin/windows |
| skaffold run | 一次性构建和部署 | high | linux/darwin/windows |
| skaffold build | 仅构建镜像 | low | linux/darwin/windows |
| skaffold delete | 删除部署 | high | linux/darwin/windows |
| tilt up | 启动Tilt开发环境 | medium | linux/darwin/windows |
| tilt down | 停止Tilt并清理 | medium | linux/darwin/windows |
| telepresence connect | 连接到Kubernetes集群 | medium | linux/darwin/windows |
| telepresence intercept | 拦截服务流量到本地 | critical | linux/darwin/windows |

---

## 使用说明

### 查询命令

```bash
# 查看所有Kubernetes生态工具分类
cmd4coder categories | grep K8s

# 列出集群管理命令
cmd4coder list "容器编排/K8s集群管理"

# 搜索Helm相关命令
cmd4coder search helm

# 查看具体命令详情
cmd4coder show "helm install"
```

### 风险等级说明

- **low**: 只读操作，无风险
- **medium**: 修改配置但可逆
- **high**: 影响运行状态或数据
- **critical**: 可能导致服务中断或数据丢失

### 平台支持

- **linux**: Linux系统
- **darwin**: macOS系统
- **windows**: Windows系统

---

## 更新日志

### v1.1.0 (2026-01-07)
- ✅ 新增12个Kubernetes生态工具分类
- ✅ 新增128条精选命令
- ✅ 覆盖集群管理、监控、网络、存储、CI/CD、安全等全栈工具
- ✅ 每条命令包含详细使用说明、选项、示例和风险评估

---

## 贡献指南

如需添加或更新Kubernetes工具命令，请参考 [CONTRIBUTING.md](../CONTRIBUTING.md)。

---

**版权所有 © 2026 cmd4coder项目**
