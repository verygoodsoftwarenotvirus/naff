package template_gen

import (
	"gitlab.com/verygoodsoftwarenotvirus/naff/forks/jennifer/jen"
	"gitlab.com/verygoodsoftwarenotvirus/naff/lib/utils"
	"gitlab.com/verygoodsoftwarenotvirus/naff/models"

	"os"
	"path/filepath"
)

const (
	packageName     = "main"
	basePackagePath = "cmd/tools/template_gen"
)

// RenderPackage renders the package
func RenderPackage(proj *models.Project) error {
	files := map[string]*jen.File{
		"main.go":            mainDotGo(proj),
		"table_configs.go":   tableConfigsDotGo(proj),
		"templates.go":       templatesDotGo(proj),
		"creator_configs.go": creatorConfigsDotGo(proj),
		"editor_configs.go":  editorConfigsDotGo(proj),
	}

	plainFiles := map[string]string{
		"templates/creator.gotpl": creatorTemplate,
		"templates/editor.gotpl":  editorTemplate,
		"templates/table.gotpl":   tableTemplate,
	}

	for path, file := range files {
		if err := utils.RenderGoFile(proj, filepath.Join(basePackagePath, path), file); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(utils.BuildTemplatePath(proj.OutputPath, filepath.Join(basePackagePath, "templates")), 0777); err != nil {
		return err
	}

	for path, file := range plainFiles {
		fp := utils.BuildTemplatePath(proj.OutputPath, filepath.Join(basePackagePath, path))

		if err := os.WriteFile(fp, []byte(file), 0644); err != nil {
			return err
		}
	}

	return nil
}

const creatorTemplate = `<div id="content" class="">
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">{{ .Title }}</h1>
    </div>
    <div class="col-md-8 order-md-1">
        <form class="needs-validation" novalidate=""{{ if ne .SubmissionURL "" }} hx-target="#content" hx-post="{{ .SubmissionURL }}"{{ end }}>{{ range $i, $field := .Fields }}
                <div class="mb3">
                    <label for="{{ $field.LabelName }}">{{ $field.StructFieldName }}</label>
                    <div class="input-group">
                        <input class="form-control" {{- if ne $field.InputType "" }} type="{{ $field.InputType }}"{{ end }} id="{{ $field.TagID }}" name="{{ $field.FormName }}" placeholder="{{ $field.InputPlaceholder }}" {{- if $field.Required }} required=""{{ end}} value="{{ print "{{ ." $field.StructFieldName " }}" }}" />
                        {{ if $field.Required }}<div class="invalid-feedback" style="width: 100%;">{{ $field.LabelName }} is required.</div>{{ end }}
                    </div>
                </div>{{ end }}
            <hr class="mb-4" />
            <button class="btn btn-primary btn-lg btn-block" type="submit">Save</button>
        </form>
    </div>
</div>`

const editorTemplate = `<div id="content" class="">
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">{{ print "{{ componentTitle . }}" }}</h1>
    </div>
    <div class="col-md-8 order-md-1">
        <form class="needs-validation" novalidate="" hx-target="#content" hx-put="{{ .SubmissionURL }}">{{ range $i, $field := .Fields }}
            <div class="mb3">
                <label for="{{ $field.LabelName }}">{{ $field.StructFieldName }}</label>
                <div class="input-group">
                    <input class="form-control" {{- if ne $field.InputType "" }} type="{{ $field.InputType }}"{{ end }} id="{{ $field.TagID }}" name="{{ $field.FormName }}" placeholder="{{ $field.InputPlaceholder }}" {{- if $field.Required }} required=""{{ end}} value="{{ print "{{ ." $field.StructFieldName " }}" }}" />
                    {{ if $field.Required }}<div class="invalid-feedback" style="width: 100%;">{{ $field.LabelName }} is required.</div>{{ end }}
                </div>
            </div>{{ end }}
            <hr class="mb-4" />
            <button class="btn btn-primary btn-lg btn-block" type="submit">Save</button>
        </form>
    </div>
</div>`

const tableTemplate = `<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    <h1 class="h2">{{ .Title }}</h1>
    {{ if .EnableSearch }}
        <div>
            Search: <input type="text">
            <button class="btn btn-secondary" hx-target="#content" hx-get="{{ .SearchURL }}">🔎</button>
            <button class="btn btn-primary" hx-target="#content" hx-trigger="keyup" hx-push-url="{{ .CreatorPagePushURL }}" hx-get="{{ .CreatorPageURL }}">New</button>
        </div>
    {{ else }}
    <button class="btn btn-primary" hx-target="#content" hx-push-url="{{ .CreatorPagePushURL }}" hx-get="{{ .CreatorPageURL }}">New</button>
    {{ end }}
</div>
<table class="table table-striped">
    <thead>
    <tr>{{ range $i, $col := .Columns }}
        <th>{{ $col }}</th>{{ end }}
    </tr>
    </thead>
    <tbody>{{ print "{{ range $i, $x := ." .RowDataFieldName " }}" }}
    <tr>
        {{ if not .ExcludeIDRow }}<td>{{ if .ExcludeLink }}{{ "{{ $x.ID }}" }}{{ else }}<button class="btn btn-sm btn-outline-dark" hx-push-url="{{ "{{ pushURL . }}" }}" hx-get="{{ "{{ individualURL . }}" }}" hx-target="#content">{{ "{{ $x.ID }}" }}</button>{{ end }}</td>{{ end }}
        {{ range $i, $x := .CellFields }}<td>{{ print "{{ $x." $x " }}" }}</td>
        {{ end }}{{ if .IncludeLastUpdatedOn }}<td>{{ "{{ relativeTimeFromPtr $x.LastUpdatedOn }}" }}</td>{{ end }}
        {{ if .IncludeCreatedOn }}<td>{{ "{{ relativeTime $x.CreatedOn }}" }}</td>{{ end }}
        {{ if .IncludeDeleteRow }}<td><button class="btn btn-sm btn-danger" hx-target="closest tr" hx-confirm="Are you sure you want to delete this?" hx-delete="{{ "{{ individualURL . }}" }}">Delete</button></td>{{ end }}
    </tr>
    {{ "{{ end }}" }}</tbody>
</table>`
