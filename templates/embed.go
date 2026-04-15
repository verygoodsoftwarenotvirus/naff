package templates

import "embed"

// FS contains all embedded template files. Post-pivot the `templates/`
// directory holds byte-for-byte copies of the upstream reference project
// under `backend/`, `proto/`, `frontend/`, and `ios/`. The pipeline walks
// this filesystem and emits every entry as a raw pass-through.
//
// The `_backend/` subtree is prefixed with `_` so Go tooling treats it as
// hidden and does NOT try to compile its `.go` files as part of the naff
// module. The `all:` prefix makes `go:embed` include underscore-prefixed
// files and directories anyway. The pipeline strips the leading `_` when
// computing output paths, so `_backend/go.mod` emits as `backend/go.mod`.
//
//go:embed all:_root all:_backend all:proto all:frontend all:ios
var FS embed.FS
