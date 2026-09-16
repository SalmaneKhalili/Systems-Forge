#!/usr/bin/env python3
"""Audit a tree for required files; print ok:/missing: lines, one per file."""
import os
import sys


def main():
    root = sys.argv[1]
    required = ["src/main.c", "src/util.h"]
    for rel in required:
        if os.path.isfile(os.path.join(root, rel)):
            print(f"ok: {rel}")
        else:
            print(f"missing: {rel}")


if __name__ == "__main__":
    main()