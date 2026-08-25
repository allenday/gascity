---
title: Gas City Identity Separator Contract — v1
description: Authoritative specification for the qualified-identity encoding shared by gascity and beads.
---

| Field | Value |
|---|---|
| Status | Authoritative specification |
| Last verified | 2026-08-25 |
| Primary implementation | `internal/agent/session_name.go` |
| Mints identities | gascity |
| Compares identities | beads |
| Concept model | [How Gas City Works](/getting-started/how-gas-city-works) — the Agent (WHO) primitive |

## What this is

gascity mints **qualified agent identities** — `rig/agent` and `city.agent`
combinations — and renders them into tmux-safe session name strings. beads
stores and compares those same strings in its own records (session-name
metadata, routing fields) without minting them itself. Two independent
codebases have to agree, byte for byte, on how the structural separators `/`
and `.` get encoded, or a qualified identity that round-trips through beads
and back can silently turn into a different identity than the one gascity
minted.

This contract exists so that agreement is written down once, instead of
re-derived by reading `internal/agent/session_name.go` from two different
repos and hoping both readings match.

## 1. The separator alphabet

A qualified identity is built from name segments joined by structural
separators. Exactly two characters are reserved as **positional separators**
in a raw (unencoded) qualified identity:

- `/` — the rig/agent boundary (`rig/agent`)
- `.` — the city/agent boundary on an imported identity (`city.agent`)

These two characters are always positional, never part of a name segment. A
literal `-` or `_`, by contrast, is an ordinary character a name segment MAY
contain: gascity's own agent-name validation permits both
(`internal/config/config.go:27`, enforced at `:4016` — a name must match
`[a-zA-Z0-9][a-zA-Z0-9_-]*`, so `hello-world` and `builder-1` are valid agent
names in their own right, not `/`-delimited pairs).

Encoding (§2) represents each positional separator as a doubled pair rather
than touching `-` or `_` on their own: a literal `-`, doubled, stands in for
`/` (as `--`); a literal `_`, doubled, stands in for `.` (as `__`). A single
`-` or single `_` is never itself a positional separator — only the exact
two-byte sequences `--` and `__` carry structural meaning under this
contract.

## 2. The two-axis encoding table

| Structural meaning | Raw character | Encoded form |
|---|---|---|
| Rig/agent boundary | `/` | `--` |
| City/agent boundary (imported identity) | `.` | `__` |

The two axes are kept distinct on purpose: `/` and `.` encode to different
two-character sequences (`--` vs `__`) so that decoding is unambiguous. A
canonicalizer implementing this rule MUST NOT collapse `--` and `__` to the
same normal form — doing so can silently merge two distinct minted
identities into one (for example, treating `hello-world/polecat` and
`hello_world.polecat` as the same key after canonicalization).

## 3. Who mints, who compares

**gascity** is the only side that **mints** qualified identities. It
composes `rig/agent` and `city.agent` strings when constructing an agent's
qualified name, and it is the only codebase that calls the encode direction
(§4) to produce a tmux-safe session name from one.

**beads** never mints and has no equivalent of the encode direction. It does,
however, independently re-derive part of the *decode* direction for
comparison purposes: `canonicalActor` (beads repo,
`internal/storage/issueops/identity.go:50`, duplicated at
`internal/validation/issue.go:112` because storage may not import
validation) special-cases an exact `--` run and decodes it to `/`, so a
rig-qualified identity compares equal regardless of which spelling it
arrived in. It does NOT restore `__` to `.` the way gascity's own
`UnsanitizeQualifiedNameFromSession` (§4) does — instead it collapses `__`,
and any other single or mixed run of `.`/`_`/`-`, to a generic `_`. This is
the independent re-derivation that let the two sides drift into being
complements before this contract existed (see the architecture ruling this
contract implements, `ga-qv1d2d`). A canonicalizer implementing this rule
MUST NOT collapse `--` and `__` to that same generic form: doing so is the
widening this contract exists to prevent.

## 4. Source of truth

The tables above describe the code; they do not replace it. The
authoritative implementation is `internal/agent/session_name.go`:

- `sessionNameQualifiedReplacer` (`internal/agent/session_name.go:15-18`) —
  the encode direction (§2, raw → encoded).
- `sessionNameQualifiedReverseReplacer` (`internal/agent/session_name.go:20-23`) —
  the decode direction (§2, encoded → raw).

If this document and `internal/agent/session_name.go` ever disagree, the
code wins. File a correction against this doc rather than reimplementing
the table anywhere else.

## 5. Worked examples

| Qualified identity | Encoded session name |
|---|---|
| `mayor` | `mayor` |
| `hello-world/polecat` | `hello-world--polecat` |
| `gastown.mayor` | `gastown__mayor` |
