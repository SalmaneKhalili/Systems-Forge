#!/usr/bin/env python3
from pathlib import Path
import sys

P = Path(sys.argv[1])

main = P / 'src' / 'main.c'
util = P / 'src' / 'util.h'

if main.exists():
    print("ok: src/main.c")
else:
    print("missing: src/main.c")

if util.exists():
    print("ok: src/util.h")
else: 
    print("missing: src/util.h")
