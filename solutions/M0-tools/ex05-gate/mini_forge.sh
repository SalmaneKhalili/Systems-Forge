#!/bin/sh
set -eu

if ! make -C sample fclean; then
    echo "FAIL"
    exit 1
fi

if ! make -C sample all; then
    echo "FAIL"
    exit 1
fi

if ! ./sample/hello | cmp -s - sample/expected.txt; then
    echo "FAIL"
    exit 1
fi

echo "PASS"
exit 0