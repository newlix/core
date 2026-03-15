# Refactor Questions

Items requiring user decisions before proceeding.

## generators/sqlc: Missing WHERE clauses for input+output methods

- **Problem**: `writeQuery` uses a priority switch — outputs-only → SELECT, inputs-only → INSERT. Methods with *both* inputs and outputs (e.g., `RemoveItem`) generate `SELECT ... FROM table;` without a WHERE clause, silently dropping input parameters. The generated query returns all rows instead of filtering.
- **Options**:
  - A: Add a `writeSelectWithWhere` path for methods with both inputs and outputs.
  - B: Treat SQL generation as scaffolding only (document that users must add WHERE clauses manually).
  - C: Rethink the SQL generation model to support richer query semantics.
- **Status**: Not yet addressed.


