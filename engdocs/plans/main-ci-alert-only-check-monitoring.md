# Main CI alert-only check monitoring plan

*Status: implementation-ready decomposition. Producing bead:
`ga-glwaz5.3`.*

## Outcome

Surface failures from explicitly enrolled, non-required checks on the default
branch as durable investigator work. The first enrollment covers Gas City's
push-only `Preflight / unit cover (noncmdgc)` context so a provider-ledger
expiry cannot remain red on main until unrelated work encounters it.

This is an alerting change, not a branch-protection change. Pull requests keep
their current required checks, and failures from unrelated non-required checks
remain outside the monitor.

## Product direction

Extend the existing local `main-ci-watcher` patrol instead of introducing a
second polling loop, GitHub workflow, persistence mechanism, or notification
channel. The watcher already owns default-branch CI-health incidents, P1
priority, duplicate suppression, and investigator routing. The new capability
adds an explicit repository-policy enrollment beside the required contexts the
watcher discovers live.

The architecture package defines that enrollment contract before tests or
implementation begin. In particular, it must preserve current filtering for
the default-branch head, push-triggered suites, latest attempts, terminal
failures, and recovered reruns.

## Work packages

| Bead | Outcome | Route | Intake label |
| --- | --- | --- | --- |
| `ga-glwaz5.3.1` | Specify the alert-only enrollment contract at the existing watcher and repository-policy boundaries. | architect | `needs-architecture` |
| `ga-glwaz5.3.2` | Author executable RED contracts for enrolled, unenrolled, recovered, duplicate, and error cases. | validator | `needs-tests` |
| `ga-glwaz5.3.3` | Extend the watcher and enroll Gas City's noncmdgc unit-cover context under the approved contract. | builder | `ready-to-build` |

Each bead carries its own measurable acceptance criteria. The architecture
package is intentionally separate because the root establishes the product
direction but leaves configuration shape and API-error semantics to the
architect.

## Dependency graph

```text
ga-glwaz5.3.1 ──> ga-glwaz5.3.2 ──> ga-glwaz5.3.3
        └──────────────────────────> ga-glwaz5.3.3
```

The validator contracts depend on the approved architecture. Implementation
depends on both the architecture and committed RED coverage, preventing the
builder from inventing either boundary while coding.

## Acceptance mapping

| Root need | Work packages |
| --- | --- |
| A red enrolled main check creates durable P1 investigator work. | `.1`, `.2`, `.3` |
| Required checks remain monitored and an empty enrollment preserves current behavior. | `.1`, `.2`, `.3` |
| Unenrolled, non-head, non-push, in-progress, and recovered failures remain excluded. | `.1`, `.2`, `.3` |
| Repeated observations do not create duplicate open incidents. | `.2`, `.3` |
| API or malformed-policy failures cannot weaken required-check monitoring or become false green. | `.1`, `.2`, `.3` |
| Gas City's live `Preflight / unit cover (noncmdgc)` context is enrolled without becoming a PR gate. | `.1`, `.2`, `.3` |
| The existing single 20-minute patrol remains the operational mechanism. | `.1`, `.3` |

## Product boundaries

- Do not add `preflight-unit-cover-noncmdgc` to `ci-required` or to pull-request
  branch protection.
- Do not edit Gas City's CI workflow merely to create an alert path for a check
  that already runs on pushes to main.
- Do not monitor every non-required check. Enrollment is explicit, exact, and
  empty by default.
- Do not hardcode repository, workflow, job, check, or agent-role names in the
  generic watcher. The initial Gas City context belongs in repository policy.
- Preserve the existing incident classifier, open-bead deduplication, priority,
  and routing behavior.
- The provider-waiver expiry mechanism and per-constructor proof work remain in
  `ga-glwaz5.1` and `ga-glwaz5.2`; this plan only closes the visibility gap.

## Risks and controls

- **False positives from unrelated Actions activity:** exact enrollment plus
  existing push-suite, head-SHA, and latest-attempt filters constrain the new
  signal.
- **Silent loss of required-check monitoring:** contracts cover configured and
  empty enrollment, plus ruleset, Actions, check-run, and malformed-policy
  errors.
- **Alert storms:** existing open-incident deduplication remains mandatory, and
  recovered latest reruns suppress older failures.
- **Accidental gate expansion:** validation must demonstrate that enrollment
  produces an alert-only incident and does not alter PR required checks.
- **Parallel-monitor drift:** all work stays at the existing watcher and
  repository-policy ownership boundaries.

## Handoff

The architecture bead is immediately actionable. The validator and builder
beads remain dependency-blocked until their predecessors close. All three are
children of `ga-glwaz5.3`, carry `source:actual-pm`, and have exactly one intake
label for their target phase.
