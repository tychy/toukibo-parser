#!/bin/bash
NUM_SAMPLE=778

SKIP_SAMPLES=(251 481 674 690 711 770)

# スクリプト全体でエラーが発生したら停止する
set -e

for ((i = 1; i <= NUM_SAMPLE; i++)); do
    # サンプル番号がスキップリストに含まれていない場合に処理を実行
    if ! [[ " ${SKIP_SAMPLES[*]} " =~ " $i " ]]; then
        echo "sample$i"
        ./bin/toukibo-parser -path="sample$i"
    fi
done