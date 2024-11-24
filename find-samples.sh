#!/bin/bash
# スクリプト全体でエラーが発生したら停止する
set -e

DATA_DIR=testdata
FIND=東方産業

for ((i = 1; i <= NUM_SAMPLE; i++)); do
    ./bin/toukibo-parser -mode="find" -path="$DATA_DIR/pdf/sample$i.pdf" -target="$FIND"
done
