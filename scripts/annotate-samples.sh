#!/bin/bash
# スクリプト全体でエラーが発生したら停止する
set -e

DATA_DIR=testdata

for ((i = 1; i <= NUM_SAMPLE; i++)); do
    echo "sample$i"
    mv $DATA_DIR/yaml/sample$i.yaml $DATA_DIR/yaml/bak_sample$i.yaml
    ./bin/toukibo-parser -mode="run" -path="$DATA_DIR/pdf/sample$i.pdf" > $DATA_DIR/yaml/sample$i.yaml
    diff ./$DATA_DIR/yaml/bak_sample$i.yaml ./$DATA_DIR/yaml/sample$i.yaml || true
    rm $DATA_DIR/yaml/bak_sample$i.yaml
done
