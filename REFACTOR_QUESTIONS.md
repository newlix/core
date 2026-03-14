# Refactor Questions

## [generators/golang/out.go, generators/swift/out.go, generators/kotlin/out.go, generators/sqlc/out.go]
- **Problem**: The `out()` function is duplicated identically in 5 packages (common + 4 generators). Each generator defines its own local copy instead of importing from `common`.
- **Options**: (A) Export `common.Out` and import everywhere — changes all call sites from `out(w, ...)` to `common.Out(w, ...)`. (B) Keep local copies for call-site convenience — accepted Go pattern for small helpers.
- **Status**: Not changed (skipped)

## [generators/golang/type.go, generators/swift/type.go, generators/kotlin/type.go]
- **Problem**: `GenerateMethodTypes` accepts a `tt []core.Type` / `ts []core.Type` parameter that is never used in the function body (golang, swift, kotlin). Only the sqlc generator actually uses it.
- **Options**: (A) Remove the unused parameter — this is a public API change that may break downstream callers. (B) Keep it for future use or API consistency across generators.
- **Status**: Not changed (skipped)
