# naff

Code generator that scaffolds a full-stack service repository (Go backend,
Proto/gRPC, TypeScript frontend, iOS) from a YAML schema. The canonical
parity target is `../../dinnerdonebetter/dinnerdonebetter` — a generated
project should reach feature and file parity with that repo when given an
equivalent `.naff.yaml`.

> **Terminology:** "DDB" in this project always refers to **Dinner Done Better**
> (the parity-target repo), never DynamoDB. If the user mentions DDB, they
> mean `../../dinnerdonebetter/dinnerdonebetter`.

## Core mental model

**naff is a single-pass scaffolder, not a re-generator.** It emits into a
fresh output directory once; the user then commits the result and evolves
the tree by hand. There is no `_generated/` marker directory, no orphan
detection, no re-run-safety layer. Do not add those.

Consequences that matter when working in this repo:

- Every emitted file must be correct on first write. There is no
  "regenerate to fix" escape hatch.
- Users will hand-edit generated code. Don't rely on regeneration to
  propagate changes — template changes benefit future generations, not
  existing trees.
- The `--clean` flag is accepted for CLI compat but is a no-op.

## Repository layout

- `cmd/naff/` — CLI entrypoint (`naff generate`, `naff validate`, etc.)
- `internal/config/` — `.naff.yaml` schema types + validation
- `internal/naming/` — Pascal/Snake/Kebab/Camel casing helpers (used
  heavily in templates via funcmap)
- `internal/renderer/` — text/template wrapper + Go source formatter +
  shared funcmap
- `internal/pipeline/` — plans + renders + writes + runs post-gen make
  targets. Start here to understand generation flow.
- `internal/builtins/` — internal helpers available during generation
- `templates/` — the template tree, `go:embed`ed into the binary
- `testdata/` — example `.naff.yaml` configs used by tests and local
  smoke runs
- `docs/parity_checklist.md` — package-level parity map against the
  dinnerdonebetter target (most surfaces are still verbatim-only)

## The templates/ tree

`templates/` is `go:embed`ed whole. Two conventions you must understand:

**Leading `_` hides a top-level directory from Go tooling.** `_backend/`
emits as `backend/`. Same for `_go.mod` → `go.mod` and `_go.sum` →
`go.sum` — nested real `go.mod` files would create a module boundary
that `go:embed` refuses to cross. These renames happen in
`rewriteOutputPath` in `internal/pipeline/pipeline.go`.

**`_root/` is a collector** for files that should emit at the output
repo root (e.g. top-level `Makefile`, `.github/`, `README.md`). The
`_root/` prefix is stripped.

**`.tmpl` suffix triggers text/template rendering.** Files without
`.tmpl` are pass-through byte copies (`Raw: true`) — safe for files
containing literal `{{` that would confuse text/template.

**`.symlink` suffix is a symlink sentinel.** The file's contents are
the symlink target string; the emitter creates a real symlink at the
de-suffixed output path. This exists because `go:embed` silently drops
actual symlinks.

**`_backend/_parameterized/` is reserved** for templates driven by the
per-domain / per-entity planner loops, not the verbatim FS walk. The
planner explicitly `SkipDir`s this subtree and emits its contents via
separate planning code. If you add a parameterized template, add it
here and wire it into the planner; do not drop it into the verbatim
tree.

## Verbatim vs parameterized

Two generation modes coexist:

1. **Verbatim walk.** `planFiles` walks `templates/` and plans one
   output per file. Templates reference the project via shared data
   (`.Project.ProjectMeta.*`). Most of the current tree is verbatim,
   which means the generated project only works if the user's schema
   *happens to match* upstream. This is the status quo, not the goal.

2. **Parameterized layer.** A second pass iterates
   `config.Domains` (and eventually `Domain.Entities`) and emits one
   file per iteration from templates under
   `_backend/_parameterized/`. Today only
   `domain_keys/keys.go.tmpl` is live. Per-entity templates exist
   under `domain_entity/` but are disabled — see the long comment in
   `pipeline.go` around `planPerEntityFiles`.

Moving a surface from verbatim to parameterized is the primary work
remaining in this repo. See `docs/parity_checklist.md`.

## How generation runs

`pipeline.Run` executes these phases:

1. **Plan** — build a `[]PlannedFile` via `planFiles` (verbatim walk +
   parameterized loops).
2. **Render + write** — parallel, diff-aware writes. `.tmpl` files go
   through text/template and (if `IsGo`) `go/format`. Raw files are
   byte-copied. Symlinks are materialized.
3. **Write `.naff.yaml`** — the resolved config is persisted at the
   output root so the project knows what it was generated from.
