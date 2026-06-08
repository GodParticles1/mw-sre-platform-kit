# 对两个现有仓库的落地修改清单

## 1. Python_Go_运维开发_365天学习仓库

建议新增：

```text
docs/projects/mw-sre-platform.md
projects/mwctl-prototype/
projects/evidence-report-generator/
projects/probe-parser-python/
```

README 里新增一段：

```markdown
## 年度主项目：mw-sre-platform

本仓库的 Python/Go 学习最终会收敛到 mw-sre-platform：一个通用中间件 SRE 诊断修复工具。Python 用于日志解析、证据包处理、报告生成和 AIOps glue；Go 用于 mwctl CLI、执行器、插件接口、并发采集和平台化 Agent。
```

## 2. K8s_AI_Ops_Python_Go_CS_DL_365天专家成长仓库

建议新增：

```text
labs/mw-sre-platform/
docs/runbooks/middleware-sre/
docs/architecture/mw-sre-platform.md
```

README 里新增：

```markdown
## Capstone 主项目：Middleware SRE + AIOps Platform

最终项目不是单纯部署 AI 模型，而是建设一个可对接 Kubernetes、Argo Workflows、OpenTelemetry、OPA、Backstage 和 AIOps 的中间件诊断修复平台。AI 推理服务作为平台上的业务场景之一纳入验证。
```

## 3. Roadmap 调整

把阶段 6 中间件学习改成插件开发阶段：

```text
每学一个中间件，必须沉淀：
1. Probe
2. Finding Rule
3. Runbook
4. Verify
5. Rollback
6. 真实案例复盘
```

把阶段 10/11 平台工程提前到阶段 8 开始穿插：

```text
Day 246 之后每周至少做一个平台化能力：
K8s Job、CRD、Argo、OPA、OTel、Backstage、AIOps API。
```
