package composite

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kirk-anchor/pggen"
	"github.com/kirk-anchor/pggen/internal/pgtest"
	"github.com/stretchr/testify/assert"
)

func TestGenerate_Go_Example_Composite(t *testing.T) {
	conn, cleanupFunc := pgtest.NewPostgresSchema(t, []string{"schema.sql"})
	defer cleanupFunc()

	tmpDir := t.TempDir()
	err := pggen.Generate(
		pggen.GenerateOptions{
			ConnString:       conn.Config().ConnString(),
			QueryFiles:       []string{"query.sql"},
			OutputDir:        tmpDir,
			GoPackage:        "composite",
			Language:         pggen.LangGo,
			InlineParamCount: 2,
			TypeOverrides: map[string]string{
				"_bool":  "[]bool",
				"bool":   "bool",
				"int8":   "int",
				"int4":   "int",
				"text":   "string",
				"citext": "github.com/jackc/pgx/v5/pgtype.Text",
			},
		})
	if err != nil {
		t.Fatalf("Generate(): %s", err)
	}

	wantQueryFile := "query.sql.go"
	gotQueryFile := filepath.Join(tmpDir, "query.sql.go")
	assert.FileExists(t, gotQueryFile,
		"Generate() should emit query.sql.go")
	wantQueries, err := os.ReadFile(wantQueryFile)
	if err != nil {
		t.Fatalf("read wanted query.go.sql: %s", err)
	}
	gotQueries, err := os.ReadFile(gotQueryFile)
	if err != nil {
		t.Fatalf("read generated query.go.sql: %s", err)
	}
	assert.Equalf(t, string(wantQueries), string(gotQueries),
		"Got file %s; does not match contents of %s",
		gotQueryFile, wantQueryFile)
}
