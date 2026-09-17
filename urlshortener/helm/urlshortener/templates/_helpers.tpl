{{/* Labels every object carries in metadata. */}}
{{- define "urlshortener.labels" -}}
{{- range $k, $v := .Values.commonLabels }}
{{ $k }}: {{ $v | quote }}
{{- end }}
{{- end }}
