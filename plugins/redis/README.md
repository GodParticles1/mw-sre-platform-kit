# Redis Plugin

Purpose: diagnose Redis local availability, authentication, master/replica role drift, and slave link health in xRocket/xBrother-style single-node or HA environments.

Lifecycle implemented:

- discover: process/port/config evidence commands
- collect: `INFO replication`, process/port and masked config summary
- diagnose: noauth, port down, both-master, slave link down, auth asymmetry
- plan/apply/verify/rollback: expressed as Runbook-as-Code; apply remains approval-gated and dry-run by default

Safety:

- all plugin collectors are read-only
- passwords are passed through process environment for local execution and masked in config evidence
- no write loop, read loop, or background process is created
