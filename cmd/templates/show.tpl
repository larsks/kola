Name: {{ .Package.Name }}
Catalog source: {{ .Package.Status.CatalogSourceDisplayName }} ({{ .Package.Status.CatalogSource }})
Publisher: {{ .Package.Status.CatalogSourcePublisher }}
Provider: {{ .Package.Status.Provider.Name }}{{ if .Package.Status.Provider.URL }} ({{ .Package.Status.Provider.URL }}){{ end }}
Keywords:
{{ range $element := .Package.GetDefaultKeywords -}}
- {{ $element }}
{{ end -}}
Channels:
{{ range .Package.GetChannels -}}
- {{ .Name }} ({{ .CurrentCSV }})
{{ end -}}
Supported install modes:
{{ range $element := .Package.GetDefaultInstallModes -}}
- {{ $element }}
{{ end }}
{{- if (gt .Verbose 0) }}
Description:
{{ .Package.GetDefaultDescription }}
{{ end }}
