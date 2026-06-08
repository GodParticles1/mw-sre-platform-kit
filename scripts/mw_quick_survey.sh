#!/usr/bin/env bash
set -u
set -o pipefail

MASTER_IP="${MASTER_IP:-10.225.20.20}"
SLAVE_IP="${SLAVE_IP:-10.225.20.21}"
VIP="${VIP:-10.225.20.22}"
SLAVE_HOST="${SLAVE_HOST:-slave}"
REDIS_PASS="${REDIS_PASS:-Xbrother*123}"
MYSQL_USER="${MYSQL_USER:-gj}"
MYSQL_PASS="${MYSQL_PASS:-xbrother}"
TIMEOUT="${TIMEOUT:-3}"

section() { echo; echo "===== $* ====="; }

tcp() {
  local host="$1" port="$2"
  timeout "$TIMEOUT" bash -c "cat < /dev/null > /dev/tcp/$host/$port" \
    && echo "OK   $host:$port" \
    || echo "FAIL $host:$port"
}

section "0. host and vip"
hostname
ip addr | egrep "$MASTER_IP|$SLAVE_IP|$VIP" || true
ssh "$SLAVE_HOST" "hostname; ip addr | egrep '$MASTER_IP|$SLAVE_IP|$VIP' || true" 2>/dev/null || echo "WARN ssh slave failed"

section "1. monit key services master"
monit summary 2>/dev/null | egrep -i 'mariadb|mysql|mongodb|mongo|redis|rabbit|mq|clickhouse|ssdb|xhouse|xbroker.v2|xacquisition|xenvoyproxy' || true

section "2. monit key services slave"
ssh "$SLAVE_HOST" "monit summary 2>/dev/null | egrep -i 'mariadb|mysql|mongodb|mongo|redis|rabbit|mq|clickhouse|ssdb|xhouse|xbroker.v2|xacquisition|xenvoyproxy' || true" 2>/dev/null || true

section "3. tcp vip"
tcp "$VIP" 3306
tcp "$VIP" 5672

section "4. local middleware tcp"
tcp 127.0.0.1 3306
tcp 127.0.0.1 27017
tcp 127.0.0.1 6379
tcp 127.0.0.1 8888
tcp 127.0.0.1 8123

section "5. ports master"
ss -lntup | egrep ':3306|:27017|:6379|:5672|:8123|:9000|:8888|:9700|:5101|:19700|:6000|:6001|:16000|:6700|:26700' || true

section "6. mysql brief"
mysql -e "SELECT @@hostname,@@server_id,@@read_only,@@log_bin; SHOW SLAVE STATUS\G;" 2>&1 \
  | egrep -i 'hostname|server_id|read_only|log_bin|Master_Host|Slave_IO_Running|Slave_SQL_Running|Seconds_Behind_Master|Last_IO_Error|Last_SQL_Error|ERROR|Access denied' || true

section "7. mongo brief"
mongo --quiet --host 127.0.0.1 --port 27017 --eval '
try {
  rs.status().members.forEach(function(m){
    print(m.name + " state=" + m.stateStr + " health=" + m.health + " msg=" + (m.lastHeartbeatMessage || ""));
  });
} catch(e) { print(e); }
' 2>&1 || true

section "8. redis brief"
redis-cli -h 127.0.0.1 -p 6379 -a "$REDIS_PASS" --no-auth-warning INFO replication 2>&1 \
  | egrep 'role:|master_host|master_port|master_link_status|connected_slaves|slave[0-9]|NOAUTH|DENIED' || true
ssh "$SLAVE_HOST" "redis-cli -h 127.0.0.1 -p 6379 -a '$REDIS_PASS' --no-auth-warning INFO replication 2>&1 | egrep 'role:|master_host|master_port|master_link_status|connected_slaves|slave[0-9]|NOAUTH|DENIED' || true" 2>/dev/null || true

section "9. rabbitmq brief"
ps -ef | egrep 'rabbitmq|beam.smp|mqproxy' | grep -v grep || true
ss -lntup | egrep ':5672|:15672|:25672|:4369|:8020' || true

section "10. clickhouse brief"
curl -sS --max-time "$TIMEOUT" 'http://127.0.0.1:8123/?query=SELECT%201' 2>&1 || true
ss -lntup | egrep ':8123|:9000|:9009' || true
if command -v clickhouse-client >/dev/null 2>&1; then
  clickhouse-client --query "SELECT cluster,host_name,host_address,is_local FROM system.clusters ORDER BY cluster,host_name LIMIT 20" 2>&1 || true
fi

section "11. xhouse brief"
ss -lntup | egrep ':9700|:5101|:19700|:8888' || true
tail -n 80 /opt/log/xhouse.log 2>/dev/null | egrep -i 'ERROR|WARN|panic|fatal|ssdb|clickhouse|bucket|refused|timeout|deadline' | tail -30 || true

section "12. south brief"
monit summary 2>/dev/null | egrep -i 'xbroker.v2|xacquisition|xsouth|xpm2' || true
ps -ef | egrep 'xbroker.v2|xacquisition|xsouth|xpm2|mqtt|snmp|xbt|mqproxy' | grep -v grep || true
ss -lntup | egrep ':6000|:6001|:16000|:1883|:6700|:26700|:8020' || true
tail -n 120 /opt/log/xbroker.v2.log 2>/dev/null | egrep -i 'PushValue|PushEvent|success|ERROR|WARN|panic|refused|no route|timeout|deadline' | tail -40 || true

section "13. interpretation hint"
cat <<'HINT'
- VIP absent on both nodes: fix keepalived/xactivestandby first.
- MySQL VIP 3306 failed: fix MariaDB/VIP before xbroker/south.
- Mongo says not running with --replSet: fix slave /etc/mongod.conf replication.replSetName.
- Redis NOAUTH: rerun with correct REDIS_PASS.
- Redis both role:master: set non-VIP node as slaveof VIP-holder physical IP.
- RabbitMQ 5672 failed: mqproxy/acquisition may fail.
- ClickHouse SELECT or cluster view failed: diagnose master startup and slave cluster view.
- SSDB 8888 failed: xhouse may fail.
- xbroker OK + PushValue success: south main link is up; page issue may be cache/check-item/device-specific.
HINT
