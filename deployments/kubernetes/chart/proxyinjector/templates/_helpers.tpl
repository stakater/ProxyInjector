{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "proxyinjector-name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" | lower -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "proxyinjector-fullname" -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "proxyinjector-labels.selector" -}}
app: {{ template "proxyinjector-name" . }}
group: {{ .Values.proxyinjector.labels.group }}
provider: {{ .Values.proxyinjector.labels.provider }}
{{- end -}}

{{- define "proxyinjector-labels.stakater" -}}
{{ template "proxyinjector-labels.selector" . }}
version: {{ .Values.proxyinjector.labels.version }}
{{- end -}}

{{- define "proxyinjector-labels.chart" -}}
chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
release: {{ .Release.Name | quote }}
heritage: {{ .Release.Service | quote }}
{{- end -}}

{{- define "proxyinjector-vol-config-name" -}}
{{- if .Values.proxyinjector.existingSecret -}}
{{ .Values.proxyinjector.existingSecret }}
{{- else -}}
{{- template "proxyinjector-name" . -}}
{{- end -}}
{{- end -}}