4. **Post-gen make targets** run in the output tree, in order:
   - `backend`: `make querier format vendor configs env_vars`
   - root: `make proto`

The order is load-bearing: `proto` consumes sqlc-generated types from
the backend, so `querier` must run first. `vendor` resolves the module
graph before formatting.

After `naff generate` finishes, the downstream workflow the user runs is
`make proto vendor generated_files` in the output tree — keep that in
mind if you're reasoning about what state is Claude's responsibility
vs. the user's.

## Dos and don'ts

### Do

- **Treat the parity target as ground truth.** When in doubt about what
  a template should emit, read the corresponding file in
  `../../dinnerdonebetter/dinnerdonebetter` and match it.
- **End-to-end test your changes.** Run `naff generate` against a
  testdata config, confirm the output compiles (`make vendor && go
  build ./...` in the generated tree), and spot-check a diff against
  upstream. Template-level unit tests catch syntax errors but not
  rendered-output regressions.
- **Extend the schema before adding feature flags.** If a template
  needs a new behavior, prefer adding a typed field to
  `internal/config/types.go` over conditional logic keyed off project
  name or similar implicit signals.
- **Keep parameterized templates in `_backend/_parameterized/`** and
  wire them through the planner explicitly. The `skipVerbatim` set in
  `planFiles` prevents collisions — use it.
- **Preserve the `_` / `.symlink` / `.tmpl` conventions.** They're the
  load-bearing part of how the FS maps to output.
- **Update `docs/parity_checklist.md`** when you move a surface from
  verbatim to parameterized.

### Don't

- **Don't add re-generation-safe scaffolding.** No orphan detection,
  no `_generated/` markers, no "is this file ours?" heuristics. Single-pass only.
- **Don't compile or leave binaries.** If you must build to verify
  something, delete the artifact immediately. (Global rule, restated
  here because this repo tempts you to `go build` the generator.)
- **Don't enable per-entity parameterized templates** (the
  `domain_entity/` subtree) until they're at parity with the
  hand-written per-entity templates. Enabling them now produces code
  that doesn't compile — see the long comment in `pipeline.go:522-533`.
- **Don't write a template that assumes `mealplanning` exists.** Any
  verbatim file referencing a specific domain name is a latent bug for
  non-ddb projects. Prefer moving it to parameterized, or gating it
  clearly in planning.
- **Don't conflate pipeline phases.** Rendering is pure; writes are
  the only side effect; make invocations run last. Don't invoke make
  from inside a template or renderer.
- **Don't hand-format Go output.** `renderer.FormatGo` runs `go/format`
  on every `IsGo: true` file. If formatting fails, the template has a
  real syntax issue — fix the template, don't paper over it by
  skipping format.

## Common tasks

**Add a new verbatim file to the generated tree.** Drop it under
`templates/_<target>/...`. Add `.tmpl` if you need template rendering,
omit it for a byte-copy. That's it — the FS walk picks it up.

**Parameterize an existing verbatim file.** (1) Move/copy it under
`_backend/_parameterized/<group>/`. (2) Add a loop in `planFiles` that
emits one `PlannedFile` per domain (or entity). (3) Add its old verbatim
path to `skipVerbatim` so the FS walk doesn't also emit it.
`domain_keys` is the reference implementation.

**Debug a render failure.** `naff generate --debug` streams post-gen
make output live. For template-level errors, add logging in
`renderer.Render` or run the specific template through `text/template`
directly in a test.

**Add a field to the schema.** Add the YAML-tagged struct field in
`internal/config/types.go`, add validation in `internal/config/validation.go`,
add any derived methods (casing, defaults), then reference it in
templates via `.Project` / `.Domain` / `.Entity`.

## Gotchas

- **`embed.FS` flattens file modes to `0o444`.** `modeFor` in
  `pipeline.go` restores `+x` for `.sh` files and anything under a
  `scripts/` directory. If you add executables elsewhere, extend
  `modeFor`.
- **Template data shape differs between layers.** Verbatim templates
  get `map[string]any{"Project": cfg}`. Parameterized templates get
  `parameterizedCtx{Project, Domain, Entity}`. Don't reference
  `.Domain` from a verbatim template.
- **The `_` prefix rule only strips the top-level segment.** A file
  deeper in the tree named `_terraform.tf` keeps its underscore.
- **Adding a nested `go.mod` inside `templates/`** will break the
  embed. Use the `_go.mod` sentinel at the directory level that needs it.

## Related docs

- `docs/parity_checklist.md` — package-level checklist of schema-derived
  surfaces vs. what's currently parameterized.
