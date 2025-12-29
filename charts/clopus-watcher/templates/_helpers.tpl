{{- define "clopus-watcher.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "clopus-watcher.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "clopus-watcher.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | quote }}
app.kubernetes.io/name: {{ include "clopus-watcher.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end -}}

{{- define "clopus-watcher.selectorLabels" -}}
app.kubernetes.io/name: {{ include "clopus-watcher.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end -}}

{{- define "clopus-watcher.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (printf "%s-sa" (include "clopus-watcher.fullname" .)) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "clopus-watcher.secretName" -}}
{{- if .Values.apiKeySecret.existingSecret -}}
{{- .Values.apiKeySecret.existingSecret -}}
{{- else -}}
{{- .Values.apiKeySecret.name -}}
{{- end -}}
{{- end -}}

{{- define "clopus-watcher.pvcName" -}}
{{- if .Values.persistence.existingClaim -}}
{{- .Values.persistence.existingClaim -}}
{{- else -}}
{{- printf "%s-data" (include "clopus-watcher.fullname" .) -}}
{{- end -}}
{{- end -}}
