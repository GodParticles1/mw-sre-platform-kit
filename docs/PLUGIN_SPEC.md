# Plugin Specification

## 1. 插件的 7 个接口

每个中间件插件都应具备：

```text
discover    发现服务是否存在
collect     采集证据
diagnose    从证据生成 finding
plan        生成修复计划
apply       执行修复动作
verify      验证是否恢复
rollback    从备份恢复
```

## 2. Probe 输出标准

```json
{
  "service": "redis",
  "probe": "redis-replication",
  "target": "slave",
  "raw": {},
  "parsed": {
    "role": "slave",
    "master_link_status": "up"
  },
  "findings": []
}
```

## 3. Finding 标准

```json
{
  "service": "redis",
  "rule_id": "redis.both_master",
  "severity": "critical",
  "summary": "Both Redis nodes are master",
  "evidence": ["master role=master", "slave role=master"],
  "recommendation": "Make the non-VIP node replicate from the VIP-holder physical IP."
}
```

## 4. Runbook 标准

Runbook 必须包含：

```text
precheck
backup
dryRun
apply
verify
rollback
riskLevel
requiresApproval
```

## 5. 插件目录规范

```text
plugins/<service>/
  probes/
  rules/
  runbooks/
  tests/
  README.md
```

## 6. 风险要求

```text
L0 只读：允许自动执行。
L1 重启：需要确认。
L2 配置修改：必须备份 + 审批。
L3 复制/副本集/集群：高级审批。
L4 数据破坏性操作：默认禁止。
```
