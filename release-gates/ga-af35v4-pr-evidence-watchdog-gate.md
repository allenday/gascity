# Release gate: PR evidence watchdog (`ga-af35v4`)

**Verdict: FAIL**

- Deploy source: `fda1be8dfec37fb9cbe4ff039249d909574e0b37`
- Latest base checked: `origin/main@08ecb0585498a0a5464e78a3b5d122236ff0ac9d`
- Deploy mode: remote; push remote would be `fork`
- Gate branch: `deploy/ga-af35v4-gate-fail-rerun-20260822` (local only; not pushed)
- Build bead: `ga-oaz41a.2`
- Review bead: `ga-vq8kck`
- Existing target PR preflight: none found for the resolved deploy source

## Checklist

| # | Criterion | Result | Evidence |
|---|---|---|---|
| 1 | Review PASS present | PASS | `ga-vq8kck` records the full review PASS at `2b0cc1d510f453431a560e8a24c411a83bdc3903`; the reviewer then independently delta-reviewed the direct descendant `fda1be8dfec37fb9cbe4ff039249d909574e0b37` and recorded `DELTA REVIEW: PASS`. |
| 2 | Acceptance criteria met | PASS | The original watchdog implementation remains unchanged from the reviewed commit. The delta adds only the mandatory untagged `scripts/prwatchdog/testenv_import_test.go` guard. Independent focused evidence confirms the exact trigger, minimum read permissions, trusted-base-only checkout, explicit PR-head lookup, bounded observation loop, fail-closed state machine, opt-in suite handling, pagination, and human-readable summary tests all pass. |
| 3 | Tests pass | **FAIL** | The documented full CI-equivalent sweep completed with 33 PASS / 7 FAIL / 0 SKIP jobs. The diff-owned watchdog and import-guard tests pass, but `TestSessionReconcilerTraceGH1654WorkRequestedStartCandidates/named_session_post-kill` failed on the candidate and passed both focused and comparable full-shard runs at the exact base, so pre-existing attribution clause (iii) is not met. Details below. |
| 3b | Policy/lint lane | PASS with attributed host findings | `make test-ci-policy`, `make fmt-check`, and `make vet` pass. Fresh-cache full lint reports only three generated dashboard `node_modules/flatted` findings; the exact same findings reproduce at the exact base and are tracked by `ga-4go623`, with no path overlap. |
| 4 | No high-severity review findings open | PASS | Reviewer recorded no blocking or HIGH security finding. Independent inspection found no PR-head checkout or execution, secret scope, write permission, or shell-interpolation surface. |
| 5 | Final branch is clean | PASS | `git status --short` was empty at the resolved deploy source immediately before this gate file was added. |
| 6 | Branch diverges cleanly from main | PASS | After the final refresh, `git merge-tree --write-tree origin/main HEAD` succeeded with tree `ba866e0948493e8b29862a239c26aa712cb37796`; `assert_deploy_ancestry_scope origin/main HEAD ga-af35v4 ga-oaz41a.2` returned 0. No self-rebase was needed. |
| 7 | Single feature theme | PASS | The commit set is one cohesive PR-evidence-watchdog feature: the trusted-base workflow, its `scripts/prwatchdog` implementation/tests, the CI-policy registration, and the required test-environment guard. |

## Criterion 3 evidence

### Required full suite — FAIL

```text
test_cmd: DOCKER_HOST=unix:///run/user/1000/podman/podman.sock TESTCONTAINERS_RYUK_DISABLED=true LOCAL_TEST_JOBS=4 make test-local-full-parallel
test_counts: 33 PASS, 7 FAIL, 0 SKIP jobs (40 total)
container_evidence: rootless Podman socket available; dolthub/dolt:2.1.7 cached
```

Hard blocker:

- `TestSessionReconcilerTraceGH1654WorkRequestedStartCandidates/named_session_post-kill` failed after 6.90 seconds with `async starts did not finish` in candidate job `cmd-gc-process-2-of-6`. Tracker `ga-hgjlhi` exists and there is no path overlap, but the exact test passes on `origin/main@08ecb0585498a0a5464e78a3b5d122236ff0ac9d` in both a focused run (0.22 seconds) and the comparable full `cmd/gc` shard 2-of-6 (1547 tests, package PASS in 76.099 seconds). Clause (iii), proven pre-existing on the exact base, is therefore missing. This occurrence is not attributable and criterion 3 remains FAIL.

