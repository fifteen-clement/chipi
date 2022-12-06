package gen

import (
	"bytes"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func data(desc string) map[string]string {
	return map[string]string{
		"description": desc,
	}
}

func dataex(desc string, ex string) map[string]string {
	return map[string]string{
		"description": desc,
		"example":     ex,
	}
}

func TestGenerator(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("Generator", func() {
		var fs *token.FileSet
		var f *dst.File

		g.BeforeEach(func() {
			var err error

			data, err := os.ReadFile("../internal/testdata/monster/monster.go")
			require.NoError(g, err)

			fs = token.NewFileSet()
			f, err = decorator.ParseFile(fs, "", data, parser.ParseComments)
			require.NoError(g, err)
		})

		g.Describe("Inspect", func() {

			g.It("should invoke callbacks for all documented fields", func() {
				pos := 0
				expected := []struct {
					parent  string
					section string
					field   string
					data    map[string]string
				}{
					{"GetMonsterRequest", "Operation", "", map[string]string{
						"tag":        "monster",
						"deprecated": "",
						"summary":    "Grab a monster and bring it to you\nknowing its Id",
					}},
					{"GetMonsterRequest", "Path", "Id", data("The _Id_ of the monster you want to\nfetch")},
					{"GetMonsterRequest", "Query", "", data("the query")},
					{"GetMonsterRequest", "Query", "Blocking", dataex("If true the request will block until\nthe monster was actually created", "ahhhhhh !")},
					{"GetMonsterRequest", "Header", "ApiKey", data("The _ApiKey_ is required to\ncheck for authorization")},
					{"GetMonsterRequest", "Header", "Something", data("This may be important")},
					{"GetMonsterRequest", "Response", "", data("what is returned")},
				}

				err := inspectStructures(f, func(parentStructName string, sectionName string, fieldName string, data map[string]string) error {
					if pos >= len(expected) {
						g.Failf("expected too short (index: %d)", pos)
					}
					dd := expected[pos]
					assert.Equal(g, dd.parent, parentStructName)
					assert.Equal(g, dd.section, sectionName)
					assert.Equal(g, dd.field, fieldName)
					assert.Equal(g, dd.data, data)

					pos++
					return nil
				})

				require.NoError(g, err)
				assert.Equal(g, 7, pos)

			})
		})

		g.Describe("GenerateAnnotations", func() {
			g.It("should generate annotations", func() {
				buffer := bytes.NewBufferString("")
				err := GenerateFieldAnnotations(buffer, f, "monster")
				require.NoError(g, err)

				// TODO: test the content
			})
		})

	})
}
