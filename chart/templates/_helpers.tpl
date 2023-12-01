{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "cert-manager-porkbun-webhook.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cert-manager-porkbun-webhook.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cert-manager-porkbun-webhook.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "cert-manager-porkbun-webhook.selfSignedIssuer" -}}
{{ printf "%s-selfsign" (include "cert-manager-porkbun-webhook.fullname" .) }}
{{- end -}}

{{- define "cert-manager-porkbun-webhook.rootCAIssuer" -}}
{{ printf "%s-ca" (include "cert-manager-porkbun-webhook.fullname" .) }}
{{- end -}}

{{- define "cert-manager-porkbun-webhook.rootCACertificate" -}}
{{ printf "%s-ca" (include "cert-manager-porkbun-webhook.fullname" .) }}
{{- end -}}

{{- define "cert-manager-porkbun-webhook.servingCertificate" -}}
{{ printf "%s-webhook-tls" (include "cert-manager-porkbun-webhook.fullname" .) }}
{{- end -}}


{{/* vim: set filetype=mustache: */}}
{{/*
Renders a value that contains template.
Usage:
{{ include "common.tplvalues.render" ( dict "value" .Values.path.to.the.Value "context" $) }}
*/}}
{{- define "common.tplvalues.render" -}}
    {{- if typeIs "string" .value }}
        {{- tpl .value .context }}
    {{- else }}
        {{- tpl (.value | toYaml) .context }}
    {{- end }}
{{- end -}}