Attributed failures satisfying all four clauses:

- `failure_attribution: TestBdFlagManifestCurrent -> ga-f0uceo + exact-base integration-tag run reproduced the same installed-bd manifest failures; the diff does not touch internal/bdflags`.
- `failure_attribution: TestGetKeyBinding_CapturesDefaultBinding -> ga-afqddr + exact-base integration-tag run reproduced the same empty tmux 3.7b default binding; the diff does not touch internal/runtime/tmux`.
- `failure_attribution: TestGetKeyBinding_CapturesDefaultBindingWithArgs -> ga-afqddr/ga-k3fxvj + exact-base integration-tag run reproduced the same empty tmux 3.7b default binding; the diff does not touch internal/runtime/tmux`.

Standing-waived raw failures:

- `TestFreshManagedBdCityInitSeedsPinnedHQDatabaseAndKeepsGCPrefix`, `TestGCLiveContract_BeadsAndEvents`, and `TestCleanInstallTutorialPath` all failed with the exact `gastownhall/beads#4566` pending-dirty-schema-migration signature. They are preserved as **FAIL-WAIVED** under the mayor's standing authorization recorded on `ga-lpfjhc` / root-cause tracker `ga-6bnc42`. The candidate changes only the PR watchdog and cannot reach Dolt schema migration or store bootstrap. This occurrence is logged on `ga-lpfjhc`; the waiver does not cover the independent hard blocker above.

### Diff-owned tests — PASS

```text
test_cmd: go test -count=1 -json ./internal/testenv/... ./scripts/prwatchdog/...
diff_tests_executed: scripts/prwatchdog 50 PASS, 0 FAIL, 0 SKIP
guard_test: TestRequiresDedicatedTestenvImportFile PASS
waiver_ref: none
```

The `scripts/prwatchdog/cmd/watchdog` binary package correctly reports no test files; its behavior is exercised from the parent watchdog package. No diff-owned test skipped.

### Policy and static lanes

- `policy_lane: make test-ci-policy` — PASS: 5 runner-policy tests, 15 CI-suite coverage tests, `scripts/cipolicy`, `scripts/prwatchdog/...`, and the four changed-static-scope tests.
- `make fmt-check` — PASS.
- `make vet` — PASS.
- `make lint-affected` — exit 0 but reported no changed worktree paths, so it was not used as sole evidence.
- Fresh-cache `make lint` — raw FAIL with exactly three findings in generated `internal/api/dashboardspa/web/node_modules/flatted/golang/pkg/flatted/flatted.go` (two `govet` inline findings and one `revive` package-comment finding). `failure_attribution: full lint generated flatted findings -> ga-4go623 + exact-base fresh-cache make lint reproduced the identical three findings; no candidate path overlap`. No watchdog lint finding was reported.

## Disposition

Gate FAIL on criterion 3. Do not push this branch, open a PR, post deploy clearance, or route a merge request. Return `ga-af35v4` to the builder with the candidate-only session-reconciler timeout and exact-base PASS evidence. The watchdog's prior diff-owned package-hygiene blocker is fixed and independently green; this reroute is solely for the remaining unattributable full-suite failure.

## Rebuild Verification

Builder: `gascity/builder-1`, in a dedicated nested worktree (`worktrees/ga-af35v4`) checked out to `builder/ga-oaz41a.2` at the same deploy source, `fda1be8dfec37fb9cbe4ff039249d909574e0b37`. No candidate source changed; this section adds attribution evidence for the criterion 3 blocker only.

### Diff-owned tests and static checks — re-confirmed green

```text
test_cmd: go test -count=1 ./internal/testenv/... ./scripts/prwatchdog/...
result: both packages ok
diff_tests_executed: scripts/prwatchdog 50 PASS, 0 FAIL, 0 SKIP
guard_test: go test -count=1 -run TestRequiresDedicatedTestenvImportFile -v ./internal/testenv/... -> PASS (0.05s)
```

`go build ./...` and `go vet ./...` both exit 0; `gofmt -l` on the changed files is empty. These match the original gate's citations exactly — the rebuild introduces no regression.

### Root cause, precisely stated

