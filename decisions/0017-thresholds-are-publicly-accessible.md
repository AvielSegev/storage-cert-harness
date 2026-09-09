# 0017. SLA thresholds are publicly accessible

- **Status:** Accepted
- **Date:** 2026-09-09

## Context

[ADR-0004](0004-threshold-secret-handling.md) treated real SLA thresholds as
sensitive: never in plaintext, code, or readable Git history, delivered only via
a dedicated encrypted folder with a shared symmetric key (mechanism still
deferred). That constraint drove the v2.0 two-file split (`catalog.json` +
`thresholds.json`, ADR-0006), the gitignored real `thresholds.json`, and the CI
secret-scan gate over numeric literals.

The thresholds are now cleared for public release. There is no longer a
requirement to keep real SLA numbers out of the repository or its history.

## Decision

Real SLA thresholds may be committed in the clear. `thresholds.json` becomes a
publishable file alongside `catalog.json`; the encrypted-folder mechanism and
shared-key custody described in ADR-0004 are no longer required. This
**supersedes ADR-0004** (which remains recorded for history).

Invariants preserved:

- The v2.0 two-file contract (ADR-0006) is unchanged: `catalog.json` +
  `thresholds.json`, joined on TR id, sharing `kb_git_commit`. Only the
  sensitivity of `thresholds.json` changes.
- Report provenance still stamps the thresholds source.
- Backend credentials remain references, never committed plaintext (ADR-0005) —
  this decision covers SLA thresholds only, not secrets in general.

## Decision notes

### Options considered

- **Publish thresholds in the clear (chosen):** removes key-management burden and
  the deferred encryption mechanism; makes partner/developer threshold retrieval
  trivial. Irreversible — once public, historical values cannot be un-published.
- **Keep encrypted-in-repo (ADR-0004):** no longer justified now that the data is
  cleared for release; retains needless key custody and rotation work.

### Consequences

- `thresholds.json` is no longer gitignored and may be committed; the fake
  `thresholds.example.json` becomes redundant but may stay as a schema-conformant
  sample.
- The CI **secret-scan** SLA-numeric diff scan is removed; the gate now covers
  only tracked run reports (`report.json` / `report.md`) and token-shaped
  credentials in editor/workspace files. `ci/config/secret-scan-allowlist.txt`
  and the scattered `secret-scan:ok` markers become no-ops (harmless; cleanup
  optional).
- `AGENTS.md` critical rule 1, the `secrets/` table row, `docs/ci.md`, and
  `.gitignore` / `.pre-commit-config.yaml` are updated to drop the secrecy
  assumption.

Related: [ADR-0004](0004-threshold-secret-handling.md),
[ADR-0006](0006-adopt-kb-export-v2-contract.md),
[ADR-0005](0005-test-plans-and-backend-config.md).
