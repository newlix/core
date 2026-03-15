# Refactor Questions

Items requiring user decisions before proceeding.

## generators/*: `log.Fatal` in library code

- **Problem**: 20+ `log.Fatal`/`log.Fatalf` calls across all generator packages (`golang`, `swift`, `kotlin`, `sqlc`). `log.Fatal` calls `os.Exit(1)`, skipping defers and preventing error recovery. Commit `4d9e30f` already removed this pattern from `InitTypes`/`InitMethods` in the core package.
- **Options**:
  - A: Change all `Generate*File` functions to return `error`, update the example `cmd/gen/main.go` caller accordingly. Also change internal helpers (`checkPackage`, `findType`, `sqlcType`, `findTableForInputs`) to return errors.
  - B: Keep `log.Fatal` in the top-level `Generate*File` functions (CLI entrypoints) but convert internal helpers to return errors.
- **Status**: Not yet addressed.

## generators/sqlc: Missing WHERE clauses for input+output methods

- **Problem**: `writeQuery` uses a priority switch — outputs-only → SELECT, inputs-only → INSERT. Methods with *both* inputs and outputs (e.g., `RemoveItem`) generate `SELECT ... FROM table;` without a WHERE clause, silently dropping input parameters. The generated query returns all rows instead of filtering.
- **Options**:
  - A: Add a `writeSelectWithWhere` path for methods with both inputs and outputs.
  - B: Treat SQL generation as scaffolding only (document that users must add WHERE clauses manually).
  - C: Rethink the SQL generation model to support richer query semantics.
- **Status**: Not yet addressed.

## core: Mutable package-level builtin type variables

- **Problem**: `core.String`, `core.Int`, `core.Bool`, `core.Float` are `var` (not `const`), so any consumer can mutate them (e.g., `core.String.Name = "oops"`), silently corrupting all downstream usage. Also a potential race condition if accessed from multiple goroutines.
- **Options**:
  - A: Convert to functions returning fresh copies (e.g., `func String() Type { ... }`). Breaking API change.
  - B: Keep as `var` but document the immutability contract.
- **Status**: Not yet addressed.

## generators/sqlc: `findTable`/`findTableForInputs` fragile name matching

- **Problem**: When a field is a builtin type, `findTable` searches all types for the first type containing a field with that name. If two types share a field name (e.g., both `item` and `order` have `id`), the wrong table may be selected.
- **Options**:
  - A: Require methods to explicitly specify the target table.
  - B: Carry the parent type context through from method definition.
- **Status**: Not yet addressed.

