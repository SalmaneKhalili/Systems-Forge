# Reading-Ladder Standard

Every exercise's `## Readings` block is a **learning path**, not a bibliography. A reader
should be able to open the subject with zero context and climb to the exercise without a
dead end. The three failures this standard exists to prevent:

1. **Cold jumps** — citing a book's deepest chapter (TLPI §21.1) while its fundamentals
   (§20) are never mentioned → the reader drowns in context they don't have.
2. **Silent gaps** — subjects with no readings at all → the module teaches nothing.
3. **Vague citations** — "TLPI Chapter 29" with no section anchor → hours of hunting.

## The rules (must hold for every `## Readings` block)

1. **Cold-start entry.** The first bullet gives a zero-context reader a way in: an inline
   glossary/primer, an external beginner source (URL), or a `man` page. Never open with the
   deepest citation.
2. **Ordered climb.** Later bullets go *fundamentals → the specific section → verification*
   (man pages, standard notes). If the primary canon assumes prerequisites, cite the
   prerequisite first (e.g. before TLPI §21.1 signal-handler design, cite §20 signals
   fundamentals).
3. **Precise anchors.** Book references use exact citations: TLPI is always
   `§N.M "Chapter title"`. Every topic that has one gets a `man` anchor.
4. **No gaps.** Every subject in every module has a Readings block. A subject without one
   is a hole rather than a shortcut.
5. **Cross-linked depth.** If the topic returns in a later module, say so ("terminal-signal
   handling and `SIGCHLD` reaping come back at M7-ex04/ex05") so the student knows the
   scope here is deliberately narrow.

## Cold-start primer

When a topic needs a handful of terms defined *before* the book makes sense (signals,
sockets, threads, mmap…), the block opens with a small glossary under an explicit marker:

```
## Readings

- **Cold-start** · signal: a software interrupt; disposition: the process-defined action
  (default / ignore / handler) for a signal; delivery: the point the disposition runs.
- The Linux Programming Interface, §20.1 "The Concept of Signals", …  (then the climb)
```

The lint requires the marker; the pedagogic content of the glossary stays a human call.

## What the machine checks — `forge lint readings`

`forge lint readings` enforces the *structural* invariants; judgment about which source is
easiest first stays human. Rules (see `forge/internal/readings/check.go`):

| rule | check |
|---|---|
| R1 | every `subjects/M*/ex*/subject.md` contains a `## Readings` section |
| R2 | any subject citing the "Linux Programming Interface"/"TLPI" also cites at least one numeric pin (`§N` / `Chapter N`) in the block |
| R3 | the block's first content line (bullet / ordered item / sub-heading) carries a cold-start marker: `https://`, `man `, `**`, `Cold-start`, `glossary`, `Reading ladder`, an italic book title, an `RFC N`, or a chapter-level citation |
| R4 | every precise `man [[sect]] name` citation resolves via `man -w` |
| R5 | (human) the remaining bullets actually climb entry → fundamentals → section → verification |

Run from the repo root. Wired into `tools/selfcheck.sh` so the check never regresses.

---

Enforced since 2026-09-11 (see PLAN.md §12, §14).