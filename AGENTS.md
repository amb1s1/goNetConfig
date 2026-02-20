# AGENTS.md

This file helps AI agents continue work in this repository with minimal ramp-up.

## Project Summary

- Project: `goNetConfig`
- Language: Go
- Purpose: Generate network config features over gRPC.
- Multi-entity support: each entity can enable different features and inject feature-specific params.

## Core Runtime Flow

1. `main.go` starts, registers feature generators, loads entity config, resolves `selected_entity`.
2. `service.Registry.SetEntity(...)` injects entity params into each generator and applies feature filtering.
3. `GetConfigGen` iterates registered generators and returns rendered configs for supported features.

## Important Paths

- `/Users/dgo/Desktop/goproject/goNetConfig/main.go`
- `/Users/dgo/Desktop/goproject/goNetConfig/service/service.go`
- `/Users/dgo/Desktop/goproject/goNetConfig/entity/entity.go`
- `/Users/dgo/Desktop/goproject/goNetConfig/configs/entity.yml`
- `/Users/dgo/Desktop/goproject/goNetConfig/configs/samples/`
- `/Users/dgo/Desktop/goproject/goNetConfig/features/aaa/`
- `/Users/dgo/Desktop/goproject/goNetConfig/features/logger/`

## Entity Config Contract

Entity config schema is in `entity.Config`:

- `params.<paramSet>.<feature>.<key>: value`
- `features.<featureSet>: [feature1, feature2, ...]`
- `entity.<entityName>.features: <featureSet>`
- `entity.<entityName>.params: <paramSet>`
- `selected_entity: <entityName>`

If no user config exists, the app falls back to `/Users/dgo/Desktop/goproject/goNetConfig/configs/entity.yml`.

## Test Strategy

- Unit tests:
  - entity parsing/selection: `/Users/dgo/Desktop/goproject/goNetConfig/entity/entity_test.go`
  - feature rendering: `/Users/dgo/Desktop/goproject/goNetConfig/features/*/*_test.go`
- Service integration tests:
  - `/Users/dgo/Desktop/goproject/goNetConfig/service/service_test.go`
  - includes sample config driven checks in `/Users/dgo/Desktop/goproject/goNetConfig/configs/samples/`

Run:

```bash
go test ./...
```

## Working Rules For Agents

- Do not remove existing generators from `allGenerator` without explicit request.
- Keep generator feature names stable (`aaa`, `logger`, etc.); entity mapping depends on them.
- When adding a feature:
  1. implement generator package under `features/`
  2. include `SetParams(map[string]interface{})`
  3. register it in `service/service.go`
  4. add entity sample params/features
  5. add or extend tests
- Avoid global-state leaks in tests; restore `service.Registry` when replacing it.