`ga-hgjlhi` describes the failing wait as a hardcoded 5-second timeout. Reading the current source (`cmd/gc/city_runtime.go`, `waitForAsyncStarts`) shows the mechanism is slightly more specific: the timeout is `cr.cfg.Daemon.ShutdownTimeoutDuration()`, and the code only falls back to a literal `5 * time.Second` when that resolves to `<= 0`. It is config-driven in principle, not an unconditional constant.

In practice that distinction does not help this test. `TestSessionReconcilerTraceGH1654WorkRequestedStartCandidates` builds its `cfg` from `writeCityTOML(t, cityDir, "trace-town", "worker")` (`cmd/gc/session_reconciler_trace_integration_test.go:616`), a minimal fixture with no daemon shutdown-timeout section, so `ShutdownTimeoutDuration()` resolves to zero and the test runs under the same fixed 5-second ceiling `ga-hgjlhi` describes. The test drives `cr.beadReconcileTick(...)` against a fake provider and then asserts `cr.waitForAsyncStarts()` returns true within that fixed wall-clock window (lines 655-658). The window does not scale with host scheduling contention, so under enough concurrent load the async-start goroutines can legitimately still be waiting for their turn on a core when the 5-second deadline fires — a false failure with no correctness defect in the watchdog or in `waitForAsyncStarts` itself.

### Fresh reproduction attempt at the exact base

To test the pre-existing/load-sensitive theory directly rather than by inference alone, a separate detached-HEAD worktree was built at the exact gate base, `08ecb0585498a0a5464e78a3b5d122236ff0ac9d`, and the full `cmd/gc` process suite was run under contention matching the candidate's original failure conditions:

```text
test_cmd: LOCAL_TEST_JOBS=4 CMD_GC_PROCESS_TOTAL=6 make test-cmd-gc-process-parallel
result: 4 of 6 shards FAILED (1, 2, 3, 4); shards 5 and 6 passed; make exited 123
```

This did not reproduce the identical assertion — `TestSessionReconcilerTraceGH1654WorkRequestedStartCandidates` is not named in any shard's failure output this run, and the string `async starts did not finish` does not appear in any shard log. Reporting that plainly rather than overstating it. What it does show:

- The exact base commit is measurably **not** reliably green under this contention profile in this exact package. Shard `cmd-gc-process-2-of-6` — the same shard number the candidate's original failure occurred in — failed again here, from a different specific test (`TestCmdWorkflowDeleteSourceStoreSelectorIgnoresLegacyRootInDifferentStore`). Shards 1 and 3 failed on `caching-store: prime FAILED: context canceled` — the same failure shape as the target flake: a fixed-deadline operation racing goroutine scheduling under CPU oversubscription, just manifesting on a different call site.
- Two separate quiet/lower-contention base reruns already on file for this gate (the deployer's focused single-test run and its comparable full shard-2 run) both came back PASS. A single clean rerun of a genuinely load-sensitive race does not disprove pre-existence — it only means that rerun did not land inside the failure window. This fresh run demonstrates the opposite result is also readily obtainable at the same commit, under the same kind of load, without touching the candidate diff at all.
- Shard 4's failure was a bare `FAIL` with no `--- FAIL:` test line, consistent with a process-level disruption under the same resource pressure rather than a test assertion; it was not investigated further since it is not the target test and does not change the conclusion either way.

### Cross-branch and same-branch history

`ga-hgjlhi` documents this exact test/subtest failing intermittently with this exact `async starts did not finish` signature across 6+ unrelated branches, always with zero path overlap between the failing diff and the reconciler-timeout code. This branch (`builder/ga-oaz41a.2`) has already hit this identical failure once before, independent of both the gate's own candidate run and this rebuild.

### Assessment

Taken together — the precise, confirmed root cause (a fixed 5-second wait with no headroom for scheduling contention, in a test fixture that leaves the timeout at its unconfigured default), the pre-existing tracker's extensive cross-branch recurrence, this branch's own prior occurrence, and a fresh demonstration that the exact base commit is genuinely unstable under matching contention in this exact package (including a repeat failure in the same shard number, from the same general timeout-under-load failure class) — this reads as a pre-existing, diff-unrelated flake. The evidence is not a literal re-trigger of the identical assertion in this rebuild; that is reported honestly above rather than implied. If the deploy gate's policy requires a byte-identical reproduction before treating clause (iii) as satisfied, this should go to the mayor as a policy question rather than be resolved unilaterally here. Returning to `gascity/deployer` with this evidence attached for that determination.
