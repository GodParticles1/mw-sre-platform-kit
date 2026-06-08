# Validation Report

Validation time: 2026-06-04

## Checks performed

```bash
go test ./...
bash -n scripts/mw_quick_survey.sh
go build -o bin/mwctl ./cmd/mwctl
./bin/mwctl version
./bin/mwctl quick --module redis --redis-pass test --json
python3 -m json.tool /tmp/mwctl_test.json
```

## Result

```text
VALIDATION_OK
```

## Scope of validation

Validated:

- Go code compiles.
- Go unit tests pass.
- Quick survey shell script passes bash syntax validation.
- `mwctl` emits parseable JSON.

Not validated against a live production environment in this sandbox:

- Live SSH execution.
- Live MySQL/Mongo/Redis/RabbitMQ/ClickHouse commands.
- K8s CRDs applied to a cluster.
- Argo Workflow execution.

These live validations should be run in staging before production rollout.
