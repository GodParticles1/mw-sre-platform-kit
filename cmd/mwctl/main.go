package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"mw-sre-platform/internal/checks"
	"mw-sre-platform/internal/core"
	"mw-sre-platform/internal/executor"
	baseplugins "mw-sre-platform/internal/plugins"
	"mw-sre-platform/internal/plugins/registry"
	"mw-sre-platform/internal/report"
)

type config struct {
	profile          string
	module           string
	target           string
	redisPass        string
	openGaussDataDir string
	jsonOut          bool
	timeout          time.Duration
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	cfg := config{profile: "generic", module: "all", target: "local", timeout: 5 * time.Second}
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	fs.StringVar(&cfg.profile, "profile", cfg.profile, "profile name, e.g. generic or xbrother-ha")
	fs.StringVar(&cfg.module, "module", cfg.module, "module to check: all,mysql,mongo,redis,opengauss,python-env,south,clickhouse")
	fs.StringVar(&cfg.target, "target", cfg.target, "target host: local or ssh host")
	fs.StringVar(&cfg.redisPass, "redis-pass", os.Getenv("REDIS_PASS"), "redis password")
	fs.StringVar(&cfg.openGaussDataDir, "opengauss-data-dir", os.Getenv("OPENGAUSS_DATA_DIR"), "OpenGauss data directory override")
	fs.BoolVar(&cfg.jsonOut, "json", false, "print JSON report")
	fs.DurationVar(&cfg.timeout, "timeout", cfg.timeout, "command timeout")

	switch cmd {
	case "quick", "check", "collect":
		_ = fs.Parse(os.Args[2:])
		r := runChecks(cfg)
		if cfg.jsonOut || cmd == "collect" {
			if err := report.WriteJSON(os.Stdout, r); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			return
		}
		printHuman(r)
	case "fix":
		fmt.Println("dry-run only in v0.1: use runbooks/*.yaml as the planned remediation contract")
		fmt.Println("repair actions must implement backup, apply, verify and rollback before being enabled")
	case "version":
		fmt.Println("mwctl v0.1.0")
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`Usage:
  mwctl quick --profile xbrother-ha --module all
  mwctl check --module redis --json
  mwctl check --module opengauss --opengauss-data-dir /opt/data/opengauss/data/dn
  mwctl check --module python-env --json
  mwctl collect --profile xbrother-ha > evidence.json
  mwctl fix --action redis-replica   # dry-run placeholder in v0.1`)
}

func runChecks(cfg config) core.CheckReport {
	started := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	exec := executor.NewLocalExecutor()
	commands := plannedCommands(cfg)
	var results []core.CommandResult
	for _, c := range commands {
		if cfg.target != "local" && cfg.target != "localhost" {
			c.Host = cfg.target
			ssh := executor.NewSSHExecutor()
			results = append(results, ssh.Run(ctx, c))
			continue
		}
		c.Host = "local"
		results = append(results, exec.Run(ctx, c))
	}

	var findings []core.Finding
	module := strings.ToLower(cfg.module)
	if module == "all" || module == "mysql" {
		findings = append(findings, checks.MySQLFindings(results)...)
	}
	if module == "all" || module == "mongo" {
		findings = append(findings, checks.MongoFindings(results)...)
	}
	for _, p := range registry.Default() {
		if baseplugins.Match(p, module) {
			findings = append(findings, p.Findings(results)...)
		}
	}
	ended := time.Now()
	return core.CheckReport{RunID: ended.Format("20060102_150405"), Profile: cfg.profile, Target: cfg.target, Module: cfg.module, StartedAt: started, EndedAt: ended, Commands: results, Findings: findings}
}

func pluginRuntimeConfig(cfg config) baseplugins.RuntimeConfig {
	return baseplugins.RuntimeConfig{
		Profile:          cfg.profile,
		Target:           cfg.target,
		Timeout:          cfg.timeout,
		RedisPass:        cfg.redisPass,
		OpenGaussDataDir: cfg.openGaussDataDir,
	}
}

func plannedCommands(cfg config) []executor.Command {
	module := strings.ToLower(cfg.module)
	var out []executor.Command
	add := func(name string, args ...string) {
		out = append(out, executor.Command{Name: name, Args: args, Timeout: int(cfg.timeout.Seconds())})
	}

	if module == "all" || module == "mysql" {
		add("mysql-slave-status", "bash", "-lc", `mysql -e "SELECT @@hostname,@@server_id,@@read_only,@@log_bin; SHOW SLAVE STATUS\\G;" 2>&1 || true`)
	}
	if module == "all" || module == "mongo" {
		add("mongo-rs-status", "bash", "-lc", `mongo --quiet --host 127.0.0.1 --port 27017 --eval 'try { rs.status().members.forEach(function(m){ print(m.name+" state="+m.stateStr+" health="+m.health+" msg="+(m.lastHeartbeatMessage||"")); }) } catch(e) { print(e); }' 2>&1 || true`)
	}
	pluginCfg := pluginRuntimeConfig(cfg)
	for _, p := range registry.Default() {
		if baseplugins.Match(p, module) {
			out = append(out, p.Commands(pluginCfg)...)
		}
	}

	if module == "all" || module == "south" {
		add("south-status", "bash", "-lc", `monit summary 2>/dev/null | egrep -i 'xbroker.v2|xacquisition|xsouth|xpm2' || true; ss -lntup | egrep ':6000|:6001|:16000|:6700|:26700' || true`)
		add("south-log", "bash", "-lc", `tail -n 120 /opt/log/xbroker.v2.log 2>/dev/null | egrep -i 'PushValue|PushEvent|success|ERROR|WARN|panic|refused|timeout' | tail -40 || true`)
	}
	if module == "all" || module == "clickhouse" {
		add("clickhouse-select", "bash", "-lc", `curl -sS --max-time 3 'http://127.0.0.1:8123/?query=SELECT%201' 2>&1 || true`)
		add("clickhouse-cluster", "bash", "-lc", `clickhouse-client --query "SELECT cluster,host_name,host_address,is_local FROM system.clusters ORDER BY cluster,host_name LIMIT 20" 2>&1 || true`)
	}
	return out
}

func printHuman(r core.CheckReport) {
	fmt.Printf("run_id=%s profile=%s module=%s target=%s\n", r.RunID, r.Profile, r.Module, r.Target)
	fmt.Println("\nCommands:")
	for _, c := range r.Commands {
		status := "OK"
		if c.ExitCode != 0 || c.Error != "" {
			status = "WARN"
		}
		fmt.Printf("- [%s] %s exit=%d duration=%s\n", status, c.Name, c.ExitCode, c.Duration)
		if strings.TrimSpace(c.Stdout) != "" {
			fmt.Println(indent(limit(c.Stdout, 1200), "  "))
		}
		if strings.TrimSpace(c.Stderr) != "" {
			fmt.Println(indent(limit(c.Stderr, 800), "  stderr: "))
		}
	}
	fmt.Println("\nFindings:")
	if len(r.Findings) == 0 {
		fmt.Println("- no rule-based findings in v0.1")
		return
	}
	for _, f := range r.Findings {
		fmt.Printf("- [%s] %s/%s: %s\n  recommendation: %s\n", f.Severity, f.Service, f.RuleID, f.Summary, f.Recommendation)
	}
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return strings.Join(lines, "\n")
}

func limit(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n...truncated..."
}

func shellQuoteSafe(s string) string {
	return strings.ReplaceAll(s, `'`, `'"'"'`)
}
