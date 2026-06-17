# OpenGauss Plugin

Purpose: diagnose OpenGauss process/port presence, `gs_ctl query` HA state, common data-directory path mistakes, permission problems, disk-full symptoms, and standby receiver loss.

Lifecycle implemented:

- discover: common root and non-root data directory candidates
- collect: `gs_ctl query`, process/port summary, bounded log tail
- diagnose: data dir missing, gs_ctl missing, HA not normal, standby receiver missing, disk full, permission denied, unhealthy startup
- plan/apply/verify/rollback: Runbook-as-Code only; no automatic database mutation is enabled

Safety:

- collectors are read-only and bounded with `head`/`tail`
- no recursive write, no data directory chown, no restart, no replication reset
- paths support `/opt/...`, `$HOME/opt/...`, and `$RUNTIME_ROOTFS/opt/...`
