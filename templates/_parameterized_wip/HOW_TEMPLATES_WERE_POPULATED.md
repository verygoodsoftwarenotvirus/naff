# How `templates/` Was Populated

The pivot on 2026-04-14 replaced the old registry-driven parameterized
templates with a verbatim mirror of the upstream `dinnerdonebetter` (DDB)
reference project. This note documents the exact commands used, so the
mirror is reproducible.

## Source

```
/Users/jeffrey/src/github.com/dinnerdonebetter/dinnerdonebetter
```

Subtrees mirrored (one rsync per subtree):

- `backend/`  →  `templates/_backend/` (prefixed with `_` so Go tooling
  skips the embedded `.go` files instead of trying to compile them as part
  of naff).
- `proto/`    →  `templates/proto/`
- `frontend/` →  `templates/frontend/`
- `ios/`      →  `templates/ios/`

## Command

```sh
DDB=/absolute/path/to/dinnerdonebetter
DST=/absolute/path/to/naff/templates

for sub in backend proto frontend ios; do
  rsync -a --exclude-from=ddb-rsync-excludes.txt "$DDB/$sub/" "$DST/_$sub/"
  # (only `backend` actually needs the `_` prefix; the rename below corrects
  # proto/frontend/ios back to their real names)
done
mv "$DST/_proto" "$DST/proto"
mv "$DST/_frontend" "$DST/frontend"
mv "$DST/_ios" "$DST/ios"

# embed can't cross a nested Go module boundary, so go.mod/go.sum in
# _backend/ get the `_` prefix too. The pipeline renames them back on emit.
mv "$DST/_backend/go.mod" "$DST/_backend/_go.mod"
mv "$DST/_backend/go.sum" "$DST/_backend/_go.sum"

# The rsync filter missed the compiled mcp binary and frontend build dirs
# because their exclude patterns were root-anchored. Remove manually:
rm -f  "$DST/_backend/mcp"
rm -rf "$DST/frontend/admin/build" "$DST/frontend/consumer/build"
```

## Exclusion filter

The `ddb-rsync-excludes.txt` file used:

```
.git/
.DS_Store
node_modules/
vendor/
.terraform/
.terraform.lock.hcl
DerivedData/
.svelte-kit/
.next/
.turbo/
.cache/
.parcel-cache/
.gradle/
.idea/
.vscode/
/backend/deploy/environments/*/terraform/.terraform/
/backend/mcp
/backend/artifacts/
/frontend/admin/build/
/frontend/consumer/build/
/ios/build/
*.xcuserdata/
*.xcuserstate
.swiftpm/
*.ipa
*.zip
*.tar
*.tar.gz
*.tgz
*.log
tmp/
coverage/
```

## Refreshing from upstream

Re-run the same `rsync` commands. The pipeline is idempotent (diff-aware
writes), and `make generate-ddb` always starts from a clean `.ddb_generated`
directory.
