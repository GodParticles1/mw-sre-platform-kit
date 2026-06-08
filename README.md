# mw-sre-platform-kit

这是一个面向 **通用中间件诊断/修复平台** 的项目对接包，目标是把现场 SOP 提炼为可复用、可审计、可验证、可回滚的 SRE 平台工具。

它不是一个只服务某次故障的脚本集合，而是按以下方向设计：

- **横向扩展**：在同一框架里扩展更多插件、更多执行器、更多输出格式、更多平台接口。
- **纵向扩展**：每个中间件不断沉淀现场问题，形成更多探针、诊断规则、Runbook、验证器和回滚器。
- **平台对接**：保留和 AIOps、Kubernetes、Argo Workflows、Backstage、OpenTelemetry、OPA 等平台能力对接的架构接口。
- **现场落地**：保留 SSH/裸机/离线环境可运行的 CLI 和快速普查脚本。

## 推荐最终形态

```text
Go CLI/Agent + Probe-as-Code + Runbook-as-Code + Evidence Bundle + K8s Operator + Argo Workflow + AIOps Policy/Reasoning
```

## 为什么不是只用 Shell

Shell 适合现场兜底，但不适合长期维护复杂状态机、回滚、证据包、平台 API、并发采集和插件体系。这个 kit 里保留了 `scripts/mw_quick_survey.sh` 作为现场兜底；主线代码以 Go CLI `mwctl` 为骨架。

## 目录

```text
cmd/mwctl/                  Go CLI 入口
internal/core/              统一数据结构
internal/executor/          local/ssh 执行器抽象
internal/checks/            检查器骨架
internal/report/            报告输出
profiles/                   现场 Profile，例如 xbrother-ha
probes/                     Probe-as-Code 样例
runbooks/                   Runbook-as-Code 样例
platform/crds/              Kubernetes CRD 样例
platform/argo/              Argo Workflow 样例
platform/backstage/         Backstage catalog 样例
scripts/                    快速普查脚本
learning-patches/           对现有 365 天学习仓库的调整建议
docs/                       架构、集成、学习计划、规范文档
```

## 快速验证

```bash
make test
bash -n scripts/mw_quick_survey.sh
```

## 当前状态

这是 v0.1 骨架，可作为你现有两个学习仓库的 Capstone 主项目引入。后续建议把它放到独立仓库：

```text
mw-sre-platform/
```

然后在两个学习仓库里把它作为全年主线项目贯穿。
