{{- define "slvlib.serviceaccount" -}}
{{- if eq .Values.serviceAccountName "" -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: slv-serviceaccount
  namespace: {{ .Release.Namespace }}
automountServiceAccountToken: true
{{- end -}}
{{- end -}}
