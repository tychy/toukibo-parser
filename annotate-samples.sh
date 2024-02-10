#!/bin/bash
NUM_SAMPLE=190

SKIP_SAMPLES=(674)

# スクリプト全体でエラーが発生したら停止する
set -e

DATA_DIR=testdata

for ((i = 1; i <= NUM_SAMPLE; i++)); do
    # サンプル番号がスキップリストに含まれていない場合に処理を実行
    if ! [[ " ${SKIP_SAMPLES[*]} " =~ " $i " ]]; then
        echo "sample$i"
    	./bin/toukibo-parser -path="$DATA_DIR/pdf/sample$i.pdf" > $DATA_DIR/yaml/sample$i.yaml
    fi
done
