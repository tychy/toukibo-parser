#!/bin/bash
NUM_SAMPLE=778

SKIP_SAMPLES=(162 192 228 251 267 298 420 473 \
 481 690 711 761 770)
SKIP_SAMPLES_GAIJI=(87 117 206 207 213 223 232 237 278 283 390 425 451 479 \
 499 518 653 668 674 684 700 716)
SKIP_SAMPLES_SHAIN=(111 127 133 493 521 704)

# スクリプト全体でエラーが発生したら停止する
set -e

for ((i = 1; i <= NUM_SAMPLE; i++)); do
    # サンプル番号がスキップリストに含まれていない場合に処理を実行
    if ! [[ " ${SKIP_SAMPLES[*]} " =~ " $i " || " ${SKIP_SAMPLES_SHAIN[*]} " =~ " $i " || " ${SKIP_SAMPLES_GAIJI[*]} " =~ " $i " ]]; then
        echo "sample$i"
        ./bin/toukibo-parser -path="sample$i"
    fi
done