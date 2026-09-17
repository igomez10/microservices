{{/* Labels every object carries in metadata. */}}
{{- define "socialapp.labels" -}}
{{- range $k, $v := .Values.commonLabels }}
{{ $k }}: {{ $v | quote }}
{{- end }}
{{- end }}

{{/* In-cluster base URL of the API Service. */}}
{{- define "socialapp.apiURL" -}}
http://socialapp.{{ .Release.Namespace }}.svc.cluster.local:{{ .Values.api.port }}
{{- end }}
