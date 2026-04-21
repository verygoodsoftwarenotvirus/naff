# Schema-Derived Package Parity Checklist

Parity target: `../../dinnerdonebetter/dinnerdonebetter`.

Currently parameterized:
- `backend/internal/domain/<domain>/keys/keys.go`
- `backend/internal/domain/<domain>/<entity>.go` (CRUD core: struct,
  input variants, `<Entity>DataManager` + `<Entity>DataService`
  interfaces, `Update`, validation)

Everything else below is either emitted verbatim (only works if the user's
schema happens to match upstream mealplanning) or not emitted at all.

## Per-domain packages

One instance per entry in `config.Domains`.

### Domain layer — `backend/internal/domain/<domain>/`

- [x] `keys/keys.go`
- [~] `<entity>.go` — CRUD core (struct, inputs, `<Entity>DataManager`,
  `Update`, validation, `Convert<Entity>ToDatabaseCreationInput`) emitted
  from schema. `<Entity>DataService` intentionally not emitted — DDB
  no longer uses this shape.
  Known drifts from hand-written upstream:
    - Collection fields on parents (`MealList.Items []*MealListItem`,
      `Meal.Components []*MealComponent`) are *not* emitted; would
      require a per-parent helper that scans sibling entities for
      matching `belongs_to`.
    - Embedded entity structs (e.g. `MealListItem.Meal Meal`) are not
      derivable — the schema has no "include parent by value" signal.
    - Bespoke types (constants, sentinel errors, `Nullable<Entity>`,
      multi-item container inputs) must be added in a sidecar file.
    - Bespoke `DataManager` methods (search indexing, custom state
      transitions, cross-entity lookups) are not emitted.
    - Argument ordering on `ArchiveMealListItem` / similar: template
      uses `(parents..., selfID)` consistently; upstream has a few
      `(selfID, parentID)` outliers.
    - `created_by_user` renders as a `CreatedByUser` field; a few
      upstream entities renamed it to `BelongsToUser` by hand.
- [ ] `<entity>_test.go` (one pair per entity — parity test still TODO)
- [ ] `<domain>.go` + `<domain>_test.go` (top-level domain types)
- [ ] `repository.go` (aggregated DataManager interface)
- [ ] `do.go` (fx providers)
- [ ] `converters/<entity>.go` + `converters/converters.go` + `converters/domain_to_creation_input.go`
- [ ] `fakes/<entity>.go` + `fakes/fake.go` + `fakes/fake_test.go`
- [ ] `mocks/repository.go`

### Postgres layer — `backend/internal/repositories/postgres/<domain>/`

- [ ] `<entity>.go` + `<entity>_test.go` (one pair per entity)
- [ ] `client.go` + `client_test.go` + `do.go`
- [ ] `sqlc_queries/<entity>.generated.sql` (input to `make querier`; `generated/` is sqlc output, not naff-emitted)

### Service layer — `backend/internal/services/<domain>/`

- [ ] `config/`
- [ ] `errors/`
- [ ] `workers/`
- [ ] `indexing/`
- [ ] `grpc/<entity>.go` + tests
- [ ] `grpc/service.go`, `grpc/do.go`, `grpc/converters/`

### Proto + generated bindings

- [ ] `proto/<domain>/<domain>_messages.proto`
- [ ] `proto/<domain>/<domain>_service.proto`
- [ ] `proto/<domain>/<domain>_service_types.proto`
- `backend/internal/grpc/generated/services/<domain>/` — protoc output, not naff-emitted

### Frontend / iOS (when targets enabled)

- [ ] `frontend/packages/api-client/src/generated/<domain>/` (proto → TS)
- [ ] `frontend/consumer` — per-entity route/page scaffolds
- [ ] `frontend/admin` — per-entity route/page scaffolds
- [ ] `ios/ios/Services/<Entity>Service.swift` + view scaffolds

## Global aggregation points

Single files that must enumerate `config.Domains` to stay consistent.

- [ ] `backend/internal/build/server/` — router wiring per service
- [ ] `backend/internal/repositories/do.go` — domain repo fx registration
- [ ] `backend/cmd/server/main.go` — wiring
- [x] `backend/cmd/localdev/server/main.go` — schema-aware for bootstrap
- [ ] `backend/Makefile` / `sqlc.yaml` — domain list drives `make querier`
- [ ] `backend/internal/grpc/` top-level wiring

## Hand-written, not schema-derived

Some per-domain files look like parity candidates but aren't realistically
derivable from the current schema — leave them as hand-written in the
generated project.

- `backend/internal/domain/<domain>/errors.go` — errors encode
  domain-specific invariants (e.g. `ErrDuplicateMealInList`) that don't
  fall out of entity field definitions. Expressing them would require
  richer constraint semantics in the schema than we have today, and
  hand-writing a handful of sentinels per domain is cheaper than that.

## Verbatim (not schema-derived)

Infrastructure domains toggled by `config.Features`, not parameterized per
user-defined entity:

`audit`, `auth`, `identity`, `oauth`, `internalops`, and feature-gated:
`comments`, `dataprivacy`, `issuereports`, `notifications`, `payments`,
`settings`, `uploadedmedia`, `waitlists`, `webhooks`.

## Suggested parity-verification order

1. `backend/internal/domain/<domain>/` top-level per-entity files — the
   parameterized templates under `templates/_backend/_parameterized/domain_entity/`
   exist but are disabled (see `internal/pipeline/pipeline.go:522-533`); bringing
   them to parity is the next logical step.
2. `converters/`, `fakes/`, `mocks/` — purely derivable from entity fields + links.
3. `repositories/postgres/<domain>/` + `sqlc_queries/` — entity fields and links
   drive SQL directly.
4. `proto/<domain>/` — messages map 1:1 to entity fields/links; service RPCs are
   per-entity CRUD.
5. `services/<domain>/grpc/` — thin handlers over repository + converters.
6. Global aggregation files listed above.
7. Frontend + iOS surfaces.
