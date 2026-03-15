# Refactor Questions

## [generators/common/out.go]
- **Problem**: `common.Out()` (and all generator write calls through it) ignores `fmt.Fprintf` errors. If the underlying `io.Writer` fails mid-generation (e.g., disk full), the error is silently lost. `GenerateFile` will report success because the `fn` callback returns `nil`.
- **Options**:
  - A: Change `Out()` to return error, update all ~100 call sites to check it. Correct but high churn.
  - B: Use an `errWriter` wrapper that captures the first write error, check it once at the end in `GenerateFile`. Low churn, still catches failures.
  - C: Accept current behavior — write errors to `*os.File` are rare in practice, and generated code is always verified by compilation.
- **Status**: Not addressed (skipped)

## [generators/golang/client.go]
- **Problem**: `io.WriteString(w, ...)` on line 41 ignores the error return. Same root cause as the `Out()` issue above — generator functions don't propagate write errors.
- **Options**: Same as above; fixing this one call without fixing `Out()` is inconsistent.
- **Status**: Not addressed (skipped, coupled to the Out() decision)
