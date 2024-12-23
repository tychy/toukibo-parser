#!/bin/bash
# スクリプト全体でエラーが発生したら停止する
set -e

DATA_DIR=testdata/pdf

rm -f $DATA_DIR/sample$IDX.pdf
for ((i = $IDX; i <= NUM_SAMPLE; i++)); do
    # ファイル名を変更する
    IDX=$((IDX + 1))
    mv $DATA_DIR/sample$IDX.pdf $DATA_DIR/sample$i.pdf
done
