package models

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/wordsmith"
)

type CodeFile struct {
	proj      *Project
	generated map[string]string
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
		"unexportedVarName": func(s string) string {
			return wordsmith.FromSingularPascalCase(s).UnexportedVarName()
		},
		"dbConditional": func(dbName, content string) string {
			if cf.proj.DatabaseIsEnabled(validDatabase(dbName)) {
				return content
			}
			return ""
		},
		"projectImport": func(path string) string {
			return fmt.Sprintf("%q", filepath.Join(append([]string{cf.proj.OutputPath}, path)...))
		},
		"projectName": func(subsequentDirectories ...string) string {
			return cf.proj.Name.Singular()
		},
		"outputPath": func() string {
			return cf.proj.OutputPath
		},
		"here": func() string {
			return fmt.Sprintf("%s/", cf.proj.OutputPath)
		},
		"generated": func(s string) string {
			return cf.generated[s]
		},
		"typeImports": func() string {
			imports := []string{}

			for _, typ := range cf.proj.DataTypes {
				imports = append(imports, fmt.Sprintf("\t%q", cf.proj.ServicePackage(typ.Name.PackageName())))
			}

			return strings.Join(imports, "\n")
		},
	}
}

func RenderCodeFile(proj *Project, rawTemplate string, generated map[string]string) string {
	cf := &CodeFile{
		proj:      proj,
		generated: generated,
	}

	tmpl := template.Must(template.New("").Funcs(cf.TemplateFunctions()).Parse(rawTemplate))

	var b bytes.Buffer
	if err := tmpl.Execute(&b, cf); err != nil {
		panic(err)
	}

	return b.String()
}
