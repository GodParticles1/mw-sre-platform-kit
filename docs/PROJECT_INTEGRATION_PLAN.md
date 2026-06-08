# 与现有学习项目的对接方案

## 1. 当前两个仓库现状

你上传的两个仓库分别适合：

1. `Python_Go_运维开发_365天学习仓库`：适合做语言、脚本、Go CLI、运维开发基本功训练。
2. `K8s_AI_Ops_Python_Go_CS_DL_365天专家成长仓库`：适合做 K8s、AI Ops、平台工程、可观测性、Operator、Capstone。

当前缺口：

```text
缺少一个贯穿两个仓库的长期主项目。
缺少平台化工具的架构骨架。
缺少 Probe-as-Code / Runbook-as-Code 规范。
缺少现场案例沉淀到插件的机制。
```

## 2. 建议新增第三个主仓库

```text
mw-sre-platform
```

它是你的长期作品集主项目。两个学习仓库作为学习路线和练习支撑，`mw-sre-platform` 作为最终产物。

## 3. 目录对接

### Python/Go 运维开发仓库

新增：

```text
docs/projects/mw-sre-platform-roadmap.md
projects/mwctl-prototype/            # Go CLI 原型
projects/probe-parser-python/        # Python 解析器实验
projects/evidence-report-generator/  # 证据包报告生成器
```

学习重点：

```text
Python:
  日志解析、报告生成、规则匹配、证据包处理。
Go:
  CLI、executor、并发采集、JSON 输出、插件接口。
```

### K8s/AI Ops 专家仓库

新增：

```text
labs/mw-sre-platform/
  operator-demo/
  argo-workflow-demo/
  otel-evidence-demo/
  opa-policy-demo/
  aiops-diagnosis-demo/

docs/runbooks/middleware-sre/
  mysql.md
  mongo.md
  redis.md
  rabbitmq.md
  clickhouse.md
  xhouse.md
  south.md
```

学习重点：

```text
K8s Operator
Argo Workflows
OpenTelemetry
OPA policy
AIOps rule engine
LLM report generation
```

## 4. 最小可行版本 MVP

MVP 不追求自动修复所有问题，只完成：

```text
1. mwctl quick
2. mwctl collect
3. mwctl check mysql/mongo/redis/clickhouse/rabbitmq
4. evidence.json
5. report.md
6. Redis 和 Mongo 两个低中风险 runbook 的 dry-run/apply/verify/rollback
```

## 5. 12 周项目里程碑

| 周 | 目标 | 产物 |
|---|---|---|
| 1 | 把这次现场复盘转成标准插件模型 | ARCHITECTURE.md、PLUGIN_SPEC.md |
| 2 | Go CLI 框架 | mwctl quick/check/collect |
| 3 | Evidence Bundle | evidence.json、report.md |
| 4 | Redis 插件 | redis check、finding、verify |
| 5 | Mongo 插件 | mongo check、finding、verify |
| 6 | MySQL 插件 | mysql check、server_id/binlog 诊断 |
| 7 | RabbitMQ / ClickHouse 插件 | rabbitmq check、clickhouse cluster diagnose |
| 8 | Runbook-as-Code | redis/mongo dry-run/apply/rollback |
| 9 | Web/API 原型 | FastAPI 或 Go HTTP API |
| 10 | K8s Job / Argo Workflow | diagnostic workflow |
| 11 | Operator CRD 原型 | DiagnosticRun CRD |
| 12 | AIOps 报告器 | 规则引擎 + Markdown 报告 + 案例库 |

## 6. 不建议的路线

```text
不建议继续堆一个超大 shell。
不建议一开始就做完整 Operator。
不建议让 LLM 直接执行修复命令。
不建议把现场 IP/密码写死到核心代码。
```

## 7. 推荐路线

```text
先 Go CLI。
再 YAML Probe/Runbook。
再 Evidence Bundle。
再 API/Operator/Argo。
最后接 AIOps。
```
