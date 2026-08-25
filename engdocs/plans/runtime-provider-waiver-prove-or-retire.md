# Runtime provider waiver prove-or-retire plan

*Status: implementation-ready decomposition. Producing bead: `ga-glwaz5.2`.
Durable governance owner: `ga-uz5t3a` (keep open).*

## Outcome

Replace the remaining bulk runtime-provider waiver debt with six independently
owned constructor decisions. Each production constructor either gains an exact,
executable `runtime.Provider` proof through
`internal/runtime/runtimetest.RunProviderTests`, or retains a short,
constructor-specific waiver whose reason describes the verified gap.

The immediate expiry/owner repair remains separate in `ga-glwaz5.1`. This plan
does not reopen that mechanism decision and does not group unrelated providers
into one implementation change.

## Grounded disposition by constructor

| Constructor boundary | Current evidence | Planned disposition | Ownership | Bead |
| --- | --- | --- | --- | --- |
| `internal/runtime/t3bridge.NewSeamBacked` | Registered by both `exact:t3bridge` and the legacy `exec:` prefix. Focused tests already use a mock T3 WebSocket bridge, but neither registration has a shared contract proof. | Prove now with one hermetic exact-constructor test and bind both ledger claims to it. | Fork-owned T3 Code integration; keep T3-specific fixtures inside `internal/runtime/t3bridge`. | `ga-uz5t3a.1` |
| `internal/runtime/k8s.NewSeamBacked` | The production path creates a real client-go configuration and client. Existing lifecycle and seam tests enter through `newProviderWithOps`, so they prove an injected fake rather than the registered constructor. | Re-waive with evidence naming the missing Kubernetes API and pod-exec contract harness. A future proof must run the exact constructor; the fake constructor cannot satisfy the ledger. | Upstream-shared K8s backend. | `ga-uz5t3a.2` |
| `internal/runtime/herdr.New` | `TestHerdrConformance` already calls the shared runner against the exact constructor, but skips in short mode and when the external `herdr` executable is absent. | Re-waive until that suite is ungated in a supported required environment; preserve the existing live coverage. | Upstream-shared herdr backend. | `ga-uz5t3a.3` |
| `internal/runtime/tmux.NewSeamBackedWithConfig` | `TestTmuxConformance` already exercises the exact seam-backed constructor on an isolated socket, but skips without tmux and uses `RunProviderTestsWithOptions`. Integration CI installs tmux. | Prove now by making the integration-owned test direct and ungated, while retaining isolated socket and cleanup ownership. | Upstream-shared tmux backend. | `ga-uz5t3a.4` |
| `internal/runtime/ssh.NewSeamBacked` | Focused tests inject `fakeRunner` through `providerWith`, bypassing both the registered constructor and its `shellRunner`. The production client invokes `ssh` from `PATH`. | Prove now with a test-owned, network-free SSH client boundary that drives the exact constructor and maintains isolated concurrent session state. | Upstream-shared SSH backend. | `ga-uz5t3a.5` |
| `cmd/gc.newHybridProvider` | The registered boundary constructs the real seam-backed tmux and K8s providers. Existing `internal/runtime/hybrid` tests compose fake backends and do not prove this `cmd/gc` boundary. | Re-evaluate after the K8s and tmux dispositions. Use a scoped proof only if it directly runs `newHybridProvider` and states which route is proved; otherwise re-waive with the unresolved component named. | Upstream-shared `cmd/gc` composition boundary. | `ga-uz5t3a.6` |

## Dependency graph

`ga-glwaz5.1` is the common blocker because it first replaces the dead waiver
owner and shared expiry with the durable owner and per-constructor dates. After
it lands, T3 bridge, K8s, herdr, tmux, and SSH can proceed independently.

| Work item | Blocked by | Why |
| --- | --- | --- |
| `ga-uz5t3a.1` through `.5` | `ga-glwaz5.1` | Start from the durable owner and independent-expiry ledger shape instead of conflicting with the in-flight repair. |
| `ga-uz5t3a.6` | `ga-glwaz5.1`, `ga-uz5t3a.2`, `ga-uz5t3a.4` | The hybrid production boundary constructs K8s and tmux; its honest scope depends on both component outcomes. |

No dependency orders T3 bridge, K8s, herdr, tmux, and SSH relative to each
other. Their provider behavior and proof fixtures are independent even though
each final change also updates the central ledger and its generated
`TESTING.md` table.

## Shared acceptance rules

Every child bead carries constructor-specific measurable criteria. Across all
six, the following rules also apply:

1. A proved claim names an exact runnable test whose factory directly returns
   the registered production constructor and whose test directly calls
   `RunProviderTests`. Focused tests or injected alternate constructors are not
   substitutes.
2. A retained waiver is owned by `ga-uz5t3a`, has its own bounded expiry, and
   names the concrete missing harness or environment gate. It must not reuse a
   bulk date or generic “no coverage” wording.
3. Both T3 bridge registrations receive the same disposition; the legacy
   `exec:` route cannot be forgotten when the exact route changes.
4. Provider package tests and
   `TestCatalogMatchesProductionWiringAndDocumentation` pass. Any ledger change
   regenerates the checked `TESTING.md` table.
5. Provider-specific behavior stays behind its existing runtime boundary. No
   T3 Code, Kubernetes, herdr, tmux, or SSH assumption moves into generic
   runtime paths merely to make a test convenient.

## Risks and coordination

- The ledger proof validator rejects pre-run skip gates, indirect constructors,
  and the wrong contract runner. A test that passes independently may still be
  ineligible as ledger evidence.
- Several children will touch `internal/testutil/providerledger/ledger.go` and
  generated `TESTING.md`. Their provider tests are parallelizable, but each
  implementation branch should refresh those central files immediately before
  review to avoid stale generated diffs.
- K8s and herdr are intentionally re-waived from current evidence. Their waivers
  remain forcing functions: a later renewal must re-check whether a supported
  required environment has appeared.
- No `upstream` Git remote is configured in this worktree. Archaeology across
  all local refs and comparison with current `origin/main` found no prior
  unlanded exact-constructor proof to port; implementations should repeat the
  path-specific history check if refs advance before work begins.

## Handoff

All six beads are `ready-to-build` and carry `source:actual-pm` plus
`runtime-provider-governance`. They are children of `ga-uz5t3a`, not of the
short-lived triage bead, so future waiver renewals and proof conversions remain
discoverable after `ga-glwaz5.2` closes.
