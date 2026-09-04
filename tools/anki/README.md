# Anki export of the review decks

Every exercise's Q&A review metadata (`answers[]` in `exercise.json`) is a spaced-
repetition review card. This tool regenerates Anki-importable decks directly from the
canon (`subjects/*/exercise.json`), so the decks never drift from the curriculum.
Quizzing lives entirely in Anki — it is not part of exercise grading.

## Regenerate

```sh
python3 tools/anki/export.py        # rewrite tools/anki/decks/*.txt (18 decks, 241 cards)
```

## Import into Anki

1. Anki → File → Import, pick `tools/anki/decks/M12-raft.txt` (repeat per module).
2. Indices and fields are already set for the default **Basic** note type:
   - Field 1 = question (`[module · exNN · title] <question>`)
   - Field 2 = answer
   - Files are **tab-separated**, so you can also paste into Anki's "Import from text";
     be sure *Fields separated by* = **Tab**.
3. Optionally rename the imported notes' deck to `systems-forge::<Module Title>` and add a
   `systems-forge` tag. Because each card's front carries the module/exercise tag, you can
   build a **filtered deck** for "this week's modules" at any time.

## Notes

- Source of truth is `subjects/` (tracked canon), never the gitignored `solutions/`/`answers/`.
- Re-importing after a curriculum change imports only new cards (questions are duplicated
  per deck); delete the old notes first to avoid drift if a question's wording changed.