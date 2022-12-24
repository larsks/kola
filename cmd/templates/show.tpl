Name: {{ .Package.Name }}
Catalog source: {{ .Package.Status.CatalogSourceDisplayName }} ({{ .Package.Status.CatalogSource }})
Publisher: {{ .Package.Status.CatalogSourcePublisher }}
Provider: {{ .Package.Status.Provider.Name }}{{ if .Package.Status.Provider.URL }} ({{ .Package.Status.Provider.URL }}){{ end }}
Keywords:
{{ range $index, $element := (index .Package.Status.Channels 0).CurrentCSVDesc.Keywords -}}
- {{ $element }}
{{ end -}}
Channels:
{{ range .Package.Status.Channels -}}
- {{ .Name }} ({{ .CurrentCSV }})
{{ end -}}
Supported install modes:
{{ range $index, $element := (index .Package.Status.Channels 0).CurrentCSVDesc.InstallModes }}
{{- if $element.Supported -}}
- {{ $element.Type }}
{{ end -}}
{{ end -}}
{{- if (gt .Verbose 0) -}}
Description:
{{ (index .Package.Status.Channels 0).CurrentCSVDesc.LongDescription }}
{{ end }}
