package golang

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/newlix/core"
	"github.com/newlix/core/generators/common"
)

type GenerateDatabaseFileConfig struct {
	Output  string
	Package string
	Types   []core.Type
}

func GenerateDatabaseFile(c GenerateDatabaseFileConfig) {
	os.MkdirAll(path.Dir(c.Output), 0o700) // Create your file
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()
	common.GenerateWarning(w)
	out(w, "package %s", PackageName(c.Package))
	out(w, `import (
	"context"
	"database/sql"
)

type QueryerContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type PreparerContext interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type ExecerContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

`)
	GenerateImports(w, c.Package, c.Types)

	GenerateDatabase(w, c.Package, c.Types)
}

// Generate writes the Go type implementations to w.
func GenerateDatabase(w io.Writer, pkg string, ts []core.Type) {
	// types
	for _, t := range ts {
		fs := core.BuiltinTypeFields(t.Fields)
		if len(fs) == 0 {
			continue
		}
		generateSelect(w, pkg, t)
		generateInsert(w, pkg, t)
		generateUpdate(w, pkg, t)
		generateUpsert(w, pkg, t)
		generateDelete(w, t)
		generateCount(w, t)
	}
}

func generateSelect(w io.Writer, pkg string, t core.Type) {
	fs := core.BuiltinTypeFields(t.Fields)
	sel := ""
	for i, f := range fs {
		sel += fmt.Sprintf("%q", f.Name)
		if i < len(fs)-1 {
			sel += ", "
		}
	}
	out(w, "const %sSelectQuery = `SELECT %s FROM %q `", t.Name, sel, t.Name)
	out(w, "")

	out(w, "func %sSelectFirst(ctx context.Context, q QueryerContext, extraQuery string, args ...any) (%s, error) {", t.CamelName, TypeGoType(pkg, t))
	out(w, "	row := q.QueryRowContext(ctx, %sSelectQuery+extraQuery, args...)", t.Name)
	out(w, "	var o %s", TypeGoType(pkg, t))
	out(w, "	err := row.Scan(")
	for _, f := range fs {
		out(w, "		&o.%s,", f.CamelName)
	}
	out(w, "	)")
	out(w, "	return o, err")
	out(w, "}")
	out(w, "")

	out(w, "func %sSelect(ctx context.Context, q QueryerContext, extraQuery string, args ...any) ([]%s, error) {", t.CamelName, TypeGoType(pkg, t))
	out(w, "	rows, err := q.QueryContext(ctx, %sSelectQuery+extraQuery, args...)", t.Name)
	out(w, "	if err != nil {")
	out(w, "		return nil, err")
	out(w, "	}")
	out(w, "	defer rows.Close()")
	out(w, "	var objs []%s", TypeGoType(pkg, t))
	out(w, "	for rows.Next() {")
	out(w, "		var obj %s", TypeGoType(pkg, t))
	out(w, "		if err := rows.Scan(")
	for _, f := range fs {
		out(w, "			&obj.%s,", f.CamelName)
	}
	out(w, "		); err != nil {")
	out(w, "			return nil, err")
	out(w, "		}")
	out(w, "		objs = append(objs, obj)")
	out(w, "	}")
	out(w, "	if err := rows.Err(); err != nil {")
	out(w, "		return nil, err")
	out(w, "	}")
	out(w, "	return objs, nil")
	out(w, "}")
	out(w, "")

	out(w, "func %sSelectWithID(ctx context.Context, q QueryerContext, id string) (%s, error) {", t.CamelName, TypeGoType(pkg, t))
	out(w, "	return %sSelectFirst(ctx, q, `WHERE \"id\" = $1`, id)", t.CamelName)
	out(w, "}")
	out(w, "")

	out(w, "func %sSelectAll(ctx context.Context, q QueryerContext) ([]%s, error) {", t.CamelName, TypeGoType(pkg, t))
	out(w, "	return %sSelect(ctx, q, \"\")", t.CamelName)
	out(w, "}")
	out(w, "")
}

