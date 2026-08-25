# Release gate: deterministic Dolt pull remote selection

- Deploy bead: `ga-nht26j`
- Build/review work: `ga-fe5cva` / `ga-mdce9d` / `ga-g04htm`
- Reviewed source: `7dbbc6d321e15ba64188dc37153bf6f4c82f3489`
- Reviewed base: `origin/main@09bae7ad1706aced67a72775d2b1d11549002cd0`
- Final base check: `origin/main@6b60167cf5891af5416ad854fe9d0542a0ade5f3`
- Deploy mode: remote; resolved push target: `origin`
- Evaluated: 2026-08-25
- Verdict: **PASS WITH ATTRIBUTED FAILURES**

`docs/PROJECT_MANIFEST.md` is not present in this worktree. This record uses
the seven criteria in the deployer contract and the sharded test commands in
`TESTING.md`.

## Gate checklist

| # | Criterion | Result | Evidence |
|---|---|---|---|
| 1 | Review PASS present | **PASS** | Closed review bead `ga-g04htm` records `verdict: pass`; both `metadata.commit` and `metadata.deploy_commit` are pinned to exact source `7dbbc6d321e15ba64188dc37153bf6f4c82f3489`. The reviewer independently re-stamped the post-rebase delta and found every branch-owned path byte-identical to the prior PASS. |
| 2 | Acceptance criteria met | **PASS** | Multi-remote pull now fails closed without `GC_DOLT_REMOTE_<DB>`; a named file remote is honored; a named non-`file://` remote additionally requires `GC_DOLT_PULL_ALLOW_REMOTE_<DB>=1`; invalid and unknown overrides fail; SQL and CLI sources are both covered. |
| 3 | Tests pass | **PASS WITH ATTRIBUTED FAILURES** | Focused package run: 363 test/subtest PASS events, 0 FAIL, 0 SKIP; all 42 diff-owned top-level tests pass by exact name. Required 40-job union: 35 PASS jobs, 5 FAIL jobs, 0 SKIP jobs. All five raw failures are non-diff-owned, tracked, logged, and mechanism-attributed below. |
| 3a | Pre-existing failures attributed | **PASS** | `TestProviderLiveClaudeKindPath`, `TestBdFlagManifestCurrent`, both tmux default-binding tests, and `TestE2E_SuspendResume_City` map to specific predating trackers. The candidate changes Dolt pull fixtures/tests and the resource-census ledger; none of those packages or mechanisms is reachable from the diff. |
| 3b | Policy and static lanes | **PASS WITH ATTRIBUTED FAILURE** | CI policy, module replacement, event-export isolation, core boundary, native DoltLite beads, docs sync, resource-census package, formatting, diff check, build, vet, and isolated-cache full affected lint all pass. The native dependency surface remains raw red at 270,245,544 bytes, but exact reviewed base is larger at 270,254,416 and current main is larger at 270,259,128; attributed to `ga-5flk3r`. |
| 4 | No high-severity review findings open | **PASS** | Reviewer reports no security, style, specification, or correctness blocker. Unresolved HIGH findings: 0. |
| 5 | Final branch is clean | **PASS** | `git status --porcelain=v1` was empty at the exact reviewed source before replacing this gate record; `git diff --check` and changed-file formatting pass. |
| 6 | Branch diverges cleanly from main | **PASS** | After a final `git fetch origin`, `git merge-tree --write-tree --messages origin/main 7dbbc6d3...` exited 0 against `origin/main@6b60167cf5` and produced tree `3da8000ea1987b1cef3789df9b116edfc875e576`. The branch is 1 behind / 9 ahead. `assert_deploy_ancestry_scope` passed for `ga-nht26j`, `ga-mdce9d`, `ga-fe5cva`, and `ga-fqi7kq`. |
| 7 | Single feature theme | **PASS** | The eight-file delta is one feature: deterministic Dolt pull remote selection, its SQL/CLI tests, parallel-safety adjustments in the shared fixtures, and the mechanically required resource-census ledger reconciliation. |

## Diff-owned test evidence

```text
test_cmd: GC_FAST_UNIT=1 scripts/go-test-observable test -- -count=1 -timeout 15m ./examples/bd/dolt/...
test_counts: 363 PASS events, 0 FAIL, 0 SKIP
full_log: /var/tmp/gascity-test.jsonl.bzywty
waiver_ref: none for diff-owned tests
```

All 42 diff-owned top-level tests passed:

- Eight new `TestPull*MultipleRemotes*` tests in
  `pull_remote_selection_test.go`.
- Both modified live-SQL pull tests in `pull_test.go`.
- All 32 modified `TestSync*` tests in `sync_test.go`.

The eight new tests directly cover ambiguous SQL and CLI refusal, explicit
file-remote selection, non-local refusal and opt-in, invalid override syntax,
unknown remote names, and named CLI selection. Every test calls
`t.Parallel()`; the shared fixture updates make that execution safe.

## Required full-suite evidence

```text
test_cmd: DOCKER_HOST=unix:///run/user/1000/podman/podman.sock TESTCONTAINERS_RYUK_DISABLED=true LOCAL_TEST_JOBS=4 GO_TEST_TIMEOUT=30m make test-local-full-parallel
test_counts: 35 PASS jobs, 5 FAIL jobs, 0 SKIP jobs (40 total)
cmd_gc_process_shards: 6 PASS, 0 FAIL, 0 SKIP
full_logs: /var/tmp/gc-local-tests.EYl9gF
```

