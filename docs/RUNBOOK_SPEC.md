# Runbook-as-Code Specification

## 1. 为什么需要 Runbook-as-Code

人工 SOP 难以平台化。Runbook-as-Code 将 SOP 变成可解析、可执行、可验证、可回滚的声明式文件。

## 2. 标准生命周期

```text
precheck -> dryRun -> approval -> backup -> apply -> verify -> rollback -> report
```

## 3. 必须字段

```yaml
apiVersion: sre.middleware/v1alpha1
kind: Runbook
metadata:
  name: redis-make-replica
spec:
  service: redis
  riskLevel: medium
  defaultMode: dryRun
  requiresApproval: true
  precheck: []
  backup: {}
  dryRun: {}
  apply: {}
  verify: []
  rollback: {}
```

## 4. 约束

```text
默认 dry-run。
apply 必须 --apply --yes 或平台审批。
每个 filePatch 前必须 backup。
每个 apply 后必须 verify。
验证失败不能默认自动进行高危 rollback，需按策略决定。
```

## 5. 平台执行状态

```text
Pending
Prechecking
DryRunReady
WaitingApproval
BackingUp
Applying
Verifying
Succeeded
Failed
RollbackReady
RollingBack
RolledBack
```
