#!/usr/bin/env bash
set -euo pipefail

bucket=${BUCKET_NAME:?BUCKET_NAME is required}
manifest=${SNAPSHOT_MANIFEST:-testdata-snapshot.json}
archive=${SAMPLE_ARCHIVE:-testdata.zip}
expected_samples=${NUM_SAMPLE:?NUM_SAMPLE is required}

if find testdata -type f -name 'bak_*' -print -quit | grep -q .; then
    echo "backup files remain under testdata; finish or restore annotation first" >&2
    exit 1
fi

pdf_samples=$(find testdata/pdf -maxdepth 1 -type f -name 'sample*.pdf' | wc -l | tr -d ' ')
yaml_samples=$(find testdata/yaml -maxdepth 1 -type f -name 'sample*.yaml' | wc -l | tr -d ' ')
if [[ "$pdf_samples" != "$expected_samples" || "$yaml_samples" != "$expected_samples" ]]; then
    echo "sample count mismatch: expected $expected_samples, PDF $pdf_samples, YAML $yaml_samples" >&2
    exit 1
fi

rm -f "$archive"
python3 - "$archive" <<'PY'
import pathlib
import shutil
import stat
import sys
import zipfile

archive = sys.argv[1]
files = sorted(
    path for path in pathlib.Path("testdata").rglob("*")
    if path.is_file() and path.name != ".DS_Store" and not path.name.startswith("bak_")
)
with zipfile.ZipFile(archive, "w", zipfile.ZIP_DEFLATED, allowZip64=True, compresslevel=6) as output:
    for path in files:
        info = zipfile.ZipInfo(path.as_posix(), (1980, 1, 1, 0, 0, 0))
        info.compress_type = zipfile.ZIP_DEFLATED
        info.external_attr = (stat.S_IFREG | path.stat().st_mode & 0o777) << 16
        with path.open("rb") as source, output.open(info, "w") as destination:
            shutil.copyfileobj(source, destination, length=1024 * 1024)
PY
sha256=$(python3 - "$archive" <<'PY'
import hashlib
import sys

h = hashlib.sha256()
with open(sys.argv[1], "rb") as f:
    for chunk in iter(lambda: f.read(1024 * 1024), b""):
        h.update(chunk)
print(h.hexdigest())
PY
)
object="snapshots/$sha256.zip"
temporary_manifest=$(mktemp "${manifest}.tmp.XXXXXX")
trap 'rm -f "$temporary_manifest"' EXIT

python3 - "$temporary_manifest" "$object" "$sha256" "$expected_samples" <<'PY'
import json
import sys

path, object_name, sha256, samples = sys.argv[1:]
with open(path, "w", encoding="utf-8") as f:
    json.dump({
        "version": 1,
        "object": object_name,
        "sha256": sha256,
        "samples": int(samples),
    }, f, ensure_ascii=False, indent=2)
    f.write("\n")
PY

wrangler r2 object put "$bucket/$object" --file "$archive" --remote
mv "$temporary_manifest" "$manifest"
trap - EXIT

echo "Uploaded sample snapshot: $object ($expected_samples samples)"
