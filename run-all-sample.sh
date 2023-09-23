#!/bin/bash
NUM_SAMPLE=778

SKIP_SAMPLES=(111 117 127 133 162 163 179 192 206 207 213 214 223 \
 228 232 236 237 239 251 267 272 278 283 298 387 390 400 420 425 451 473 479 \
 481 493 499 518 521 543 653 668 674 675 684 690 700 704 711 716 761 770)
SKIP_SAMPLES_GAIJI=(87)

# スクリプト全体でエラーが発生したら停止する
set -e

for ((i = 1; i <= NUM_SAMPLE; i++)); do
    # サンプル番号がスキップリストに含まれていない場合に処理を実行
    if ! [[ " ${SKIP_SAMPLES[*]} " =~ " $i " || " ${SKIP_SAMPLES_GAIJI[*]} " =~ " $i " ]]; then
        echo "sample$i"
        ./bin/toukibo-parser -path="sample$i"
    fi
done