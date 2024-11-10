#!/bin/bash
# スクリプト全体でエラーが発生したら停止する
set -e

FILE_PATH=/Users/ryu/Downloads/20241110
IDX=1245

for file in $FILE_PATH/*; do
    if [ -f "$file" ]; then
        mv "$file" "$FILE_PATH/sample$(printf '%04d' $IDX).pdf"
        ((IDX++))
    fi
done
