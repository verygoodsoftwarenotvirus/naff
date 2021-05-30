package models

import (
	"bytes"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
	"path/filepath"
	"strings"
	"text/template"
)

type CodeFile struct {
	proj *Project
}

func (cf *CodeFile) TemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"clean":     strings.TrimSpace,
		"lowercase": strings.ToLower,
		"uppercase": strings.ToUpper,
		"kebabFmt": func(s string) string {
			return wordsmith.FromSingularPascalCase(s).KebabName()
		},
		"route_fmt": func(s string) string {
			return wordsmith.FromSingularPascalCase(s).RouteName()
		},
		"projectImport": func(path string) string {
			return filepath.Join(append([]string{cf.proj.OutputPath}, path)...)
		},
		"projectName": func(subsequentDirectories ...string) string {
			return cf.proj.Name.Singular()
		},
		"outputPath": func() string {
			return cf.proj.OutputPath
		},
	}
}

func RenderCodeFile(proj *Project, rawTemplate string) string {
	cf := &CodeFile{
		proj: proj,
	}

	tmpl := template.Must(template.New("").Funcs(cf.TemplateFunctions()).Parse(rawTemplate))

	var b bytes.Buffer
	if err := tmpl.Execute(&b, cf); err != nil {
		panic(err)
	}

	return b.String()
}
