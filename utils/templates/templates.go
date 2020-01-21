package templatesUtils

import (
	"github.com/flosch/pongo2"
	"path"
)

func Render(file string, ctx pongo2.Context) string {

	fullPath := path.Join("templates", file)
	engine, err := pongo2.FromFile(fullPath)
	tpl := pongo2.Must(engine, err)
	outData, err := tpl.Execute(ctx)
	if err != nil {

		return ""
	}
	return outData
}
