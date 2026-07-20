#!/usr/bin/env bash
set -euo pipefail

DATA_DIR=testdata

for ((i = 1; i <= NUM_SAMPLE; i++)); do
    echo "sample$i"
    yaml="$DATA_DIR/yaml/sample$i.yaml"
    temporary=$(mktemp "$DATA_DIR/yaml/sample$i.yaml.tmp.XXXXXX")
    trap 'rm -f "$temporary"' EXIT

    ./bin/toukibo-parser -mode=run -path="$DATA_DIR/pdf/sample$i.pdf" > "$temporary"
    diff "$yaml" "$temporary" || true
    mv "$temporary" "$yaml"
    trap - EXIT
done
