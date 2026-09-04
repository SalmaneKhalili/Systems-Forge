#!/usr/bin/env python3
"""Export the curriculum's Q&A review decks (subjects/*/exercise.json answers[])
as Anki-compatible card files.

One card = one question/answer pair from exercise.json. Cards are written to
tools/anki/decks/<MODULE>.txt in Anki import format (Front<TAB>Back), one per
line. Import in Anki via File > Import (tab-separated; default "Basic" note type
works, or use a Cloze-less basic). Each card carries a tag line on the Front:
[<module> · ex<nn> · <title>] so you can review per-module with a filtered deck.

Re-run whenever the curriculum changes: `python3 tools/anki/export.py`.
"""
import glob
import json
import os

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))  # repo root
SUBS = os.path.join(ROOT, "subjects")
OUT = os.path.join(ROOT, "tools", "anki", "decks")

# Map module dir -> human deck name; module.md titles are the source, keep tidy.
TITLES = {
    "M0-tools": "Tools of the Trade",
    "M1-clib": "The C gate",
    "M2-procs": "Processes",
    "M3-memory": "Memory",
    "M4-concurrency": "Concurrency",
    "M5-filesio": "Files & I/O",
    "M6-networking": "Networking",
    "M7-resilience": "Resilience",
    "M8-messaging": "Messaging + the Switch",
    "M9-time": "Time & Ordering",
    "M10-commit": "Commit & Consensus",
    "M11-log": "Replicated Log & Consistency",
    "M12-raft": "Raft (flagship capstone)",
    "M13-shard": "Sharding",
    "M14-membership": "Membership",
    "M15-tx": "Transactions & Chaos",
    "M16-storage": "Storage Engines",
    "M17-observability": "Observability",
}


def ex_key(text: str) -> str:
    """Convert 'M12-ex06' style id to 'ex06' tag fragment."""
    return text.split("-ex")[-1].split("-")[0]


def main():
    os.makedirs(OUT, exist_ok=True)
    by_mod = {}
    total = 0
    for ex in sorted(glob.glob(os.path.join(SUBS, "*", "ex*", "exercise.json"))):
        d = json.load(open(ex))
        mod = d.get("module", "?")
        qas = d.get("answers", [])
        exid = d.get("id", "?")
        title = d.get("title", "")
        exn = ex_key(exid)
        for qa in qas:
            q, a = qa.get("q", ""), qa.get("a", "")
            for ch in ("\n", "\t", "\r"):
                q = q.replace(ch, " ")
                a = a.replace(ch, " ")
            front = f"[{mod} · {exn} · {title}] {q}"
            by_mod.setdefault(mod, []).append(f"{front}\t{a}\n")
            total += 1

    for mod, lines in sorted(by_mod.items()):
        deck = mod  # filename-safe module id (avoid '/' etc. in titles)
        with open(os.path.join(OUT, deck + ".txt"), "w") as f:
            f.writelines(lines)
        print(f"{mod:16} {len(lines):3} cards  -> {deck}.txt (import under 'systems-forge::{TITLES.get(mod,mod)}')")
    print(f"total {total} cards across {len(by_mod)} decks")


if __name__ == "__main__":
    main()