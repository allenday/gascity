---
title: Gas City Identity Separator Contract — v1
description: Authoritative specification for the qualified-identity encoding shared by gascity and beads.
---

| Field | Value |
|---|---|
| Status | Authoritative specification |
| Last verified | 2026-08-18 |
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
separators. Three characters are reserved as **positional separators**:

- `/` — the rig/agent boundary (`rig/agent`)
- `.` — the city/agent boundary on an imported identity (`city.agent`)
- `_` — reserved alongside them because the encoded form uses it (§2)

These three characters are always positional separators, never part of a name.
A name segment does not itself contain `/`, `.`, or a double underscore.
Encoding only has to disambiguate `/` and `.`, because those are the two
separators that appear in a raw qualified identity; `_` appears solely as
half of an encoded pair (`--` or `__`), never on its own.

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

**beads** only ever **compares** already-minted strings. It stores the
encoded form verbatim — for example as session-name metadata on a bead —
and matches it byte-for-byte against what gascity re-derives. beads does not
independently encode or decode a qualified identity, and MUST NOT implement
its own copy of the separator table in §2: a second implementation is
exactly how the two sides drift apart.

## 4. Source of truth

The tables above describe the code; they do not replace it. The
authoritative implementation is `internal/agent/session_name.go`:

- `sessionNameQualifiedReplacer` — the encode direction (§2, raw → encoded).
- `sessionNameQualifiedReverseReplacer` — the decode direction (§2, encoded → raw).

If this document and `internal/agent/session_name.go` ever disagree, the
code wins. File a correction against this doc rather than reimplementing
the table anywhere else.

## 5. Worked examples

| Qualified identity | Encoded session name |
|---|---|
| `mayor` | `mayor` |
| `hello-world/polecat` | `hello-world--polecat` |
| `gastown.mayor` | `gastown__mayor` |
