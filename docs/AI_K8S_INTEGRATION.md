# AIOps and Kubernetes Platform Integration

## 1. 平台对接目标

本工具要能同时服务：

```text
裸机现场 SSH
Kubernetes Job
Kubernetes Operator
Argo Workflow
AIOps 平台
ChatOps/工单系统
Backstage/内部开发者平台
```

## 2. Kubernetes 对接

使用 CRD 抽象：

```text
DiagnosticRun      一次诊断任务
RemediationPlan    一次 dry-run 修复计划
RemediationApply   一次审批后的执行
EvidenceBundle     一次证据包索引
```

Operator 只负责调度、状态回写和审计，不在 controller 里写大量业务排查逻辑。真正执行仍由 mwctl/mw-agent 完成。

## 3. Argo Workflow 对接

诊断天然是 DAG：

```text
collect topology
  -> check mysql/mongo/redis/rabbitmq/clickhouse/ssdb
  -> diagnose
  -> generate plan
  -> approval
  -> apply
  -> verify
  -> report
```

## 4. AIOps 对接

AIOps 不直接执行高危命令，只做：

```text
告警关联
证据解释
相似案例检索
Runbook 推荐
风险分级
报告生成
```

执行必须走：

```text
Runbook Schema -> Policy -> Dry-run -> Approval -> Apply -> Verify -> Audit
```

## 5. OpenTelemetry

每次运行建议输出：

```text
run_id
service
module
duration
result
finding_count
critical_count
action
risk_level
```

未来可以通过 OTLP 上报到可观测平台。

## 6. OPA Policy

示例策略：

```text
禁止自动执行 L4 数据破坏动作。
L3 复制/副本集/ClickHouse replica 变更必须高级审批。
非工作时间禁止 apply，除非 emergency=true。
未生成备份禁止 filePatch。
未定义 verify 禁止 apply。
```
