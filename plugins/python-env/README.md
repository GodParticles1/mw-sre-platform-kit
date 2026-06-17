# Python Environment Plugin

Purpose: diagnose driver/Python/Console-script failures caused by `/etc/profile` or user profile pollution, non-zero `source`, noisy profile output, missing Python runtime, broken standard-library imports and suspicious Python-related env variables.

Lifecycle implemented:

- discover: current Python executable and selected environment
- collect: quiet-source check for `/etc/profile`, bounded profile risk lines, Python import smoke test
- diagnose: source non-zero, profile output, python3 missing, import failure, possible path pollution, TIMEOUT/TMOUT baseline hints
- plan/apply/verify/rollback: runbook stays dry-run; profile edits require manual review and backup

Safety:

- collectors are read-only
- no profile file is modified automatically
- output is bounded and avoids dumping all environment variables
