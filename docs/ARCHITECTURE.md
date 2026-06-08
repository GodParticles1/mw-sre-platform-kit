# Middleware SRE Platform Architecture

## 1. 最终定位

本项目定位为 **通用中间件诊断与修复平台内核**，不是一次性脚本。它把现场经验沉淀为可版本化、可审计、可回滚、可验证的 Probe 和 Runbook。

## 2. 推荐技术选型

| 层级 | 选型 | 原因 |
|---|---|---|
| 现场 CLI / Agent | Go | 单二进制、并发、SSH/K8s/API 生态好、适合长期维护 |
| 现场兜底 | Shell | 离线可执行、复制方便、只用于 quick survey 和 bootstrap |
| 分析插件 | Python | 日志分析、统计、报告、LLM/RAG glue 更方便 |
| 编排 | Argo Workflows | 多步骤 DAG、并发、重试、K8s 原生执行 |
| 平台控制面 | Kubernetes Operator + CRD | 声明式、可审计、适合平台化 |
| 策略 | OPA/Rego | 高危动作准入、审批规则、租户隔离 |
| 可观测 | OpenTelemetry | 统一日志、指标、链路、事件输出 |
| 门户 | Backstage/自研 Portal | Docs-as-Code、服务目录、Runbook 入口 |
| AIOps | 规则引擎 + 案例检索 + LLM 辅助 | 先确定性规则，再让 AI 解释和推荐，不裸执行高危动作 |

## 3. 逻辑架构

```text
AIOps / Portal / ChatOps
          |
          v
SRE API / Policy / Approval / Evidence Store
          |
          +--> Argo Workflow Engine
          |
          +--> Kubernetes Operator
          |
          v
mwctl / mw-agent Execution Core
          |
          +--> local executor
          +--> ssh executor
          +--> k8s Job executor
          +--> agent executor
          |
          v
Middleware Targets: MySQL / Mongo / Redis / RabbitMQ / ClickHouse / SSDB / Product Services
```

## 4. 横向扩展

横向扩展指框架能力扩展：

```text
executor: local -> ssh -> k8s-job -> agent -> winrm
output: text -> json -> markdown -> evidence bundle -> OpenTelemetry
policy: local confirm -> RBAC -> OPA -> approval workflow
platform: CLI -> API -> Operator -> Portal -> ChatOps
```

## 5. 纵向扩展

纵向扩展指每个中间件持续沉淀现场问题：

```text
MySQL:
  server_id_equal
  binlog_missing_1236
  access_denied
  slave_sql_error
  vip_unreachable

Mongo:
  not_running_with_replset
  member_unhealthy
  keyfile_auth_failed
  startup2_too_long

Redis:
  noauth
  both_master
  link_down
  config_not_persisted

ClickHouse:
  master_down
  slave_cannot_reach_master
  cluster_missing_master
  remote_servers_inconsistent
  macros_duplicate_replica

RabbitMQ:
  port_down
  auth_failed
  queue_backlog
  disk_alarm
```

## 6. 执行生命周期

每个修复动作必须具备：

```text
precheck -> dry-run plan -> backup -> apply -> verify -> rollback -> report
```

禁止直接执行高风险动作。

## 7. 风险分级

| 等级 | 动作 | 默认策略 |
|---|---|---|
| L0 | 只读采集 | 自动允许 |
| L1 | 安全重启 | 可审批后执行 |
| L2 | 配置补丁 | 必须备份 + 审批 |
| L3 | 复制/副本集/集群变更 | 高级审批，默认 dry-run |
| L4 | 数据目录、DROP、重建副本 | 默认禁止自动执行 |

