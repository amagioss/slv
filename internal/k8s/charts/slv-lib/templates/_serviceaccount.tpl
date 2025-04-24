{{- define "slvlib.serviceaccount" -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .Values.config.serviceAccountName | default "slv-serviceaccount"}}
  namespace: {{ .Release.Namespace }}
automountServiceAccountToken: true
{{- end -}}
