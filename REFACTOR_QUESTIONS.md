# Refactor Questions

Items requiring user decisions before proceeding.

## generators/sqlc: Missing WHERE clauses for input+output methods

- **Problem**: `writeQuery` uses a priority switch — outputs-only → SELECT, inputs-only → INSERT. Methods with *both* inputs and outputs (e.g., `RemoveItem`) generate `SELECT ... FROM table;` without a WHERE clause, silently dropping input parameters. The generated query returns all rows instead of filtering.
- **Options**:
  - A: Add a `writeSelectWithWhere` path for methods with both inputs and outputs.
  - B: Treat SQL generation as scaffolding only (document that users must add WHERE clauses manually).
  - C: Rethink the SQL generation model to support richer query semantics.
- **Status**: Not yet addressed.


## generators/sqlc: `findTable`/`findTableForInputs` fragile name matching

- **Problem**: When a field is a builtin type, `findTable` searches all types for the first type containing a field with that name. If two types share a field name (e.g., both `item` and `order` have `id`), the wrong table may be selected.
- **Options**:
  - A: Require methods to explicitly specify the target table.
  - B: Carry the parent type context through from method definition.
- **Status**: Not yet addressed.