func generateInsert(w io.Writer, pkg string, t core.Type) {
	fs := core.BuiltinTypeFields(t.Fields)
	fmt.Fprintf(w, "const %sInsertQuery = `INSERT INTO %q (", t.Name, t.Name)
	for i, f := range fs {
		fmt.Fprintf(w, "%q", f.Name)
		if i < len(fs)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	fmt.Fprintf(w, ") ")
	fmt.Fprintf(w, "VALUES (")
	for i := range fs {
		fmt.Fprintf(w, "$%d", i+1)
		if i < len(fs)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	out(w, ");`")
	out(w, "")
	out(w, "func %sInsert(ctx context.Context, e ExecerContext, obj %s) error {", t.CamelName, TypeGoType(pkg, t))
	out(w, "	_, err := e.ExecContext(ctx, %sInsertQuery,", t.Name)
	for _, f := range fs {
		out(w, "		obj.%s,", f.CamelName)
	}
	out(w, "	)")
	out(w, "	return err")
	out(w, "}")
	out(w, "")
}

func generateUpdate(w io.Writer, pkg string, t core.Type) {
	fs := core.BuiltinTypeFields(t.Fields)
	var assigns []string
	i := 2
	for _, f := range fs {
		if f.Name != "id" {
			assigns = append(assigns, fmt.Sprintf("%q = $%d", f.Name, i))
			i += 1
		}
	}
	out(w, "const %sUpdateQuery = `UPDATE %q", t.Name, t.Name)
	out(w, "SET %s", strings.Join(assigns, ", "))
	out(w, "WHERE id = $1;`")
	out(w, "")
	out(w, "func %sUpdate(ctx context.Context, e ExecerContext, obj %s) error {", t.CamelName, TypeGoType(pkg, t))
	out(w, "	_, err := e.ExecContext(ctx, %sUpdateQuery,", t.Name)
	out(w, "		obj.ID,")
	for _, f := range fs {
		if f.Name == "id" {
			continue
		}
		out(w, "		obj.%s,", f.CamelName)
	}
	out(w, "	)")
	out(w, "	return err")
	out(w, "}")
	out(w, "")
}

func generateUpsert(w io.Writer, pkg string, t core.Type) {
	fs := core.BuiltinTypeFields(t.Fields)
	fmt.Fprintf(w, "const %sUpsertQuery = `INSERT INTO %q (", t.Name, t.Name)
	for i, f := range fs {
		fmt.Fprintf(w, "%q", f.Name)
		if i < len(fs)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	fmt.Fprintf(w, ") ")
	fmt.Fprintf(w, "VALUES (")

	for i := range fs {
		fmt.Fprintf(w, "$%d", i+1)
		if i < len(fs)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	fmt.Fprintf(w, ")\n")
	out(w, "ON CONFLICT (\"id\")")
	fmt.Fprintf(w, "DO UPDATE SET ")
	for i, f := range fs {
		fmt.Fprintf(w, "%q = $%d", f.Name, i+1)
		if i < len(fs)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	out(w, ";`")
	out(w, "")
	out(w, "func %sUpsert(ctx context.Context, e ExecerContext, obj %s) error {", t.CamelName, TypeGoType(pkg, t))
	out(w, "	_, err := e.ExecContext(ctx, %sUpsertQuery,", t.Name)
	for _, f := range fs {
		out(w, "		obj.%s,", f.CamelName)
	}
	out(w, "	)")
	out(w, "	return err")
	out(w, "}")
	out(w, "")
}

func generateDelete(w io.Writer, t core.Type) {
	out(w, "const %sDeleteQuery = `DELETE FROM %q WHERE \"id\" = $1;`", t.Name, t.Name)
	out(w, "")
	out(w, "func %sDelete(ctx context.Context, e ExecerContext, id string) error {", t.CamelName)
	out(w, "	_, err := e.ExecContext(ctx, %sDeleteQuery, id)", t.Name)
	out(w, "	return err")
	out(w, "}")
	out(w, "")
}

func generateCount(w io.Writer, t core.Type) {
	out(w, "const %sCountQuery = \"SELECT COUNT(*) FROM %s \"", t.Name, t.Name)
	out(w, "")
	out(w, "func %sCount(ctx context.Context, q QueryerContext, extraQuery string, args ...any) (int, error) {", t.CamelName)
	out(w, "	row := q.QueryRowContext(ctx, %sCountQuery+extraQuery, args...)", t.Name)
	out(w, "	var c int")
	out(w, "	err := row.Scan(&c)")
	out(w, "	return c, err")
	out(w, "}")
}
