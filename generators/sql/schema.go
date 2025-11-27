package sql

import (
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/newlix/core"
)

type Dialect string

const (
	Cockroachdb Dialect = "cockroachdb"
	Sqlite              = "sqlite"
)

type GenerateSchemaFileConfig struct {
	Output  string
	Types   []core.Type
	Dialect Dialect
}

func GenerateSchemaFile(c GenerateSchemaFileConfig) {
	os.MkdirAll(path.Dir(c.Output), 0o700)
	w, err := os.Create(c.Output)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	GenerateSchema(w, c.Types, c.Dialect)
}

// Generate writes the Go type implementations to w.
func GenerateSchema(w io.Writer, tt []core.Type, dialect Dialect) {
	// types
	for _, t := range tt {
		ff := core.BuiltinTypeFields(t.Fields)
		if len(ff) == 0 {
			continue
		}
		out(w, "-- %s", t.Description)
		out(w, "CREATE TABLE IF NOT EXISTS %q (", t.Name)
		out(w, "  id   text PRIMARY KEY")
		out(w, ");")
		writeFields(w, t.Name, ff, dialect)
	}

}

// writeFields to writer.
func writeFields(w io.Writer, table string, ff []core.Field, dialect Dialect) {
	for _, f := range ff {
		if f.Name == "id" {
			continue
		}
		ttype := sqlType(f, dialect)
		out(w, "ALTER TABLE %q ADD COLUMN IF NOT EXISTS %q %s %s;", table, f.Name, ttype, sqlDefault(ttype))
		if !strings.HasPrefix(f.Type.GoType, "*") {
			out(w, "ALTER TABLE %q ALTER COLUMN %q SET NOT NULL;", table, f.Name)
		}
	}
}

func sqlType(f core.Field, dialect Dialect) string {
	switch dialect {
	case Cockroachdb:
		return f.Type.CockroachdbType
	case Sqlite:
		return f.Type.SqliteType
	default:
		log.Fatal("unhandled dialect: " + dialect)
	}
	return ""
}

func sqlDefault(ttype string) string {
	switch ttype {
	case core.String.CockroachdbType:
		return "DEFAULT ''"
	default:
		return ""
	}
}
