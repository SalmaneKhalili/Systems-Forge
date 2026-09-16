#!/bin/bash

set -euo pipefail

declare -i sum=0
for arg in "$@"
do 
	sum=$sum+$arg
done
echo $sum