### Raw failures and attribution

The raw failures remain visible; none was rewritten green.

- `TestProviderLiveClaudeKindPath` -> `ga-fh1flg` / `ga-zron27`.
  It failed in `unit-core` with the exact tracked `agent_pane_busy` on
  `w1:p1` signature. The candidate has no `internal/runtime/herdr` path or
  pane-allocation mechanism. **FAIL-WAIVED** under
  `mayor-2026-08-20-herdr-pane-standing`.
- `TestBdFlagManifestCurrent` -> `ga-f0uceo`. It compared the installed
  `bd` binary with `internal/bdflags`; neither is changed or reachable from
  the Dolt pull/test-policy diff.
- `TestGetKeyBinding_CapturesDefaultBinding` and
  `TestGetKeyBinding_CapturesDefaultBindingWithArgs` -> `ga-afqddr`.
  Both observed the tracked empty host tmux 3.7b default-binding table. The
  candidate has no `internal/runtime/tmux` path or mechanism.
- `TestE2E_SuspendResume_City` -> `ga-yc0e3a`. It timed out waiting for
  `citysus.report`, the exact tracked lifecycle-fixture signature. The
  candidate does not execute city suspend/resume.

Each occurrence, exact SHA, mechanism proof, and log path was appended to its
tracker before scoring this gate.

```text
failure_attribution: TestProviderLiveClaudeKindPath -> ga-fh1flg / ga-zron27 + mayor-2026-08-20-herdr-pane-standing + separate-package no-mechanism proof
failure_attribution: TestBdFlagManifestCurrent -> ga-f0uceo + separate-package no-mechanism proof
failure_attribution: TestGetKeyBinding_CapturesDefaultBinding{,WithArgs} -> ga-afqddr + separate-package no-mechanism proof
failure_attribution: TestE2E_SuspendResume_City -> ga-yc0e3a + separate-package no-mechanism proof
```

### Pre-push hook evidence

`GIT_TERMINAL_PROMPT=0 git push --dry-run origin HEAD` ran the active
pre-push hook. Nine fast jobs passed and only `unit-core` failed:
`TestProviderLiveClaudeKindPath` reproduced the same
`agent_pane_busy` / `w1:p1` signature already recorded above. The
diff-owned `examples/bd/dolt` package passed in 43.974 seconds. Log:
`/var/tmp/gc-local-tests.KtYRiK/unit-core.log`.

This exact hook failure is covered by
`mayor-2026-08-20-herdr-pane-standing`. It is preserved as
**FAIL-WAIVED**. The specific deploy-branch push uses the standing
authorization's one-time `--no-verify` path instead of retrying a known,
non-diff-owned failure until it happens to turn green.

## Policy, static, and baseline evidence

```text
make test-ci-policy — PASS
make check-gomod-replace — PASS
make check-eventexport-isolation — PASS
make check-core-boundary — PASS
make test-native-doltlite-beads — PASS
make check-docs — PASS
go test -count=1 ./internal/testpolicy/resourcecensus — PASS
go build ./... — PASS
go vet ./... — PASS
make fmt-check-changed LINT_CHANGED_SCOPE=tracked LINT_CHANGED_REF=09bae7ad... — PASS
git diff --check 09bae7ad...HEAD — PASS
GOLANGCI_LINT_CACHE=<isolated-on-disk-cache> make lint-affected LINT_CHANGED_SCOPE=tracked LINT_CHANGED_REF=09bae7ad... — PASS, 0 issues
```

The first affected-lint attempt used the shared golangci-lint cache and emitted
616 diagnostics whose paths pointed into deleted, unrelated
`/var/tmp/ga-*-gate` worktrees. Repeating the identical full affected-package
closure with an isolated on-disk linter cache produced `0 issues`; the shared
Go build cache was not changed or cleared.

`make check-native-dependency-surface` remains raw red and is attributed to
open tracker `ga-5flk3r`:

```text
candidate 7dbbc6d3:        270245544 bytes
exact base 09bae7ad:       270254416 bytes
current main 6b60167c:     270259128 bytes
ceiling:                   270000000 bytes
```

The candidate is 8,872 bytes smaller than its exact base and adds no dependency
to the `gc` runtime graph. This is baseline threshold debt, not candidate
growth; the occurrence and both A/B measurements are preserved on
`ga-5flk3r`.

## Pre-flight and provenance

- GitHub commit-to-PR lookup returned `[]`; no existing PR carries the
  reviewed source.
- `origin/builder/ga-nht26j-footprint-fix` resolves to the exact reviewed
  source.
- The source tree was clean, the ancestry scope guard passed, and the isolated
  deploy branch was derived mechanically as `deploy/ga-nht26j-gate`.
- The gate was evaluated at the exact reviewed source; no contributor branch
  was modified or pushed.

## Release decision

The candidate is ready for an isolated deploy-branch push and pull request.
Attributed baseline failures remain explicit in this record. Merge authority
remains with the mayor / maintainer-review workflow.
