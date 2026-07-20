#!/usr/bin/env bash
set -euo pipefail

manifest=${SNAPSHOT_MANIFEST:-testdata-snapshot.json}
url=${SAMPLE_URL:?SAMPLE_URL is required}
archive=${SAMPLE_ARCHIVE:-testdata.zip}

IFS=$'\t' read -r object expected_sha256 expected_samples < <(python3 - "$manifest" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as f:
    manifest = json.load(f)
for key in ("object", "sha256", "samples"):
    if key not in manifest:
        raise SystemExit(f"snapshot manifest is missing {key}")
print(manifest["object"], manifest["sha256"], manifest["samples"], sep="\t")
PY
)

curl --fail --location --silent --show-error --output "$archive" "$url/$object"
actual_sha256=$(python3 - "$archive" <<'PY'
import hashlib
import sys

h = hashlib.sha256()
with open(sys.argv[1], "rb") as f:
    for chunk in iter(lambda: f.read(1024 * 1024), b""):
        h.update(chunk)
print(h.hexdigest())
PY
)
if [[ "$actual_sha256" != "$expected_sha256" ]]; then
    echo "sample snapshot checksum mismatch: expected $expected_sha256, got $actual_sha256" >&2
    exit 1
fi

unzip -q "$archive"
actual_samples=$(find testdata/pdf -maxdepth 1 -type f -name 'sample*.pdf' | wc -l | tr -d ' ')
if [[ "$actual_samples" != "$expected_samples" ]]; then
    echo "sample count mismatch: expected $expected_samples, got $actual_samples" >&2
    exit 1
fi

echo "Downloaded sample snapshot: $object ($actual_samples samples)"
